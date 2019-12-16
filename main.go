package skl

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"github.com/parnurzeal/gorequest"
	"github.com/robertkrimen/otto"

)

type SKL struct {
	user string
	pass string
	token string
	app *gorequest.SuperAgent
}

var js *otto.Otto

func init() {
	desUrl := "http://cas.hdu.edu.cn/cas/comm/desJs/des.desJs"

	filepath := "des.desJs"
	var desJs []byte

	 if exist(filepath) {
	 	desJs,_ = ioutil.ReadFile(filepath)
	 }else{
	 	resp,body,_ := gorequest.New().Get(desUrl).End()
	 	if resp == nil ||resp.StatusCode != 200 {
	 		panic("des.desJs Not Available")
		}
		desJs = []byte(body)
		err := ioutil.WriteFile(filepath, desJs,0666)
		if err != nil {
			panic(err)
		}
	 }
	 js := otto.New()
	 _,err := js.Run(desJs)
	 if err != nil {
	 	panic(err)
	 }
}

func exist(filepath string)bool {
	_,err := os.Stat(filepath)
	if err != nil{
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func encStr(user, pass, lt string) string {

	ret,err := js.Call("strEnc",nil,user+pass+lt,"1","2","3")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", ret)

	return fmt.Sprintf("%v",ret)
}

func (s *SKL)New(user,pass string) *SKL {
	return &SKL{
		user:  user,
		pass: pass,
		token: "",
		app: gorequest.New(),
	}
}

func (s *SKL)getParams() (map[string]string, error) {
	req := gorequest.New()
	url := "https://cas.hdu.edu.cn/cas/login?state=&service=https%3A%2F%2Fskl.hdu.edu.cn%2Fapi%2Fcas%2Flogin%3Findex%3D"
	resp, body, errs := req.Get(url).End()
	if resp == nil || resp.StatusCode != 200 {
		log.Panic(errs)
	}

	ok, err := regexp.MatchString("(<span.*?id=\"errormsghide\".*?>|<p.*?class=\"unauthorise_p\">)(.*?)<", body)
	if err != nil {
		log.Panic(err)
	}
	if !ok {
		return nil, errors.New("Error matched")
	}

	r := regexp.MustCompile("<input.*?name=\"lt\".*?value=\"(.*?)\"")
	lt := r.FindStringSubmatch(body)
	r = regexp.MustCompile("<input.*?name=\"execution\".*?value=\"(.*?)\"")
	execution := r.FindStringSubmatch(body)

	//fmt.Println(resp.Header["Set-Cookie"])
	var cookie string
	for _, t := range resp.Header["Set-Cookie"] {
		if strings.Contains(t, "JSESSIONID") {
			cookie = t
		}
	}

	return map[string]string{
		"lt":        lt[1],
		"execution": execution[1],
		"cookie":    cookie,
	}, nil
}

func (s *SKL)login(user, pass string) error {

	data, err := s.getParams()
	if err != nil {
		return err
	}

	urlquery := "https://cas.hdu.edu.cn/cas/login?state=&service=https%3A%2F%2Fskl.hdu.edu.cn%2Fapi%2Fcas%2Flogin%3Findex%3D"

	s.app.Post(urlquery)
	s.app.Set("Origin", "https://cas.hdu.edu.cn").
		Set("Content-Type", "application/x-www-form-urlencoded").
		Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)").
		Set("Referer", urlquery).
		Set("Cookie", data["cookie"])
	s.app.Send(map[string]string{
		"rsa":       encStr(user, pass, data["lt"]),
		"ul":        fmt.Sprintf("%v", len(user)),
		"pl":        fmt.Sprintf("%v", len(pass)),
		"lt":        data["lt"],
		"execution": data["execution"],
		"_eventId":  "submit",
	})
	s.app.RedirectPolicy(func(req gorequest.Request, via []gorequest.Request) error {
		return nil
	})

	resp, _, errs := s.app.End()
	if resp == nil || resp.StatusCode != 200 {
		return errs[0]
	}

	retUrl := resp.Request.URL.String()
	if strings.Contains(retUrl, "token=") {
		pos := strings.Index(retUrl, "token=")
		s.token = retUrl[pos+6:]
		s.app.Set("")
	} else {
		return errors.New("Token not found")
	}
	return nil
}

