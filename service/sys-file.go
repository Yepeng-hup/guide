package service

import (
	"encoding/json"
	"fmt"
	"guide/core"
	"guide/global"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mholt/archiver"
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
		Href:          strings.TrimRight(c.Request.URL.Path, "/") + "/..",
		//IpPort: ipAndPort,
	})
	for _, file := range files {
		href := strings.ReplaceAll(c.Request.URL.Path+"/"+file.Name(), "//", "/")
		times := file.ModTime()
		if file.IsDir() {
			dirList = append(dirList, DirectoryAnchor{
				DirectoryName: file.Name(),
				Href:          href,
				Size:          file.Size() / 1024 / 1024,
				Time:          times.Format("2006-01-02 15:04:05"),
				Power:         file.Mode(),
				//IpPort: ipAndPort,
			})

		} else {
			fileList = append(fileList, FileAnchor{
				FileName: file.Name(),
				Href:     href,
				Size:     file.Size() / 1024 / 1024,
				Time:     times.Format("2006-01-02 15:04:05"),
				Power:    file.Mode(),
				//IpPort: ipAndPort,
			})
		}
	}
	c.HTML(http.StatusOK, "file.tmpl", gin.H{
		"dirList":  dirList,
		"fileList": fileList,
		//"rootDir":  c.Request.URL.Path,
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

func CatFile(c *gin.Context) {
	var fileTailNameList = []string{"go", "sh", "txt", "py", "yaml", "yml", "md", "java", "c", "json", "env", "dockerfile", "conf", "js", "html", "css", "ts",
		"tmpl", "sql", "bat", "ps1", "php", "tmp", "xml", "ini", "jenkinsfile"}
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
		} else {
			//add if file size
			fileInfo, err := os.Stat(global.SaveDataDir + filePath + "/" + fileList[0])
			if err != nil {
				log.Println("ERROR: show file info fail,", err.Error())
				return
			}
			if fileInfo.Size()/1024/1024 > 1 {
				log.Printf("ERROR: only allow viewing files below 1M, the is file size -> [%v]M", fileInfo.Size()/1024/1024)
				return
			}
			//cat file
			fileContents, err := ioutil.ReadFile(global.SaveDataDir + filePath + "/" + fileList[0])
			if err != nil {
				log.Printf("ERROR: not is read file, %v\n", err.Error())
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"code":     http.StatusOK,
				"fileText": string(fileContents),
				"fileName": fileList[0],
			})
		}

	} else {
		log.Println("ERROR: Not in character [.] .")
		return
	}
}

func UpdateFile(c *gin.Context) {
	u := Update{
		FileName: c.PostForm("file"),
		Centent:  c.PostForm("content"),
		FilePath: c.PostForm("path"),
	}
	fileList := strings.Fields(u.FileName)
	fileWritePath := global.SaveDataDir + "/" + u.FilePath + "/" + fileList[0]
	err := os.WriteFile(fileWritePath, []byte(u.Centent), 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}

func CreateDir(c *gin.Context) {
	f := CreateDirs{
		DirName: c.PostForm("name"),
		DirPath: c.PostForm("path"),
	}

	createDirPath := global.SaveDataDir + f.DirPath + "/" + f.DirName
	err := os.Mkdir(createDirPath, 0755)
	if err != nil {
		log.Println("ERROR: create dir fail.", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "ERROR: create dir fail." + err.Error(),
		})
	} else {
		log.Printf("INFO: create dir success ---> [%v].", f.DirName)
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
		})
	}
}

func CreateFile(c *gin.Context) {
	f := CreateFiles{
		FileName: c.PostForm("name"),
		FilePath: c.PostForm("path"),
	}
	createFilePath := global.SaveDataDir + f.FilePath + "/" + f.FileName
	file, err := os.Create(createFilePath)
	if err != nil {
		log.Println("ERROR: create file fail.", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "ERROR: create file fail." + err.Error(),
		})
	} else {
		log.Printf("INFO: create  file success ---> [%v].", f.FileName)
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
		})
	}
	defer file.Close()
}

func copyDir(src string, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !srcInfo.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	// _, err = os.Open(dst)
	// if !os.IsNotExist(err) {
	// 	return fmt.Errorf("destination already exists")
	// }

	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			// Perform file copy
			err = copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	inFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, inFile)
	if err != nil {
		return err
	}

	fileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, fileInfo.Mode())
}

func DeleteDirAndFile(c *gin.Context) {
	d := Global{
		FileDirName: c.PostForm("FDname"),
		FileDirPath: c.PostForm("FDpath"),
	}
	fileAndDirList := strings.Fields(d.FileDirName)
	if fileAndDirList[0] == ".." || fileAndDirList[0] == "." {
		log.Println("ERROR: cannot delete root directory.")
		return
	}
	folderFilePath := global.SaveDataDir + d.FileDirPath + "/" + fileAndDirList[0]
	v, err := os.Stat(folderFilePath)
	if err != nil {
		log.Println("ERROR: show dir and file info fail.", err.Error())
		return
	}
	if v.IsDir() {
		if err := copyDir(folderFilePath, global.Hs); err != nil {
			log.Println("ERROR: mv dir and file to [hs] fail.", err.Error())
			return
		}

		if err := os.RemoveAll(folderFilePath); err != nil {
			log.Println("ERROR: delete dir fail.", err.Error())
			return
		}
		log.Printf("INFO: delete dir success. ---> [%s]", fileAndDirList[0])
		return
	} else {
		if err := core.CopyFile(folderFilePath, global.Hs); err != nil {
			log.Println("ERROR: mv file to [hs] fail.", err.Error())
			return
		}

		err := os.Remove(folderFilePath)
		if err != nil {
			log.Println("ERROR: delete file fail.", err.Error())
			return
		}
		log.Printf("INFO: delete file success. ---> [%s]", fileAndDirList[0])
		return
	}

}

