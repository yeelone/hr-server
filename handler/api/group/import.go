package group

import (
	"errors"
	"fmt"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

func Import(c *gin.Context) {
	log.Info("Import Group from excel to create function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println("file err ", err)
		h.SendResponse(c, errno.ErrTemplateInvalid, err.Error())
		return
	}

	filename, subffix := util.ExtractFileName(file.Filename)

	newFilename := "upload/import/" + filename + "-" + time.Now().Format("20060102150405") + subffix
	if err := c.SaveUploadedFile(file, newFilename); err != nil {
		h.SendResponse(c, errno.ErrImport, err.Error())
		return
	}
	newFile := "/export/importGroupResult.xlsx"
	if errs, err := model.ImportGroupFromExcel(newFilename); len(errs) > 0 {
		if err != nil {
			fmt.Println("OpenFile", err)
			h.SendResponse(c, errors.New("导入数据库之后发现错误，请下载错误文件"), CreateResponse{File: "", Error: "无法将错误信息写入文件"})
			return
		}
		xlsx := excelize.NewFile()
		for i, err := range errs {
			xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(i+1), err)
		}
		err = xlsx.SaveAs("." + newFile)
		if err != nil {
			h.SendResponse(c, errno.ErrImport, CreateResponse{File: "", Error: err.Error()})
			return
		}
		h.SendResponse(c, errno.ErrImport, CreateResponse{File: "importGroupResult.xlsx", Error: ""})
		return
	}

	rsp := CreateResponse{}
	model.CreateOperateRecord(c, fmt.Sprintf("批量导入机构信息, 文件名: %s ", newFile))
	h.SendResponse(c, nil, rsp)
}
