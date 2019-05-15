package template

import (
	"fmt"
	"hrgdrc/util"
	"strings"
)

// ComparedTemplate : 双比两个模板配置文件有哪些异同
// @params ： oldFile , newFile 指定要对比的两个文件
// @return ： 返回发现变更的字段以及变更的内容
// 主要对比内容包括： 有无新增项目，项目类型是否有变化， 公式变化，公式依赖参数变化 ，函数调用变化 ，函数参数变化 ，固定值
func ComparedTemplate(oldFile, newFile string) (result string, err error) {
	oldTemplate, err := ResolveTemplate(oldFile)
	if err != nil {
		return "", err
	}
	newTemplate, err := ResolveTemplate(newFile)
	if err != nil {
		return "", err
	}

	oldFields := []Field{}
	oldFields = append(oldFields, oldTemplate.Base...)
	oldFields = append(oldFields, oldTemplate.Input...)
	oldFields = append(oldFields, oldTemplate.Coefficient...)
	oldFields = append(oldFields, oldTemplate.Upload...)
	oldFields = append(oldFields, oldTemplate.Buildin...)
	oldFields = append(oldFields, oldTemplate.Formula...)
	oldFields = append(oldFields, oldTemplate.Related...)
	newFields := []Field{}
	newFields = append(newFields, newTemplate.Base...)
	newFields = append(newFields, newTemplate.Input...)
	newFields = append(newFields, newTemplate.Coefficient...)
	newFields = append(newFields, newTemplate.Upload...)
	newFields = append(newFields, newTemplate.Buildin...)
	newFields = append(newFields, newTemplate.Formula...)
	newFields = append(newFields, newTemplate.Related...)

	changeContent := comparedField(oldFields, newFields)
	if len(changeContent) > 0 {
		changeContent = "模板发生了变更:" + oldFile + ";" + changeContent
	}
	return changeContent, err
}

func comparedField(oldFields, newFields []Field) (change string) {
	idList := make(map[string]struct{})
	oldFieldMap := make(map[string]Field)
	newFieldMap := make(map[string]Field)
	for _, field := range oldFields {
		oldFieldMap[field.ID] = field
		idList[field.ID] = struct{}{}
	}

	for _, field := range newFields {
		newFieldMap[field.ID] = field
		idList[field.ID] = struct{}{}
	}

	for id := range idList {
		if _, ok := oldFieldMap[id]; !ok { //如果在旧的模板不存在，则表示是新增的
			newField := newFieldMap[id]
			change += "新增项目:key=" + newField.Key + ", 名=" + newField.Name + ", 别名=" + newField.Alias + ";"
			continue
		}

		if _, ok := newFieldMap[id]; !ok {
			field := oldFieldMap[id]
			change += "删除项目:key=" + field.Key + ", 名=" + field.Name + ", 别名=" + field.Alias + ";"
			continue
		}
		oldField := oldFieldMap[id]
		newField := newFieldMap[id]
		if oldField.Key != newField.Key {
			change += "Key发生变更: 原key=" + oldField.Key + ", 新key=" + newField.Key + ",当Key发生变更时，请仔细检查所有关联到此Key的计算项目;"
		}
		if oldField.Name != newField.Name {
			change += "字段名字发生变更: 原名=" + oldField.Name + ", 新名=" + newField.Name + ";"
		}
		if oldField.Invalid != newField.Invalid {
			change += "字段冻结状态变更: 原状态=" + fmt.Sprintf("%s", oldField.Invalid) + ", 新状态=" + fmt.Sprintf("%s", newField.Invalid) + ";"
		}
		if oldField.Alias != newField.Alias {
			change += "字段别名发生变更: 原名=" + oldField.Alias + ", 新名=" + newField.Alias + ";"
		}
		if oldField.Type != newField.Type {
			change += "字段类型发生变更: 原类型=" + oldField.Type + ", 新类型=" + newField.Type + ";"
		}
		if !util.StringSliceEqualBCE(oldField.Require, newField.Require) {
			change += "依赖数据发生变更: 原依赖的数据有=" + strings.Join(oldField.Require, ",") + ", 现依赖的数据有=" + strings.Join(newField.Require, ",") + ";"
		}
		if oldField.Call != newField.Call {
			change += "函数调用发生变更: 原调用函数=" + oldField.Call + ", 现调用函数=" + newField.Call + ";"
		}
		if util.Strip(oldField.Formula) != util.Strip(newField.Formula) {
			change += "公式调用发生变更: 原公式=" + util.Strip(oldField.Formula) + ", 现公式=" + util.Strip(newField.Formula) + ";"
		}
		if !util.StringSliceEqualBCE(oldField.Params, newField.Params) {
			change += "公式或者函数参数发生变更: 原参数=" + strings.Join(oldField.Params, ",") + ", 现参数=" + strings.Join(newField.Params, ",") + ";"
		}
		if oldField.Value != newField.Value {
			change += "固定输入值变更: 原值=" + fmt.Sprint(oldField.Value) + ", 新值=" + fmt.Sprint(newField.Value) + ";"
		}
		if oldField.ShouldTax != newField.ShouldTax {
			change += "关于是否属于扣税项: 原值=" + fmt.Sprint(oldField.ShouldTax) + ", 新值=" + fmt.Sprint(newField.ShouldTax) + ";"
		}
		if oldField.IsIncome != newField.IsIncome {
			change += "关于是否属于收入项: 原值=" + fmt.Sprint(oldField.IsIncome) + ", 新值=" + fmt.Sprint(newField.IsIncome) + ";"
		}
		if oldField.IsDeduct != newField.IsDeduct {
			change += "关于是否属于扣减项: 原值=" + fmt.Sprint(oldField.ShouldTax) + ", 新值=" + fmt.Sprint(newField.ShouldTax) + ";"
		}
		if oldField.FitIntoMonth != newField.FitIntoMonth {
			change += "工资项目纳入其它月份: 原值=" + FitIntoMonth[oldField.FitIntoMonth] + ", 新值=" + FitIntoMonth[newField.FitIntoMonth] + ";"
		}
		if oldField.RelateTemplate != newField.RelateTemplate {
			change += "工资计算项目的值来源于其它计算模板发生变更 : 原模板=" + oldField.RelateTemplate + ", 新模板=" + newField.RelateTemplate + ";"
		}
		if oldField.RelateKey != newField.RelateKey {
			change += "工资计算项目的值来源于其它计算模板的指定项目发生变更 : 原关联项目=" + oldField.RelateKey + ", 新关联项目=" + newField.RelateKey + ";"
		}
		if oldField.FixedData != newField.FixedData {
			change += "该项数据是否来自于设定模板时所设置的固定上传文件:原=" + fmt.Sprint(oldField.FixedData) + ", 新关联项目=" + fmt.Sprint(newField.FixedData) + ";"
		}
		if oldField.MustRounding != newField.MustRounding {
			change += "该项是否进行四舍五入保持整数:原关联项目=" + fmt.Sprint(oldField.MustRounding) + ", 新关联项目=" + fmt.Sprint(newField.MustRounding) + ";"
		}
	}
	return change
}