func DecompressionZipTar(c *gin.Context) {
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
			if err != nil {
				log.Println("ERROR: unarchive zip fail,", err.Error())
				c.JSON(http.StatusOK, gin.H{
					"code":    http.StatusInternalServerError,
					"message": fmt.Sprint("ERROR: unarchive zip fail,", err.Error()),
				})
				return
			}
		case "gz":
			g := archiver.NewTarGz()
			err := g.Unarchive(global.SaveDataDir+f.FileDirPath+"/"+fileList[0], global.SaveDataDir+f.FileDirPath)
			if err != nil {
				log.Println("ERROR: unarchive tar.gz fail,", err.Error())
				err := core.UnGz(global.SaveDataDir + f.FileDirPath + "/" + fileList[0])
				if err != nil {
					log.Println("ERROR: unarchive gz fail,", err.Error())
					c.JSON(http.StatusOK, gin.H{
						"code":    http.StatusInternalServerError,
						"message": fmt.Sprint("ERROR: unarchive gz fail,", err.Error()),
					})
					return
				}
			}
		case "tar":
			t := archiver.NewTar()
			err := t.Unarchive(global.SaveDataDir+f.FileDirPath+"/"+fileList[0], global.SaveDataDir+f.FileDirPath)
			if err != nil {
				log.Printf("ERROR: unarchive tar fail %s\n", err.Error())
				c.JSON(http.StatusOK, gin.H{
					"code":    http.StatusInternalServerError,
					"message": fmt.Sprint("ERROR: unarchive tar fail,", err.Error()),
				})
				return
			}
		default:
			log.Println("ERROR: Invalid decompression format.")
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "ERROR: Invalid decompression format.",
			})
			return
		}
	} else {
		log.Println("ERROR: Not in character [.] .")
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "ERROR: Not in character [.] .",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "解压成功.",
	})
}

func CompressZipTar(c *gin.Context) {
	f := Global{
		c.PostForm("fileName"),
		c.PostForm("filePath"),
	}
	fileList := strings.Fields(f.FileDirName)
	if fileList[0] == ".." || fileList[0] == "." {
		log.Println("ERROR: cannot compress zip root directory.")
		return
	}
	err := archiver.Archive([]string{global.SaveDataDir + f.FileDirPath + "/" + fileList[0]}, global.SaveDataDir+f.FileDirPath+"/"+fileList[0]+".zip")
	if err != nil {
		log.Printf("ERROR: zip file and dir fail, %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"message": fmt.Sprint("ERROR: zip file and dir fail,", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "压缩成功.",
	})
}

func ShowRecycle(c *gin.Context) {
	fileDirSlice := make([]string, 0)

	d, err := os.Open(global.Hs)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"data": err.Error(),
		})
		return
	}
	defer d.Close()

	// -1 read all file
	files, err := d.Readdir(-1)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"data": err.Error(),
		})
		return
	}

	for _, file := range files {
		fileDirSlice = append(fileDirSlice, file.Name())
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": fileDirSlice,
	})

}

func DeleteRecycleFile(c *gin.Context) {
	fileName := c.PostForm("fileName")
	fileList := strings.Fields(fileName)
	deleteFilePath := global.Hs + "/" + fileList[0]
	err := os.Remove(deleteFilePath)
	if err != nil {
		if err := os.RemoveAll(deleteFilePath); err != nil {
			log.Println("delete file or dir fail -> ", deleteFilePath, err.Error())
			return
		}
	}
	log.Println("delete file or dir ok -> ", deleteFilePath)
}

func listFilesAndDirs(root *string, searchStr string, ssPath string) ([]FileAnchor, error) {
	fileList := make([]FileAnchor, 0)
	err := filepath.Walk(*root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		// 检查是否是文件及是否包含搜索关键字的文件
		if !info.IsDir() && strings.Contains(info.Name(), searchStr) {
			href := ssPath + "/" + info.Name()
			fileList = append(fileList, FileAnchor{
				FileName: info.Name(),
				Href:     href,
				Size:     info.Size() / 1024 / 1024,
				Time:     info.ModTime().Format("2006-01-02 15:04:05"),
				Power:    info.Mode(),
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return fileList, nil
}

func FileSearch(c *gin.Context) {
	f := SsText{
		c.PostForm("fileName"),
		c.PostForm("filePath"),
	}

	if f.SsFilePath == "/" {
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}

	folderFilePath := global.SaveDataDir + f.SsFilePath
	fileList, err := listFilesAndDirs(&folderFilePath, f.SsFile, f.SsFilePath)
	if err != nil {
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"fileList": fileList,
		"rootDir":  folderFilePath,
	})
}

func RootDirCheck(c *gin.Context) {
	data, _ := c.GetRawData()
	var body map[string]string
	_ = json.Unmarshal(data, &body)
	password := body["password"]
	if password == global.DirAccessPwd {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusInternalServerError,
	})
	return
}
