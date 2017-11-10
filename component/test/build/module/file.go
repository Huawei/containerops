package module

import (
	"io/ioutil"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func ReadFile(readf string) []byte {
	b, err := ioutil.ReadFile(readf)
	check(err)

	return b
}

func WriteFile(outputb []byte, outputf string) {

	err := ioutil.WriteFile(outputf, outputb, 0666)
	check(err)

}
