package plugins

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//todo 祖安话
var minlevelApiurl = "https://nmsl.shadiao.app/api.php?level=min&lang=zh_cn"
var apiurl = "https://nmsl.shadiao.app/api.php?lang=zh_cn"

func GetZuanResult(level string) string {
	url := apiurl
	if level == "min" {
		url = minlevelApiurl
	}
	reqest, err := http.NewRequest("GET", url, nil) //建立一个请求
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(0)
	}
	//Add 头协议
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")
	client := &http.Client{}
	response, err := client.Do(reqest) //提交
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	cookies := response.Cookies() //遍历cookies
	for _, cookie := range cookies {
		fmt.Println("cookie:", cookie)
	}

	body, err1 := ioutil.ReadAll(response.Body)
	if err1 != nil {
		// handle error
	}
	fmt.Println(string(body)) //网页源码

	return string(body)

}
