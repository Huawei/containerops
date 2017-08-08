package main

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	model "github.com/Huawei/containerops/singular/model"
	log "github.com/Sirupsen/logrus"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// A deployment has one related infrastructure and multiple installed components
type Deployment struct {
	ID          int64
	InfraName   string
	InfraLogo   string
	InfraLog    template.HTML
	Components  []Component
	StatusIcon  string
	StatusColor string
}
type Component struct {
	Name          string
	Log           template.HTML
	ImageSrc, Alt string
	Width, Height int
}

var cfgFile string

func main() {
	// Init cmd
	cmd := &cobra.Command{
		Use:   "singular_website",
		Short: "singular_website is the tool to serve the website that shows singular's deployments",
		Run:   runCmd,
	}

	cobra.OnInitialize(initConfig)

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Configuration file path")
	viper.BindPFlag("config", cmd.Flags().Lookup("config"))

	if err := cmd.Execute(); err != nil {
		log.Fatalf("Failed to run command: %s", err.Error())
		os.Exit(1)
	}

}

func runCmd(cmd *cobra.Command, args []string) {
	// if err := cmd.ParseFlags(args); err != nil {
	//     log.Fatalln("Failed to parse flags: %s", err.Error())
	//     os.Exit(1)
	// }
}

func initConfig() {
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	dbConfig := viper.GetStringMapString("database")
	webConfig := viper.GetStringMapString("web")

	driver := dbConfig["driver"]
	host := dbConfig["host"]
	port := dbConfig["port"]
	user := dbConfig["user"]
	password := dbConfig["password"]
	db := dbConfig["db"]

	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True&loc=Local", user, password, host, port, db)
	if DB, err = gorm.Open(driver, dsn); err != nil {
		log.Fatal("Initlization database connection error.", err)
		os.Exit(1)
	} else {
		DB.DB()
		DB.DB().Ping()
		DB.DB().SetMaxIdleConns(10)
		DB.DB().SetMaxOpenConns(100)
		DB.SingularTable(true)
		// DB.LogMode(true)
	}

	var server *http.Server

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("public"))
	mux.Handle("/public/", http.StripPrefix("/public/", fs))

	mux.HandleFunc("/", ServeListTemplate)
	mux.HandleFunc("/detail", ServeDetailTemplate)

	webMode := webConfig["mode"]
	address := webConfig["address"]

	switch webMode {

	case "https":
		port, _ := strconv.Atoi(webConfig["port"])
		cert := webConfig["cert"]
		key := webConfig["key"]
		listenAddr := fmt.Sprintf("%s:%d", address, port)
		server = &http.Server{Addr: listenAddr, TLSConfig: &tls.Config{MinVersion: tls.VersionTLS10}, Handler: mux}
		if err := server.ListenAndServeTLS(cert, key); err != nil {
			log.Fatalf("Start Dockyard https service error: %s\n", err.Error())
		}

		break
	case "unix":
		listenAddr := fmt.Sprintf("%s", address)

		_, err := os.Stat(listenAddr)
		if err == nil || os.IsExist(err) {
			os.Remove(listenAddr)
		}

		fmt.Println("listening on ", listenAddr)
		if listener, err := net.Listen("unix", listenAddr); err != nil {
			log.Fatalf("Start Dockyard unix socket error: %s\n", err.Error())
		} else {
			server = &http.Server{Handler: mux}
			if err := server.Serve(listener); err != nil {
				log.Fatalf("Start Dockyard unix socket error: %s\n", err.Error())
			}
		}
		break
	default:
		log.Fatalf("Invalid listen mode: %s\n", webMode)
		os.Exit(1)
		break
	}
}

func ServeListTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid method"))
		return
	}

	funcs := template.FuncMap{
		"component_names": stringifyComponentsNames,
	}

	deployments, err := getDeploymentList()
	if err != nil {
		log.Errorf("Failed to get deployment list: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// The template's name should be the same with file name
	listTmpl, err := template.New("list.template").Funcs(funcs).ParseFiles("./templates/list.template")
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = listTmpl.Execute(w, deployments)
	if err != nil {
		log.Error(err)
		// w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func ServeDetailTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid method"))
		return
	}
	// Get the deployment information
	queryArgs := r.URL.Query()
	deployment_ids := queryArgs["deployment_id"]
	if len(deployment_ids) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid argument of deployment_id in query string"))
		return
	}

	deployment_id := deployment_ids[0]
	deploymentID, _ := strconv.Atoi(deployment_id)
	deployment := getDeploymentDetail(deploymentID)
	if deployment == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Deployment not found"))
		return
	}

	// The template's name should be the same with file name
	listTmpl, err := template.New("detail.template").ParseFiles("./templates/detail.template")
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = listTmpl.Execute(w, deployment)
	if err != nil {
		log.Error(err)
		// w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func getDeploymentDetail(deploymentID int) *Deployment {
	var dep model.DeploymentV1
	var infra model.InfraV1
	var comps []model.ComponentV1

	// Get the infra and components
	err := DB.Where("id=?", deploymentID).First(&dep).Error
	if err != nil && err.Error() == "record not found" {
		return nil
	}
	DB.Where("deployment_v1=?", dep.ID).First(&infra)
	DB.Where("infra_v1=?", infra.ID).Find(&comps)

	components := []Component{}
	for i := 0; i < len(comps); i++ {
		components = append(components, convertComponent(&comps[i]))
	}
	deployment := Deployment{
		InfraName:  getInfraName(infra.Name),
		InfraLogo:  getInfraLogo(infra.Name),
		InfraLog:   convertToBr(infra.Log),
		Components: components,
	}

	return &deployment
}

