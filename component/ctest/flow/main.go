package main

import (
	"fmt"
	"os"
)

func main() {

	//Get the CO_DATA from environment parameter "CO_DATA"
	data := os.Getenv("CO_DATA")
	if len(data) == 0 {
		fmt.Fprintf(os.Stderr, "[COUT] %s\n", "The CO_DATA value is null.")
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")

}
