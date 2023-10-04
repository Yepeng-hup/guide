package core

import (
	"bufio"
	"guide/global"
	"log"
	"os"
	"strings"
)

func ShowUrl()([]News,){
	var structSlice []News
	file, err := os.Open(global.FilePath)
	if err != nil {
		log.Println("ERROR: ",err.Error())
		return nil
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		slice := strings.Split(line, " ")
		txtStruct := News{
			UName:   slice[0],
			Url: slice[1],
			Notes: slice[2],
		}
		structSlice = append(structSlice, txtStruct)
	}

	if err := scanner.Err(); err != nil {
		log.Println("ERROR: ",err.Error())
		return nil
	}
	return structSlice
}


//func BitConvert(bit int64)(int64, string){
//	bitsPerKb := int64(1024)
//	bitsPerMb := bitsPerKb * 1024
//	bitsPerGb := bitsPerMb * 1024
//
//	switch {
//	case bit < bitsPerKb:
//		return int64(bit), "bit"
//	case bit < bitsPerMb:
//		return int64(bit) / int64(bitsPerKb), "Kb"
//	case bit < bitsPerGb:
//		return int64(bit) / int64(bitsPerMb), "Mb"
//	default:
//		return int64(bit) / int64(bitsPerGb), "Gb"
//	}
//}

