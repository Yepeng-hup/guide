package service

import (
	"github.com/gin-gonic/gin"
	"guide/core"
	"guide/global"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)


func CutDirAndFile(c *gin.Context, fullPath *string) {
	files, _ := ioutil.ReadDir(*fullPath)
	dirList := make([]DirectoryAnchor, 0)
	fileList := make([]FileAnchor, 0)
	ipAndPort, err := core.ShowLocalIp(&global.InterfaceName)
	if err != nil {
		log.Println(err.Error())
	}
	dirList = append(dirList, DirectoryAnchor{
		DirectoryName: "..",
		Href: strings.TrimRight(c.Request.URL.Path, "/") + "/..",
		IpPort: ipAndPort,
	})
	for _, file := range files {
		href := strings.ReplaceAll(c.Request.URL.Path +"/"+ file.Name(), "//", "/")
		times := file.ModTime()
		if file.IsDir() {
			dirList = append(dirList, DirectoryAnchor{
				DirectoryName: file.Name(),
				Href:          href,
				Size: file.Size()/1024/1024,
				Time: times.Format("2006-01-02 15:04:05"),
				Power: file.Mode(),
				IpPort: ipAndPort,
			})

		} else {
			fileList = append(fileList, FileAnchor{
				FileName: file.Name(),
				Href:     href,
				Size: file.Size()/1024/1024,
				Time: times.Format("2006-01-02 15:04:05"),
				Power: file.Mode(),
				IpPort: ipAndPort,
			})
		}
	}
	c.HTML(http.StatusOK, "file.tmpl", gin.H{
		"dirList": dirList,
		"fileList": fileList,
		"rootDir": c.Request.URL.Path,
	})
}


func UploadData(c *gin.Context) {
	file, _ := c.FormFile("file")
	if file == nil {
		log.Println("WARN: file is nil.")
		c.Redirect(http.StatusFound, "/")
		return
	}
	filename := filepath.Base(file.Filename)
	if strings.Contains(c.PostForm("path"), "..") {
		c.Redirect(http.StatusFound, "/")
		return
	}
	savePath := filepath.Join(global.SaveDataDir, c.PostForm("path"), filename)
	err := c.SaveUploadedFile(file, savePath)
	if err != nil {
		log.Println("ERROR: file save fail,", err.Error())
		c.Redirect(http.StatusFound, "/")
		return
	}
	log.Printf("INFO: file push success ---> [%s]", filename)
	c.Redirect(http.StatusFound, c.PostForm("path"))
}


func DownloadData(c *gin.Context, p *string) {
	c.Writer.WriteHeader(200)
	//提示客户端这是个二进制文件而非普通文本格式
	c.Header("Content-Type", "application/octet-stream")
	c.File(*p)
}


func CatFile(c *gin.Context){
	// 定义格式支持
	var fileTailNameList = []string{"go","sh","txt","py","yaml","yml","md","java","c","json","env","dockerfile","conf","js","html","css","ts",
		"tmpl","sql","bat","ps1","php"}
	fileNmae := c.Query("fileName")
	filePath := c.Query("filePath")
	fileList := strings.Fields(fileNmae)
	lastIndex := strings.LastIndex(fileList[0], ".")
	if lastIndex != -1 && lastIndex+1 < len(fileList[0]) {
		// 最后一个点的后一个位置开始截取字符串
		fName := fileList[0][lastIndex+1:]
		if !core.SliceCheck(fileTailNameList, fName) {
			log.Println("ERROR: This is not a file.")
			return
		}else {
			//cat file
			fileContents, err := ioutil.ReadFile(global.SaveDataDir+filePath+"/"+fileList[0])
			if err != nil {
				log.Printf("ERROR: not is read file, %v\n", err.Error())
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"fileText": string(fileContents),
			})
		}

	} else {
		log.Println("ERROR: Not in character [.] .")
		return
	}
}


func CreateDir(c *gin.Context){
	f := Creates{
		DirName: c.PostForm("name"),
		DirPath: c.PostForm("path"),
	}
	createDirPath := global.SaveDataDir+f.DirPath+"/"+f.DirName
	err := os.Mkdir(createDirPath, 0755)
	if err != nil {
		log.Println("ERROR: create dir and file fail.", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"message": "目录创建失败.",
		})
	}
	log.Printf("INFO: create dir and file success ---> [%v].", f.DirName)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message": "目录创建成功.",
	})
}


func DeleteDirAndFile(c *gin.Context){
	d := Deletes{
		FileDirName: c.PostForm("FDname"),
		FileDirPath: c.PostForm("FDpath"),
	}
	fileAndDirList := strings.Fields(d.FileDirName)
	if fileAndDirList[0] == ".."||fileAndDirList[0] == "."{
		log.Println("ERROR: cannot delete root directory.")
		return
	}
	folderFilePath := global.SaveDataDir+d.FileDirPath+"/"+fileAndDirList[0]
	v, err := os.Stat(folderFilePath)
	if err != nil {
		log.Println("ERROR: show dir and file info fail.", err.Error())
		return
	}
	if v.IsDir() {
		err := os.RemoveAll(folderFilePath)
		if err != nil {
			log.Println("ERROR: delete dir fail.", err.Error())
			return
		}
		log.Printf("INFO: delete dir success. ---> [%s]", fileAndDirList[0])
		return
	}else {
		err := os.Remove(folderFilePath)
		if err != nil {
			log.Println("ERROR: delete file fail.", err.Error())
			return
		}
		log.Printf("INFO: delete file success. ---> [%s]", fileAndDirList[0])
		return
	}

}


