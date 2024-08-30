package service

import (
	"bufio"
	"fmt"
	_ "fmt"
	"guide/core"
	_ "guide/core"
	"log"
	"net/http"
	"os"
	"regexp"
	_ "strconv"
	"strings"
	_ "time"

	"github.com/gin-gonic/gin"
	_ "github.com/shirou/gopsutil/cpu"
	_ "github.com/shirou/gopsutil/mem"
)

// const cpuAll int = 100
// const gm uint64 = 1074000000

const (
	messagesPath = "/var/log/messages"
	securePath = "/var/log/secure"
)


// func CpuValueWDB() {
// 	for {
// 		time.Sleep(30 * time.Second)
// 		go func() {
// 			c2, _ := cpu.Percent(time.Duration(time.Second), false)
// 			cpuFrees := cpuAll - int(c2[0])
// 			// w to database
// 			if err := core.InsertActCPU(cpuFrees); err != nil {
// 				log.Println(err.Error())
// 			}
// 		}()

// 		go func() {
// 			m, _ := mem.VirtualMemory()
// 			if m.Total < gm {
// 				free := float64(m.Free/1024/1024) / float64(1024)
// 				strNum := strconv.FormatFloat(free, 'f', 2, 64)
// 				memFloatNum, err := strconv.ParseFloat(strNum, 64)
// 				if err != nil {
// 					log.Println("ERROR: mem 1G num str to float fail,", err.Error())
// 					return
// 				}
// 				if err := core.InsertActMEM(memFloatNum); err != nil {
// 					log.Println(err.Error())
// 				}
// 			} else {
// 				free := float64(m.Free) / 1024 / 1024 / 1024
// 				strNum := strconv.FormatFloat(free, 'f', 2, 64)
// 				memFloatNum, err := strconv.ParseFloat(strNum, 64)
// 				if err != nil {
// 					log.Println("ERROR: mem num str to float fail,", err.Error())
// 					return
// 				}
// 				if err := core.InsertActMEM(memFloatNum); err != nil {
// 					log.Println(err.Error())
// 				}
// 			}
// 		}()
// 	}
// }

// func showCpu() ([]int, error) {
// 	cpuList := make([]int, 0)
// 	rel, err := core.SelectCPU(fmt.Sprintf("SELECT cpunum FROM cpu ORDER BY id DESC LIMIT 50"))
// 	if err != nil {
// 		return nil, fmt.Errorf(err.Error())
// 	}
// 	for _, v := range rel {
// 		cpuList = append(cpuList, v.CpuNum)
// 	}
// 	return cpuList, nil
// }


// func showMem()([]float64, error) {
// 	memList := make([]float64, 0)
// 	rel, err := core.SelectMEM(fmt.Sprintf("SELECT memnum FROM mem ORDER BY id DESC LIMIT 50"))
// 	if err != nil {
// 		return nil, fmt.Errorf(err.Error())
// 	}
	
// 	for _, v := range rel {
// 		memList = append(memList, v.MemNum)
// 	}
// 	return memList, nil
// }


func SecurityIndex(c *gin.Context){
	// cpuList, err := showCpu()
	// if err != nil {log.Println(err); c.JSON(http.StatusOK, gin.H{
	// 	"code": http.StatusBadGateway,
	// 	"mag": err,
	// })}
	// memList, err := showMem()
	// if err != nil {log.Println(err); c.JSON(http.StatusOK, gin.H{
	// 	"code": http.StatusBadGateway,
	// 	"mag": err,
	// })}

	// c.HTML(http.StatusOK, "security.tmpl", gin.H{
	// 	"CPU": cpuList,
	// 	"MEM": memList,
	// })

	MessagesIpSlice, err := showMessagesLog()
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusOK, "security.tmpl", gin.H{
			"messageError": "error:"+err.Error(),
		})
		return
	}
	SecuresIpSlice, err := showSecureLog()
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusOK, "security.tmpl", gin.H{
			"securesError": "error:"+err.Error(),
		})
		return
	}
	if err := deleteFile("/tmp/ip-messages.txt"); err != nil {
		log.Println("ERROR: delete ip-messages.txt fail.")
	}
	if err := deleteFile("/tmp/ip-secure.txt"); err != nil {
		log.Println("ERROR: delete ip-secure.txt fail.")
	}
	
	c.HTML(http.StatusOK, "security.tmpl", gin.H{
		"Messages": MessagesIpSlice,
		"Secure": SecuresIpSlice,
	})
}

