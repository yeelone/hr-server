package group

import (
	"errors"
	"fmt"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"
	"os"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

func ImportTags(c *gin.Context) {
	log.Info("Import tags into group from excel function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println("file err ", err)
		h.SendResponse(c, errno.ErrTemplateInvalid, err.Error())
		return
	}

	if !util.Exists("upload/temporary/") {
		os.MkdirAll("upload/temporary/", os.ModePerm) //创建文件
	}

	filename, subffix := util.ExtractFileName(file.Filename)
	newFilename := "upload/temporary/" + filename + "-" + time.Now().Format("20060102150405") + subffix
	if err := c.SaveUploadedFile(file, newFilename); err != nil {
		h.SendResponse(c, errno.ErrImport, err.Error())
		return
	}
	newFile := "/export/importGroupTagsRelationshipResult.xlsx"
	if errs, err := model.ImportGroupTagRelationshipFromExcel(newFilename); len(errs) > 0 {
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
		h.SendResponse(c, errno.ErrImport, CreateResponse{File: "importGroupTagsRelationshipResult.xlsx", Error: ""})
		return
	}

	rsp := CreateResponse{}
	model.CreateOperateRecord(c, fmt.Sprintf("导入机构与标签映射关系表, 文件名: %s ", newFile))
	h.SendResponse(c, nil, rsp)
}
