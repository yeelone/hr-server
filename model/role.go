package model

import (
	"errors"
	"hr-server/pkg/constvar"
	"hr-server/util"
	"strconv"
	"strings"
)

type Role struct {
	BaseModel
	Name        string        `json:"name" gorm:"name"`
	Users       []User        `json:"users" gorm:"many2many:user_roles;"`
	Permissions []Permissions `json:"permissions" gorm:"many2many:permissions_roles;"`
}

const RoleTableName = "tb_roles"

// TableName :
func (r *Role) TableName() string {
	return RoleTableName
}

// Create creates a new Role.
func (r *Role) Create() error {
	return DB.Self.Create(&r).Error
}

// DeleteRole deletes the role by the user identifier.
func DeleteRole(id uint64) error {
	role := Role{}
	role.ID = id
	return DB.Self.Delete(&role).Error
}

// Update updates an user Role information.
func (r *Role) Update() (err error) {
	tx := DB.Self.Begin()
	if err := tx.Model(&r).Update(map[string]interface{}{"name": r.Name}).Error; err != nil {
		tx.Rollback()
		return errors.New("无法更新")
	}
	tx.Commit()

	return err
}

//AddRoleUsers :
func AddRoleUsers(rid uint64, uids []uint64) (err error) {
	r := &Role{}

	if r, err = GetRole(rid, false); err != nil {
		return errors.New("User Role is not existed!")
	}

	tx := DB.Self.Begin()

	var users []User
	for _, id := range uids {
		u := User{}
		u.ID = id
		users = append(users, u)
		tx.Model(&u).Association("Roles").Clear()
	}

	err = tx.Model(&r).Association("Users").Append(users).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//RemoveUserGroupUsers :
func RemoveRoleUsers(rid uint64, idList []uint64) (err error) {
	r := &Role{}
	if r, err = GetRole(rid, false); err != nil {
		return errors.New("Role is not existed!")
	}

	tx := DB.Self.Begin()

	uids := make([]string, len(idList))

	for i, id := range idList {
		uids[i] = util.Uint2Str(id)
	}
	err = tx.Model(&r).Exec(" delete from user_roles where user_id in (" + strings.Join(uids, ",") + ") and role_id = " + util.Uint2Str(rid) + " ;").Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// GetRoles :
func GetRole(id uint64, withUsers bool) (result *Role, err error) {
	r := &Role{}
	if id == 0 {
		return result, errors.New("cannot find Role by id " + util.Uint2Str(id))
	}
	err = DB.Self.Select("id,name").First(&r, id).Error
	if withUsers {
		DB.Self.Model(&result).Select("id").Association("Users").Find(&r.Users)
	}
	return r, err
}

//ListRoles :
func ListRoles(offset, limit int, where string, whereKeyword string) (rs []*Role, total int, err error) {
	r := &Role{}
	if limit == 0 {
		limit = constvar.DefaultLimit
	}
	fieldsStr := "id,name"
	if len(where) > 0 {
		if err := DB.Self.Select(fieldsStr).Where(where+" = ?", whereKeyword).Offset(offset).Limit(limit).Find(&rs).Error; err != nil {
			return rs, 0, errors.New("cannot get Role list by where " + where + " and keyword " + whereKeyword)
		}

		if err := DB.Self.Model(r).Where(where+" = ?", whereKeyword).Count(&total).Error; err != nil {
			return rs, 0, errors.New("cannot fetch count of the row")
		}
	} else {
		if err := DB.Self.Select(fieldsStr).Offset(offset).Limit(limit).Find(&rs).Error; err != nil {
			return rs, 0, errors.New("cannot get Role list ")
		}
		if err := DB.Self.Model(r).Count(&total).Error; err != nil {
			return rs, 0, errors.New("cannot fetch count of the row")
		}
	}

	return rs, total, nil

}

//GetRoleRelatedUsers :
func GetRoleRelatedUsers(rid uint64, offset, limit int) (users []*User, total int, err error) {
	if limit == 0 {
		limit = constvar.DefaultLimit
	}
	r := &Role{}
	r.ID = rid

	uids := []uint64{}

	selectSql := ""
	countSql := ""
	if rid == 0 {
		selectSql = "SELECT user_id from user_roles offset " + strconv.Itoa(offset) + " limit " + strconv.Itoa(limit)
	} else {
		selectSql = "SELECT user_id from user_roles where role_id = " + util.Uint2Str(rid) + " offset " + strconv.Itoa(offset) + " limit " + strconv.Itoa(limit)
		countSql = "SELECT  count(user_id) from user_roles where role_id = " + util.Uint2Str(rid)
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

	if rid == 0 {
		DB.Self.Model(User{}).Count(&total)
	} else {
		rows, _ := DB.Self.Raw(countSql).Rows()
		for rows.Next() {
			rows.Scan(&total)
		}
	}

	return users, total, nil
}
