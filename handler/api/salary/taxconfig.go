package salary

import (
	"fmt"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"
	"io/ioutil"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
	yaml "gopkg.in/yaml.v2"
)

func TaxSetting(c *gin.Context) {
	log.Info("TemplateConfig Create function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r TaxRequest
	if err := c.Bind(&r); err != nil {
		fmt.Println(err)
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}
	filename := "conf/salary.yaml"
	if util.Exists(filename) {
		err := util.MoveFile(filename, "conf/old/salary-"+time.Now().Format("20060102-150405")+".yaml")
		if err != nil {
			fmt.Println("cannot move file to new directory" + err.Error())
		}
	}

	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		fmt.Println("create file error:", err.Error())
	}

	d, err := yaml.Marshal(&r.TaxConf)
	if err != nil {
		fmt.Printf("error: %v", err)
	} else {
		_, err := f.Write(d)
		if err != nil {
			fmt.Println("cannot create the file ")
		}
	}
	record := model.Record{}
	record.Object = "salary"
	record.Body = "个税起征点有发生变更,请仔细检查！"
	if err := record.Create(); err != nil {

	}
	h.SendResponse(c, nil, nil)
}

func GetTaxSetting(c *gin.Context) {
	yamlFile, err := ioutil.ReadFile("conf/salary.yaml")

	m := model.TaxConf{}
	err = yaml.Unmarshal(yamlFile, &m)
	if err != nil {
		fmt.Println(err)
	}

	resp := &TaxResponse{}
	resp.Conf = m

	h.SendResponse(c, nil, resp)
}
