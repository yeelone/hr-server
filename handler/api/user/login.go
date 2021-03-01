package user

import (
	"encoding/csv"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/auth"
	"hr-server/pkg/errno"
	"hr-server/pkg/token"
	"hr-server/util"
	"io"
	"os"
	"strings"

	"github.com/axgle/mahonia"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
)

// @Summary Login generates the authentication token
// @Produce  json
// @Param username body string true "Username"
// @Param password body string true "Password"
// @Param CaptchaId body string true "CaptchaId"
// @Param CaptchaValue body string true "CaptchaValue"
// @Success 200 {string} json "{"code":0,"message":"OK","data":{"token":"yourtoken"}}"
// @Router /api/login [post]
func Login(c *gin.Context) {
	// Binding the data with the user struct.
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}

	if r.CaptchaId == "" || r.CaptchaValue == "" {
		h.SendResponse(c, errno.ErrCaptcha, "验证码错误")
		return
	}

	if !captcha.VerifyString(r.CaptchaId, r.CaptchaValue) {
		h.SendResponse(c, errno.ErrCaptcha, "验证码错误")
		return
	}

	// Get the user information by the login username.
	d, err := model.GetUserByName(r.Username)
	if err != nil {
		log.Warnf("Login function called. ID: %d | username: %s | info: %s  ", 0, r.Username, "user trying to login")
		h.SendResponse(c, errno.ErrUserNotFound, err.Error())
		return
	}
	// Compare the login password with the user password.
	if err := auth.Compare(d.Password, r.Password); err != nil {
		log.Warnf("Login function called. ID: %d | username: %s | info: %s  ", d.ID, d.Username, "password incorrect")
		h.SendResponse(c, errno.ErrPasswordIncorrect, err.Error())
		return
	}

	role := model.Role{}
	if len(d.Roles) > 0 {
		role = d.Roles[0]
	}
	// Sign the json web token.
	t, err := token.Sign(c, token.Context{ID: d.ID, Username: d.Username, Role: role.Name}, "")
	if err != nil {
		log.Warnf("Login function called. ID: %d | username: %s | info: %s  ", d.ID, d.Username, "token error")
		h.SendResponse(c, errno.ErrToken, err.Error())
		return
	}

	log.Infof("Login function called. ID: %d | username: %s | info: login success", d.ID, d.Username)
	//返回给客户端之前把密码抹除
	d.Password = ""
	c.Set("CurrentUsername", d.Username)
	c.Set("CurrentUserID", d.ID)

	record := &model.OperateRecord{
		Body: d.Username + " 登录系统",
	}

	if err := record.Create(); err != nil {
	}

	permissions := getRolePermissionFromCSVFile(role.Name)

	h.SendResponse(c, nil, model.Token{
		Token:       t,
		User:        d,
		Permissions: getPermissionFieldsFromConf(role.Name, permissions),
	})
}

func Logout(c *gin.Context) {
	h.SendResponse(c, nil, "Successfully logged out")
}

func getPermissionFieldsFromConf(subject string, permissions map[string]model.Resource) map[string]map[string]model.Resource {

	runtimeViper := viper.New()
	runtimeViper.AddConfigPath("conf/permission") // 如果没有指定配置文件，则解析默认的配置文件
	runtimeViper.SetConfigName("permission")

	runtimeViper.SetConfigType("yaml")                  // 设置配置文件格式为YAML
	if err := runtimeViper.ReadInConfig(); err != nil { // viper解析配置文件
		return nil
	}

	keys := make(map[string]map[string]model.Resource)

	for _, key := range runtimeViper.AllKeys() {
		s := strings.Split(key, ".")
		if len(s) > 0 {
			if _, ok := keys[s[0]]; !ok {
				keys[s[0]] = make(map[string]model.Resource)
			}

			if _, ok := keys[s[0]][s[1]]; ok {
				continue
			}
			resource := model.Resource{}
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

func getRolePermissionFromCSVFile(name string) (permissions map[string]model.Resource) {
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

	permissions = make(map[string]model.Resource)
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			//log.Fatal("getRolePermissionFromCSVFile error" + error.Error())
		}
		s := strings.Replace(name+","+line[2]+","+line[3], " ", "", -1)
		permissions[s] = model.Resource{Checked: true}
	}

	return permissions
}
