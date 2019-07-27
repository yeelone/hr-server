package buildinfunc

import (
	"fmt"
	"hr-server/model"
	"hr-server/pkg/formula"
	"hr-server/util"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/lexkong/log"
	yaml "gopkg.in/yaml.v2"
)

// 因为Addup 需要从数据库上取得以前的数据，如果第一次调用Addup查询一次数据库，会比较浪费。在职工数量不是太大的情况下，可以这样设置缓存 。之后 直接从缓存中取
// map 的格式为： 缓存key => profile_id => month => SalaryField
// 其中缓存key 为 key +"-"+year+"-"+fromMonth+"-"+toMonth
var fieldCache = make(map[string]map[uint64]map[string]model.SalaryField)

// 因为LastMonthValue 需要从数据库上取得上个月的数据，如果第一次调用Addup查询一次数据库，会比较浪费。在职工数量不是太大的情况下，可以这样设置缓存 。之后 直接从缓存中取
// map 的格式为： 缓存key => profile_id  => SalaryField
// 其中缓存key 为 key +"-"+year+"-"+ month
var lastMonthValueCache = make(map[string]map[uint64]model.SalaryField)

// CalcuWorkingAge
// @params date:string , 2018-07
// @return year:number
func (b *BuildinFunc) WorkingAge(date string) float64 {
	newDate := strings.Split(date, "-")
	old, err := strconv.Atoi(newDate[0])
	if err != nil {
		return -1.00
	}
	year := time.Now().Year()

	return float64(year - old + 1.00)
}

func (b *BuildinFunc) BaseSalary() float64 {
	base, err := model.GetSalaryConfig()
	if err != nil {
		log.Info("内置函数调用出错，函数名：GetBaseSalary,出错信息:" + err.Error())
		return 0.00
	}
	return base.Base
}

func (b *BuildinFunc) Taxable(salaryBeforeTax float64) float64 {
	yamlFile, err := ioutil.ReadFile("conf/salary.yaml")

	m := model.TaxConf{}
	err = yaml.Unmarshal(yamlFile, &m)
	if err != nil {
		fmt.Println(err)
	}

	payment := salaryBeforeTax - m.Threshold
	tax := 0.00
	switch {
	case payment < 0.00:
		tax = 0.00
	case payment <= m.Level[0]:
		tax = payment * m.Rating[0] / 100
	case payment <= m.Level[1]:
		tax = payment*m.Rating[1]/100 - m.Deduction[1]
	case payment <= m.Level[2]:
		tax = payment*m.Rating[2]/100 - m.Deduction[2]
	case payment <= m.Level[3]:
		tax = payment*m.Rating[3]/100 - m.Deduction[3]
	case payment <= m.Level[4]:
		tax = payment*m.Rating[4]/100 - m.Deduction[4]
	case payment <= m.Level[5]:
		tax = payment*m.Rating[5]/100 - m.Deduction[5]
	case payment > m.Level[5]:
		tax = payment*m.Rating[6]/100 - m.Deduction[6]
	}
	return tax
}

func (b *BuildinFunc) Taxable2019(salaryBeforeTax float64) float64 {
	yamlFile, err := ioutil.ReadFile("conf/pre_deduction_rate.yaml")

	m := model.TaxConf{}
	err = yaml.Unmarshal(yamlFile, &m)
	if err != nil {
		fmt.Println(err)
	}

	payment := salaryBeforeTax
	tax := 0.00
	switch {
	case payment < 0.00:
		tax = 0.00
	case payment <= m.Level[0]:
		tax = payment * m.Rating[0] / 100
	case payment <= m.Level[1]:
		tax = payment*m.Rating[1]/100 - m.Deduction[1]
	case payment <= m.Level[2]:
		tax = payment*m.Rating[2]/100 - m.Deduction[2]
	case payment <= m.Level[3]:
		tax = payment*m.Rating[3]/100 - m.Deduction[3]
	case payment <= m.Level[4]:
		tax = payment*m.Rating[4]/100 - m.Deduction[4]
	case payment <= m.Level[5]:
		tax = payment*m.Rating[5]/100 - m.Deduction[5]
	case payment > m.Level[5]:
		tax = payment*m.Rating[6]/100 - m.Deduction[6]
	}
	return tax
}

