package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/Huawei/containerops/singular/init_config"
)

//log "github.com/Sirupsen/logrus"
//"github.com/kubespray/kargo-cli/common"

// "github.com/fatih/color"
// "github.com/spf13/cobra"
// "github.com/spf13/viper"

var (
	// RedPrint      = color.New(color.FgHiRed).SprintFunc()
	// YellowPrint   = color.New(color.FgHiYellow).SprintFunc()
	// GreenPrint    = color.New(color.FgHiGreen).SprintFunc()
	CfgFile       string
	KargoPath     string
	InventoryPath string
	LogFile       string
	ClusterName   string
	etcdCount     uint16
	masterCount   uint16
	nodeCount     uint16
	Log           *log.Logger
)

type SSHCommander struct {
	User string
	IP   string
}

//
func RestartSvc(svcArr []string) error {
	for _, svc := range svcArr {
		args := []string{"restart", svc}
		ExecCMDparams(svc, args)
		args = []string{"enable", svc}
		ExecCMDparams(svc, args)
	}
	args := []string{"daemon-reload"}
	//_, err := exec.Command("systemctl", args...).Output()
	err := ExecCMDparams("systemctl", args)
	return err
}
func Reload() error {
	args := []string{"daemon-reload"}
	//_, err := exec.Command("systemctl", args...).Output()

	return ExecCMDparams("systemctl", args)

}

func ExecCommand(service string) error {
	args := []string{"start", service}
	//	_, err := exec.Command("systemctl", args...).Output()

	return ExecCMDparams("systemctl", args)
}

func ServiceStart(service string) error {
	args := []string{"start", service}
	//	_, err := exec.Command("systemctl", args...).Output()

	return ExecCMDparams("systemctl", args)
}

func ServiceStop(service string) error {
	args := []string{"stop", service}
	//	_, err := exec.Command("systemctl", args...).Output()

	return ExecCMDparams("systemctl", args)
}

func ServiceExists(service string) bool {
	args := []string{"status", service}
	outBytes, _ := exec.Command("systemctl", args...).Output()
	ExecCMDparams("systemctl", args)
	output := string(outBytes)
	if strings.Contains(output, "Loaded: not-found") {
		return false
	}
	return true
}

func ServiceIsEnabled(service string) bool {
	args := []string{"is-enabled", service}
	//	_, err := exec.Command("systemctl", args...).Output()
	ExecCMDparams("systemctl", args)

	// if err != nil {
	// 	return false
	// }
	return true
}

// ServiceIsActive will check is the service is "active". In the case of
// crash looping services (kubelet in our case) status will return as
// "activating", so we will consider this active as well.
func ServiceIsActive(service string) bool {
	args := []string{"is-active", service}
	// Ignoring error here, command returns non-0 if in "activating" status:
	outBytes, _ := exec.Command("systemctl", args...).Output()
	ExecCMDparams("systemctl", args)
	output := strings.TrimSpace(string(outBytes))
	if output == "active" || output == "activating" {
		return true
	}
	return false
}

// inner cmmd

func ExecShCommandEcho(txtContet string, targetName string) error {
	//	_, err := exec.Command("sh", "-c", "echo 456 /n 123  >>/etc/hosts").Output()
	return ExecCMDparams("echo", []string{"-e", txtContet + " " + targetName, ">>", "/etc/hosts"})
}

func ExecCPparams(sourceName string, targetName string) error {

	return ExecCMDparams("cp", []string{sourceName, targetName})
}

func ExecCMDparams(commandName string, params []string) error {

	cmdstr := []string{init_config.User + "@" + init_config.TargetIP}
	fmt.Println(cmdstr)
	cmdstr = append(cmdstr, commandName)
	fmt.Println(cmdstr)

	for _, item := range params {
		cmdstr = append(cmdstr, item)
	}
	cmd := exec.Command("ssh", cmdstr...)
	//show cmds
	fmt.Println(cmd.Args)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return err
	}

	cmd.Start()

	reader := bufio.NewReader(stdout)

	// show content of stream in time
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
	}

	cmd.Wait()
	return err
}

func LocalExecCMDparams(commandName string, params []string) error {

	cmdstr := []string{}
	fmt.Println(cmdstr)
	cmdstr = append(cmdstr, commandName)
	fmt.Println(cmdstr)

	for _, item := range params {
		cmdstr = append(cmdstr, item)
	}
	cmd := exec.Command(commandName, params...)

	//show cmds
	fmt.Println(cmd.Args)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return err
	}

	cmd.Start()

	reader := bufio.NewReader(stdout)

	// show content of stream in time
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
	}

	cmd.Wait()
	return err
}
