package service

import (
	"fmt"
	"guide/core"
	"log"
	"net/http"
	"strconv"
	"time"
	"guide/core/cmd"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

const cpuAll int = 100
const gm uint64 = 1074000000

type (
	IpUse interface{

	}
)

func CpuValueWDB() {
	for {
		time.Sleep(30 * time.Second)
		go func() {
			c2, _ := cpu.Percent(time.Duration(time.Second), false)
			cpuFrees := cpuAll - int(c2[0])
			// w to database
			if err := core.InsertActCPU(cpuFrees); err != nil {
				log.Println(err.Error())
			}
		}()

		go func() {
			m, _ := mem.VirtualMemory()
			if m.Total < gm {
				free := float64(m.Free/1024/1024) / float64(1024)
				strNum := strconv.FormatFloat(free, 'f', 2, 64)
				memFloatNum, err := strconv.ParseFloat(strNum, 64)
				if err != nil {
					log.Println("ERROR: mem 1G num str to float fail,", err.Error())
					return
				}
				if err := core.InsertActMEM(memFloatNum); err != nil {
					log.Println(err.Error())
				}
			} else {
				free := float64(m.Free) / 1024 / 1024 / 1024
				strNum := strconv.FormatFloat(free, 'f', 2, 64)
				memFloatNum, err := strconv.ParseFloat(strNum, 64)
				if err != nil {
					log.Println("ERROR: mem num str to float fail,", err.Error())
					return
				}
				if err := core.InsertActMEM(memFloatNum); err != nil {
					log.Println(err.Error())
				}
			}
		}()
	}
}

func showCpu() ([]int, error) {
	cpuList := make([]int, 0)
	rel, err := core.SelectCPU(fmt.Sprintf("SELECT cpunum FROM cpu ORDER BY id DESC LIMIT 50"))
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	for _, v := range rel {
		cpuList = append(cpuList, v.CpuNum)
	}
	return cpuList, nil
}


func showMem()([]float64, error) {
	memList := make([]float64, 0)
	rel, err := core.SelectMEM(fmt.Sprintf("SELECT memnum FROM mem ORDER BY id DESC LIMIT 50"))
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	
	for _, v := range rel {
		memList = append(memList, v.MemNum)
	}
	return memList, nil
}


func SecurityIndex(c *gin.Context){
	cpuList, err := showCpu()
	if err != nil {log.Println(err); c.JSON(http.StatusOK, gin.H{
		"code": http.StatusBadGateway,
		"mag": err,
	})}
	memList, err := showMem()
	if err != nil {log.Println(err); c.JSON(http.StatusOK, gin.H{
		"code": http.StatusBadGateway,
		"mag": err,
	})}

	c.HTML(http.StatusOK, "security.tmpl", gin.H{
		"CPU": cpuList,
		"MEM": memList,
	})

}


// Only applicable to linux
func ShowSysAll(c *gin.Context){
	lCmd := ""
	if err := cmd.UseCmd(lCmd); err != nil {
		log.Println(err)
		return
	}


}


func DisabledIp(c *gin.Context){

}


func ReleaseIp(c *gin.Context){

}