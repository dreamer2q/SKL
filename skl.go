package skl

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
)

var recCaptchaUrl = "http://localhost:5000/captcha"

func RegisterUrl(url string) {
	recCaptchaUrl = url
}

//isTokenValid
func (s *User) TokenValid() error {
	return s.LoadInfo()
}

func (s *User) LoadInfo() error {

	url := "https://skl.hdu.edu.cn/api/userinfo"
	resp, body, err := s.get(url)
	if resp == nil || resp.StatusCode != 200 {
		return fmt.Errorf("loadInfo: %v", err)
	}
	userInfo := &UserInfo{}
	err = json.Unmarshal([]byte(body), userInfo)
	if err != nil {
		return err
	}

	s.UserInfo = *userInfo
	return nil
}

func (s *User) get(url string) (gorequest.Response, string, error) {
	req := app.Get(url)
	req.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.79 Mobile Safari/537.36")
	req.Set("X-Auth-Token", s.Token)
	req.Set("Origin", "https://skl.hduhelp.com")
	req.Set("Referer", "https://skl.hdu.edu.cn/")

	resp, body, errs := req.End()
	return resp, body, fmt.Errorf("%v", errs)
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
		if err, ok := err.(codeError); ok {
			if err.statusCode == 400 {
				//签到码识别错误，retry
				continue
			}
		}
		return err
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
	if resp.StatusCode != 200 {
		return codeError{statusCode: resp.StatusCode}
	}
	return nil
}

func (s *User) SKLStatus(d DateRange) (*SKLCheckData, error) {
	from, to := getDataRange(d)
	url := fmt.Sprintf("https://skl.hdu.edu.cn/api/stat/stu/user?startDate=%s&endDate=%s", from, to)

	resp, body, err := s.get(url)
	if resp == nil || resp.StatusCode != 200 {
		return nil, err
	}
	if body == "[]" { //没有数据
		return nil, nil
	}
	ret := &SKLCheckData{}
	err = json.Unmarshal([]byte(body), ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *User) SKLCheckList(r DateRange) (*SKLCheckListStruct, error) {
	from, to := getDataRange(r)
	url := fmt.Sprintf("https://skl.hdu.edu.cn/api/check-in-student-detail/my?startDate=%s&endDate=%s", from, to)

	resp, body, err := s.get(url)
	if resp == nil || resp.StatusCode != 200 {
		return nil, err
	}
	if body == "[]" { //没有数据
		return nil, nil
	}
	ret := &SKLCheckListStruct{}
	err = json.Unmarshal([]byte(body), ret)
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
	resp, _, err := req.End()
	if resp == nil || resp.StatusCode != 200 {
		return fmt.Errorf("comment error: %v", err)
	}
	return nil
}
