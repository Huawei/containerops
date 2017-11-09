package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"github.com/Huawei/containerops/component/ctest/build/module"
	"path/filepath"
    "fmt"
	"flag"
	"strings"
)

var tag = "latest"
var registry = "hub.opshub.sh"
var namespace = "containerops"
var api_url = "https://build.opshub.sh/assembling/build?"


func main() {
	flag.Parse()
	//root := flag.Arg(0)

	var root string //to do with out tar
	flag.StringVar(&root, "path", "tar", "dir path")
	getFilelist(root)
	
	//return
	
	// var image string
	// flag.StringVar(&image, "image", "test_image", "image name")


	// flag.Parse()

	//Print result
	fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "true")
	os.Exit(0)
}



func getFilelist(path string) {
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
			if ( f == nil ) {return err}
			if f.IsDir() {return nil}
			println("origin:",path)
			gpath(path)
			return nil
	})
	if err != nil {
			fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}

func gpath(imgpath string){
	
	image :=strings.Replace(imgpath, "tar", "", -1)
	image =strings.Replace(image, "/", "", 1)
	image =strings.Replace(image, ".", "", 1)
	
	fmt.Println("image:", image)
	url := api_url
	buf := bytes.NewBufferString(url)
	buf.Write([]byte("image="))
	buf.Write([]byte(image))
	buf.Write([]byte("&tag="))
	buf.Write([]byte(tag))
	buf.Write([]byte("&tag="))
	buf.Write([]byte(tag))
	buf.Write([]byte("&registry="))
	buf.Write([]byte(registry))
	buf.Write([]byte("&namespace="))
	buf.Write([]byte(namespace))
	fmt.Println(buf.String())

	//module.CreateYMLwihtURL(image,path,"hub.opshub.sh/containerops/component-python-coala-workflow111:latest")
	
	fmt.Println("imgpath:", imgpath)	
	UploadBinaryFile(image, imgpath, buf.String())
	
}


// Upload binary file to the service.
func UploadBinaryFile(name, filePath, url string) error {
	
	if f, err := os.Open(filePath); err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Parse the CO_DATA error: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
		return err
	} else {
		defer f.Close()
		if req, err := http.NewRequest(http.MethodPost,
			url, f); err != nil {
			fmt.Fprintf(os.Stderr, "[COUT] Parse the CO_DATA error: %s\n", err.Error())
			fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
			return err
		} else {
			req.Header.Set("Content-Type", "text/plain")

			client := &http.Client{}
			if resp, err := client.Do(req); err != nil {
				return err
			} else {
				fmt.Println("UploadBinaryFile:", "imgpath")	
				
				defer resp.Body.Close()

				switch resp.StatusCode {
				case http.StatusOK:
					{
						body := &bytes.Buffer{}
						_, err := body.ReadFrom(resp.Body)
						if err != nil {
							log.Fatal(err)
							fmt.Fprintf(os.Stderr, "[COUT] Parse the CO_DATA error: %s\n", err.Error())
							fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
						}
						resp.Body.Close()
						fmt.Println(resp.StatusCode)
						fmt.Println(resp.Header)
						fmt.Println(body)
						var jsonobj module.Image
						jsonobj = module.Json2obj(body.String())
						fmt.Println(jsonobj)
						// module.Buildyml(jsonobj.Endpoint)
						//hub.opshub.sh/containerops/component-python-coala-workflow111:latest
						//module.CreateYMLwihtURL(name,filePath,"hub.opshub.sh/containerops/component-python-coala-workflow111:latest")
						
						//TODO move to flow componnent 
						//module.CreateYMLwihtURL(name, filePath, jsonobj.Endpoint)
						return nil
					}
				case http.StatusBadRequest:
					{

						return fmt.Errorf("Binary upload failed.")

					}
				case http.StatusUnauthorized:
					return fmt.Errorf("Action unauthorized.")
				default:
					return fmt.Errorf("Unknown error.")
				}

			}
		}
	}

	return nil
}
