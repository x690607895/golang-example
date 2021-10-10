package fund

import "log"

type paiXuShuJu struct {
	fundCode string  //基金代码
	value    float64 //需要排序的值
}

type paiXuShuJuShuZu []paiXuShuJu

func NewPaiXuShuJuShuZuByDrawdown(jiJinJianYaoXinXi *map[string]zhiShuXiaJiJinJianYaoXinXi) paiXuShuJuShuZu {
	result := make(paiXuShuJuShuZu, 0, len(*jiJinJianYaoXinXi))
	for _, value := range *jiJinJianYaoXinXi {
		if value.Drawdown[3] <= 0 {
			continue
		}
		result = append(result, paiXuShuJu{value.FundCode, value.Drawdown[3]})
	}
	return result
}

func NewPaiXuShuJuShuZuByTrackError(jiJinJianYaoXinXi *map[string]zhiShuXiaJiJinJianYaoXinXi) paiXuShuJuShuZu {
	result := make(paiXuShuJuShuZu, 0, len(*jiJinJianYaoXinXi))
	for _, value := range *jiJinJianYaoXinXi {
		if value.TrackError <= 0 || value.FundType != 1 {
			continue
		}
		result = append(result, paiXuShuJu{value.FundCode, value.TrackError})
	}
	return result
}

func NewPaiXuShuJuShuZuByExcessEarnings(jiJinJianYaoXinXi *map[string]zhiShuXiaJiJinJianYaoXinXi) paiXuShuJuShuZu {
	result := make(paiXuShuJuShuZu, 0, len(*jiJinJianYaoXinXi))
	for _, value := range *jiJinJianYaoXinXi {
		if value.ExcessEarnings <= 0 || value.FundType != 2 {
			continue
		}
		result = append(result, paiXuShuJu{value.FundCode, value.ExcessEarnings})
	}
	return result
}

func (this paiXuShuJuShuZu) Len() int {
	return len(this)
}

func (this paiXuShuJuShuZu) Less(i, j int) bool {
	return this[i].value <= this[j].value
}

func (this paiXuShuJuShuZu) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this paiXuShuJuShuZu) PrintLn(jiJinJianYaoXinXi map[string]zhiShuXiaJiJinJianYaoXinXi) {
	for _, value := range this {
		object := jiJinJianYaoXinXi[value.fundCode]
		log.Println(object.String())
	}
}
