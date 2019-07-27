package profile

import (
	"errors"
	"fmt"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/errno"
	"hr-server/util"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

func ListProfileTags(c *gin.Context) {
	log.Info("ListProfileTags function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})

	pid, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.SendResponse(c, errno.ErrQueryParams, nil)
		return
	}
	profile, err := model.GetProfileWithTags(pid)
	if err != nil {
		h.SendResponse(c, errno.ErrProfileNotFound, nil)
		return
	}
	tags, err := covertTagsToMap(profile.Tags)
	if err != nil {
		h.SendResponse(c, errno.ErrTagNoFount, nil)
		return
	}
	profile.Tags = nil //把tags 去掉，避免重复
	h.SendResponse(c, nil, ListTagsResponse{Profile: profile, Tags: tags})
	return
}

func covertTagsToMap(tags []model.Tag) (result []TagResponse, err error) {
	fmt.Println("tags", tags)

	tids := make([]uint64, len(tags))

	for i, tag := range tags {
		tids[i] = tag.Parent
	}

	parent, err := model.GetTagsByIDList(tids)
	if err != nil {
		return nil, errors.New("cannot fetch tag with special children id ")
	}

	result = make([]TagResponse, len(parent))
	for i, p := range parent {
		response := TagResponse{}
		response.Tag = p
		for _, tag := range tags {
			if tag.Parent == p.ID {
				response.Children = append(response.Children, tag)
			}
		}
		result[i] = response
	}

	return result, err
}
