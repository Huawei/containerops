package module

import (
	"encoding/json"
	"fmt"
)

type Image struct {
	Endpoint string
}

func Json2obj(jsonstring string) (obj Image) {
	
	Obj := Image{}
	
	json.Unmarshal([]byte(jsonstring), &Obj)
		return  Obj
}

func Obj2Json(obj Image) (jsonstring string) {
	s, err := json.Marshal(obj)
	if err != nil {
		fmt.Println("error:", err)
	}	
	return string(s)
}
