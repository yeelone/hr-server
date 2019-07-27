package salary

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/pkg/template"
	"hr-server/util"
	"io/ioutil"
	"os"
	"strings"
)

// TemplateConfig : 配置模板
func TemplateConfig(c *gin.Context) {
	log.Info("TemplateConfig Create function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		fmt.Println("bind error ", err)
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}

	m := model.Template{
		Name:       r.Name,
		Type:       r.Type,
		Months:     r.Months,
		Startup:    false,
		InitData:   r.InitData,
		UserID:     1,
		AuditState: model.AuditStateWaiting,
	}

	////要判断key是唯一的
	// remark: 之前考虑到key是全系统唯一，因为key要从配置文件中去解析，后来想了想，改成在查询明细时从数据库里解析key ，这样就不会存在key冲突的问题，一个key是跟相应有日期绑定在一起的
	//allKeys := findKeys(r.Name)
	//errMsg := []string{}
	//for _, f := range r.Body {
	//	for _, key := range allKeys {
	//		if f.Key == key {
	//			errMsg = append(errMsg, f.Key+"已存在于其它模板中，请确定key是全系统唯一的标识.")
	//		}
	//	}
	//}
	//
	//if len(errMsg) > 0 {
	//	h.SendResponse(c, errno.ErrCreateTemplateAccount, errMsg)
	//	return
	//}

	//创建一个新的模板，待审核通过之后，会迁移到这个模板。删除旧模板
	if err := m.Create(); err != nil {
		fmt.Println("create error", err)
		h.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	//不管是创建还是更新，都会创建一个 模板名 + ID ，例如 "绩效模板-12.yaml" 的新配置文件，如审核通过，会将这个文件改名为 "绩效模板.yaml" 成为最终可用模板
	filename := "conf/templates/" + r.Name + "-" + util.Uint2Str(m.ID) + ".yaml"

	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		log.Error("audit create function , create file error:", err)
	}

	d, err := yaml.Marshal(&r.Body)
	if err != nil {
		log.Error("audit create function , yaml.Marshal error:", err)
	} else {
		_, err := f.Write(d)
		if err != nil {
			log.Error("audit create function , cannot create the file error:", err)
		}
	}

	//同时需要将旧模板的状态更新为AuditStateWaiting，待审核
	if r.ID > 0 {
		old, err := model.GetTemplate(r.ID)
		if err != nil {
			h.SendResponse(c, errno.ErrDatabase, nil)
			return
		}

		if err := old.UpdateState(model.AuditStateWaiting); err != nil {
			h.SendResponse(c, errno.ErrDatabase, nil)
			return
		}
	}
	//创建的同时需同时创建审核条目
	audit := &model.Audit{}
	userid, _ := c.Get("userid")
	audit.OperatorID = userid.(uint64)
	audit.Object = model.TemplateAuditObject
	if r.ID > 0 {
		audit.Action = model.AUDITUPDATEACTION
		change, _ := template.ComparedTemplate(r.Name, r.Name+"-"+util.Uint2Str(m.ID))
		audit.Body = "描述:更新模板;" + change
	} else {
		audit.Action = model.AUDITCREATEACTION
		audit.Body = "描述:创建模板;" +
			"档案名:" + m.Name + "; "
	}
	audit.OrgObjectID = []int64{int64(r.ID)}
	audit.DestObjectID = []int64{int64(m.ID)}
	audit.State = model.AuditStateWaiting

	audit.Remark = r.Remark

	if err := audit.Create(); err != nil {
		log.Error("audit create error", err)
		h.SendResponse(c, errno.ErrDatabase, err.Error())
		return
	}
	model.CreateOperateRecord(c, fmt.Sprintf("配置模板,模板名: %s", r.Name))
	h.SendResponse(c, nil, nil)
}

func findKeys(excludeFile string) (keys []string) {
	folder := `conf/templates`
	// 获取所有文件
	files, _ := ioutil.ReadDir(folder)
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			runtimeViper := viper.New()
			runtimeViper.AddConfigPath("conf/templates") // 如果没有指定配置文件，则解析默认的配置文件

			filenameOnly := strings.TrimSuffix(file.Name(), ".yaml")
			if filenameOnly == excludeFile { //不需要对比同一个模板
				continue
			}
			runtimeViper.SetConfigName(filenameOnly)

			runtimeViper.SetConfigType("yaml")                  // 设置配置文件格式为YAML
			if err := runtimeViper.ReadInConfig(); err != nil { // viper解析配置文件
				fmt.Println("err", err)
				continue
			}

			for _, key := range runtimeViper.AllKeys() {
				if strings.Contains(key, ".type") {
					t := runtimeViper.GetString(key)
					// Base 基本信息
					// coefficient 系数表
					// Related 关联其它模板
					// 以上这三样，都允许Key可重复。
					if t == "Base" || t == "Coefficient" || t == "Related" {
						continue
					} else {
						s := strings.Split(key, ".")
						keys = append(keys, s[0])
					}
				}
			}
		}

	}
	return keys
}
