package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	data := os.Getenv("CO_DATA")
	if len(data) == 0 {
		fmt.Fprintf(os.Stderr, "[COUT] %s\n", "The CO_DATA value is null.")
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
		os.Exit(1)
	}

	target, hub, namespace, repo, tag, binary, err := parseEnv(data)

	if err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Parse the CO_DATA: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = false\n")
		os.Exit(1)
	}

	if err := clone(); err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Failed to clone git repo: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = false\n")
		os.Exit(1)
		return
	}

	localFile, err := build(target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Failed to build project: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = false\n")
		os.Exit(1)
		return
	}

	// push
	url, err := push(localFile, hub, namespace, repo, tag, binary)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Failed to upload binary: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = false\n")
		os.Exit(1)
		return
	}

	fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT=true\n")
	fmt.Fprintf(os.Stdout, "[COUT] CO_URL=%s\n", url)
}

func parseEnv(env string) (target, hub, namespace, repo, tag, binary string, err error) {
	files := strings.Fields(env)
	if len(files) == 0 {
		err = fmt.Errorf("CO_DATA value is null\n")
		return
	}

	for _, v := range files {
		s := strings.Split(v, "=")
		key, value := s[0], s[1]

		switch key {
		case "target":
			target = value
		case "hub":
			hub = value
		case "namespace":
			namespace = value
		case "repo":
			repo = value
		case "tag":
			tag = value
		case "binary":
			binary = value
		default:
			fmt.Fprintf(os.Stdout, "[COUT] Unknown Parameter: [%s]\n", s)
		}
	}
	return
}

func push(filePath, hub, namespace, repo, tag, binary string) (string, error) {
	url := fmt.Sprintf("https://%s/binary/v1/%s/%s/binary/%s/%s", hub, namespace, repo, tag, binary)
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	client := http.Client{}
	req, _ := http.NewRequest(http.MethodPut, url, file)
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	// res, err := http.Post(url, "binary/octet-stream", file)
	// if err != nil {
	//     return err
	// }
	// defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to upload binary file, status: %d", res.StatusCode)
	}

	return url, nil
}
