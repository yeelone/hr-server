package permission

import (
	"encoding/csv"
	"fmt"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/axgle/mahonia"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Permission struct {
	Subject string `json:"subject"`
	Object  string `json:"object"`
	Action  string `json:"action"`
}

func Get(c *gin.Context) {
	rid, _ := strconv.Atoi(c.Param("id"))
	role, err := model.GetRole(uint64(rid), false)
	if err != nil {
		h.SendResponse(c, errno.ErrUserBelongRoles, nil)
		return
	}

	resp := &CreateResponse{}
	permissions := getRolePermissionFromCSVFile(role.Name)
	resp.Fields = getPermissionFieldsFromConf(role.Name, permissions)
	// fmt.Println(util.PrettyJson(resp.Fields))
	h.SendResponse(c, nil, resp)
}

func getPermissionFieldsFromConf(subject string, permissions map[string]Resource) map[string]map[string]Resource {

	runtimeViper := viper.New()
	runtimeViper.AddConfigPath("conf/permission") // 如果没有指定配置文件，则解析默认的配置文件
	runtimeViper.SetConfigName("permission")

	runtimeViper.SetConfigType("yaml")                  // 设置配置文件格式为YAML
	if err := runtimeViper.ReadInConfig(); err != nil { // viper解析配置文件
		fmt.Println(err)
		return nil
	}

	keys := make(map[string]map[string]Resource)

	for _, key := range runtimeViper.AllKeys() {
		s := strings.Split(key, ".")
		if len(s) > 0 {
			if _, ok := keys[s[0]]; !ok {
				keys[s[0]] = make(map[string]Resource)
			}

			if _, ok := keys[s[0]][s[1]]; ok {
				continue
			}
			resource := Resource{}
			resource.ID = runtimeViper.GetString(s[0] + "." + s[1] + ".resource")
			str := subject + "," + runtimeViper.GetString(s[0]+"."+s[1]+".object") + "," + runtimeViper.GetString(s[0]+"."+s[1]+".action")
			if _, ok := permissions[str]; ok {
				resource.Checked = true
				keys[s[0]][s[1]] = resource
			} else {
				resource.Checked = false
				keys[s[0]][s[1]] = resource
			}
		}
	}
	return keys
}

func getRolePermissionFromCSVFile(name string) (permissions map[string]Resource) {
	orgFilename := "conf/permission/" + name + ".csv"
	if !util.Exists(orgFilename) {
		return
	}
	file, err := os.Open(orgFilename)
	if err != nil {
	}
	defer file.Close()

	decoder := mahonia.NewDecoder("utf8")            // 把原来ANSI格式的文本文件里的字符，用utf8进行解码。
	reader := csv.NewReader(decoder.NewReader(file)) // 这样，最终返回的字符串就是utf-8了。（go只认utf8）

	permissions = make(map[string]Resource)
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal("getRolePermissionFromCSVFile error" + error.Error())
		}
		s := strings.Replace(name+","+line[2]+","+line[3], " ", "", -1)
		permissions[s] = Resource{Checked: true}
	}

	return permissions
}
