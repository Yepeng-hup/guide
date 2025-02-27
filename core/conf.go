package core

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

var Cfg *Config = nil

type (
	Config struct {
		FileDataDir          string `json:"file_data_dir"`
		ListenHost           string `json:"listen_host"`
		ListenPort           string `json:"listen_port"`
		StartWhiteList       string `json:"start_white_list"`
		WhiteList            string `json:"white_list"`
		PasswdAdminWhiteList string `json:"passwd_admin_white_list"`
		DebugMode            string `json:"debug_mode"`
		AccessDir            string `json:"access_dir"`
		//DirAccessPwd         string `json:"dir_access_pwd"`
		HsDir         string `json:"hs_dir"`
		InterfaceName string `json:"interface_name"`
		StartMon      string `json:"start_mon"`
	}
)

func ReadJson(jsonFilePath string) {
	file, err := os.Open(jsonFilePath)
	if err != nil {
		log.Fatal("open json file: ", err.Error())
	}
	defer file.Close()
	f := bufio.NewReader(file)
	configObj := json.NewDecoder(f)
	if err = configObj.Decode(&Cfg); err != nil {
		log.Fatal(err.Error())
		return
	}
	return
}
