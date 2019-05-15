package template

import (
	"fmt"
)

func GetFields(name string) (fields []Field) {
	tp, err := ResolveTemplate(name)
	if err != nil {
		fmt.Println("模板解析失败，请检查模板格式是否正确", err.Error())
		return nil
	}

	fields = make([]Field, 0)
	fields = append(fields, tp.Base...)
	fields = append(fields, tp.Buildin...)
	fields = append(fields, tp.Coefficient...)
	fields = append(fields, tp.Formula...)
	fields = append(fields, tp.Input...)
	fields = append(fields, tp.Related...)
	fields = append(fields, tp.Upload...)
	return fields
}