type (
	Messages struct {
		Ip string
		IpNum int
	}
	Secures struct {
		Ip string
		IpNum int
	}
)

func deleteFile(folderFilePath string)(error){
	err := os.Remove(folderFilePath)
		if err != nil {
			return fmt.Errorf("ERROR: delete file fail, %s", err.Error())
		}
		return nil
}

func getLogIpWriteFile(logPath string, title string)error{
	if err := core.WriteFile(logPath, "127.0.0.1"); err != nil {
		return err
	}
	file, err := os.Open(logPath)
    if err != nil {
        return err
    }
    defer file.Close()
 
    scanner := bufio.NewScanner(file)
    re := regexp.MustCompile(`((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}`)
    for scanner.Scan() {
        matches := re.FindAllString(scanner.Text(), -1)
        for _, match := range matches {
			wLogPath := "/tmp/ip-"+title+".txt"
			file, err := os.OpenFile(wLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = file.WriteString(match+"\n")
			if err != nil {
				return err
			}
        }
    }
 
    if err := scanner.Err(); err != nil {
        return err
    }	
	return nil
}

func showMessagesLog()([]Messages, error){
	err := getLogIpWriteFile("/var/log/messages", "messages")
	if err != nil {
		log.Println(err.Error())
		return nil, fmt.Errorf(err.Error())
	}
	ipCounts := make(map[string]int)
	ipSlice := make([]Messages, 0)
	file, err := os.Open("/tmp/ip-messages.txt")
	if err != nil {  
		return nil, fmt.Errorf(err.Error())
	}  
	defer file.Close() 
	scanner := bufio.NewScanner(file)  
	for scanner.Scan() {  
		// ipLines := strings.Split(scanner.Text(), " ")
		if scanner.Text() == "" || strings.TrimSpace(scanner.Text()) == "" {  
			continue
		}  
		ipCounts[scanner.Text()]++
	}  
	if err := scanner.Err(); err != nil {  
		return nil, fmt.Errorf(err.Error())
	} 
	for k, v := range ipCounts {  
		ipSlice = append(ipSlice, Messages{  
			Ip:   k,  
		 	IpNum: v,  
		})  
	}
	return ipSlice, nil
}


func showSecureLog()([]Secures, error){
	err := getLogIpWriteFile("/var/log/secure", "secure")
	if err != nil {
		log.Println(err.Error())
		return nil, fmt.Errorf(err.Error())
	}
	ipCounts := make(map[string]int)
	ipSlice := make([]Secures, 0)
	file, err := os.Open("/tmp/ip-secure.txt")
	if err != nil {  
		return nil, fmt.Errorf(err.Error())
	}  
	defer file.Close() 
	scanner := bufio.NewScanner(file)  
	for scanner.Scan() {  
		if scanner.Text() == "" || strings.TrimSpace(scanner.Text()) == "" {  
			continue
		}  
		ipCounts[scanner.Text()]++
	}  
	if err := scanner.Err(); err != nil {  
		return nil, fmt.Errorf(err.Error())
	} 
	for k, v := range ipCounts {  
		ipSlice = append(ipSlice, Secures{  
			Ip:   k,  
		 	IpNum: v,  
		})  
	}
	return ipSlice, nil
}


func  AddBlacklistIp(c *gin.Context){

}


func MoveBlacklistIp(c *gin.Context){

}

