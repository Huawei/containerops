package unittest

// //	initsystem "github.com/Huawei/containerops_he3io/singular/initsystem"

// func Test() {
// 	for k, item := range nodes {

// 		fmt.Printf("k=%v, item=%v item[1]=%v \n", k, item[0], item[1])

// 		// if item.kes == "centos-master" {
// 		// 	nodes.Deploymaster(nodes, item[0][1])
// 		// }
// 		// if item[1][1] == "centos-minion" {
// 		// 	nodes.Deploynode(nodes, item[0][1])
// 		// }
// 	}

// 	// for item := range m {

// 	// 	if item.kes == "centos-master" {
// 	// 		nodes.Deploymaster(nodes, item[0][1])
// 	// 	}
// 	// 	if item[1][1] == "centos-minion" {
// 	// 		nodes.Deploynode(nodes, item[0][1])
// 	// 	}
// 	// }

// 	//1
// 	// command := "/usr/bin/systemctl stop etcd"
// 	// params := []string{"-i"}
// 	// //执行cmd命令: ls -l
// 	// execCommand(command, params)

// 	//2
// 	// path, err := exec.LookPath("systemctl")
// 	// if err != nil {
// 	// 	log.Fatal("installing fortune is in your future %q", err)
// 	// }
// 	// fmt.Printf("fortune is available at %s\n", path)

// 	//3
// 	// cmd := exec.Command("tr", "a-z", "A-Z")
// 	// cmd := exec.Command("/usr/bin/systemct", " stop", " etcd")

// 	// //cmd.Stdin = strings.NewReader("some input")
// 	// var out1 bytes.Buffer
// 	// cmd.Stdout = &out1
// 	// err := cmd.Run()
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// fmt.Printf("in all caps: %q\n", out1.String())

// 	//initsystem.ServiceStart("")
// 	//4
// 	// var ErrNotFound = errors.New("executable file not found in $PATH")
// 	// fmt.Printf("2 in all caps: %q\n", ErrNotFound)

// 	//5
// 	// cmd := exec.Command("sleep", "5")
// 	// err := cmd.Start()
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// log.Printf("Waiting for command to finish...")
// 	// err = cmd.Wait()
// 	// log.Printf("Command finished with error: %v", err)

// 	//6
// 	//	initsystem.Test()
// 	// a := initsystem.SystemdInitSystem{}
// 	// a.ServiceStart("etcd")

// 	//ServiceStart("etcd")
// }

// // type SystemdInitSystem struct{}

// // // getInitSystem returns an InitSystem for the current system, or nil
// // // if we cannot detect a supported init system for pre-flight checks.
// // // This indicates we will skip init system checks, not an error.
// // func GetInitSystem() (InitSystem, error) {
// // 	// Assume existence of systemctl in path implies this is a systemd system:
// // 	_, err := exec.LookPath("systemctl")
// // 	if err == nil {
// // 		return &SystemdInitSystem{}, nil
// // 	}
// // 	return nil, fmt.Errorf("no supported init system detected, skipping checking for services")
// // }

// // type InitSystem interface {
// // 	// ServiceStart tries to start a specific service
// // 	ServiceStart(service string) error

// // 	// ServiceStop tries to stop a specific service
// // 	ServiceStop(service string) error

// // 	// ServiceExists ensures the service is defined for this init system.
// // 	ServiceExists(service string) bool

// // 	// ServiceIsEnabled ensures the service is enabled to start on each boot.
// // 	ServiceIsEnabled(service string) bool

// // 	// ServiceIsActive ensures the service is running, or attempting to run. (crash looping in the case of kubelet)
// // 	ServiceIsActive(service string) bool
// // }

// // func execCommand(commandName string, params []string) bool {
// // 	cmd := exec.Command(commandName, params...)

// // 	//显示运行的命令
// // 	fmt.Println(cmd.Args)

// // 	stdout, err := cmd.StdoutPipe()

// // 	if err != nil {
// // 		fmt.Println(err)
// // 		return false
// // 	}

// // 	cmd.Start()

// // 	reader := bufio.NewReader(stdout)

// // 	//实时循环读取输出流中的一行内容
// // 	for {
// // 		line, err2 := reader.ReadString('\n')
// // 		if err2 != nil || io.EOF == err2 {
// // 			break
// // 		}
// // 		fmt.Println(line)
// // 	}

// // 	cmd.Wait()
// // 	return true
// // }

// func restartSvc(svcArr []string) error {
// 	for _, svc := range svcArr {
// 		args := []string{"restart", svc}
// 		cmd.execCMDparams(svc, args)
// 		args = []string{"enable", svc}
// 		cmd.execCMDparams(svc, args)
// 	}
// 	args := []string{"daemon-reload"}
// 	_, err := exec.Command("systemctl", args...).Output()
// 	return err
// }
// func reload() error {
// 	args := []string{"daemon-reload"}
// 	_, err := exec.Command("systemctl", args...).Output()
// 	return err
// }

// func execCommand(service string) error {
// 	args := []string{"start", service}
// 	_, err := exec.Command("systemctl", args...).Output()
// 	return err
// }
// func cmd.execCMDparams(commandName string, params []string) bool {
// 	cmd := exec.Command(commandName, params...)

// 	//显示运行的命令
// 	fmt.Println(cmd.Args)

// 	stdout, err := cmd.StdoutPipe()

// 	if err != nil {
// 		fmt.Println(err)
// 		return false
// 	}

// 	cmd.Start()

// 	reader := bufio.NewReader(stdout)

// 	//实时循环读取输出流中的一行内容
// 	for {
// 		line, err2 := reader.ReadString('\n')
// 		if err2 != nil || io.EOF == err2 {
// 			break
// 		}
// 		fmt.Println(line)
// 	}

// 	cmd.Wait()
// 	return true
// }

// func ServiceStart(service string) error {
// 	args := []string{"start", service}
// 	_, err := exec.Command("systemctl", args...).Output()
// 	return err
// }

// func ServiceStop(service string) error {
// 	args := []string{"stop", service}
// 	_, err := exec.Command("systemctl", args...).Output()
// 	return err
// }

// func ServiceExists(service string) bool {
// 	args := []string{"status", service}
// 	outBytes, _ := exec.Command("systemctl", args...).Output()
// 	output := string(outBytes)
// 	if strings.Contains(output, "Loaded: not-found") {
// 		return false
// 	}
// 	return true
// }

// func ServiceIsEnabled(service string) bool {
// 	args := []string{"is-enabled", service}
// 	_, err := exec.Command("systemctl", args...).Output()
// 	if err != nil {
// 		return false
// 	}
// 	return true
// }

// // ServiceIsActive will check is the service is "active". In the case of
// // crash looping services (kubelet in our case) status will return as
// // "activating", so we will consider this active as well.
// func ServiceIsActive(service string) bool {
// 	args := []string{"is-active", service}
// 	// Ignoring error here, command returns non-0 if in "activating" status:
// 	outBytes, _ := exec.Command("systemctl", args...).Output()
// 	output := strings.TrimSpace(string(outBytes))
// 	if output == "active" || output == "activating" {
// 		return true
// 	}
// 	return false
// }
