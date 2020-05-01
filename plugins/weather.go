package plugins

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	iotqq "zuanGirl/model"
)

//todo 天气
func GetWeather(cm []string,qq string) string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}} //如果需要测试自签名的证书 这里需要设置跳过证书检测 否则编译报错
	client := &http.Client{Transport: tr}
	var req *http.Request
	tmp := iotqq.GetCook()
	req, _ = http.NewRequest("GET", "https://weather.mp.qq.com/cgi-bin/rich?g_tk="+tmp.Gtk+"&city="+url.PathEscape(cm[1]), nil)
	c1 := &http.Cookie{Name: "uin", Value: qq, Path: "/", Domain: ".weather.mp.qq.com"}
	c2 := &http.Cookie{Name: "skey", Value: tmp.Skey, Path: "/", Domain: ".weather.mp.qq.com"}
	req.AddCookie(c1)
	req.AddCookie(c2)
	req.Header.Add("Referer", "http://weather.mp.qq.com/ark")
	req.Header.Add("User-Agent", "PostmanRuntime/7.20.1")
	req.Header.Add("Accept", "PostmanRuntime/7.20.1")
	req.Header.Add("Content-Type", "text/json: charset=UTF-8")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	var weather iotqq.Weather
	err = json.Unmarshal([]byte(string(b)), &weather)
	if err != nil {
		fmt.Println("反序列化出错,info:", err)
	}
	m := "{\"app\":\"com.tencent.weather\",\"desc\":\"天气\",\"view\":\"RichInfoView\",\"ver\":\"1.0.0.217\",\"prompt\":\"[应用]天气\",\"meta\":{\"richinfo\":{\"adcode\":\"%s\",\"air\":\"%s\",\"city\":\"%s\",\"date\":\"%s\",\"max\":\"%s\",\"min\":\"%s\",\"ts\":\"1554951408\",\"type\":\"%s\",\"wind\":\"%s\"}},\"config\":{\"forward\":1,\"autosize\":1,\"type\":\"card\"}}"
	n := fmt.Sprintf(m, weather.Data.Adcode, weather.Data.Air, weather.Data.City, weather.Data.Date, weather.Data.Max, weather.Data.Min, weather.Data.Type, weather.Data.Wind)
	log.Println(n)
	return n
}
