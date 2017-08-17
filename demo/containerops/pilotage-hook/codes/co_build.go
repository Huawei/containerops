package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
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

func getGitRepoParentPath() (string, error) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return "", fmt.Errorf("Failed to get $GOPATH")
	}
	repoParentPath := fmt.Sprintf("%s/src/github.com/Huawei", gopath)
	return repoParentPath, nil
}

func clone() error {
	repoParentPath, err := getGitRepoParentPath()
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "clone", "https://github.com/Huawei/containerops")
	cmd.Dir = repoParentPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
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

func build(target string) (string, error) {
	if err := installDependencies(); err != nil {
		return "", err
	}

	repoParentPath, err := getGitRepoParentPath()
	if err != nil {
		return "", err
	}

	buildCmd := exec.Command("go", "build")
	buildCmd.Stderr = os.Stderr
	buildCmd.Stdout = os.Stdout
	buildCmd.Dir = fmt.Sprintf("%s/containerops/%s", repoParentPath, target)
	// var buf bytes.Buffer
	// buildCmd.Stdout = buf
	// buildCmd.Stderr = buf
	if err := buildCmd.Run(); err != nil {
		return "", err
	}

	localFile := fmt.Sprintf("%s/%s", buildCmd.Dir, target)
	return localFile, nil
}

func installDependencies() error {
	deps := [][]string{
		{"get", "-v", "github.com/Sirupsen/logrus"},
		{"get", "-v", "github.com/cloudflare/cfssl/cli"},
		{"get", "-v", "github.com/cloudflare/cfssl/cli/genkey"},
		{"get", "-v", "github.com/cloudflare/cfssl/cli/sign"},
		{"get", "-v", "github.com/cloudflare/cfssl/csr"},
		{"get", "-v", "github.com/cloudflare/cfssl/initca"},
		{"get", "-v", "github.com/cloudflare/cfssl/signer"},
		{"get", "-v", "github.com/digitalocean/godo"},
		{"get", "-v", "github.com/fernet/fernet-go"},
		{"get", "-v", "github.com/jinzhu/gorm"},
		{"get", "-v", "github.com/jinzhu/gorm/dialects/mysql"},
		{"get", "-v", "github.com/logrusorgru/aurora"},
		{"get", "-v", "github.com/mitchellh/go-homedir"},
		{"get", "-v", "github.com/pkg/sftp"},
		{"get", "-v", "github.com/spf13/cobra"},
		{"get", "-v", "github.com/spf13/viper"},
		{"get", "-v", "golang.org/x/crypto/ssh"},
		{"get", "-v", "golang.org/x/net/context"},
		{"get", "-v", "golang.org/x/oauth2"},
		{"get", "-v", "gopkg.in/macaron.v1"},
		{"get", "-v", "gopkg.in/yaml.v2"},
	}
	for i := 0; i < len(deps); i++ {
		args := deps[i]
		cmd := exec.Command("go", args...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
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
