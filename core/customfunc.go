package core

import "strings"

func CheckFileTailStr(fileName string, tailStr ...string) bool {
	for _, suffix := range tailStr {
		if strings.HasSuffix(fileName, suffix) {
			return true
		}
	}
	return false
}


func IfFileSize(){

}


