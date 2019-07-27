package service

import (
	"hr-server/model"
)

// ListUserCoefficient
// @description  list all the groups , because groups would not be too large
func ListUserCoefficient() map[string]*model.Group {
	groups, _, _ := model.ListGroup(0, 5000, "", "")

	gmap := make(map[string]*model.Group)
	for _, g := range groups {
		gmap[g.Name] = g
	}

	return gmap
}
