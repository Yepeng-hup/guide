package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"guide/core"
	"guide/core/cmd"
	"net/http"
	"strings"
)

var (
	cronMap    = make(map[string]cron.EntryID)
	globalCron *cron.Cron
)

func init() {
	globalCron = cron.New()
}

func CoutomCron(c *gin.Context) {
	f := CronsFrom{
		Cname:  c.PostForm("cname"),
		Ctime:  c.PostForm("ctime"),
		Ccode:  c.PostForm("ccode"),
		Cnotes: c.PostForm("cnotes"),
	}
	if f.Cname == "" || f.Cnotes == "" || f.Ccode == "" || f.Ctime == "" {
		mlog.Warn("add element is nil.")
		c.Redirect(http.StatusFound, "/cron/index")
		return
	}
	// check from security
	s := core.Security{}
	if s.CheckForm(f.Ccode) {
		mlog.Error(fmt.Sprintf("detected code containing dangerous characters --> [%v]", f.Ccode))
		c.Redirect(http.StatusFound, "/cron/index")
		return
	}

	v, err := core.SelectAct("cronName", f.Cname, false)
	if err != nil {
		mlog.Error(err.Error())
	}

	if len(v) < 1 {
		core.InsertAct(f.Cname, f.Ctime, f.Ccode, f.Cnotes)
		go func() {
			entryId, err := globalCron.AddFunc(f.Ctime, func() {
				mlog.Info(fmt.Sprintf("corn use success, use name -> [%s].", f.Cname))
				// Timed task logic code
				err := cmd.UseCmd(f.Ccode)
				if err != nil {
					mlog.Error(err.Error())
					return
				}

			})
			/*
				If persistence can be written to storage, I won't use persistence here.
				Restarting will result in all failures.
			*/
			cronMap[f.Cname] = entryId
			mlog.Info(fmt.Sprintf("cron Job add with ID [%v] for name [%s]", entryId, f.Cname))
			if err != nil {
				mlog.Error(fmt.Sprintf("cron add function error,%v", err.Error()))
			}
			defer globalCron.Start()
			//select {}
		}()
	} else {
		mlog.Error("we already have the same cron name.")
		return
	}
	c.Redirect(http.StatusFound, "/cron/index")
}

func ShowCron(c *gin.Context) {
	cronList, err := core.SelectAct("", "", true)
	if err != nil {
		mlog.Error(err.Error())
		return
	}
	c.HTML(http.StatusOK, "cronnum.tmpl", gin.H{
		"cronList": cronList,
	})
}

// Delete cron func
func deleteMaptoCronjob(id cron.EntryID, jobName string) {
	delete(cronMap, jobName)
	globalCron.Remove(id)
	mlog.Info(fmt.Sprintf("delete cron job [%s] ok.", jobName))
}

// Delete cron route
func DelCron(c *gin.Context) {
	crons := c.PostForm("cron")
	cronList := strings.Fields(crons)
	v := cronMap[cronList[0]]
	deleteMaptoCronjob(v, cronList[0])
	err := core.DeleteAct(cronList[0])
	if err != nil {
		mlog.Error(fmt.Sprintf("delete cron fail -> [%s].", cronList[0]))
		return
	}
	c.Redirect(http.StatusFound, "/cron/list")
}
