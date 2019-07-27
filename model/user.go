package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"

	"hr-server/pkg/auth"
	"hr-server/pkg/constvar"

	validator "gopkg.in/go-playground/validator.v9"
)

const (
	USERSTATEFREEZE = 1 //冻结状态
	USERSTATEACTIVE = 0 //激活状态
)

// User : User represents a registered user.
type User struct {
	ID        uint64      `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	CreatedAt time.Time   `gorm:"column:createdAt" json:"-"`
	UpdatedAt time.Time   `gorm:"column:updatedAt" json:"-"`
	Email     string      `json:"email" gorm:"column:email;"`
	Username  string      `json:"username" gorm:"column:username;not null" binding:"required" validate:"min=1,max=32"`
	Nickname  string      `json:"nickname" gorm:"column:nichname;not null" `
	IDCard    string      `json:"id_card" gorm:"not null;unique" binding:"required"`
	Password  string      `json:"-" gorm:"column:password;not null" `
	IsSuper   bool        `json:"is_super"`
	Picture   string      `json:"picture"`
	State     int         `json:"state"`
	Groups    []UserGroup `json:"groups" gorm:"many2many:user_usergroups;"`
	Roles     []Role      `json:"roles" gorm:"many2many:user_roles"`
	Profile   Profile     `json:"profile" gorm:"-"`
}

// TableName :
func (u *User) TableName() string {
	return "tb_users"
}

// Create creates a new user account.
func (u *User) Create() error {
	return DB.Self.Create(&u).Error
}

// DeleteUser deletes the user by the user identifier.
func DeleteUser(id uint64) error {
	user := User{}
	user.ID = id
	return DB.Self.Delete(&user).Error
}

// Update updates an user account information.
func (u *User) Update() (err error) {
	tx := DB.Self.Begin()
	if err := tx.Model(&u).Update(map[string]interface{}{"username": u.Username, "nickname": u.Nickname, "id_card": u.IDCard}).Error; err != nil {
		tx.Rollback()
		return errors.New("无法更新:" + err.Error())
	}
	tx.Commit()

	return err
}

// Save :create or update account information
func (u User) Save() (User, error) {
	tx := DB.Self.Begin()

	if u.IsSuper {
		u.IsSuper = true
	} else {
		u.IsSuper = false
	}
	if u.ID > 0 {
		tx.Model(&u).Where("id = ?", u.ID).Updates(u)
	} else if len(u.Email) > 0 {
		u.Password, _ = auth.Encrypt(u.Password)

		if err := tx.Create(&u).Error; err != nil {
			tx.Rollback()
			return u, err
		}
	}

	tx.Commit()

	return u, nil
}

// GetUser gets an user by the user identifier.
func GetUserByName(username string) (*User, error) {
	u := &User{}
	d := DB.Self.Where("username = ?", username).Preload("Roles").First(&u)

	if len(u.IDCard) > 0 {
		// 查询profile
		profile, err := GetProfileByIDCard(u.IDCard)
		if err == nil {
			u.Profile = profile
		}
	}

	return u, d.Error
}

// GetUser gets an user by the user identifier.
func GetUser(id uint64) (*User, error) {
	u := &User{}
	u.ID = id
	d := DB.Self.Preload("Roles").First(&u)

	return u, d.Error
}

// ListUser List all users
func ListUser(username string, offset, limit int) ([]*User, uint64, error) {
	if limit == 0 {
		limit = constvar.DefaultLimit
	}

	users := make([]*User, 0)
	var count uint64

	where := fmt.Sprintf("username like '%%%s%%'", username)
	if err := DB.Self.Model(&User{}).Where(where).Count(&count).Error; err != nil {
		return users, count, err
	}

	if err := DB.Self.Where(where).Offset(offset).Limit(limit).Order("id desc").Preload("Groups").Find(&users).Error; err != nil {
		return users, count, err
	}

	return users, count, nil
}

// Compare with the plain text password. Returns true if it's the same as the encrypted one (in the `User` struct).
func (u *User) Compare(pwd string) (err error) {
	err = auth.Compare(u.Password, pwd)
	return
}

// Encrypt the user password.
func (u *User) Encrypt() (err error) {
	u.Password, err = auth.Encrypt(u.Password)
	return
}

func FreezeUsers(uids []uint64) (err error) {
	tx := DB.Self.Begin()
	if err := tx.Model(&User{}).Where("id in (?)", uids).Update(map[string]interface{}{"state": USERSTATEFREEZE}).Error; err != nil {
		tx.Rollback()
		return errors.New("无法冻结")
	}
	tx.Commit()

	return err
}

func ActiveUsers(uids []uint64) (err error) {
	tx := DB.Self.Begin()
	if err := tx.Model(&User{}).Where("id in (?)", uids).Update(map[string]interface{}{"state": USERSTATEACTIVE}).Error; err != nil {
		tx.Rollback()
		return errors.New("无法激活")
	}
	tx.Commit()

	return err
}

func ResetUsersPassword(uids []uint64) (err error) {
	password, err := auth.Encrypt(viper.GetString("default_password"))
	tx := DB.Self.Begin()
	if err := tx.Model(&User{}).Where("id in (?)", uids).Update(map[string]interface{}{"password": password}).Error; err != nil {
		tx.Rollback()
		return errors.New("无法重置密码")
	}
	tx.Commit()

	return err
}

func ChangeUsersPassword(id uint64, password string) (err error) {
	newPassword, err := auth.Encrypt(password)
	tx := DB.Self.Begin()
	if err := tx.Model(&User{}).Where("id = ?", id).Update(map[string]interface{}{"password": newPassword}).Error; err != nil {
		tx.Rollback()
		return errors.New("无法修改密码")
	}
	tx.Commit()

	return err
}

// Validate the fields.
func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
