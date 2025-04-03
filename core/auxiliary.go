package core

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"guide/global"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func BackupPrefix() string {
	t := time.Now()
	return t.Format("2006-01-02")
}

func ShowUrl() []News {
	var structSlice []News
	file, err := os.Open(global.FilePath)
	if err != nil {
		err := os.Mkdir("text", 0755)
		if err != nil {
			mlog.Error(fmt.Sprintf("create dir fail,%s", err.Error()))
			return nil
		}
		return nil
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		slice := strings.Split(line, " ")
		txtStruct := News{
			UName: slice[0],
			Url:   slice[1],
			Notes: slice[2],
		}
		structSlice = append(structSlice, txtStruct)
	}

	if err := scanner.Err(); err != nil {
		mlog.Error(err.Error())
		return nil
	}
	return structSlice
}

func ShowLocalIp(interfaceName *string) (string, error) {
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
				return ipnet.IP.String() + ":" + Cfg.ListenPort, nil
			}
		}

	}

	//interfaces, err := net.Interfaces()
	//if err != nil {
	//	return "", fmt.Errorf("ERROR: Unable to obtain the list of network interfaces, %s", err.Error())
	//}
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
	//				return ipNet.IP.String()+":"+Cfg.ListenPort, nil
	//			}
	//		}
	//	}
	//}
	return "127.0.0.1" + ":" + Cfg.ListenPort, nil
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

func SliceCheck(slices []string, targets string) bool {
	for _, value := range slices {
		if value == targets {
			return true
		}
	}
	return false
}

func IfElement(slice []string, element string) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

func UnGz(gzSrcPath string) error {
	//open file
	gzFile, err := os.Open(gzSrcPath)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer gzFile.Close()

	// new gzip Reader
	gzReader, err := gzip.NewReader(gzFile)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer gzReader.Close()
	outFile, err := os.Create(strings.TrimSuffix(gzSrcPath, ".gz"))
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer outFile.Close()

	// gzip data cp to outfile
	_, err = io.Copy(outFile, gzReader)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}

func PasswordEncryption(p, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("ERROR: New key fail,%s", err.Error())
	}
	ciphertext := make([]byte, aes.BlockSize+len(p))
	iv := ciphertext[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return "", fmt.Errorf(err.Error())
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(p))
	encryptedPassword := base64.URLEncoding.EncodeToString(ciphertext)
	return encryptedPassword, nil
}

func PasswordDecrypt(p, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("ERROR: New key fail,%s", err.Error())
	}
	ciphertext, err := base64.URLEncoding.DecodeString(p)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ERROR: %s", "The ciphertext character is too short.")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	decryptedPassword := string(ciphertext)
	return decryptedPassword, nil
}

func ShowSys() string {
	return runtime.GOOS
}

func LinuxC(code string) error {
	cmd := exec.Command("/bin/bash", "-c", code)
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("use command error,%s", err)
	}
	fmt.Println("****************************************************************************************************************************\n", string(out))
	fmt.Println("****************************************************************************************************************************")
	return nil
}

func WinC(code string) error {
	cmd := exec.Command("cmd", "/c", code)
	out, err := cmd.Output()
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

func WriteFile(filePath, fileContent string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer file.Close()

	_, err = file.WriteString(fileContent + "\n")
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	return nil
}

func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("unable to open source file: %v", err)
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst + "/" + filepath.Base(src))
	if err != nil {
		return fmt.Errorf("unable to create destination file: %v", err)
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %v", err)
	}

	// sync disk
	err = destinationFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync file: %v", err)
	}

	return nil
}

//  泛型去重切片元素

func DeduplicateGeneric[T comparable](slice []T) []T {
	seen := make(map[T]struct{})
	result := []T{}

	for _, item := range slice {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
