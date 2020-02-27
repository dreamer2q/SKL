package skl

import (
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"regexp"
	"strings"
)

type User struct {
	UserID string
	Token  string
}

//login in with HDU, and get a token for SKL
func (u *User) Login(user, pass string) error {

	req := gorequest.New()
	url := "https://cas.hdu.edu.cn/cas/login?state=&service=https%3A%2F%2Fskl.hdu.edu.cn%2Fapi%2Fcas%2Flogin%3Findex%3D"
	resp, body, errs := req.Get(url).End()

	if resp == nil || resp.StatusCode != 200 {
		return fmt.Errorf("login: %v", errs)
	}

	r := regexp.MustCompile("(<span.*?id=\"errormsghide\".*?>|<p.*?class=\"unauthorise_p\">)(.*?)<")
	if r.MatchString(body) {
		return errors.New(fmt.Sprintf("login: %v", r.FindString(body)))
	}

	r = regexp.MustCompile("<input.*?name=\"lt\".*?value=\"(.*?)\"")
	lt := r.FindStringSubmatch(body)
	r = regexp.MustCompile("<input.*?name=\"execution\".*?value=\"(.*?)\"")
	execution := r.FindStringSubmatch(body)

	if len(lt) != 2 || len(execution) != 2 {
		return fmt.Errorf("login: hdu page error")
	}

	req.Post(url)
	req.Type("form")

	rsa, err := encStr_i(user, pass, lt[1])
	//fast method failed, try another but slow method
	if err != nil {
		rsa = encStr(user, pass, lt[1])
	}

	req.Send(map[string]string{
		"rsa":       rsa,
		"ul":        fmt.Sprintf("%d", len(user)),
		"pl":        fmt.Sprintf("%d", len(pass)),
		"lt":        lt[1],
		"execution": execution[1],
		"_eventId":  "submit",
	})
	req.RedirectPolicy(func(req gorequest.Request, via []gorequest.Request) error {
		return nil
	})

	resp, body, errs = req.End()
	if resp == nil || resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("login: %v", errs))
	}

	r = regexp.MustCompile("(<span.*?id=\"errormsghide\".*?>|<p.*?class=\"unauthorise_p\">)(.*?)<")
	if r.MatchString(body) {
		return errors.New(fmt.Sprintf("login: %v\n", r.FindString(body)))
	}

	retUrl := resp.Request.URL.String()
	if strings.Contains(retUrl, "token=") {
		pos := strings.Index(retUrl, "token=")
		u.Token = retUrl[pos+6:]
	} else {
		return errors.New("login: token not found")
	}

	return u.LoadInfo()
}

//login in with Token
func (u *User) LoginByToken(token string) error {
	u.Token = token
	return u.LoadInfo()
}
