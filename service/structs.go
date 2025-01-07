package service

import (
	"os"
)

type (
	UrlType struct {
		TypeName string
	}

	From2 struct {
		UrlName   string
		UrlAddr string
		UrlType string
		UrlNotes string
	}

	From struct {
		Url   string
		UName string
		Notes string
	}
	Name struct {
		Name string
	}

	DirectoryAnchor struct {
		DirectoryName string
		Href          string
		Size int64
		Time string
		Power os.FileMode
		//IpPort string
	}

	FileAnchor struct {
		FileName string
		Href     string
		Size int64
		Time string
		Power os.FileMode
		//IpPort string
	}

	CreateDirs struct {
		DirName string
		DirPath string
	}

	CreateFiles struct {
		FileName string
		FilePath string
	}

	Global struct {
		FileDirName string
		FileDirPath string
	}

	CronsFrom struct {
		Cname string
		Ctime string
		Ccode string
		Cnotes string
	}

	SvcFrom struct {
		SvcName string
		SvcCmd string
		SvcNotes string
	}

	Update struct {
		FileName string
		Centent string
		FilePath string
	}

	SsText struct {
		SsFile string
		SsFilePath string
	}
)
