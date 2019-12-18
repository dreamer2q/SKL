package skl

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"regexp"
	"strings"
)

type UserInfo struct {
	ID string `json:"id"`
	UserName string `json:"userName"`
	UserType int `json:"userType"`
	UnitID string `json:"unitId"`
	UnitCode string `json:"unitCode"`
	UnitName string `json:"unitName"`
	Grade string `json:"grade"`
	ClassNo string `json:"classNo"`
	Sex string `json:"sex"`
	Major string `json:"major"`
	//RoleList interface{} `json:"roleList"`
}

type User struct {
	ID  string
	User *UserInfo
	Token string
	Group string
	app   *gorequest.SuperAgent
}


//get a *User with empty params, login is required to continue next operations.
func NewUser(id,group string) *User {
	user := &User{
		ID:    id,
		User:  nil,
		Token: "",
		Group: group,
		app:   gorequest.New(),
	}
	user.app.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.79 Mobile Safari/537.36")
	return user
}

//login in with HDU, and get a token for SKL
func (s *User) Login(user, pass string) error {

	url := "https://cas.hdu.edu.cn/cas/login?state=&service=https%3A%2F%2Fskl.hdu.edu.cn%2Fapi%2Fcas%2Flogin%3Findex%3D"
	resp, body, errs := s.app.Get(url).End()
	if resp == nil || resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Errors: %s \n", errs))
	}

	r := regexp.MustCompile("(<span.*?id=\"errormsghide\".*?>|<p.*?class=\"unauthorise_p\">)(.*?)<")
	ok := r.MatchString(body)
	if ok {
		return errors.New(fmt.Sprintf("Error: %s\n", r.FindString(body)))
	}

	r = regexp.MustCompile("<input.*?name=\"lt\".*?value=\"(.*?)\"")
	lt := r.FindStringSubmatch(body)
	r = regexp.MustCompile("<input.*?name=\"execution\".*?value=\"(.*?)\"")
	execution := r.FindStringSubmatch(body)

	if len(lt) != 2 || len(execution) != 2 {
		return errors.New("Missing some params ")
	}

	//urlquery := "https://cas.hdu.edu.cn/cas/login?state=&service=https%3A%2F%2Fskl.hdu.edu.cn%2Fapi%2Fcas%2Flogin%3Findex%3D"
	//urlquery := "http://116.62.36.144/"
	s.app.Post(url)
	s.app.Type("form") 	//s.app.Set("Content-Type", "application/x-www-form-urlencoded")

	s.app.Send(map[string]string{
		"rsa":       encStr(user, pass, lt[1]),	//this encStr is very slow...
		"ul":        fmt.Sprintf("%d", len(user)),
		"pl":        fmt.Sprintf("%d", len(pass)),
		"lt":        lt[1],
		"execution": execution[1],
		"_eventId":  "submit",
	})

	s.app.RedirectPolicy(func(req gorequest.Request, via []gorequest.Request) error {
		return nil
	})

	resp, body, errs = s.app.End()
	if resp == nil || resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Errors: %s \n", errs))
	}

	r = regexp.MustCompile("(<span.*?id=\"errormsghide\".*?>|<p.*?class=\"unauthorise_p\">)(.*?)<")
	ok = r.MatchString(body)
	if ok {
		return errors.New(fmt.Sprintf("Error: %s\n", r.FindString(body)))
	}

	retUrl := resp.Request.URL.String()
	if strings.Contains(retUrl, "token=") {
		pos := strings.Index(retUrl, "token=")
		s.Token = retUrl[pos+6:]
	} else {
		return errors.New("Token not found")
	}

	return s.loadInfo()
}

//login in with Token
func (s *User) LoginToken(token string) error {
	s.Token = token
	s.app.Set("X-Auth-Token",s.Token)
	return s.loadInfo()
}

func (s *User) loadInfo() error {
	//s.app.Set("X-Auth-Token",s.Token)
	url := "https://skl.hdu.edu.cn/api/userinfo"
	resp,body,errs := s.app.Get(url).End()
	if resp == nil || resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("%v",errs))
	}
	userInfo := &UserInfo{}
	err := json.Unmarshal([]byte(body),userInfo)
	if err != nil {
		return err
	}
	return nil
}

func (s *User) Test() {
	s.loadInfo()
}

func (* User) CheckCode(code string) {

}