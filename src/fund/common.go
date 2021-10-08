package fund

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

//check err信息
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// httpGet请求
func httpGet(httpUrl string, params map[string]string) []byte {
	httpGetParams := url.Values{}
	requestUrl, err := url.Parse(httpUrl)
	checkErr(err)
	if len(params) > 0 {
		for k, v := range params {
			httpGetParams.Set(k, v)
		}
	}
	requestUrl.RawQuery = httpGetParams.Encode()
	resp, err := http.Get(requestUrl.String())
	checkErr(err)
	defer resp.Body.Close()
	respStr, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	return respStr
}

// 字符串转日期
func strToDate(str string) time.Time {
	date, err := time.ParseInLocation("2006-01-02", str, time.Local)
	checkErr(err)
	return date
}

func strToFloat64(str string) float64 {
	if str == "" {
		return 0
	}
	floatData, err := strconv.ParseFloat(str, 64)
	checkErr(err)
	return floatData
}
