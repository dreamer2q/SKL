package skl

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"os"
	"time"
)

var js *otto.Otto

type DateRange int

const (
	Today DateRange = iota
	ThisWeek
	LastWeek
	ThisMonth
	HalfYear
)

func init() {
	desUrl := "http://cas.hdu.edu.cn/cas/comm/js/des.js"
	filepath := "des.js"
	var desJs []byte

	if exist(filepath) {
		desJs, _ = ioutil.ReadFile(filepath)
	} else {
		resp, body, _ := gorequest.New().Get(desUrl).End()
		if resp == nil || resp.StatusCode != 200 {
			panic("des.desJs Not Available")
		}
		desJs = []byte(body)
		err := ioutil.WriteFile(filepath, desJs, 0666)
		if err != nil {
			panic(err)
		}
	}
	js = otto.New()
	_, err := js.Run(desJs)
	if err != nil {
		panic(err)
	}
}

func exist(filepath string) bool {
	_, err := os.Stat(filepath)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func encStr(user, pass, lt string) string {

	ret, err := js.Call("strEnc", nil, user+pass+lt, "1", "2", "3")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", ret)

	return fmt.Sprintf("%v", ret)
}

func wton(w time.Weekday) (t int) {
	t = int(w)
	t = (t + 6) % 7
	return
}

func getDataRange(r DateRange) (string, string) {
	curr := time.Now()
	var from, to time.Time
	switch r {
	case Today:
		from = curr
		to = curr
	case ThisWeek:
		w := wton(curr.Weekday())
		from = curr.AddDate(0, 0, -w)
		to = curr.AddDate(0, 0, 6-w)
	case LastWeek:
		w := wton(curr.Weekday())
		curr = curr.AddDate(0, 0, -7)
		from = curr.AddDate(0, 0, -w)
		to = curr.AddDate(0, 0, 6-w)
	case ThisMonth:
		from = curr.AddDate(0, -1, 0)
		to = curr
	case HalfYear:
		from = curr.AddDate(0, 6, 0)
		to = curr
	}
	format := "2006-01-02"
	return from.Format(format), to.Format(format)
}

func GetCheckStatus(key string) string {
	return map[string]string{
		"3": "到",
		"0": "旷",
		"2": "迟",
		"4": "早",
		"1": "假",
	}[key]
}
func GetUpdateMode(key string) string {
	return map[string]string{
		"0": "签到码",
		"1": "手动点名",
	}[key]
}
