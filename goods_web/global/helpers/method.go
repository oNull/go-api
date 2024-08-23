package helpers

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// HttpGet Get 网络请求
func HttpGet(apiURL string, params url.Values) (rs []byte, err error) {
	var Url *url.URL
	Url, err = url.Parse(apiURL)
	if err != nil {
		fmt.Printf("解析url错误:\r\n%v", err)
		return nil, err
	}
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	resp, err := http.Get(Url.String())
	if err != nil {
		fmt.Println("err:", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("err:", err)
		}
	}(resp.Body)

	return ioutil.ReadAll(resp.Body)
}

// HttpPost Post  网络请求 ,params 是url.Values类型
func HttpPost(apiURL string, params url.Values) (rs []byte, err error) {
	resp, err := http.PostForm(apiURL, params)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("err:", err)
		}
	}(resp.Body)
	return ioutil.ReadAll(resp.Body)
}
