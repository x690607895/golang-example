package main

// export GOPATH=/Users/qiuqiu/go:/Users/qiuqiu/store/golang-example/
import (
	"fund"
	"log"
	"strconv"
	"time"
)

type MeiRiXinXi struct {
	Jzzzl   string `json:"JZZZL"`
	Dwjz    string `json:"DWJZ"`
	Navtype string `json:"NAVTYPE"`
	Fsrq    string `json:"FSRQ"`
	Ljjz    string `json:"LJJZ"`
}
type RespData struct {
	Data    []MeiRiXinXi `json:"data"`
	Success bool         `json:"success"`
}

func main() {
	zhiShuJiBenXinXi := fund.HuoQuZhiShuJiBenXinXiShuJu("SH000905")
	log.Print(zhiShuJiBenXinXi.String())
	fund.GetZhiShuXiaJiJinJianYaoXinXi("SH000905")
	// zuiDaHuiCheFunc()
}

func zuiDaHuiCheFunc() {
	// fmt.Println("请输入定投金额：")
	// var dingTouJinE int
	// fmt.Scan(&dingTouJinE)

	// fmt.Println("请选择定投方式：")
	// fmt.Println("-----1：每周定投X元")
	// fmt.Println("-----2：每天定投X元")
	// var dingTouFangShi int
	// fmt.Scan(&dingTouFangShi)

	// if dingTouFangShi == 1 {
	// 	log.Printf("每周定投%d元。", dingTouJinE)
	// } else {
	// 	log.Printf("每天定投%d元。", dingTouJinE)
	// }

	// dingTouJinE := 1000

	// // 计算普通定投收益
	// var dingTouShouYiData dingTouShouYi
	// dingTouShouYiData.maiRuMingXi = make(map[float64]float64)
	// now := time.Now()
	// startDate := time.Date(now.Year()-1, now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	// lastZhou := startDate.Month()
	// for i := 0; i < len(s.Data); i++ {
	// 	data := s.Data[i]
	// 	nowDate := strToDate(data.Fsrq)
	// 	// 记录开始时间前的最后一个交易日的净值
	// 	if nowDate.Before(startDate) || nowDate == startDate {
	// 		qiChuJinZhi := strToFloat64(data.Ljjz)
	// 		if qiChuJinZhi <= 0 {
	// 			continue
	// 		}
	// 		dingTouShouYiData.qiChuJinZhi = qiChuJinZhi
	// 	}

	// 	if nowDate.After(startDate) {
	// 		week := nowDate.Month()
	// 		if lastZhou == week {
	// 			continue
	// 		}
	// 		lastZhou = week
	// 		// 普通定投
	// 		dangRiJinZhi := strToFloat64(data.Dwjz)
	// 		gouMaiFengE := float64(dingTouJinE) / dangRiJinZhi
	// 		// gouMaiFengE = decimal(gouMaiFengE)
	// 		dingTouShouYiData.leiJiTouRu += float64(dingTouJinE)
	// 		dingTouShouYiData.leiJiFengE += gouMaiFengE
	// 		dingTouShouYiData.maiRuMingXi[dangRiJinZhi] += gouMaiFengE
	// 		dingTouShouYiData.qiMoJinZhi = dangRiJinZhi
	// 	}
	// }
	// dingTouShouYiData.qiMoZongZiChan = dingTouShouYiData.leiJiFengE * dingTouShouYiData.qiMoJinZhi
	// dingTouShouYiData.qiMoShouYi = dingTouShouYiData.qiMoZongZiChan - dingTouShouYiData.leiJiTouRu
	// log.Printf("过去一年累计投入%.0f元，累计购买%.2f份额，截止今日总资产%.2f元，收益为%.2f元", dingTouShouYiData.leiJiTouRu, dingTouShouYiData.leiJiFengE, dingTouShouYiData.qiMoZongZiChan, dingTouShouYiData.qiMoShouYi)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

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

type zuiDaHuiChe struct {
	quJian      int
	zuiDaHuiChe float64
	kaiShiRiQi  time.Time
}

type dingTouShouYi struct {
	leiJiTouRu     float64
	leiJiFengE     float64
	qiMoJinZhi     float64
	qiChuJinZhi    float64
	qiMoShouYi     float64
	qiMoZongZiChan float64
	maiRuMingXi    map[float64]float64
}
