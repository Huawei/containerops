package main

import (
	"fmt"
	"os"
	"os/exec"
)

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

func getGitRepoParentPath() (string, error) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return "", fmt.Errorf("Failed to get $GOPATH")
	}
	repoParentPath := fmt.Sprintf("%s/src/github.com/Huawei", gopath)
	return repoParentPath, nil
}
