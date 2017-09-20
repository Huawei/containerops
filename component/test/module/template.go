package module

import (
	"os"
	"text/template"
	"bytes"
	"fmt"
)

	var filename = "./output.yml"
	var tphead = "./ext/head.yml"
	var tpaction = "./ext/action.yml"
	var tpfoot = "./ext/foot.yml"

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
		tmpl, err := template.New("test").Parse(bugString) 
		if err != nil {
			panic(err)
		}
		err = tmpl.Execute(os.Stdout, value)
		if err != nil {
			panic(err)
		}

		fmt.Println("bugString:", bugString)	
		WriteFile([]byte(bugString),poutput)
  
//TO DO: upload to assebling
}
