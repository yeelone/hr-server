package template

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"github.com/xcltapestry/xclpkg/algorithm"
)

var FitIntoMonth = map[string]string{
	"LASTMONTH": "上个月",
	"LASTYEAR":  "上一年份",
}

type Field struct {
	ID                    string      `json:"id" yaml:"id"` //生成一个唯一的ID，主要用于在以后的新旧模板对比中判断是否属于同一个Field
	Key                   string      `json:"key" yaml:"key"`
	Invalid               bool        `json:"invalid" yaml:"invalid"`
	Type                  string      `json:"type" yaml:"type"`
	Name                  string      `json:"name" yaml:"name"`
	Alias                 string      `json:"alias" yaml:"alias"`
	Require               []string    `json:"require" yaml:"require"`
	Call                  string      `json:"call" yaml:"call"`
	Formula               string      `json:"formula" yaml:"formula"`
	From                  string      `json:"from" yaml:"from"`
	Params                []string    `json:"params" yaml:"params"`
	Value                 interface{} `json:"value" yaml:"value"`
	ShouldTax             bool        `json:"should_tax" yaml:"should_tax"`
	Order                 int         `json:"order" yaml:"order"`
	IsIncome              bool        `json:"is_income" yaml:"is_income"`
	IsDeduct              bool        `json:"is_deduct" yaml:"is_deduct"`
	FitIntoMonth          string      `json:"fit_into_month" yaml:"fit_into_month"`
	RelateTemplateAccount string      `json:"related_templateaccount" yaml:"related_templateaccount"`
	RelateTemplate        string      `json:"related_template" yaml:"related_template"`
	RelateKey             string      `json:"related_key" yaml:"related_key"`
	RelateYear            string      `json:"related_year" yaml:"related_year"`
	RelateMonth           string      `json:"related_month" yaml:"related_month"`
	Description           string      `json:"description" yaml:"description"`
	FixedData             bool        `json:"fixed_data" yaml:"fixed_data"`       //判断是否固定上传的数据 ，如果是，则需要从初始化数据中取得数据并进行填充
	MustRounding          bool        `json:"must_rounding" yaml:"must_rounding"` //必须四舍五入，例如住房公积金
	Visible               bool        `json:"visible" yaml:"visible"`             //在excel中是否可见，默认为可见
}

type Template struct {
	Name          string
	File          string
	All           []Field
	Base          []Field
	Input         []Field
	Coefficient   []Field
	NeedCalculate []Field
	Buildin       []Field
	Formula       []Field
	Upload        []Field
	Related       []Field
}

var runtimeViper = viper.New()

var keys map[string]int   //为了去重
var all []Field           //基础要素，只从用户表里取相关数据
var base []Field          //基础要素，只从用户表里取相关数据
var input []Field         //值来自于用户设定要素
var coefficient []Field   //系数，从群组里取得用户相关群组，并获取该群组的系数
var wantCalculate []Field //系数，从群组里取得用户相关群组，并获取该群组的系数
var otherTemplate []Field //数据从其它模板中取
var buildin []Field       //内建函数，调用系统所提供的内建函数，会利用反射来实现
var calculate []Field     //计算要素，会依赖于其它要素
var fileImport []Field    //计算要素，会依赖于其它要素

func ResolveTemplate(name string) (t Template, err error) {
	runtimeViper = viper.New()
	runtimeViper.AddConfigPath("conf/templates") // 如果没有指定配置文件，则解析默认的配置文件
	runtimeViper.SetConfigName(name)

	runtimeViper.SetConfigType("yaml")                  // 设置配置文件格式为YAML
	if err := runtimeViper.ReadInConfig(); err != nil { // viper解析配置文件
		return t, err
	}

	t = readTemplate()
	resolveDependency(&t)
	return t, nil
}

