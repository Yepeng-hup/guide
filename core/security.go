package core

import (
	"strings"
)

type Security struct {

}

func dangerKeywords(formInput string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(formInput, keyword) {
			return true
		}
	}
	return false
}

func (s *Security)CheckForm(formStrList ...string) bool{
	dangerKeywordsList := []string{"rm", "delete", ">", "<", "|", "&", "&&", "||", ";", "..", "!", "^", "$", "@", "#", "*"}

	for _, v := range formStrList {
		if dangerKeywords(v, dangerKeywordsList) {
			return true
		}
	}
	return false
}
