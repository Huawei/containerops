package module

import (
	"os"
	"text/template"
	"bytes"
	"fmt"
)

	var poutput = "./output.yml"
	var tphead = "./module/head.yml"
	var tpaction = "./module/action.yml"
	var tpfoot = "./module/foot.yml"

func Buildtp(value string) {
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
