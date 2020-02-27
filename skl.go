package skl

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
)

var recCaptchaUrl = "http://localhost:5000/captcha"

func (u *User) TokenValid() error {
	return u.LoadInfo()
}

func (u *User) LoadInfo() error {

	url := "https://skl.hdu.edu.cn/api/userinfo"
	resp, body, err := u.get(url)
	if resp == nil || resp.StatusCode != 200 {
		return fmt.Errorf("loadInfo: %v", err)
	}
	userInfo := &UserInfo{}
	err = json.Unmarshal([]byte(body), userInfo)
	if err != nil {
		return err
	}
	u.UserID = userInfo.UserID
	return nil
}

func (u *User) get(url string) (gorequest.Response, string, error) {
	if u.Token == "" {
		return nil, "", fmt.Errorf("get: token is empty")
	}

	req := gorequest.New().Get(url)
	req.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.79 Mobile Safari/537.36")
	req.Set("X-Auth-Token", u.Token)
	req.Set("Origin", "https://skl.hduhelp.com")
	req.Set("Referer", "https://skl.hdu.edu.cn/")

	resp, body, errs := req.End()
	return resp, body, fmt.Errorf("%v", errs)
}

func (u *User) CheckCode(code string, maxRetry int) error {

	for i := 0; i < maxRetry; i++ {
		url := fmt.Sprintf("https://skl.hdu.edu.cn/api/checkIn/code-check-in?code=%s", code)
		resp, body, err := u.get(url)
		if resp == nil || resp.StatusCode != 200 {
			return err
		}
		url2 := "https://skl.hdu.edu.cn/api/checkIn/create-code-img"
		resp, body, err = u.get(url2)
		if resp == nil || resp.StatusCode != 200 {
			return err
		}

		b6Img := base64.StdEncoding.EncodeToString([]byte(body))
		retCode, err := u.recCaptcha(string(b6Img))
		if err != nil {
			return err
		}
		err = u.validCode(retCode)
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
	return errors.New("checkCode: reach the max retry times")
}

func (u *User) recCaptcha(img string) (string, error) {

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

func (u *User) validCode(code string) error {
	url := fmt.Sprintf("https://skl.hdu.edu.cn/api/checkIn/valid-code?code=%u", code)
	resp, _, err := u.get(url)
	if resp == nil {
		return err
	}
	if resp.StatusCode != 200 {
		return codeError{statusCode: resp.StatusCode}
	}
	return nil
}

func (u *User) Status(d DateRange) (*CheckData, error) {
	from, to := getDataRange(d)
	url := fmt.Sprintf("https://skl.hdu.edu.cn/api/stat/stu/user?startDate=%u&endDate=%u", from, to)

	resp, body, err := u.get(url)
	if resp == nil || resp.StatusCode != 200 {
		return nil, err
	}
	if body == "[]" { //没有数据
		return nil, nil
	}
	ret := &CheckData{}
	err = json.Unmarshal([]byte(body), ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (u *User) CheckList(r DateRange) (*CheckListData, error) {
	from, to := getDataRange(r)
	url := fmt.Sprintf("https://skl.hdu.edu.cn/api/check-in-student-detail/my?startDate=%u&endDate=%u", from, to)

	resp, body, err := u.get(url)
	if resp == nil || resp.StatusCode != 200 {
		return nil, err
	}
	if body == "[]" { //没有数据
		return nil, nil
	}
	ret := &CheckListData{}
	err = json.Unmarshal([]byte(body), ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (u *User) Comment(mark string) error {
	url := "https://skl.hdu.edu.cn/api/teacher-mark/save"
	req := gorequest.New()
	req.Post(url)
	req.Set("X-Auth-Token", u.Token).
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

func (u *User) Info() (*UserInfo, error) {
	url := "https://skl.hdu.edu.cn/api/userinfo"
	resp, body, err := u.get(url)
	if resp == nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("info: %v", err)
	}
	userInfo := &UserInfo{}
	err = json.Unmarshal([]byte(body), userInfo)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}
