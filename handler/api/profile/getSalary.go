package profile

import (
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetSalary(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	_, err = model.GetSalaryFieldByProfileAndMonth("2020", "02", id)
	if err != nil {
		h.SendResponse(c, errno.ErrUserNotFound, err.Error())
		return
	}
	h.SendResponse(c, nil, CreateResponse{})
}