var infraNames map[string]string = map[string]string{
	"digital_ocean": "Digital Ocean",
	"gke":           "Google Container Engine",
	"aws":           "AWS EC2",
	"azure":         "Microsoft Azure",
}
var infraLogos map[string]string = map[string]string{
	"digital_ocean": "./public/icons/digital-ocean.svg",
	"gke":           "./public/icons/google-cloud.svg",
	"aws":           "./public/icons/aws-ec2.svg",
	"azure":         "./public/icons/azure.svg",
}

func getInfraName(key string) string {
	name := infraNames[key]
	if name == "" {
		name = "Unsupported Cloud Service Provider"
	}
	return name
}

func getInfraLogo(name string) string {
	logo := infraLogos[name]
	if logo == "" {
		logo = "./public/icons/google-cloud.svg"
	}
	return logo
}

// Convert \n to <br />
func convertToBr(src string) template.HTML {
	replaced := strings.Replace(src, "\n", "\n<br />", -1)
	return template.HTML(replaced)
	// return strings.Join(strings.Split(src, "\n"), "<br />")
}

func stringifyComponentsNames(args ...interface{}) (string, error) {
	v := reflect.ValueOf(args)
	numArgs := v.Len()
	if numArgs != 1 {
		fmt.Fprintf(os.Stderr, "component_names function expect 1 argument, but got %d", numArgs)
		return "", fmt.Errorf("Expect 1 argument")
	}

	if components, ok := v.Index(0).Interface().([]Component); !ok {
		fmt.Fprintln(os.Stderr, "function component_names receives an argument which is not []Component")
		return "", fmt.Errorf("Argument is not []Component!")
	} else {
		s := ""
		for i := 0; i < len(components); i++ {
			c := components[i]
			s += c.Name + ", "
		}
		return s[:len(s)-1], nil
	}
}

var DB *gorm.DB

func convertComponent(input *model.ComponentV1) Component {
	c := Component{
		Name: input.Name,
		Log:  convertToBr(input.Log),
	}
	// Got the size and icon
	switch c.Name {
	case "kubernetes":
		c.Width = 40
		c.Height = 40
		c.ImageSrc = "./public/icons/kubernetes.svg"
		break
	case "docker":
		c.Width = 40
		c.Height = 40
		c.ImageSrc = "./public/icons/docker.svg"
		break
	case "flannel":
		c.Width = 40
		c.Height = 40
		c.ImageSrc = "./public/icons/flannel.svg"
		break
	case "etcd":
		c.Width = 40
		c.Height = 40
		c.ImageSrc = "./public/icons/etcd.svg"
		break
	default:
		break
	}

	return c
}

func getDeploymentList() ([]Deployment, error) {
	var deployments []model.DeploymentV1
	err := DB.Find(&deployments).Error
	if err != nil {
		return nil, err
	}

	var retDeployments []Deployment
	// Find the corresponding infrastructure & components info
	sig := make(chan error)
	for i := 0; i < len(deployments); i++ {
		go func(deployment *model.DeploymentV1) {
			var infra model.InfraV1
			DB.Where("deployment_v1=?", deployment.ID).First(&infra)
			var components []model.ComponentV1
			DB.Where("infra_v1=?", infra.ID).Find(&components)

			statusIcon, statusColor := "clear", "red"
			if deployment.Result == true {
				statusIcon, statusColor = "check", "green"
			}
			retDeployment := Deployment{
				// Name: deployment.Name,
				ID:          deployment.ID,
				InfraName:   infra.Name,
				StatusIcon:  statusIcon,
				StatusColor: statusColor,
			}
			retComponents := []Component{}

			for j := 0; j < len(components); j++ {
				retComponents = append(retComponents, convertComponent(&components[j]))
			}

			retDeployment.Components = retComponents

			retDeployments = append(retDeployments, retDeployment)

			sig <- nil
		}(&deployments[i])
	}

	for i := 0; i < len(deployments); i++ {
		if err := <-sig; err != nil {
			return nil, err
		}
	}
	return retDeployments, nil
}