//读取模板文件
func readTemplate() Template {
	keys = make(map[string]int)      //为了去重
	all = make([]Field, 0)           //基础要素，只从用户表里取相关数据
	base = make([]Field, 0)          //基础要素，只从用户表里取相关数据
	input = make([]Field, 0)         //值来自于用户设定要素
	coefficient = make([]Field, 0)   //系数，从群组里取得用户相关群组，并获取该群组的系数
	wantCalculate = make([]Field, 0) //系数，从群组里取得用户相关群组，并获取该群组的系数
	otherTemplate = make([]Field, 0) //数据从其它模板中取
	buildin = make([]Field, 0)       //内建函数，调用系统所提供的内建函数，会利用反射来实现
	calculate = make([]Field, 0)     //计算要素，会依赖于其它要素
	fileImport = make([]Field, 0)    //计算要素，会依赖于其它要素

	for _, key := range runtimeViper.AllKeys() {
		s := strings.Split(key, ".")
		keys[s[0]] = 0
	}
	t := Template{}
	field := Field{}
	for key := range keys {
		field.Key = key
		field.ID = runtimeViper.GetString(key + ".id")
		field.Invalid = runtimeViper.GetBool(key + ".invalid")
		field.Type = runtimeViper.GetString(key + ".type")
		field.Name = runtimeViper.GetString(key + ".name")
		field.Value = runtimeViper.GetFloat64(key + ".value")
		field.Alias = runtimeViper.GetString(key + ".alias")
		field.Require = runtimeViper.GetStringSlice(key + ".require")
		field.Call = runtimeViper.GetString(key + ".call")
		field.Formula = runtimeViper.GetString(key + ".formula")
		field.From = runtimeViper.GetString(key + ".from")
		field.Order = runtimeViper.GetInt(key + ".order")
		field.RelateTemplateAccount = runtimeViper.GetString(key + ".related_templateaccount")
		field.RelateTemplate = runtimeViper.GetString(key + ".related_template")
		field.RelateKey = runtimeViper.GetString(key + ".related_key")
		field.RelateYear = runtimeViper.GetString(key + ".related_year")
		field.RelateMonth = runtimeViper.GetString(key + ".related_month")
		field.ShouldTax = runtimeViper.GetBool(key + ".should_tax")
		field.Params = runtimeViper.GetStringSlice(key + ".params")
		field.IsIncome = runtimeViper.GetBool(key + ".is_income")
		field.IsDeduct = runtimeViper.GetBool(key + ".is_deduct")
		field.FitIntoMonth = runtimeViper.GetString(key + ".fit_into_month")
		field.FixedData = runtimeViper.GetBool(key + ".fixed_data")
		field.MustRounding = runtimeViper.GetBool(key + ".must_rounding")
		field.Visible = runtimeViper.GetBool(key + ".visible")

		switch t := runtimeViper.GetString(key + ".type"); t {
		case "Base":
			base = append(base, field)
			all = append(all, field)
		case "Coefficient":
			coefficient = append(coefficient, field)
			all = append(all, field)
		case "Input":
			input = append(input, field)
			all = append(all, field)
		case "Buildin":
			buildin = append(buildin, field)
			wantCalculate = append(wantCalculate, field)
			all = append(all, field)
		case "Calculate":
			calculate = append(calculate, field)
			wantCalculate = append(wantCalculate, field)
			all = append(all, field)
		case "Upload":
			fileImport = append(fileImport, field)
			all = append(all, field)
		case "Related":
			otherTemplate = append(otherTemplate, field)
			all = append(all, field)
		}
	}
	t.Base = base
	t.Coefficient = coefficient
	t.Upload = fileImport
	t.Buildin = buildin
	t.Formula = calculate
	t.Related = otherTemplate
	t.Input = input
	t.All = all
	return t
}

// resolveDependency :
// 根据buildin 和calculate类型的元素所拥有的依赖关系 ，利用拓扑排序算法进行依赖排序，对这些元素重新排列，之后才可以按序进行计算
func resolveDependency(t *Template) {
	//去除 Base Coefficient Related Upload,
	var gMap = make(map[string][]string)
	for _, field := range wantCalculate {
		required := runtimeViper.GetStringSlice(field.Key + ".require")
		temp := make([]string, 0)
		for _, k := range required {
			if runtimeViper.GetString(k+".type") == "Base" || runtimeViper.GetString(k+".type") == "Coefficient" ||
				runtimeViper.GetString(k+".type") == "Related" || runtimeViper.GetString(k+".type") == "Upload" {
				continue
			}

			temp = append(temp, k)
		}
		gMap[field.Key] = temp

	}
	sortList := dependSort(gMap)

	newList := make([]Field, 0)
	for _, key := range sortList {
		for _, f := range wantCalculate {
			if f.Key == key {
				newList = append(newList, f)
			}
		}
	}
	t.NeedCalculate = newList
}

