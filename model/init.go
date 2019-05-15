package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

type Database struct {
	Self   *gorm.DB
	Docker *gorm.DB
}

var DB *Database

var TableNames = map[string]string{"Profile": "tb_profile", "Template": "tb_template"}

func openDB(username, password, addr, name string) *gorm.DB {

	config := fmt.Sprintf("host=%s dbname=%s user=%s  password=%s sslmode=disable",
		addr,
		name,
		username,
		password,
	)

	db, err := gorm.Open("postgres", config)
	if err != nil {
		//log.Errorf(err, "Database connection failed. Database name: %s", name)
	}

	// set for db connection
	setupDB(db)

	return db
}

func setupDB(db *gorm.DB) {
	db.LogMode(viper.GetBool("gormlog"))
	//db.DB().SetMaxOpenConns(20000) // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	db.DB().SetMaxIdleConns(0) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
}

// InitSelfDB ; used for cli
func InitSelfDB() *gorm.DB {
	return openDB(viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.addr"),
		viper.GetString("db.name"))
}

//GetSelfDB :
func GetSelfDB() *gorm.DB {
	return InitSelfDB()
}

//Init :
func (db *Database) Init() {
	DB = &Database{
		Self: GetSelfDB(),
	}
	initTable()
}

//Close :
func (db *Database) Close() {
	DB.Self.Close()
}

//InitDatabaseTable :
func initTable() {
	var user User
	var profile Profile
	var g Group
	var template Template
	var s Salary
	var t Tag
	var gt GroupTransfer
	var sf SalaryField
	var sc SalaryConfig
	var tas TemplateAccount
	var audit Audit
	var role Role
	var permissions Permissions
	var usergroup UserGroup
	var record Record
	DB.Self.AutoMigrate(&user, &g, &template, &gt, &record, &profile, &tas, &audit, &usergroup, &permissions, &role, &sf, &s, &t, &tas, &sc)
	initAdmin()
	initDefaultGroup()
	initDefaultUserGroup()
	initDefaultRole()
}

//initAdmin: 初始化管理员账号
func initAdmin() {
	u := User{}
	//查看账号是否存在
	email := viper.GetString("admin.email")
	err := DB.Self.Where("email = ?", email).First(&u).Error

	if err != nil {
		u.Email = email
		u.Username = viper.GetString("admin.username")
		u.IDCard = "000000"
		u.Password = viper.GetString("admin.password")
		u.Save()
	}
}

func initDefaultUserGroup() {
	//查看是否存在

	err := DB.Self.Where("name = ?", viper.GetString("company.name")).First(&UserGroup{}).Error
	if err != nil {
		m := UserGroup{
			Name:   viper.GetString("company.name"),
			Parent: 0,
		}
		m.Create()
	}
}

func initDefaultGroup() {
	//查看是否存在
	depart := Group{}
	err := DB.Self.Where("name = ?", "部门").First(&depart).Error
	if err != nil {
		m := Group{
			Name:        "部门",
			Code:        0,
			Parent:      0,
			Coefficient: 0,
			Locked:      true,
		}
		m.Create()
	}

	err = DB.Self.Where("name = ?", "空部门").First(&Group{}).Error
	if err != nil {
		m := Group{
			Name:        "空部门",
			Code:        0,
			Parent:      depart.ID,
			Coefficient: 0,
			Locked:      true,
			IsDefault:   true,
		}
		m.Create()
	}

	g := Group{}
	err = DB.Self.Where("name = ?", "学历").First(&g).Error
	if err != nil {
		m := Group{
			Name:        "学历",
			Code:        0,
			Parent:      0,
			Coefficient: 0,
			Locked:      true,
		}
		m.Create()
	}

	err = DB.Self.Where("name = ?", "无学历").First(&Group{}).Error
	if err != nil {
		m := Group{
			Name:        "无学历",
			Code:        0,
			Parent:      g.ID,
			Coefficient: 0,
			Locked:      true,
			IsDefault:   true,
		}
		m.Create()
	}

	g = Group{}
	err = DB.Self.Where("name = ?", "岗位").First(&g).Error
	if err != nil {
		m := Group{
			Name:        "岗位",
			Code:        0,
			Parent:      0,
			Coefficient: 0,
			Locked:      true,
		}
		m.Create()
	}

	err = DB.Self.Where("name = ?", "空岗位").First(&Group{}).Error
	if err != nil {
		m := Group{
			Name:        "空岗位",
			Code:        0,
			Parent:      g.ID,
			Coefficient: 0,
			Locked:      true,
			IsDefault:   true,
		}
		m.Create()
	}

	g = Group{}
	err = DB.Self.Where("name = ?", "职称").First(&g).Error
	if err != nil {
		m := Group{
			Name:        "职称",
			Code:        0,
			Parent:      0,
			Coefficient: 0,
			Locked:      true,
		}
		m.Create()
	}
	err = DB.Self.Where("name = ?", "无职称").First(&Group{}).Error
	if err != nil {
		m := Group{
			Name:        "无职称",
			Code:        0,
			Parent:      g.ID,
			Coefficient: 0,
			Locked:      true,
			IsDefault:   true,
		}
		m.Create()
	}

	g = Group{}
	err = DB.Self.Where("name = ?", "状态").First(&g).Error
	if err != nil {
		m := Group{
			Name:        "状态",
			Code:        0,
			Parent:      0,
			Coefficient: 0,
			Locked:      true,
		}
		m.Create()
	}
	err = DB.Self.Where("name = ?", "无状态").First(&Group{}).Error
	if err != nil {
		m := Group{
			Name:        "无状态",
			Code:        0,
			Parent:      g.ID,
			Coefficient: 0,
			Locked:      true,
			IsDefault:   true,
		}
		m.Create()
	}
}

func initDefaultRole() {
	//查看是否存在
	err := DB.Self.Where("name = ?", "操作岗").First(&Role{}).Error
	if err != nil {
		m := Role{
			Name: "操作岗",
		}
		m.Create()
	}

	err = DB.Self.Where("name = ?", "复核岗").First(&Role{}).Error
	if err != nil {
		m := Role{
			Name: "复核岗",
		}
		m.Create()
	}

	err = DB.Self.Where("name = ?", "主管岗").First(&Role{}).Error
	if err != nil {
		m := Role{
			Name: "主管岗",
		}
		m.Create()
	}

	err = DB.Self.Where("name = ?", "管理岗").First(&Role{}).Error
	if err != nil {
		m := Role{
			Name: "管理岗",
		}
		m.Create()
	}

	err = DB.Self.Where("name = ?", "查询岗").First(&Role{}).Error
	if err != nil {
		m := Role{
			Name: "查询岗",
		}
		m.Create()
	}
}
