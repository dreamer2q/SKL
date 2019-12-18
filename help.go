package skl

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"os"
)

var js *otto.Otto

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
