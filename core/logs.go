package core

import "fmt"

type Logs interface {
	Info(string)
	Error(string)
	Warn(string)
}

type Wlogs struct {
	Date string
}

func (w Wlogs) Info(l string){
	log := fmt.Sprintf("[GUIDE] "+"INFO "+"%s "+l, w.Date)
	fmt.Println(log)
}

func (w Wlogs) Error(l string){
	log := fmt.Sprintf("[GUIDE] "+"ERROR "+"%s "+l, w.Date)
	fmt.Println(log)
}

func (w Wlogs) Warn(l string){
	log := fmt.Sprintf("[GUIDE] "+"WARN "+"%s "+l, w.Date)
	fmt.Println(log)
}