package profile

import (
	"github.com/gin-gonic/gin"
	h "hrgdrc/handler"
	"hrgdrc/model"
	"hrgdrc/pkg/errno"
	"strconv"
)

func GetTransfer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	transfer, err := model.GetProfileTransfer(id)
	if err != nil {
		h.SendResponse(c, errno.ErrUserNotFound, err.Error())
		return
	}

	h.SendResponse(c, nil, TransferResponse{Transfer: transfer})
}