// Addup : 累计计算
// @params : key 即salary fields key ,表示需要累计的key
// @params : fromMonth , toMonth 在指定的月份之间进行合计，如果fromMonth,toMonth 为空字符串，则默认合计从当年1月份到核算当月的前一个月
// @params : currentMonth ,这个是为了构建缓存
func (b *BuildinFunc) Addup(profile uint64, currentMonth, key, fromMonth, toMonth string) float64 {
	if currentMonth == "1" || currentMonth == "01" {
		return 0.00
	}
	if len(fromMonth) == 1 {
		fromMonth = "0" + fromMonth
	}
	if len(toMonth) == 1 {
		toMonth = "0" + toMonth
	}

	currentTime := time.Now()
	year := strconv.Itoa(currentTime.Year())
	cacheKey := key + "-" + currentMonth + "-" + year + "-" + fromMonth + "-" + toMonth

	if _, ok := fieldCache[cacheKey]; !ok {
		fieldCache[cacheKey] = make(map[uint64]map[string]model.SalaryField)
		_, lastMonth := util.LastMonth(year, currentMonth)
		fields, err := model.GetFieldByKeyAndMonth(key, year, fromMonth, lastMonth)
		if err != nil {
			return 0.00
		}
		for _, f := range fields {
			if _, ok := fieldCache[cacheKey][f.ProfileID]; !ok {
				fieldCache[cacheKey][f.ProfileID] = make(map[string]model.SalaryField)
			}
			fieldCache[cacheKey][f.ProfileID][f.Month] = f
		}
	}

	value := 0.00
	for _, field := range fieldCache[cacheKey][profile] {
		value += field.Value
	}

	// 转换成map来方便取数

	return value
}

//Range : 区间计算器，根据设定的上限和下限，不允许指定金额超出区间
func (b *BuildinFunc) Range(salary, upper, lower float64) float64 {
	var eps = 0.000001
	if math.Max(salary, lower) == lower && math.Abs(salary-upper) > eps {
		return lower
	}

	if math.Max(salary, upper) == salary && math.Abs(salary-upper) > eps {
		return upper
	}
	return salary
}

//Range : 时间范围计算器，在设定的时间范围内可返回设定的值,比如从5-9月
func (b *BuildinFunc) TimeRange(startMonth, endMonth, value float64) float64 {
	month := int(time.Now().Month())
	if month < int(endMonth) && month > int(startMonth) {
		return value
	}

	return 0.00
}

//IF : 判断，条件成立，返回相应的数值，不成立则返回不成立的数值
func (b *BuildinFunc) IF(sentence string, trueSentence string, falseSentence string, valueMap map[string]float64) float64 {

	//先去除所有空格
	sentence = util.Strip(sentence)
	trueSentence = util.Strip(trueSentence)
	falseSentence = util.Strip(falseSentence)
	var eps = 0.000001
	result := true
	i := strings.Index(sentence, ">=")
	if i > 0 {
		params := strings.Split(sentence, ">=")
		a, b := resolveParams(params, valueMap)
		if (math.Max(a, b) == a && math.Abs(a-b) > eps) || a-b == 0 {
			result = true
		} else {
			result = false
		}
	}

	j := strings.Index(sentence, ">") //有可能错误匹配到 >=
	if j > 0 && i < 0 {
		params := strings.Split(sentence, ">")
		a, b := resolveParams(params, valueMap)
		if math.Max(a, b) == a && math.Abs(a-b) > eps && a-b != 0 {
			result = true
		} else {
			result = false
		}
	}

	i = strings.Index(sentence, "<=")
	if i > 0 {
		params := strings.Split(sentence, "<=")
		a, b := resolveParams(params, valueMap)
		if (math.Max(a, b) == b && math.Abs(a-b) > eps) || a-b == 0 {
			result = true
		} else {
			result = false
		}
	}

	j = strings.Index(sentence, "<")
	if j > 0 && i < 0 {
		params := strings.Split(sentence, "<")
		a, b := resolveParams(params, valueMap)
		if math.Max(a, b) == b && math.Abs(a-b) > eps && a-b != 0 {
			result = true
		} else {
			result = false
		}
	}

	i = strings.Index(sentence, "!=")
	if i > 0 {
		params := strings.Split(sentence, "!=")
		a, b := resolveParams(params, valueMap)
		if a != b {
			result = true
		} else {
			result = false
		}
	}

	if result {
		return util.Decimal(formula.Resolve(trueSentence+" * 1", valueMap))
	} else {
		//falseSentence 手工再 * 1 ，避免构不成表达式。
		return util.Decimal(formula.Resolve(falseSentence+" * 1", valueMap))
	}

	return 0.00
}

