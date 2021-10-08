package fund

import (
	"encoding/json"
	"fmt"
	"log"
)

type zhuShuJiBenXinXi struct {
	Id        int     `json:"id"`
	IndexCode string  `json:"index_code"`    //基金代码
	Name      string  `json:"name"`          //基金名称
	Pe        float64 `json:"pe"`            //pe
	PePer     float64 `json:"pe_percentile"` //pe百分位(百分比)
	Peg       float64 `json:"peg"`           //预测pe
	Pb        float64 `json:"pb"`            //pb
	PbPer     float64 `json:"pb_percentile"` //pb百分位(百分比)
	Roe       float64 `json:"roe"`           //roe(百分比)
	Yeild     float64 `json:"yeild"`         //股息率(百分比)
}

type zhiShuJiBenXinXiFanHuiShuJu struct {
	Data        zhuShuJiBenXinXi `json:"data"`
	Result_code int              `json:"result_code"`
}

func HuoQuZhiShuJiBenXinXiShuJu(zhiShuDaiMa string) zhuShuJiBenXinXi {
	RequestUrl := fmt.Sprintf("https://danjuanapp.com/djapi/index_eva/detail/%s", zhiShuDaiMa)
	responseDataByteArr := httpGet(RequestUrl, make(map[string]string))
	var result zhiShuJiBenXinXiFanHuiShuJu
	err := json.Unmarshal(responseDataByteArr, &result)
	checkErr(err)

	// 部分百分比数据x100
	log.Println(result.Data)
	result.Data.PePer *= 100
	result.Data.PbPer *= 100
	result.Data.Roe *= 100
	result.Data.Yeild *= 100
	return result.Data
}

func (this *zhuShuJiBenXinXi) String() string {
	var result string
	result += fmt.Sprintln()
	result += fmt.Sprintln("以下是该代码基本信息：")
	result += fmt.Sprintln("代码：", this.IndexCode)
	result += fmt.Sprintln("名字：", this.Name)
	result += fmt.Sprintln("PE：", this.Pe)
	result += fmt.Sprintf("PE百分位：%.2f%%\r\n", this.PePer)
	result += fmt.Sprintln("预测PE：", this.Peg)
	result += fmt.Sprintln("PB：", this.Pb)
	result += fmt.Sprintf("PB百分位：%.2f%%\r\n", this.PbPer)
	result += fmt.Sprintln("ROE", this.Roe, "%")
	result += fmt.Sprintln("股息率：", this.Yeild, "%")
	if this.Pe <= 30 {
		result += fmt.Sprintln("估值区间：低估")
	} else if this.Pe > 30 && this.Pe < 70 {
		result += fmt.Sprintln("估值区间：正常")
	} else {
		result += fmt.Sprintln("估值区间：高估")
	}
	return result
}
