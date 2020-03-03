package statistics

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"
	"strconv"
)

//查询指定月份之间所有收入项
func Income(c *gin.Context) {

}

func EmployeeAnnualIncome(c *gin.Context) {
	log.Info("EmployeeAnnualIncome function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}
	var result []model.Statistics
	var err error
	if result, err = model.FetchAnnualIncome(r.Year, r.ProfileID); err != nil {
		h.SendResponse(c, errno.ErrAnnulIncome, nil)
		return
	}

	filename := writeExcel(result)

	h.SendResponse(c, nil, CreateResponse{
		File: filename,
	})
}

func DepartmentAnnualIncome(c *gin.Context) {
	log.Info("DepartmentAnnualIncome function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		fmt.Println("annual error", err)
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}
	h.SendResponse(c, nil, nil)
}
func writeExcel(sf []model.Statistics) string {
	xlsx := excelize.NewFile()

	//cols := map[string]int{
	//	"profile": 1 ,
	//	"id_card": 2 ,
	//	"department": 3 ,
	//	"post": 4 ,
	//	"year": 5 ,
	//	"month": 6 ,
	//	"field": 7 ,
	//	"value": 8,
	//}

	xlsx.SetCellValue("Sheet1", "A1", "姓名")
	xlsx.SetCellValue("Sheet1", "B1", "身份证")
	xlsx.SetCellValue("Sheet1", "C1", "部门")
	xlsx.SetCellValue("Sheet1", "D1", "时间")
	xlsx.SetCellValue("Sheet1", "E1", "收入")

	row := 2

	for _, s := range sf {
		xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(row), s.Profile)
		xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(row), s.IDCard)
		xlsx.SetCellValue("Sheet1", "C"+strconv.Itoa(row), s.Department)
		xlsx.SetCellValue("Sheet1", "D"+strconv.Itoa(row), s.Year)
		xlsx.SetCellValue("Sheet1", "E"+strconv.Itoa(row), s.Total)
		row++
	}

	// Save xlsx file by the given path.
	filename := "annual_income.xlsx"
	err := xlsx.SaveAs("./export/" + filename)
	if err != nil {
		fmt.Println(err)
	}
	return filename
}
