package fund

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"sync"
	"time"
)

type zhiShuXiaJiJinJianYaoXinXi struct {
	FundCode       string          `json:"fund_code"`       //基金代码
	FundName       string          `json:"fund_name"`       //基金名称
	NavGrwR1m      float64         `json:"nav_grw_r1m"`     //过去一月收益（百分比）
	NavGrwRiy      float64         `json:"nav_grw_r1y"`     //过去一年收益（百分比）
	TrackError     float64         `json:"track_error"`     //跟踪误差（百分比）（ 被动指数型才有）
	ExcessEarnings float64         `json:"excess_earnings"` //超额收益（百分比）（增强指数型才有）
	Drawdown       map[int]float64 //最大回撤（百分比） key为年，value为值
	FundType       int             //基金类型（1:被动指数型，2:指数增强型）
}

type JiJinMeiRiXinXi struct {
	Jzzzl   string `json:"JZZZL"`   //不清楚
	Dwjz    string `json:"DWJZ"`    //单位净值
	Navtype string `json:"NAVTYPE"` //不清楚
	Fsrq    string `json:"FSRQ"`    //发生日期
	Ljjz    string `json:"LJJZ"`    //累计净值
}
type JiJinMeiRiXinXiFanHuiShuJu struct {
	Data    []JiJinMeiRiXinXi `json:"data"`    //数据流
	Success bool              `json:"success"` //是否成功
}

type zuiDaHuiCheShuJu struct {
	quJian     int       //区间：（1年、3年、5年）
	kaiShiRiQi time.Time //开始时间
}

func GetZhiShuXiaJiJinJianYaoXinXi(zhiShuDaiMa string) map[string]zhiShuXiaJiJinJianYaoXinXi {
	wg := sync.WaitGroup{}
	data := make(chan zhiShuXiaJiJinJianYaoXinXi, 100)
	url := "https://danjuanapp.com/djapi/fundx/base/index/traces"
	params := make(map[string]string)
	params["symbol"] = zhiShuDaiMa
	respStr := httpGet(url, params)
	var result map[string]interface{}
	err := json.Unmarshal(respStr, &result)
	checkErr(err)
	items := result["data"].(map[string]interface{})["items"].([]interface{})
	result2 := make(map[string]zhiShuXiaJiJinJianYaoXinXi)
	for _, item := range items {
		fundTypeStr := item.(map[string]interface{})["fund_type"].(string)
		fundTypeInt := 0
		if fundTypeStr == "被动指数型" {
			fundTypeInt = 1
		} else if fundTypeStr == "增强指数型" {
			fundTypeInt = 2
		}

		fundList := item.(map[string]interface{})["funds"].([]interface{})
		for _, fund := range fundList {
			wg.Add(1)
			go jiSuanJianYaoXinXi(fundTypeInt, fund, &wg, data)
		}
	}
	wg.Wait()

	close(data)
	for {
		jiJinJianYaoXinXi, ok := <-data
		if !ok {
			break
		}
		result2[jiJinJianYaoXinXi.FundCode] = jiJinJianYaoXinXi
	}
	return result2
}

func jiSuanJianYaoXinXi(fundType int, fund interface{}, wg *sync.WaitGroup, dataChannel chan zhiShuXiaJiJinJianYaoXinXi) {
	defer wg.Done()
	var jiJinJianYaoXinXi zhiShuXiaJiJinJianYaoXinXi
	fundInfo := fund.(map[string]interface{})
	jiJinJianYaoXinXi.FundCode = fundInfo["fund_code"].(string)
	jiJinJianYaoXinXi.FundName = fundInfo["fund_name"].(string)
	jiJinJianYaoXinXi.FundType = fundType
	if fundInfo["nav_grw_r1m"] != nil {
		jiJinJianYaoXinXi.NavGrwR1m = fundInfo["nav_grw_r1m"].(float64)
	}
	if fundInfo["nav_grw_r1y"] != nil {
		jiJinJianYaoXinXi.NavGrwRiy = fundInfo["nav_grw_r1y"].(float64)
	}
	switch fundType {
	case 1:
		if fundInfo["track_error"] != nil {
			jiJinJianYaoXinXi.TrackError = fundInfo["track_error"].(float64)
		}
	case 2:
		if fundInfo["excess_earnings"] != nil {
			jiJinJianYaoXinXi.ExcessEarnings = fundInfo["excess_earnings"].(float64)
		}
	default:
		log.Fatal(fundInfo)
	}
	jiJinJianYaoXinXi.jiSuanZuiDaHuiChe()
	dataChannel <- jiJinJianYaoXinXi
}

