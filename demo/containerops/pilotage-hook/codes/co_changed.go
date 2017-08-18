package main

import (
	"bytes"
	"fmt"
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

	target, err := parseEnv(data)
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

	changed, err := diff(target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Failed to compare merges: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = false\n")
		os.Exit(1)
		return
	}

	fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT=true")
	fmt.Fprintf(os.Stdout, fmt.Sprintf("[COUT] CO_CHANGED=%t", changed))
}

func diff(target string) (bool, error) {
	repoParentPath, err := getGitRepoParentPath()
	if err != nil {
		return false, err
	}

	logCmd := exec.Command("git", "log", "--merges", "--oneline") //, "|", "head", "n2")
	logCmd.Dir = fmt.Sprintf("%s/containerops", repoParentPath)

	if err != nil {
		return false, err
	}

	headCmd := exec.Command("head", "-n2")

	headCmd.Stdin, err = logCmd.StdoutPipe()
	if err != nil {
		return false, err
	}
	var buf bytes.Buffer
	headCmd.Stdout = &buf
	headCmd.Start()
	logCmd.Run()
	headCmd.Wait()
	results := strings.Split(buf.String(), "\n")
	if len(results) < 2 {
		return false, fmt.Errorf("No comparisons")
	}

	lastMerge := results[0]
	prevMerge := results[1]

	parts := strings.Split(lastMerge, " ")
	lastMergeCommit := parts[0]

	parts = strings.Split(prevMerge, " ")
	prevMergeCommit := parts[0]

	diffCmd := exec.Command("git", "diff", fmt.Sprintf("%s..%s", prevMergeCommit, lastMergeCommit), "--name-only")
	diffCmd.Dir = fmt.Sprintf("%s/containerops/%s/", repoParentPath, target)
	buf.Reset()
	diffCmd.Stdout = &buf
	diffCmd.Stderr = &buf
	if err := diffCmd.Run(); err != nil {
		return false, err
	}

	results = strings.Split(buf.String(), "\n")
	// TODO Test the senarios that folders or files are renamed
	singularCodeChanged := false
	targetFolder := fmt.Sprintf("%s/", target)
	for i := 0; i < len(results); i++ {
		changedFile := results[i]
		if strings.HasPrefix(changedFile, targetFolder) /* && strings.HasSuffix(changedFile, ".go")  */ {
			singularCodeChanged = true
		}
	}

	return singularCodeChanged, nil
}

func parseEnv(env string) (target string, err error) {
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
		default:
			fmt.Fprintf(os.Stdout, "[COUT] Unknown Parameter: [%s]\n", s)
		}
	}
	return
}
