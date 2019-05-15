package permission

import (
	"bufio"
	"fmt"
	h "hrgdrc/handler"
	"hrgdrc/pkg/auth"
	"hrgdrc/pkg/errno"
	"hrgdrc/util"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
	"github.com/spf13/viper"
)

func Create(c *gin.Context) {
	log.Info("permisstion Create function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		fmt.Println("err", err.Error())
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}

	saveCSV(r)
	auth.RefreshEnforcer()
	rsp := CreateResponse{}

	// Show the user information.
	h.SendResponse(c, nil, rsp)

}

func saveCSV(r CreateRequest) {
	file, err := os.Create("conf/permission/" + r.RoleName + ".csv")
	checkError("Cannot create file", err)
	defer file.Close()

	runtimeViper := viper.New()
	runtimeViper.AddConfigPath("conf/permission") // 如果没有指定配置文件，则解析默认的配置文件
	runtimeViper.SetConfigName("permission")

	runtimeViper.SetConfigType("yaml")                  // 设置配置文件格式为YAML
	if err := runtimeViper.ReadInConfig(); err != nil { // viper解析配置文件
		fmt.Println(err)
		return
	}

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for k, fields := range r.Fields {
		for k2, resource := range fields {
			if resource.Checked {
				//注意这里要加个空格在前面
				subject := r.RoleName
				object := runtimeViper.GetString(k + "." + k2 + ".object")
				action := runtimeViper.GetString(k + "." + k2 + ".action")
				s := "p" + ", " + subject + ", " + object + ", " + action + "\n"
				writer.Write([]byte(s))
			}
		}
	}
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
