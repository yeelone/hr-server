package salary

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/pkg/template"
	"hr-server/util"
	"io/ioutil"
	"strconv"
)

func SalaryProfileConfig(c *gin.Context) {
	log.Info("SalaryProfileConfig Create function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateProfileConfigCreateRequest
	if err := c.Bind(&r); err != nil {
		fmt.Println("err", err)
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}

	m := model.SalaryProfileConfig{
		ProfileId:       r.ProfileID,
		TemplateFieldId: r.TemplateFieldID,
		Operate:         r.Operate,
		Value:           r.Value,
		Description:     r.Description,
	}

	if err := m.Create(); err != nil {
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}
	h.SendResponse(c, nil, nil)
}

func GetSalaryProfileConfigList(c *gin.Context) {
	log.Info("SalaryProfileConfig Create function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})

	list, err := model.GetSalaryProfileConfig()
	if err != nil {
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}

	//
	files, _ := ioutil.ReadDir("conf/templates")
	fields := make(map[string]template.Field)
	for _, f := range files {
		name, _ := util.ExtractFileName(f.Name())
		t, _ := template.ResolveTemplate(name)

		for _, f := range t.All {
			fields[f.ID] = f
		}
	}

	for i, item := range list {
		list[i].TemplateField = fields[item.TemplateFieldId]
	}
	h.SendResponse(c, nil, SalaryProfileConfigResponse{
		ConfigList: list,
	})
}

func DeleteSalaryProfileConfig(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	m := model.SalaryProfileConfig{}
	m.ID = uint64(id)

	if err := m.Delete(); err != nil {
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}
	h.SendResponse(c, nil, nil)

}
