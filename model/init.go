package model

import (
	"bufio"
	"fmt"
	"github.com/lexkong/log"
	"hr-server/pkg/auth"
	"os"
	"strings"

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
	var operate OperateRecord
	var salaryProfileConfig SalaryProfileConfig
	var message Message
	var messageText MessageText
	DB.Self.AutoMigrate( &message,&messageText, &role, &record, &salaryProfileConfig, &t, &operate, &user, &g, &template, &gt, &profile, &tas, &audit, &usergroup, &permissions, &sf, &s, &tas, &sc)
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

	//查看是否存在
	role := Role{}
	err = DB.Self.Where("name = ?", viper.GetString("admin.role")).First(&role).Error
	if err != nil {
		m := Role{
			Name: viper.GetString("admin.role"),
		}
		m.Create()

		role = m
	}

	users := []uint64{u.ID}
	if err = AddRoleUsers(role.ID, users); err != nil {
		log.Fatal("cannot add user to role ", err)
	}

	// 超级管理员拥有所有的权限
	filename := "conf/permission/" + role.Name + ".csv"
	err = os.Remove(filename) //删除文件test.txt
	if err != nil {
		//如果删除失败则输出 file remove Error!
		fmt.Println("file remove Error!  ", err)
	} else {
		//如果删除成功则输出 file remove OK!
		fmt.Println("file remove OK!")
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Cannot create permission file ", err)
	}
	defer file.Close()

	runtimeViper := viper.New()
	runtimeViper.AddConfigPath("conf/permission") // 如果没有指定配置文件，则解析默认的配置文件
	runtimeViper.SetConfigName("permission")

	runtimeViper.SetConfigType("yaml")                  // 设置配置文件格式为YAML
	if err := runtimeViper.ReadInConfig(); err != nil { // viper解析配置文件
		fmt.Println(err)
		return
	}

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	uni := make(map[string]struct{})
	for _, key := range runtimeViper.AllKeys() {
		keyStr := strings.Split(key, ".")
		if len(keyStr) == 3 {
			subject := role.Name
			object := runtimeViper.GetString(keyStr[0] + "." + keyStr[1] + ".object")
			action := runtimeViper.GetString(keyStr[0] + "." + keyStr[1] + ".action")
			s := "p" + ", " + subject + ", " + object + ", " + action + "\n"

			if _, ok := uni[s]; !ok {
				uni[s] = struct{}{}
				writer.Write([]byte(s))
			}
		}
	}

	writer.Flush()

	auth.MergePermission("./conf/permission")
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
		depart = Group{
			Name:        "部门",
			Code:        0,
			Parent:      0,
			Coefficient: 0,
			Locked:      true,
		}
		depart.Create()
	}

	g := Group{}
	err = DB.Self.Where("name = ?", "学历").First(&g).Error
	if err != nil {
		g = Group{
			Name:        "学历",
			Code:        0,
			Parent:      0,
			Coefficient: 0,
			Locked:      true,
		}
		g.Create()
	}

	g1 := Group{}
	err = DB.Self.Where("name = ?", "博士后").First(&g1).Error
	if err != nil {
		g1 = Group{
			Name:        "博士后",
			Code:        0,
			Parent:      g.ID,
			Coefficient: 0,
			Locked:      false,
		}
		g1.Create()
	}

	g1 = Group{}
	err = DB.Self.Where("name = ?", "博士").First(&g1).Error
	if err != nil {
		g1 = Group{
			Name:        "博士",
			Code:        0,
			Parent:      g.ID,
			Coefficient: 0,
			Locked:      false,
		}
		g1.Create()
	}

	g1 = Group{}
	err = DB.Self.Where("name = ?", "研究生").First(&g1).Error
	if err != nil {
		g1 = Group{
			Name:        "研究生",
			Code:        0,
			Parent:      g.ID,
			Coefficient: 0,
			Locked:      false,
		}
		g1.Create()
	}

	g1 = Group{}
	err = DB.Self.Where("name = ?", "本科").First(&g1).Error
	if err != nil {
		g1 = Group{
			Name:        "本科",
			Code:        0,
			Parent:      g.ID,
			Coefficient: 0,
			Locked:      false,
			IsDefault:   true,
		}
		g1.Create()
	}

	g1 = Group{}
	err = DB.Self.Where("name = ?", "大专").First(&g1).Error
	if err != nil {
		g1 = Group{
			Name:        "大专",
			Code:        0,
			Parent:      g.ID,
			Coefficient: 0,
			Locked:      false,
		}
		g1.Create()
	}

	g1 = Group{}
	err = DB.Self.Where("name = ?", "高中").First(&g1).Error
	if err != nil {
		g1 = Group{
			Name:        "高中",
			Code:        0,
			Parent:      g.ID,
			Coefficient: 0,
			Locked:      false,
		}
		g1.Create()
	}

	g1 = Group{}
	err = DB.Self.Where("name = ?", "中专").First(&g1).Error
	if err != nil {
		g1 = Group{
			Name:        "中专",
			Code:        0,
			Parent:      g.ID,
			Coefficient: 0,
			Locked:      false,
		}
		g1.Create()
	}

	g1 = Group{}
	err = DB.Self.Where("name = ?", "初中").First(&g1).Error
	if err != nil {
		g1 = Group{
			Name:        "初中",
			Code:        0,
			Parent:      g.ID,
			Coefficient: 0,
			Locked:      false,
		}
		g1.Create()
	}

	g1 = Group{}
	err = DB.Self.Where("name = ?", "小学").First(&g1).Error
	if err != nil {
		g1 = Group{
			Name:        "小学",
			Code:        0,
			Parent:      g.ID,
			Coefficient: 0,
			Locked:      false,
		}
		g1.Create()
	}

	g = Group{}
	err = DB.Self.Where("name = ?", "岗位").First(&g).Error
	if err != nil {
		g = Group{
			Name:        "岗位",
			Code:        0,
			Parent:      0,
			Coefficient: 0,
			Locked:      true,
		}
		g.Create()
	}

	g = Group{}
	err = DB.Self.Where("name = ?", "职称").First(&g).Error
	if err != nil {
		g = Group{
			Name:        "职称",
			Code:        0,
			Parent:      0,
			Coefficient: 0,
			Locked:      true,
		}
		g.Create()
	}

	g = Group{}
	err = DB.Self.Where("name = ?", "状态").First(&g).Error
	if err != nil {
		g = Group{
			Name:        "状态",
			Code:        0,
			Parent:      0,
			Coefficient: 0,
			Locked:      true,
		}
		g.Create()
	}
	err = DB.Self.Where("name = ?", "在职").First(&Group{}).Error
	if err != nil {
		m := Group{
			Name:        "在职",
			Code:        0,
			Parent:      g.ID,
			Coefficient: 0,
			Locked:      true,
			IsDefault:   true,
		}
		m.Create()
	}

	err = DB.Self.Where("name = ?", "离职").First(&Group{}).Error
	if err != nil {
		m := Group{
			Name:        "离职",
			Code:        0,
			Parent:      g.ID,
			Coefficient: 0,
			Locked:      true,
		}
		m.Create()
	}

	err = DB.Self.Where("name = ?", "退休").First(&Group{}).Error
	if err != nil {
		m := Group{
			Name:        "退休",
			Code:        0,
			Parent:      g.ID,
			Coefficient: 0,
			Locked:      true,
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
