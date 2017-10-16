package module

import (
	"os"
	"text/template"
	"bytes"
	"fmt"
	yaml "gopkg.in/yaml.v2"
)


var flowfile = "./output1.yml"

// Flow is DevOps orchestration flow struct.
type Flow struct {
	ID        int64      `json:"-" yaml:"-"`
	Model     string     `json:"-" yaml:"-"`
	URI       string     `json:"uri" yaml:"uri"`
	Number    int64      `json:",omitempty" yaml:",omitempty"`
	Title     string     `json:"title" yaml:"title"`
	Version   int64      `json:"version" yaml:"version"`
	Tag       string     `json:"tag" yaml:"tag"`
	Timeout   int64      `json:"timeout" yaml:"timeout"`
	Status    string     `json:"status,omitempty" yaml:"status,omitempty"`
	Logs      []string   `json:"logs,omitempty" yaml:"logs,omitempty"`
	Stages    []Stage    `json:"stages,omitempty" yaml:"stages,omitempty"`
	Receivers []Receiver `json:"receivers,omitempty" yaml:"receivers,omitempty"`
}
type Stage struct {
	ID         int64    `json:"-" yaml:"-"`
	T          string   `json:"type" yaml:"type"`
	Name       string   `json:"name" yaml:"name"`
	Title      string   `json:"title" yaml:"title"`
	Sequencing string   `json:"sequencing,omitempty" yaml:"sequencing,omitempty"`
	Status     string   `json:"status,omitempty" yaml:"status,omitempty"`
	Logs       []string `json:"logs,omitempty" yaml:"logs,omitempty"`
	Actions    []Action `json:"actions,omitempty" yaml:"actions,omitempty"`
}

// Receiver receives the flow execution result
type Receiver struct {
	Type    string `json:"type" yaml:"type"`
	Address string `json:"address" yaml:"address"`
}

// Action is
type Action struct {
	ID     int64    `json:"-" yaml:"-"`
	Name   string   `json:"name" yaml:"name"`
	Title  string   `json:"title" yaml:"title"`
	Status string   `json:"status,omitempty" yaml:"status,omitempty"`
	Jobs   []Job    `json:"jobs,omitempty" yaml:"jobs,omitempty"`
	Logs   []string `json:"logs,omitempty" yaml:"logs,omitempty"`
}

// Job is
type Job struct {
	ID            int64               `json:"-" yaml:"-"`
	T             string              `json:"type" yaml:"type"`
	Name          string              `json:"name" yaml:"name,omitempty"`
	Kubectl       string              `json:"kubectl" yaml:"kubectl"`
	Endpoint      string              `json:"endpoint" yaml:"endpoint"`
	Timeout       int64               `json:"timeout" yaml:"timeout"`
	Status        string              `json:"status,omitempty" yaml:"status,omitempty"`
	Resources     Resource            `json:"resources" yaml:"resources"`
	Logs          []string            `json:"logs,omitempty" yaml:"logs,omitempty"`
	Environments  []map[string]string `json:"environments" yaml:"environments"`
	Outputs       []string            `json:"outputs,omitempty" yaml:"outputs,omitempty"`
	Subscriptions []map[string]string `json:"subscriptions,omitempty" yaml:"subscriptions,omitempty"`
}
type Resource struct {
	CPU    string `json:"cpu" yaml:"cpu"`
	Memory string `json:"memory" yaml:"memory"`
}

	var poutput = "./output.yml"
	var tphead = "./module/head.yml"
	var tpaction = "./module/action.yml"
	var tpfoot = "./module/foot.yml"

//discard
func Buildyml(value string) {
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

func UnmarshlYML(flowfile string)  Flow{
	//fmt.Println("Hello World!")
	yml := ReadFile(flowfile)
	//fmt.Println("yml %s\n", yml)
	//yaml.Unmarshal([]byte("name: 1\nb: 2"), &t)
	var flow Flow
	yaml.Unmarshal(yml, &flow)
	//fmt.Println("object %s\n", flow)
	//fmt.Println("stages %s", flow.Name) //*
	fmt.Println("stages %s", flow.Title) //* title: Components For Python
	fmt.Println("stages %s", flow.Tag) //* tag: latest	
	
	// todo diff for  stages , build test and flow 
	fmt.Println("stages %s", flow.Stages[1].Title) //** 	
	fmt.Println("stages %s", flow.Stages[1].Name) //**
	fmt.Println("stages %s", flow.Stages[1].Actions[0].Title) //*
	fmt.Println("stages %s", flow.Stages[1].Actions[0].Name) //*
	fmt.Println("stages %s", flow.Stages[1].Actions[0].Jobs[0].Endpoint) //*
	fmt.Println("stages %s", flow.Stages[1].Actions[0].Jobs[0].Environments[0]["CO_DATA"]) //*
	for k, v := range flow.Stages[1].Actions[0].Jobs[0].Environments[0] {
		fmt.Println(k, v)
	}
	//fmt.Println("stages %s", flow.Stages.Type)
	// todo rebuild yml from ex-ymal file
	return flow 
}

func (f *Flow) YAML() ([]byte, error) {
	return yaml.Marshal(f)
}

func  CreateYMLwihtURL(name ,path,imgUrl string){
	flow :=UnmarshlYML(path+"/"+name+".yml")
	//flow.Stages[1].Actions[0].Jobs[0].Environments[0]["CO_DATA"] = imgUrl
	flow.Stages[1].Actions[0].Jobs[0].Endpoint = imgUrl
	flow.Stages[1].Actions[0].Name = name 
	s,_:=flow.YAML()
	WriteFile(s,name+".yml") 
}