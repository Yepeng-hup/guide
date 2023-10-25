package cmd

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os/exec"
	"runtime"
)

func showSys()string{
	return runtime.GOOS
}

func linuxC(code string)error{
	cmd := exec.Command("/bin/bash", "-c", code)
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("use command error,%s", err)
	}
	fmt.Println("****************************************************************************************************************************\n", string(out))
	fmt.Println("****************************************************************************************************************************")
	return nil
}


func winC(code string)error{
	cmd := exec.Command("cmd", "/c", code)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("use command error,%s", err)
	}
	reader := transform.NewReader(bytes.NewReader(out), simplifiedchinese.GBK.NewDecoder())
	output, err := ioutil.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("byte encoding conversion error,%s", err.Error())
	}
	fmt.Println("****************************************************************************************************************************\n", string(output))
	fmt.Println("****************************************************************************************************************************")
	return nil
}

func UseCmd(code string)(error,){
	osType := showSys()
	switch osType {
	case "linux":
		err := linuxC(code)
		if err != nil {
			return fmt.Errorf("ERROR: %s", err)
		}
		return nil
	case "windows":
		err := winC(code)
		if err != nil {
			return fmt.Errorf("ERROR: %s", err)
		}
		return nil
	default:
		return fmt.Errorf("%s","WARN: unsupported operating system.")
	}
}