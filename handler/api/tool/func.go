package tool

import (
	h "hrgdrc/handler"
	"hrgdrc/pkg/buildinfunc"

	"github.com/gin-gonic/gin"
)

func ListFunc(c *gin.Context) {
	b := &buildinfunc.BuildinFunc{}
	data, err := b.ListFunc()

	if err != nil {
		h.SendResponse(c, err, nil)
		return
	}

	list := make([]map[string]interface{}, 0)
	for _, value := range data {
		list = append(list, value)
	}
	h.SendResponse(c, nil, CreateResponse{List: list})
	return
}
