package skl

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"regexp"
	"strings"
)

type UserInfo struct {
	ID       string `json:"id"`
	UserName string `json:"userName"`
	UserType int    `json:"userType"`
	UnitID   string `json:"unitId"`
	UnitCode string `json:"unitCode"`
	UnitName string `json:"unitName"`
	Grade    string `json:"grade"`
	ClassNo  string `json:"classNo"`
	Sex      string `json:"sex"`
	Major    string `json:"major"`
	//RoleList interface{} `json:"roleList"`
}

type User struct {
	ID    string
	User  *UserInfo
	Token string
	Group string
	app   *gorequest.SuperAgent
}

var (
	recCaptchaUrl = "http://localhost:5000/captcha"
)

//get a *User with empty params, login is required to continue next operations.
func NewUser(id, group string) *User {
	return &User{
		ID:    id,
		User:  nil,
		Token: "",
		Group: group,
		app:   gorequest.New(),
	}
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

	s.app.Post(url)
	s.app.Type("form")

	rsa, err := encStr_i(user, pass, lt[1])
	if err != nil {
		rsa = encStr(user, pass, lt[1]) //this method is very slow...
	}
	s.app.Send(map[string]string{
		"rsa":       rsa,
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
	return s.loadInfo()
}

//isTokenValid
func (s *User) TokenValid() error {
	return s.loadInfo()
}

func (s *User) loadInfo() error {

	url := "https://skl.hdu.edu.cn/api/userinfo"
	resp, body, err := s.get(url)
	if resp == nil || resp.StatusCode != 200 {
		return errors.New(body)
	}
	userInfo := &UserInfo{}
	err = json.Unmarshal([]byte(body), userInfo)
	if err != nil {
		return err
	}
	s.User = userInfo
	return nil
}

func (s *User) get(url string) (gorequest.Response, string, error) {
	req := s.app.Get(url)
	req.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.79 Mobile Safari/537.36")
	req.Set("X-Auth-Token", s.Token)
	req.Set("Origin", "https://skl.hduhelp.com")
	req.Set("Referer", "https://skl.hdu.edu.cn/")

	resp, body, errs := req.End()
	return resp, body, errors.New(fmt.Sprintf("%v", errs))
}

func (s *User) CheckCode(code string, maxRetry int) error {

	for i := 0; i < maxRetry; i++ {
		url := fmt.Sprintf("https://skl.hdu.edu.cn/api/checkIn/code-check-in?code=%s", code)
		resp, body, err := s.get(url)
		if resp == nil || resp.StatusCode != 200 {
			return err
		}

		b6Img := base64.StdEncoding.EncodeToString([]byte(body))
		retCode, err := s.recCaptcha(string(b6Img))
		if err != nil {
			return err
		}
		err = s.validCode(retCode)
		if err == nil {
			return nil
		}
		if strings.Contains(err.Error(), "401") {
			return errors.New("签到码无效")
		}
		if !strings.Contains(err.Error(), "400") {
			//未知错误
			return err
		}
		//验证码识别错误，重试
	}
	return errors.New("reach the max retry times")
}

func (s *User) recCaptcha(img string) (string, error) {

	req := gorequest.New().Post(recCaptchaUrl)
	req.Type("form")
	data := map[string]string{
		"data": img,
	}
	req.Send(data)
	resp, body, errs := req.End()
	if resp == nil || resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("%v", errs))
	}
	var ret = &struct {
		Err  int    `json:"err"`
		Msg  string `json:"msg"`
		Data string `json:"data"`
	}{}
	err := json.Unmarshal([]byte(body), ret)
	if err != nil {
		return "", err
	}
	if ret.Err != 0 {
		return "", errors.New(ret.Msg)
	}
	return ret.Data, nil
}

func (s *User) validCode(code string) error {
	url := fmt.Sprintf("https://skl.hdu.edu.cn/api/checkIn/valid-code?code=%s", code)
	resp, _, err := s.get(url)
	if resp == nil {
		return err
	}

	if resp.StatusCode != 200 { //400 验证码不正确，401，签到码无效，200 ？
		return errors.New(fmt.Sprintf("%d", resp.StatusCode))
	}
	return nil
}

func RegisterUrl(url string) {
	recCaptchaUrl = url
}

func (s *User) SKLStatus(d DateRange) (*SKLCheckData, error) {
	from, to := getDataRange(d)
	url := fmt.Sprintf("https://skl.hdu.edu.cn/api/stat/stu/user?startDate=%s&endDate=%s", from, to)
	fmt.Println(url)

	resp, body, _ := s.get(url)
	if resp == nil || resp.StatusCode != 200 {
		return nil, errors.New("get response error")
	}
	if body == "[]" { //没有数据
		return nil, nil
	}
	ret := &SKLCheckData{}
	err := json.Unmarshal([]byte(body), ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *User) SKLCheckList(r DateRange) (*SKLCheckListStruct, error) {
	from, to := getDataRange(r)
	url := fmt.Sprintf("https://skl.hdu.edu.cn/api/check-in-student-detail/my?startDate=%s&endDate=%s", from, to)

	resp, body, _ := s.get(url)
	if resp == nil || resp.StatusCode != 200 {
		return nil, errors.New("get response error")
	}
	if body == "[]" { //没有数据
		return nil, nil
	}
	ret := &SKLCheckListStruct{}
	err := json.Unmarshal([]byte(body), ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *User) SKLComment(mark string) error {
	url := "https://skl.hdu.edu.cn/api/teacher-mark/save"
	req := gorequest.New()
	req.Post(url)
	req.Set("X-Auth-Token", s.Token).
		Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Mobile Safari/537.36").
		Type("json")
	req.Send(&map[string]interface{}{
		"totalMark": 5,
		"mark1":     5,
		"mark2":     5,
		"comments":  mark,
	})
	resp, _, _ := req.End()
	if resp == nil || resp.StatusCode != 200 {
		return errors.New("Comments failed")
	}
	return nil
}
