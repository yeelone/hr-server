package model

import (
	"errors"
	"fmt"
	"hr-server/pkg/constvar"
	"hr-server/util"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

//项目比较急，一开始设想的时候将Group用于管理职员档案了，现在需要一个新的群组Model来专门管理User，只能将model命名为UserGroup了
// UserGroup : 用于User管理的群组
type UserGroup struct {
	BaseModel
	Name   string `json:"name" gorm:"column:name;not null"`
	Users  []User `json:"users" gorm:"many2many:user_usergroups;"`
	Parent uint64 `json:"parent" gorm:"column:parent;"`
	Levels string `json:"levels" gorm:"column:levels"` //保存父子层级关系图,例如 pppid.ppid.pid.id
}

// TableName :
func (g *UserGroup) TableName() string {
	return "tb_usergroups"
}

// Create : Create a new UserGroup
func (g *UserGroup) Create() error {
	pm := &UserGroup{}
	g.Levels = "0."
	if g.Parent != 0 {
		pm.BaseModel.ID = g.Parent
		if err := DB.Self.First(&pm).Error; err != nil {
			return errors.New("找不到父目录")
		}
		g.Levels = pm.Levels + util.Uint2Str(g.Parent) + "."
	}

	err := DB.Self.Create(&g).Error
	return err
}

// Update updates an UserGroup information.
// only update name and coefficient
func (g *UserGroup) Update() error {
	_, err := GetUserGroup(g.ID, false)
	if err != nil {
		return err
	}

	tx := DB.Self.Begin()
	if err := tx.Model(&g).Update(map[string]interface{}{"name": g.Name, "parent": g.Parent}).Error; err != nil {
		tx.Rollback()
		return errors.New("无法更新")
	}
	tx.Commit()
	return nil
}

func CountUserGroup() (count int, err error) {
	err = DB.Self.Model(&UserGroup{}).Count(&count).Error
	return count, err
}

//GetAllUserGroup :
func ListUserGroup(offset, limit int, where string, whereKeyword string) (gs []*UserGroup, total int, err error) {
	g := &UserGroup{}
	if limit == 0 {
		limit = constvar.DefaultLimit
	}
	fieldsStr := "id,name,parent,levels"
	if len(where) > 0 {
		if err := DB.Self.Select(fieldsStr).Where(where+" = ?", whereKeyword).Offset(offset).Limit(limit).Find(&gs).Error; err != nil {
			return gs, 0, errors.New("cannot get UserGroup list by where " + where + " and keyword " + whereKeyword)
		}

		if err := DB.Self.Model(g).Where(where+" = ?", whereKeyword).Count(&total).Error; err != nil {
			return gs, 0, errors.New("cannot fetch count of the row")
		}
	} else {
		if err := DB.Self.Select(fieldsStr).Offset(offset).Limit(limit).Find(&gs).Error; err != nil {
			return gs, 0, errors.New("cannot get UserGroup list ")
		}
		if err := DB.Self.Model(g).Count(&total).Error; err != nil {
			return gs, 0, errors.New("cannot fetch count of the row")
		}
	}

	return gs, total, nil

}

//GetUserGroupRelatedUsers :
func GetUserGroupRelatedUsers(id uint64, offset, limit int) (users []*User, total int, err error) {
	if limit == 0 {
		limit = constvar.DefaultLimit
	}
	gs := []UserGroup{}
	if err := DB.Self.Where("levels LIKE ? OR id = ?", "%."+util.Uint2Str(id)+".%", id).Find(&gs).Error; err != nil {
		fmt.Println(err)
		return nil, 0, err
	}

	gids := make([]string, len(gs))
	for i, g := range gs {
		gids[i] = util.Uint2Str(g.ID)
	}
	uids := []uint64{}

	selectSql := ""
	countSql := ""
	if id == 0 {
		selectSql = "SELECT user_id from user_usergroups offset " + strconv.Itoa(offset) + " limit " + strconv.Itoa(limit)
	} else {
		selectSql = "SELECT user_id from user_usergroups where user_group_id in (" + strings.Join(gids, ",") + ")" + " offset " + strconv.Itoa(offset) + " limit " + strconv.Itoa(limit)
		countSql = "SELECT  count(user_id) from user_usergroups where user_group_id in (" + strings.Join(gids, ",") + ")"
	}
	rows, _ := DB.Self.Raw(selectSql).Rows() // Note: Ignoring errors for brevity

	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			return nil, 0, err
		}
		uids = append(uids, id)
	}

	if err := DB.Self.Where(" id in (?)", uids).Find(&users).Error; err != nil {
		return users, 0, err
	}

	if id == 0 {
		DB.Self.Model(User{}).Count(&total)
	} else {
		rows, _ := DB.Self.Raw(countSql).Rows()
		for rows.Next() {
			rows.Scan(&total)
		}
	}

	return users, total, nil
}

func AddUserToDefaultGroup(uid uint64) (err error) {
	gname := viper.GetString("company.name")
	g := &UserGroup{}
	if err := DB.Self.Where("name = ?", gname).First(g).Error; err != nil {
		return err
	}

	idlist := []uint64{uid}
	err = AddUserGroupUsers(g.ID, idlist)
	return err

}

//AddUserGroupUsers :
func AddUserGroupUsers(gid uint64, uids []uint64) (err error) {

	g := &UserGroup{}

	if g, err = GetUserGroup(gid, false); err != nil {
		return errors.New("User Group is not existed!")
	}

	tx := DB.Self.Begin()

	var users []User
	for _, id := range uids {
		u := User{}
		u.ID = id
		users = append(users, u)
		tx.Model(&u).Association("Groups").Clear()
	}

	err = tx.Model(&g).Association("Users").Append(users).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//RemoveUserGroupUsers :
func RemoveGroupUsers(gid uint64, idList []uint64) (err error) {
	g := &UserGroup{}
	if g, err = GetUserGroup(gid, false); err != nil {
		return errors.New("User Group is not existed!")
	}

	tx := DB.Self.Begin()

	uids := make([]string, len(idList))

	for i, id := range idList {
		uids[i] = util.Uint2Str(id)
	}
	err = tx.Model(&g).Exec(" delete from user_usergroups where user_id in (" + strings.Join(uids, ",") + ") and user_group_id = " + util.Uint2Str(gid) + " ;").Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// GetUserGroup :
func GetUserGroup(id uint64, withUsers bool) (result *UserGroup, err error) {
	g := &UserGroup{}
	if id == 0 {
		return result, errors.New("cannot find UserGroup by id " + util.Uint2Str(id))
	}
	err = DB.Self.Select("id,name,parent,levels").First(&g, id).Error
	if withUsers {
		DB.Self.Model(&result).Select("id").Association("Users").Find(&g.Users)
	}
	return g, err
}

// GetUserGroupByName :
func GetUserGroupByName(name string) (result *UserGroup, err error) {
	g := &UserGroup{}
	err = DB.Self.Select("id,name,parent,levels").Where("name = ?", name).First(&g).Error
	return g, err
}

// DeleteUserGroup : delete children UserGroup when parent had deleted
func DeleteUserGroup(id uint64) error {

	group, err := GetUserGroup(id, false)
	if err != nil {
		return err
	}
	cat := &UserGroup{}
	cat.ID = group.ID
	tx := DB.Self.Begin()
	if err := tx.Where("levels LIKE ? OR id = ?", "%."+util.Uint2Str(id)+".%", id).Delete(&cat).Error; err != nil {
		tx.Rollback()
		return errors.New("无法删除")
	}
	tx.Commit()

	return nil
}
