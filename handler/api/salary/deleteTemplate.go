package salary

import (
	"fmt"
	h "hrgdrc/handler"
	"hrgdrc/model"
	"hrgdrc/pkg/errno"
	"hrgdrc/util"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func DeleteTemplate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	t, err := model.GetTemplate(uint64(id))

	if err != nil {
		h.SendResponse(c, errno.ErrDatabase,  err.Error())
		return
	}

	if err := model.DeleteTemplate(t.ID); err != nil {
		h.SendResponse(c, errno.ErrDatabase,  err.Error())
		return
	}

	filename := "conf/templates/" + t.Name + ".yaml"
	if util.Exists(filename) {
		err := util.MoveFile(filename, "templates/old/"+t.Name+"-delete-"+time.Now().Format("20060102-150405")+".yaml")
		if err != nil {
			fmt.Println("cannot delete file " + t.Name + err.Error())
		}
	}

	h.SendResponse(c, nil, nil)

}
