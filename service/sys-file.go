package service

import (
	"github.com/gin-gonic/gin"
	"github.com/mholt/archiver"
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
	//ipAndPort, err := core.ShowLocalIp(&global.InterfaceName)
	//if err != nil {
	//	log.Println(err.Error())
	//}
	dirList = append(dirList, DirectoryAnchor{
		DirectoryName: "..",
		Href: strings.TrimRight(c.Request.URL.Path, "/") + "/..",
		//IpPort: ipAndPort,
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
				//IpPort: ipAndPort,
			})

		} else {
			fileList = append(fileList, FileAnchor{
				FileName: file.Name(),
				Href:     href,
				Size: file.Size()/1024/1024,
				Time: times.Format("2006-01-02 15:04:05"),
				Power: file.Mode(),
				//IpPort: ipAndPort,
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
	f, _ := c.MultipartForm()
	files := f.File["file"]
	for _, file := range files {
		// init route path
		filename := filepath.Base(file.Filename)
		// 注释部分为支持中文url，我这里是不支持，如果要想支持中文url，取消下面注释即可。
		//decodePath, err0 := url.PathUnescape(c.PostForm("path"))
		//if err0 != nil {
		//	log.Println("ERROR: decod path fail,", err0.Error())
		//	return
		//}
		//savePath := filepath.Join(global.SaveDataDir, decodePath, filename)
		savePath := filepath.Join(global.SaveDataDir, c.PostForm("path"), filename)
		// save file
		err := c.SaveUploadedFile(file, savePath)
		if err != nil {
			log.Println("ERROR: file save fail,", err.Error())
			return
		}
		log.Printf("INFO: file push success ---> [%s]", filename)
	}
}


func DownloadData(c *gin.Context, p *string) {
	c.Writer.WriteHeader(200)
	//prompt the client that this is a binary file rather than a regular text format
	c.Header("Content-Type", "application/octet-stream")
	c.File(*p)
}


func CatFile(c *gin.Context){
	var fileTailNameList = []string{"go","sh","txt","py","yaml","yml","md","java","c","json","env","dockerfile","conf","js","html","css","ts",
		"tmpl","sql","bat","ps1","php","tmp","xml","ini","jenkinsfile",}
	fileNmae := c.Query("fileName")
	filePath := c.Query("filePath")
	fileList := strings.Fields(fileNmae)
	lastIndex := strings.LastIndex(fileList[0], ".")
	if lastIndex != -1 && lastIndex+1 < len(fileList[0]) {
		// starting from the last position of the last point, truncate the string
		fName := fileList[0][lastIndex+1:]
		if !core.SliceCheck(fileTailNameList, fName) {
			log.Println("ERROR: This is not a file.")
			return
		}else {
			//add if file size
			fileInfo, err := os.Stat(global.SaveDataDir+filePath+"/"+fileList[0])
			if err != nil{
				log.Println("ERROR: show file info fail,", err.Error())
				return
			}
			if fileInfo.Size()/1024/1024 > 1 {
				log.Printf("ERROR: only allow viewing files below 1M, the is file size -> [%v]M", fileInfo.Size()/1024/1024)
				return
			}
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
	f := CreateDirs{
		DirName: c.PostForm("name"),
		DirPath: c.PostForm("path"),
	}
	createDirPath := global.SaveDataDir+f.DirPath+"/"+f.DirName
	err := os.Mkdir(createDirPath, 0755)
	if err != nil {
		log.Println("ERROR: create dir fail.", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"message": "目录创建失败.",
		})
	}
	log.Printf("INFO: create dir success ---> [%v].", f.DirName)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message": "目录创建成功.",
	})
}


func CreateFile(c *gin.Context){
	f := CreateFiles{
		FileName: c.PostForm("name"),
		FilePath: c.PostForm("path"),
	}
	createFilePath := global.SaveDataDir+f.FilePath+"/"+f.FileName
	//err := os.Mkdir(createDirPath, 0755)
	file, err := os.Create(createFilePath)
	defer file.Close()
	if err != nil {
		log.Println("ERROR: create file fail.", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"message": "文件创建失败.",
		})
	}
	log.Printf("INFO: create  file success ---> [%v].", f.FileName)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message": "文件创建成功.",
	})
}




func DeleteDirAndFile(c *gin.Context){
	d := Global{
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


func DecompressionZipTar(c *gin.Context){
	f := Global{
		c.PostForm("fileName"),
		c.PostForm("filePath"),
	}
	fileList := strings.Fields(f.FileDirName)
	lastIndex := strings.LastIndex(fileList[0], ".")
	if lastIndex != -1 && lastIndex+1 < len(fileList[0]) {
		fName := fileList[0][lastIndex+1:]
		switch fName {
			case "zip":
				err := archiver.Unarchive(global.SaveDataDir+f.FileDirPath+"/"+fileList[0], global.SaveDataDir+f.FileDirPath)
				if err != nil{
					log.Println("ERROR: unarchive zip fail,", err.Error())
					return
				}
			case "gz":
				g := archiver.NewTarGz()
				err := g.Unarchive(global.SaveDataDir+f.FileDirPath+"/"+fileList[0], global.SaveDataDir+f.FileDirPath)
				if err != nil {
					log.Println("ERROR: unarchive tar.gz fail,", err.Error())
					err := core.UnGz(global.SaveDataDir+f.FileDirPath+"/"+fileList[0])
					if err != nil{
						log.Println("ERROR: unarchive gz fail,", err.Error())
						return
					}
				}
			case "tar":
				t := archiver.NewTar()
				err := t.Unarchive(global.SaveDataDir+f.FileDirPath+"/"+fileList[0], global.SaveDataDir+f.FileDirPath)
				if err != nil {
					log.Printf("ERROR: unarchive tar fail %s\n",err.Error())
					return
			}
		default:
			log.Println("ERROR: Invalid decompression format.")
			return
		}
	}else {
		log.Println("ERROR: Not in character [.] .")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message": "解压成功.",
	})
}



func CompressZipTar(c *gin.Context){
	f := Global{
		c.PostForm("fileName"),
		c.PostForm("filePath"),
	}
	fileList := strings.Fields(f.FileDirName)
	if fileList[0] == ".."||fileList[0] == "."{
		log.Println("ERROR: cannot compress zip root directory.")
		return
	}
	err := archiver.Archive([]string{global.SaveDataDir+f.FileDirPath+"/"+fileList[0]}, global.SaveDataDir+f.FileDirPath+"/"+fileList[0]+".zip")
	if err != nil {
		log.Printf("ERROR: zip file and dir fail, %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"message": "压缩失败.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message": "压缩成功.",
	})
}