func (this *zhiShuXiaJiJinJianYaoXinXi) String() string {
	var result string
	result += fmt.Sprintf("当前基金名称：%s，", this.FundName)
	result += fmt.Sprintf("当前基金代码：%s，", this.FundCode)
	if this.FundType == 1 {
		result += "当前基金类型：被动指数型，"
		result += fmt.Sprintf("跟踪误差：%.2f%%，", this.TrackError)
	} else if this.FundType == 2 {
		result += "当前基金类型：增强指数型，"
		result += fmt.Sprintf("超额收益%.2f%%，", this.ExcessEarnings)
	}
	result += fmt.Sprintf("过去一月收益：%.2f%%，", this.NavGrwR1m)
	result += fmt.Sprintf("过去一年收益：%.2f%%，", this.NavGrwRiy)
	for year, drawDown := range this.Drawdown {
		result += fmt.Sprintf("最近%d年最大回撤率：%.2f%%。", year, drawDown*100)
	}
	return result
}

func (this *zhiShuXiaJiJinJianYaoXinXi) jiSuanZuiDaHuiChe() {
	url := "https://uni-fundts.1234567.com.cn/dataapi/fund/FundNetDiagram2"
	params := make(map[string]string)
	params["CODE"] = this.FundCode
	params["FCODE"] = this.FundCode
	params["RANGE"] = "ln"
	params["CustomerNo"] = ""
	params["UserId"] = ""
	params["Uid"] = ""
	params["CToken"] = ""
	params["UToken"] = ""
	params["MobileKey"] = ""
	params["DATES"] = ""
	params["POINTCOUNT"] = ""
	params["deviceid"] = "5C1164BA-3D4A-4953-A488-00E47222D4BB"
	params["plat"] = "Iphone"
	params["AppType"] = "Iphone"
	params["product"] = "EFund"
	params["version"] = "6.2.5"
	params["Serverversion"] = "6.2.5"
	params["appversion"] = "6.2.5"

	resp := httpGet(url, params)
	var respData JiJinMeiRiXinXiFanHuiShuJu
	err := json.Unmarshal(resp, &respData)
	checkErr(err)
	zuiDaHuiCheMap := getDrawDownStaticData()
	this.Drawdown = make(map[int]float64)

	i := 0
	for _, v := range respData.Data {
		nowDate := strToDate(v.Fsrq)
		i++
		for _, zuiDaHuiCheXinXi := range zuiDaHuiCheMap {
			if nowDate.After(zuiDaHuiCheXinXi.kaiShiRiQi) && !zuiDaHuiCheXinXi.kaiShiRiQi.IsZero() {
				dangRiLeiJiJinZhi := strToFloat64(v.Ljjz)
				if dangRiLeiJiJinZhi <= 0 {
					continue
				}
				zuiXiaoJinZhi := 999.99
				for j := i; j < len(respData.Data); j++ {
					dangRiLeiJiJinZhi2 := strToFloat64(respData.Data[j].Ljjz)
					if dangRiLeiJiJinZhi2 <= 0 {
						continue
					}
					zuiXiaoJinZhi = math.Min(zuiXiaoJinZhi, dangRiLeiJiJinZhi2)
				}
				zuiDaHuiChe := math.Dim(dangRiLeiJiJinZhi, zuiXiaoJinZhi) / dangRiLeiJiJinZhi
				if zuiDaHuiChe <= 0 {
					continue
				}
				this.Drawdown[zuiDaHuiCheXinXi.quJian] = math.Max(this.Drawdown[zuiDaHuiCheXinXi.quJian], zuiDaHuiChe)
			}
		}
	}
}

func getDrawDownStaticData() []zuiDaHuiCheShuJu {
	now := time.Now()
	zuiDaHuiChe := make([]zuiDaHuiCheShuJu, 3)
	zuiDaHuiChe[0] = zuiDaHuiCheShuJu{
		quJian:     1,
		kaiShiRiQi: time.Date(now.Year()-1, now.Month(), now.Day(), 0, 0, 0, 0, time.Local),
	}
	zuiDaHuiChe[1] = zuiDaHuiCheShuJu{
		quJian:     3,
		kaiShiRiQi: time.Date(now.Year()-3, now.Month(), now.Day(), 0, 0, 0, 0, time.Local),
	}
	zuiDaHuiChe[2] = zuiDaHuiCheShuJu{
		quJian:     5,
		kaiShiRiQi: time.Date(now.Year()-5, now.Month(), now.Day(), 0, 0, 0, 0, time.Local),
	}
	return zuiDaHuiChe
}
