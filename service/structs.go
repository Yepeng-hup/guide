package service

import (
	"os"
)

type (
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
		IpPort string
	}

	FileAnchor struct {
		FileName string
		Href     string
		Size int64
		Time string
		Power os.FileMode
		IpPort string
	}

	Creates struct {
		DirName string
		DirPath string
	}

	Deletes struct {
		FileDirName string
		FileDirPath string
	}
)
