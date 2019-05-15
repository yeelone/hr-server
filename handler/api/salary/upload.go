package salary

import (
	h "hrgdrc/handler"
	"hrgdrc/pkg/errno"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	// Source
	file, err := c.FormFile("file")
	if err != nil {
		h.SendResponse(c, errno.ErrTemplateInvalid, nil)
		return
	}
	filename, subffix := renameFile(file.Filename)
	if subffix != ".xlsx" {
		h.SendResponse(c, errno.ErrUploadFileTypeInvalid, nil)
		return
	}

	filepath := "upload/" + filename
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		h.SendResponse(c, errno.ErrTemplateInvalid, nil)
		return
	}

	h.SendResponse(c, nil, UploadResponse{UploadFile: filepath})
}

func renameFile(filename string) (name, subffix string) {
	t := strconv.FormatInt(time.Now().UnixNano(), 10)

	var filenameWithSuffix string
	filenameWithSuffix = path.Base(filename)

	var fileSuffix string
	fileSuffix = path.Ext(filenameWithSuffix) //获取文件后缀

	var filenameOnly string
	filenameOnly = strings.TrimSuffix(filenameWithSuffix, fileSuffix) //获取文件名

	return filenameOnly + t + fileSuffix, fileSuffix

}
