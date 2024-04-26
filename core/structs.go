package core

type (
	News struct {
		UName string
		Url   string
		Notes string
	}

	Cron struct {
		CronId      string
		CronNewDate string
		CronName    string
		CronTime    string
		CronCode    string
		CronNotes   string
	}

	ServiceTools struct {
		Id           string
		ServiceName  string
		ServiceDate  string
		StartCmd     string
		ServiceNotes string
	}

	UserPwd struct {
		Id          string
		ServiceName string
		User        string
		Passwd      string
		Notes       string
	}

	User struct {
		User string
	}

	UserAll struct {
		Id          string
		UserName    string
		NewUserDate string
		Password    string
	}

	ErrorLog struct {
		Id         string
		Date       string
		LogType    string
		LogContent string
	}
)
