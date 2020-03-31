package buildinfunc

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"hr-server/util"
	"io/ioutil"
	"reflect"
	"strings"
)

type Data struct {
	Entry []Entry `yaml:"entries"`
}

type Entry map[string]string

type Buildin interface {
	CalcuWorkingAge(map[string]interface{}) int
}

type BuildinFunc struct{}

//将调用同一个函数相同参数的进行缓存
var funcCache = make(map[string][]interface{})

//https://play.golang.org/p/BQpOfDWBzqn
func (b *BuildinFunc) Call(name string, args []interface{}) (values []interface{}, err error) {

	defer func() {
		if r := recover(); r != nil {
			// 这里可以对异常进行一些处理和捕获
			err = errors.New("无法调用函数，请仔细检查函数名称是否正确,函数名为" + name)
		}
	}()

	//先从cache中取值
	hashName := hashFunc(name, args)
	if cacheValues, err := getCache(hashName); err == nil {
		return cacheValues, nil
	}

	s := &BuildinFunc{}
	v := reflect.ValueOf(s)

	params := []reflect.Value{}

	for _, arg := range args {
		params = append(params, reflect.ValueOf(arg))
	}
	for _, v := range v.MethodByName(name).Call(params) {
		values = append(values, v.Interface())

		if name == "Addup" || name == "LastMonthValue" { // Addup 没办法这样简单的进行缓存
			continue
		}
		setCache(hashName, values)
	}

	return values, err
}

func setCache(name string, values []interface{}) {
	funcCache[name] = values
}

func getCache(name string) (values []interface{}, err error) {
	if val, ok := funcCache[name]; ok {
		return val, nil
	}
	return values, errors.New("cannot find any cache")
}

func hashFunc(name string, params []interface{}) (hashStr string) {
	var paramSlice []string
	for _, param := range params {
		paramSlice = append(paramSlice, fmt.Sprintf("%v", param))
	}
	//has := md5.Sum([]byte(strings.Join(paramSlice, "_")))
	has := util.Strip(strings.Join(paramSlice, "_"))
	return fmt.Sprintf("%x", has) //将[]byte转成16进制
}

func (b *BuildinFunc) ListFunc() (data map[string]map[string]interface{}, err error) {
	out := make(map[string]map[string]interface{})
	yamlFile, err := ioutil.ReadFile("conf/func_list.yaml")
	if err != nil {
		return data, err
	}

	err = yaml.Unmarshal(yamlFile, out)
	if err != nil {
		return data, err
	}
	return out, nil

}