// https://play.golang.org/p/tVmbJ3akrhB
func dependSort(G map[string][]string) []string {
	var inMap = make(map[string]int, 0)

	//第一步，先把所有的节点入度初始为0
	for key, _ := range G {
		inMap[key] = 0
	}
	//第二步，遍历依赖节点，为节点增加入度
	for _, g := range G {
		for _, key := range g {
			inMap[key] += 1
		}
	}

	//第三步，将入度为0的提取出来
	var qList = make([]string, 0)
	for key, _ := range G {
		if inMap[key] == 0 {
			qList = append(qList, key)
		}
	}

	//第四步，
	var sList = make([]string, 0)
	for len(qList) > 0 {
		//pop 操作
		key := qList[len(qList)-1]
		qList = qList[:len(qList)-1]

		sList = append(sList, key)
		for _, v := range G[key] {
			inMap[v] -= 1

			if inMap[v] == 0 {
				qList = append(qList, v)
			}
		}
	}
	return sList
}

func ResolveFormula(str string) {
	exp, err := ExpConvert(covertFormulaToArray(str))
	if err != nil {
		fmt.Println("中序表达式转后序表达式失败! ", err)
	} else {
		Exp(exp)
	}
}

// covertFormulaToArray : 将字符串
func covertFormulaToArray(str string) (sarr []string) {
	//先去除所有的空格
	str = strings.Replace(str, " ", "", -1)
	//同时去除换行符
	str = strings.Replace(str, "\n", "", -1)
	word := ""
	for i, s := range str {
		ch := string(s)
		if IsOperator(ch) {
			if word != "" {
				sarr = append(sarr, word)
			}
			sarr = append(sarr, ch)
			word = ""
			continue
		}

		word += ch

		if i == len(str)-1 {
			sarr = append(sarr, word)
		}

	}

	return sarr
}

func ExpConvert(str []string) ([]string, error) {

	var result []string
	stack := algorithm.NewStack()
	for _, s := range str {
		ch := s
		if IsOperator(ch) { //是运算符
			if stack.Empty() || ch == "(" {
				stack.Push(ch)
			} else {
				if ch == ")" { //处理括号
					for {
						if stack.Empty() {
							return nil, errors.New("表达式有问题! 没有找到对应的\"(\"号")
						}
						if stack.Top().(string) == "(" {
							break
						}
						result = append(result, stack.Top().(string))
						stack.Pop()
					}

					//弹出"("
					stack.Top()
					stack.Pop()
				} else { //非括号
					//比较优先级
					for {
						if stack.Empty() {
							break
						}
						m := stack.Top().(string)
						if Priority(ch) > Priority(m) {
							break
						}
						result = append(result, m)
						stack.Pop()
					}
					stack.Push(ch)
				}
			}

		} else { //非运算符
			result = append(result, ch)
		} //end IsOperator()

	} //end for range str

	for {
		if stack.Empty() {
			break
		}
		result = append(result, stack.Top().(string))
		stack.Pop()
	}
	return result, nil
}

func Exp(str []string) {
	stack := algorithm.NewStack()
	for _, s := range str {
		ch := s
		if IsOperator(ch) { //是运算符
			if stack.Empty() {
				break
			}
			stack.Print()
			b := stack.Top().(string)
			stack.Pop()

			a := stack.Top().(string)
			stack.Pop()

			sv := fmt.Sprintf("%f", Calc(ch, a, b))

			stack.Push(sv)
		} else {
			stack.Push(ch)
		} //end IsOperator
	}
	if !stack.Empty() {
		stack.Pop()
	}
}

func IsOperator(op string) bool {
	switch op {
	case "(", ")", "+", "-", "*", "/":
		return true
	default:
		return false
	}
}

func Priority(op string) int {
	switch op {
	case "*", "/":
		return 3
	case "+", "-":
		return 2
	case "(":
		return 1
	default:
		return 0
	}
}

func covert2Int(a string) int {
	i, _ := strconv.Atoi(a)
	return i
}

func covert2Float(a string) float64 {
	i, _ := strconv.ParseFloat(a, 64)
	return i
}

func findKey(n string) (float64, bool) {
	return 0.00, false
}

func Calc(op string, a, b string) float64 {
	var ia float64 = -1.00
	var ib float64 = -1.00

	if val, ok := findKey(a); ok {
		ia = val
	}

	if val, ok := findKey(b); ok {
		ib = val
	}
	if ia == -1 {
		ia = covert2Float(a)
	}

	if ib == -1 {
		ib = covert2Float(b)
	}
	switch op {
	case "*":
		return (ia * ib)
	case "/":
		return (ia / ib)
	case "+":
		return (ia + ib)
	case "-":
		return (ia - ib)
	default:
		return 0
	}
}
