package service

import (
	"bufio"
	"fmt"
	"guide/core"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/shirou/gopsutil/cpu"
	_ "github.com/shirou/gopsutil/mem"
)

const (
	// messagesPath = "/var/log/messages"
	// securePath = "/var/log/secure"
	// cpuAll int = 100
	// gm uint64 = 1074000000
	ipMaxNum = 20
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

func SecurityIndex(c *gin.Context) {
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
		mlog.Error(err.Error())
		c.HTML(http.StatusOK, "security.tmpl", gin.H{
			"messageError": "error:" + err.Error(),
		})
		return
	}
	SecuresIpSlice, err := showSecureLog()
	if err != nil {
		mlog.Error(err.Error())
		c.HTML(http.StatusOK, "security.tmpl", gin.H{
			"securesError": "error:" + err.Error(),
		})
		return
	}
	if err := deleteFile("/tmp/ip-messages.txt"); err != nil {
		mlog.Error("delete ip-messages.txt fail.")
	}
	if err := deleteFile("/tmp/ip-secure.txt"); err != nil {
		mlog.Error("delete ip-secure.txt fail.")
	}

	c.HTML(http.StatusOK, "security.tmpl", gin.H{
		"Messages": MessagesIpSlice,
		"Secure":   SecuresIpSlice,
	})
}

type (
	Messages struct {
		Ip    string
		IpNum int
	}
	Secures struct {
		Ip    string
		IpNum int
	}
)

func deleteFile(folderFilePath string) error {
	err := os.Remove(folderFilePath)
	if err != nil {
		return fmt.Errorf("ERROR: delete file fail, %s", err.Error())
	}
	return nil
}

func getLogIpWriteFile(logPath string, title string) error {
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
			wLogPath := "/tmp/ip-" + title + ".txt"
			file, err := os.OpenFile(wLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = file.WriteString(match + "\n")
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

func showMessagesLog() ([]Messages, error) {
	err := getLogIpWriteFile("/var/log/messages", "messages")
	if err != nil {
		mlog.Error(err.Error())
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
		if v > ipMaxNum {
			ipSlice = append(ipSlice, Messages{
				Ip:    k,
				IpNum: v,
			})
		}
	}
	return ipSlice, nil
}

func showSecureLog() ([]Secures, error) {
	err := getLogIpWriteFile("/var/log/secure", "secure")
	if err != nil {
		mlog.Error(err.Error())
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
		// ip次数超过x次
		if v > ipMaxNum {
			ipSlice = append(ipSlice, Secures{
				Ip:    k,
				IpNum: v,
			})
		}
	}
	return ipSlice, nil
}

type (
	Ip struct {
		BlacklistIp string `json:"blacklistIp"`
	}
)

func AddBlacklistIp(c *gin.Context) {
	var ip Ip
	if err := c.ShouldBindJSON(&ip); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "json analysis parameter error.",
		})
		return
	}

	blacklistIpStr := fmt.Sprintf("all:%s:deny", ip.BlacklistIp)
	fmt.Println(blacklistIpStr)

	// 先查询db此ip是否已经被拉黑
	sql := fmt.Sprintf("select * from blacklist where ip=\"%s\"", ip.BlacklistIp)
	ipList, err := core.SelectBlacklistIp(sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadGateway,
			"msg": err.Error(),
		})
		return
	}

	if len(ipList) <= 0 {
		if err := core.WriteFile("/etc/hosts.deny", blacklistIpStr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadGateway,
				"msg": err.Error(),
			})
			return
		}

		// 写入db
		if err := core.InsertActBlacklist(ip.BlacklistIp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadGateway,
				"msg": err.Error(),
			})
			return
		}

		mlog.Info(fmt.Sprintf("add ip [%s] blacklist Ok", ip.BlacklistIp))
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg": "add ip blacklist OK.",
		})
		return
	}

	mlog.Warn(fmt.Sprintf("add ip [%s] blacklist already present.", ip.BlacklistIp))
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg": "add ip blacklist already present.",
	})

}

func MoveBlacklistIp(c *gin.Context) {

	const (
		filePath  = "/etc/hosts.deny"
		separator = ":"
	)

	var ip Ip
	if err := c.ShouldBindJSON(&ip); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "json analysis parameter error.",
		})
		return
	}

	sql := fmt.Sprintf("select * from blacklist where ip=\"%s\"", ip.BlacklistIp)
	ipList, err := core.SelectBlacklistIp(sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadGateway,
			"msg": err.Error(),
		})
		return
	}

	if len(ipList) <= 0 {
		mlog.Warn(fmt.Sprintf("add ip fail, find this IP address -> [%s]", ip.BlacklistIp))
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadGateway,
			"msg": "I couldn't find this IP address.",
		})
		return
	}

	content, _ := os.ReadFile(filePath)
	lines := strings.Split(string(content), "\n")

	var filteredLines []string
	for _, line := range lines {
		// 跳过空行和注释行
		if strings.HasPrefix(line, "#") || line == "" {
			filteredLines = append(filteredLines, line)
			continue
		}

		parts := strings.Split(line, separator)
		if len(parts) < 3 || parts[1] != ip.BlacklistIp {
			filteredLines = append(filteredLines, line)
		}
	}
	

	newContent := strings.Join(filteredLines, "\n")
	if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
		mlog.Error(err.Error()+" delete file is ip fail.")
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadGateway,
			"msg": err.Error(),
		})
		return
	}else {
		if err := core.DeleteBlacklistIp(ip.BlacklistIp); err != nil {
			mlog.Error(fmt.Sprintf("delete blacklist ip fail. ip -> [%s]", ip.BlacklistIp))
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusBadGateway,
				"msg": fmt.Sprintf("delete blacklist ip fail. ip -> [%s]", ip.BlacklistIp),
			})
			return
		}

	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg": "delete blacklist ip OK.",
	})
	
}


func AddNetworkBlacklistIp(c *gin.Context){

}


func MoveNetworkBlacklistIp(c *gin.Context){
	
}


func ShowDbBlacklistIp(c *gin.Context){
	sql := "select * from blacklist"
	ipList, err := core.SelectBlacklistIp(sql)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadGateway,
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"blackList": ipList,
	})
}