//IFGroupTransferMonthEqualCurrentMonth : 根据员工加入群组的时间跟当前时间对比，如果属于同一年月，则返回真值，否则返回假值
// @params profile : profile id
// @params group : group id
// @params trueSentence : 真值
// @params falseSentence: 假值
// valueMap : 依赖与值的映射 ，主要用于 trueSentence and falseSentence 的表达式应用
func (b *BuildinFunc) IFTransferMonthEqualCurrentMonth(profile uint64, group uint64, trueSentence string, falseSentence string, date string, valueMap map[string]float64) (result float64) {
	transfer, err := model.GetTransferByNewGroupAndProfile(group, profile)
	if err != nil {
		return 0.00
	}

	gDate := transfer.CreatedAt.Format("2006-01")
	currentDate := ""
	if len(date) > 0 {
		currentDate = date
	} else {
		currentDate = fmt.Sprint(time.Now().Year()) + "-" + fmt.Sprint(int(time.Now().Month()))
	}

	if gDate == currentDate {
		result, _ = strconv.ParseFloat(trueSentence, 64)
		return result
	} else {
		result, _ = strconv.ParseFloat(falseSentence, 64)
		return result
	}
	return result
}

func (b *BuildinFunc) IFGroupEqual(profile uint64, group string, trueSentence string, falseSentence string, valueMap map[string]float64) (result float64) {
	isProfileInGroup := model.IfGroupHaveProfile(group, profile)
	if isProfileInGroup {
		result = util.Decimal(formula.Resolve(trueSentence+" * 1", valueMap))
	} else {
		result = util.Decimal(formula.Resolve(falseSentence+" * 1", valueMap))
	}
	return result
}

func (b *BuildinFunc) IFTagEqual(profile uint64, tag string, trueSentence string, falseSentence string, valueMap map[string]float64) (result float64) {
	isProfileInGroup := model.IfTagHaveProfile(tag, profile)
	if isProfileInGroup {
		result = util.Decimal(formula.Resolve(trueSentence+" * 1", valueMap))
	} else {
		result = util.Decimal(formula.Resolve(falseSentence+" * 1", valueMap))
	}
	return result
}

func (b *BuildinFunc) LastMonthValue(profile uint64, currentMonth, key string) float64 {
	if currentMonth == "1" || currentMonth == "01" {
		return 0.00
	}
	currentTime := time.Now()
	year := strconv.Itoa(currentTime.Year())
	cacheKey := key + "-" + year + "-" + currentMonth

	if _, ok := lastMonthValueCache[cacheKey]; !ok {
		lastMonthValueCache[cacheKey] = make(map[uint64]model.SalaryField)
		_, lastMonth := util.LastMonth(year, currentMonth)
		fields, err := model.GetFieldLastMonthValueByKey(key, year, lastMonth)
		fmt.Println(util.PrettyJson(fields))
		if err != nil {
			return 0.00
		}
		for _, f := range fields {
			lastMonthValueCache[cacheKey][f.ProfileID] = f
		}
	}

	value := 0.00
	for _, field := range lastMonthValueCache[cacheKey] {
		fmt.Println(field.Key, field.Value)
		value = field.Value
	}

	// 转换成map来方便取数
	return value
}

// 解析参数，如果是形式于 []
func resolveParams(params []string, valueMap map[string]float64) (a, b float64) {
	if len(params) != 2 {
		return 0.00, 0.00
	}
	if _, ok := valueMap[params[0]]; ok {
		a = valueMap[params[0]]
	} else {
		a, _ = strconv.ParseFloat(params[0], 64)
	}

	if _, ok := valueMap[params[1]]; ok {
		b = valueMap[params[1]]
	} else {
		b, _ = strconv.ParseFloat(params[1], 64)
	}
	return a, b
}
