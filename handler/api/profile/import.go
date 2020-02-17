package profile

import (
	"fmt"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"
	"os"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

func Import(c *gin.Context) {
	log.Info("Import Profiles from excel to create function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	file, err := c.FormFile("file")
	if err != nil {
		h.SendResponse(c, errno.ErrTemplateInvalid, nil)
		return
	}

	userid, ok := c.Get("userid")
	if !ok {
		h.SendResponse(c, errno.StatusUnauthorized, nil)
		return
	}


	filename, subffix := util.ExtractFileName(file.Filename)

	// if subffix != ".csv" {
	// 	// 只支持csv格式
	// 	h.SendResponse(c, errno.ErrUploadFormatInvalid, nil)
	// 	return
	// }

	if !util.Exists("upload/import/") {
		os.MkdirAll("upload/import/",os.ModePerm) //创建文件
	}

	newFilename := "upload/import/" + filename + "-" + time.Now().Format("20060102150405") + subffix
	if err := c.SaveUploadedFile(file, newFilename); err != nil {
		log.Error("上传文件出现错误:", err)
		h.SendResponse(c, errno.ErrTemplateInvalid, nil)
		return
	}

	if f, err := model.ImportProfileFromExcel(newFilename,userid.(uint64)); err != nil {
		log.Error("导入失败:", err)
		h.SendResponse(c, errno.ErrImport, CreateResponse{File: f, Error: err.Error()})
		return
	}

	rsp := CreateResponse{}
	model.CreateOperateRecord(c, fmt.Sprintf("导入员工, 文件：[ %s ]", newFilename))
	// Show the user information.
	h.SendResponse(c, nil, rsp)
}

//handleCsvFile:
//在这一步要把csv第一行取出来，根据ProfileI18nMap来构造查询字符串，供postgresql /copy 命令入参。
//然后把第一行删除掉生成一个新的文件。
// func handleCsvFile(filename string) (fields string, newfilename string, err error) {
// 	orgFilename := "upload/import/" + filename
// 	newCsvfileName := "upload/import/csv/" + filename

// 	file, err := os.Open(orgFilename)
// 	if err != nil {
// 		return "", "", nil
// 	}
// 	defer file.Close()

// 	decoder := mahonia.NewDecoder("gbk")             // 把原来ANSI格式的文本文件里的字符，用gbk进行解码。
// 	reader := csv.NewReader(decoder.NewReader(file)) // 这样，最终返回的字符串就是utf-8了。（go只认utf8）

// 	contents, err := reader.ReadAll()
// 	sz := len(contents)
// 	if err != nil {
// 		return "", "", err
// 	}

// 	cols := make([]string, len(contents[0]))
// 	for index, field := range contents[0] {
// 		cols[index] = model.ProfileI18nMap[field]
// 	}

// 	buf := new(bytes.Buffer)
// 	writer := csv.NewWriter(buf)
// 	for i := 1; i < sz; i++ {
// 		writer.Write(contents[i])
// 		writer.Flush()
// 	}

// 	fout, err := os.Create(newCsvfileName)
// 	defer fout.Close()
// 	if err != nil {
// 		fmt.Println(newCsvfileName, err)
// 		return
// 	}

// 	_, err = fout.WriteString(buf.String())
// 	if err != nil {
// 		return "", "", err
// 	}

// 	return strings.Join(cols, ","), newCsvfileName, nil
// }

func handleUploadedExcel(filename string) (fields string, err error) {
	//分析第一行
	xlsx, err := excelize.OpenFile(filename)
	if err != nil {
		fmt.Println("OpenFile", err)
		return "", err
	}

	rows,_ := xlsx.GetRows("Sheet1")
	cols := make([]string, len(rows[0]))
	for index, colCell := range rows[0] {
		cols[index] = model.ProfileI18nMap[colCell]
	}

	return strings.Join(cols, ","), err
}

// func handleUploadedExcel(filename, subffix string) (fields string, newfilename string, err error) {
// 	filepath := "upload/import/" + filename + subffix
// 	//分析第一行
// 	xlsx, err := excelize.OpenFile(filepath)
// 	if err != nil {
// 		fmt.Println("OpenFile", err)
// 		return "", "", err
// 	}
// 	fmt.Println("read excel")

// 	rows := xlsx.GetRows("Sheet1")
// 	fmt.Println(len(rows))
// 	cols := make([]string, len(rows[0]))
// 	for index, colCell := range rows[0] {
// 		cols[index] = model.ProfileI18nMap[colCell]
// 		fmt.Println("colCell", colCell)
// 	}

// 	fmt.Print(strings.Join(cols, ","))

// 	// //将excel转为csv
// 	// csvfileName := "upload/import/csv/" + filename + ".csv"
// 	// buf := new(bytes.Buffer)
// 	// r2 := csv.NewWriter(buf)

// 	// for _, row := range rows[1:] {
// 	// 	fmt.Println(row)
// 	// 	r2.Write(row)
// 	// 	r2.Flush()
// 	// }
// 	// fout, err := os.Create(csvfileName)
// 	// defer fout.Close()
// 	// if err != nil {
// 	// 	fmt.Println(csvfileName, err)
// 	// 	return
// 	// }

// 	// fout.WriteString(buf.String())

// 	return strings.Join(cols, ","), csvfileName, err
// }
