package core

import (
	"bufio"
	"fmt"
	"guide/global"
	"log"
	"net"
	"os"
	"strings"
)

func ShowUrl()([]News,){
	var structSlice []News
	file, err := os.Open(global.FilePath)
	if err != nil {
		log.Println("ERROR: ",err.Error())
		return nil
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		slice := strings.Split(line, " ")
		txtStruct := News{
			UName:   slice[0],
			Url: slice[1],
			Notes: slice[2],
		}
		structSlice = append(structSlice, txtStruct)
	}

	if err := scanner.Err(); err != nil {
		log.Println("ERROR: ",err.Error())
		return nil
	}
	return structSlice
}


func ShowLocalIp(interfaceName *string)(string, error){
	iface, err := net.InterfaceByName(*interfaceName)
	if err != nil {
		return "", fmt.Errorf("Unable to obtain network interface %s：%v\n", *interfaceName, err.Error())
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return "", fmt.Errorf("Unable to obtain network interface -> [%s] ip address：%v\n", *interfaceName, err.Error())
	}
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
			if ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String()+":"+global.Port, nil
				}
			}


	}

	// 获取所有网络接口
	//interfaces, err := net.Interfaces()
	//if err != nil {
	//	return "", fmt.Errorf("ERROR: Unable to obtain the list of network interfaces, %s", err.Error())
	//}

	// 遍历网络接口并获取IP地址
	//for _, ifaceEth := range interfaces {
	//	addrs, err := ifaceEth.Addrs()
	//	if err != nil {
	//		return "", fmt.Errorf("ERROR: Unable to obtain network interface address, %s", err.Error())
	//	}
	//
	//	for _, addr := range addrs {
	//		ipNet, ok := addr.(*net.IPNet)
	//		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
	//			if ifaceEth.Name == "eth0" || ifaceEth.Name == "WLAN" || ifaceEth.Name == "ens33"{
	//				//fmt.Printf("接口：%s，IP地址：%s\n", ifaceEth.Name, ipNet.IP.String())
	//				return ipNet.IP.String()+":"+global.Port, nil
	//			}
	//		}
	//	}
	//}
	return "127.0.0.1"+":"+global.Port, nil
}

//func BitConvert(bit int64)(int64, string){
//	bitsPerKb := int64(1024)
//	bitsPerMb := bitsPerKb * 1024
//	bitsPerGb := bitsPerMb * 1024
//
//	switch {
//	case bit < bitsPerKb:
//		return int64(bit), "bit"
//	case bit < bitsPerMb:
//		return int64(bit) / int64(bitsPerKb), "Kb"
//	case bit < bitsPerGb:
//		return int64(bit) / int64(bitsPerMb), "Mb"
//	default:
//		return int64(bit) / int64(bitsPerGb), "Gb"
//	}
//}

