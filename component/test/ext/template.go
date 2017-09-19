package ext

//https://studygolang.com/articles/464
import (
	"os"
	"text/template"
	"bytes"
	"fmt"
)
//curl -i -X POST -H 'Content-type':'application/yaml' --data-binary @output.yml  45.55.29.141:8080/flow/v1/containerops/python_analysis_coala/flow/latest/yaml
//http://blog.csdn.net/sryan/article/details/52353937

	//var wireteString = "测试n"
	var filename = "./output.yml"
	var tphead = "./ext/head.yml"
	var tpaction = "./ext/action.yml"
	var tpfoot = "./ext/foot.yml"
//func Etemplate(tpString string,value string) (string) {
		
func Buildtp(value string) {
		var poutput = "./output.yml"
		var tphead = "./head.yml"
		var tpaction = "./action.yml"
		var tpfoot = "./foot.yml"
		bhead := ReadFile(tphead)
		baction := ReadFile(tpaction)
		bfoot := ReadFile(tpfoot)

		buf := bytes.NewBuffer(bhead)
		buf.Write(baction)
		buf.Write(bfoot)
		var bugString =	buf.String()
		tmpl, err := template.New("test").Parse(bugString) //建立一个模板，内容是"hello, {{.}}"
		if err != nil {
			panic(err)
		}
		err = tmpl.Execute(os.Stdout, value) //将string与模板合成，变量name的内容会替换掉{{.}}
		//合成结果放到os.Stdout里
		if err != nil {
			panic(err)
		}

		fmt.Println("bugString:", bugString)	
		WriteFile([]byte(bugString),poutput)
  
//TO DO: upload to assebling
//return tpString
}
