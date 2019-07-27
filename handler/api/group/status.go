package group

import (
	"strconv"

	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

//Lock
func Lock(c *gin.Context) {
	log.Info("Lock function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	// Get the user id from the url parameter.
	groupID, _ := strconv.Atoi(c.Param("id"))

	// Save changed fields.
	if err := model.LockGroupOrNot(uint64(groupID), true); err != nil {
		h.SendResponse(c, errno.ErrLockGroup, nil)
		return
	}

	h.SendResponse(c, nil, nil)
}

//UnLock
func UnLock(c *gin.Context) {
	log.Info("UnLock function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	// Get the user id from the url parameter.
	groupID, _ := strconv.Atoi(c.Param("id"))

	// Save changed fields.
	if err := model.LockGroupOrNot(uint64(groupID), false); err != nil {
		h.SendResponse(c, errno.ErrUnLockGroup, nil)
		return
	}

	h.SendResponse(c, nil, nil)
}

//Invalid
func Invalid(c *gin.Context) {
	log.Info("Invalid function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	// Get the user id from the url parameter.
	groupID, _ := strconv.Atoi(c.Param("id"))

	// Save changed fields.
	if err := model.InvalidGroupOrNot(uint64(groupID), true); err != nil {
		h.SendResponse(c, errno.ErrInvalidGroup, nil)
		return
	}

	h.SendResponse(c, nil, nil)
}

//Valid
func Valid(c *gin.Context) {
	log.Info("Valid function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	// Get the user id from the url parameter.
	groupID, _ := strconv.Atoi(c.Param("id"))

	// Save changed fields.
	if err := model.InvalidGroupOrNot(uint64(groupID), false); err != nil {
		h.SendResponse(c, errno.ErrValidGroup, nil)
		return
	}

	h.SendResponse(c, nil, nil)
}
