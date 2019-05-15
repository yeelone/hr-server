package grouptransfer

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	h "hrgdrc/handler"
	"hrgdrc/model"
	"hrgdrc/pkg/errno"
	"strconv"
	"time"
)

func List(c *gin.Context) {
	var r ListRequest
	if err := c.ShouldBind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}
	transfer, err := model.GetAllProfileTransfer()
	if err != nil {
		h.SendResponse(c, err, nil)
		return
	}

	xlsx := excelize.NewFile()
	xlsx.SetCellValue("Sheet1", "A1", "姓名")
	xlsx.SetCellValue("Sheet1", "B1", "从")
	xlsx.SetCellValue("Sheet1", "C1", "到")
	xlsx.SetCellValue("Sheet1", "D1", "日期")
	for i, t := range transfer {
		xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), t.ProfileName)
		xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), t.OldGroupName)
		xlsx.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), t.NewGroupName)
		xlsx.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), t.CreatedAt.Format("20060102"))
	}
	// Save xlsx file by the given path.
	filename := "调动记录表 " + time.Now().Format("20060102150405") + ".xlsx"
	err = xlsx.SaveAs("./export/" + filename)
	if err != nil {
		fmt.Println(err)
		h.SendResponse(c, errno.ErrCreateFile, err.Error())
		return
	}

	h.SendResponse(c, nil, ListResponse{
		File: filename,
	})
}
