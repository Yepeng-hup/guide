package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"guide/core"
	"net/http"
	"strings"
)

func writeUrlDb(urlName, urlAddr, urlType, urlNotes string) error {
	if err := core.InsertActUrl(urlName, urlAddr, urlType, urlNotes); err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}

func ShowDbUrl(c *gin.Context) {
	urlTypeRel, err := core.SelectUrlType("SELECT urlType FROM url_type")
	if err != nil {
		mlog.Error(err.Error())
	}
	c.HTML(http.StatusOK, "url2.tmpl", gin.H{
		"urlTypeList": urlTypeRel,
	})
}

type UrlTypeName struct {
	UrlType string `json:"urltype"`
}

func ShowTypeUrl(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Error getting raw data: " + err.Error(),
		})
		return
	}
	var bodys UrlTypeName
	err = json.Unmarshal(data, &bodys)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Error unmarshaling JSON: " + err.Error(),
		})
		return
	}

	urlRel, err := core.SelectUrl(fmt.Sprintf("SELECT urlName,urlAddress,urlType,urlNotes FROM url where urlType=\"%s\"", bodys.UrlType))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadGateway,
			"data": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": urlRel,
	})
}

func RewriteUrl2(c *gin.Context) {
	f := From2{
		UrlName:  c.PostForm("url-name"),
		UrlAddr:  c.PostForm("url-addr"),
		UrlType:  c.PostForm("url-type"),
		UrlNotes: c.PostForm("url-txt"),
	}
	if f.UrlName == "" && f.UrlAddr == "" {
		mlog.Warn("add element is nil.")
	} else {
		if f.UrlNotes == "" {
			f.UrlNotes = "无备注"
		}
		err := writeUrlDb(f.UrlName, f.UrlAddr, f.UrlType, f.UrlNotes)
		if err != nil {
			mlog.Error(err.Error())
		}
	}
	c.Redirect(http.StatusFound, "/url/index")
}

func DelUrl2(c *gin.Context) {
	n := Name{
		Name: c.PostForm("u-name"),
	}
	if n.Name == "" {
		mlog.Warn("del element is nil.")
	} else {
		err := core.DeleteUrl(n.Name)
		if err != nil {
			mlog.Error(err.Error())
		}
	}
	c.Redirect(http.StatusFound, "/url/index")
}

type UrlInfo2 struct {
	UrlName string `json:"urlName"`
	Url     string `json:"url"`
	UrlType string `json:"type"`
	Notes   string `json:"notes"`
}

func UpdateUrlInfo2(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Error getting raw data: " + err.Error(),
		})
		return
	}

	var body UrlInfo2
	err = json.Unmarshal(data, &body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Error unmarshaling JSON: " + err.Error(),
		})
		return
	}

	if err := core.DeleteUrl(body.UrlName); err != nil {
		mlog.Error(err.Error())
		return
	} else {
		if err := core.InsertActUrl(body.UrlName, body.Url, body.UrlType, body.Notes); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusBadGateway,
				"data": err.Error(),
			})
		} else {
			mlog.Info(fmt.Sprintf("update url ok. name -> [%s].", body.UrlName))
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
			})
		}
	}

}

func CreateType(c *gin.Context) {
	urlType := UrlType{
		TypeName: c.PostForm("type-name"),
	}
	if err := core.InsertActUrlType(urlType.TypeName); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "type create fail.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "type create ok.",
	})
}

func ShowUrlType(c *gin.Context) {
	urlTypeRel, err := core.SelectUrlType("SELECT urlType FROM url_type")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadGateway,
			"data": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": urlTypeRel,
	})
}

func DelUrlType(c *gin.Context) {
	urlType := UrlType{
		TypeName: c.PostForm("type-name"),
	}
	strList := strings.Split(urlType.TypeName, "	")
	urlRel, err := core.SelectUrl(fmt.Sprintf("SELECT urlName,urlAddress,urlType,urlNotes FROM url where urlType=\"%s\"", strList[1]))
	if err != nil {
		mlog.Error(err.Error())
	}
	urlNum := len(urlRel)
	if urlNum <= 0 && strList[1] != "other" {
		if err := core.DeleteUrlType(strList[1]); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusBadGateway,
				"data": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
			})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadGateway,
			"data": "errror: This type also has URL data.",
		})
	}

}
