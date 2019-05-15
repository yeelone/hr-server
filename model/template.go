package model

import (
	"errors"
	"fmt"
	"hrgdrc/util"
	"strings"
)

var (
	Normal_Templte      = "normal"
	Instalment_Template = "instalment"
)

const TemplateAuditObject = "Template"

type Template struct {
	BaseModel
	Name             string            `json:"name"  gorm:"not null;"`
	Type             string            `json:"type" `                          //可以是普通模板和分摊模板
	Months           int               `json:"months"`                         //如果分摊模板，可以设置分几个月摊完
	Startup          bool              `json:"startup"`                        //是否已经开始
	InitData         string            `json:"file"`                           //分推模板需要有初始数据
	Order            int               `json:"order"`                          //模板的顺序
	Groups           []Group           `json:"groups" `                        //指定群组的人进行计算
	UserID           uint64            `json:"user_id" gorm:"not null"`        //模板创建人
	AuditState       int               `json:"audit_state" gorm:"audit_state"` //审核结果
	TemplateAccounts []TemplateAccount `gorm:"many2many:templateaccount_templates"`
}

const TemplateTableName = "tb_template"

// TableName :
func (t *Template) TableName() string {
	return TemplateTableName
}

// Create creates a new base salary.
func (t *Template) Create() (err error) {
	tx := DB.Self.Begin()
	err = tx.Create(&t).Error
	tx.Commit()
	return err
}

// Create creates a new base salary.
func (t *Template) Save() (err error) {
	tx := DB.Self.Begin()

	if t.ID > 0 {
		err = tx.Model(&t).Where("id = ?", t.ID).Updates(t).Error
	} else {
		err = tx.Create(&t).Error
	}

	tx.Commit()

	return err
}

func (t *Template) UpdateState(state int) (err error) {
	tx := DB.Self.Begin()
	if err = tx.Model(&t).Update(map[string]interface{}{"audit_state": state}).Error; err != nil {
		tx.Rollback()
		return errors.New("无法更新")
	}
	tx.Commit()
	return err
}

func GetTemplate(id uint64) (*Template, error) {
	w := &Template{}
	w.ID = id
	if err := DB.Self.First(&w).Error; err != nil {
		return w, err
	}

	return w, nil
}

// GetAccountTemplates  : 获取账套当前所关联的所有的模板
func GetAccountTemplates(account uint64) (ts []Template, err error) {
	tac := &TemplateAccount{}
	tac.ID = account
	if err = DB.Self.Debug().Preload("Templates").First(&tac).Error; err != nil {
		return ts, errors.New("cannot query template from database")
	}
	fmt.Println(util.PrettyJson(tac))
	return tac.Templates, nil
}

func ListTemplates() (ts []Template, err error) {
	if err = DB.Self.Where("audit_state =? ", AuditStatePermit).Order("order").Find(&ts).Error; err != nil {
		return ts, errors.New("cannot query template from database")
	}

	return ts, nil
}

func DeleteTemplate(id uint64) error {
	t := Template{}
	t.ID = id
	return DB.Self.Delete(&t).Error
}

//UpdateTemplateOrder  更新所有模板的顺序
func UpdateTemplateOrder(orders map[uint64]int) (err error) {
	t := &Template{}
	tx := DB.Self.Begin()

	// postgresql 批量更新的方法
	//  update test set info=tmp.info from (values (1,'new1'),(2,'new2'),(6,'new6')) as tmp (id,info) where test.id=tmp.id;
	sql := `update ` + TemplateTableName + ` set  "order"=tmp.order from (values `
	values := []string{}
	for id, order := range orders {
		values = append(values, `(`+fmt.Sprint(id)+`,`+fmt.Sprint(order)+`)`)
	}
	sql += strings.Join(values, ",") + `) as tmp(id,"order") where ` + TemplateTableName + `.id=tmp.id`
	fmt.Println("sql", sql)
	if err = tx.Model(&t).Exec(sql).Error; err != nil {
		tx.Rollback()
		return errors.New("无法更新")
	}
	tx.Commit()
	return err
}
