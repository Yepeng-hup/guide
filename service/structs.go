package service

import (
	"os"
	"time"
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
		Time time.Time
		Power os.FileMode
	}

	FileAnchor struct {
		FileName string
		Href     string
		Size int64
		Time time.Time
		Power os.FileMode
	}
)
