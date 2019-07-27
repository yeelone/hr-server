package model_test

import (
	"fmt"
	"hr-server/model"
	"hr-server/util"
	"testing"
)

func TestGetFieldByKeyAndMonth(t *testing.T) {
	result, err := model.GetFieldByKeyAndMonth("医药补贴", "2018", "01", "12")
	fmt.Println(util.PrettyJson(result))
	if err != nil {
		t.Error("GetFieldByKeyAndMonth failed!")
	}

	t.Log("GetFieldByKeyAndMonth test pass")
}
