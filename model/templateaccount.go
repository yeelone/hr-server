package model

import (
	"errors"
	"hrgdrc/pkg/constvar"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
)

type TemplateAccount struct {
	ID        uint64        `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	CreatedAt time.Time     `gorm:"column:createdAt" json:"-"`
	Name      string        `json:"name"  gorm:"not null;unique"`
	Order     pq.Int64Array `json:"order" gorm:"type:integer[]"`
	Groups    []Group       `json:"groups" gorm:"many2many:templateaccount_groups;"` //指定群组的人进行计算
	Templates []Template    `json:"templates" gorm:"many2many:templateaccount_templates;"`
}

// TableName :
func (t *TemplateAccount) TableName() string {
	return "tb_template_account"
}

// Create creates a new template account.
func (t *TemplateAccount) Save() (err error) {
	tx := DB.Self.Begin()

	if t.ID > 0 {
		err = tx.Model(&t).Where("id = ?", t.ID).Updates(t).Error
	} else {
		err = tx.Create(&t).Error
	}

	tx.Commit()

	return err
}

func GetTemplateAccount(id uint64) (*TemplateAccount, error) {
	w := &TemplateAccount{}
	w.ID = id
	if err := DB.Self.Preload("Groups").Preload("Templates").First(&w).Error; err != nil {
		return w, err
	}
	return w, nil
}

func GetTemplateAccountByName(name string) (*TemplateAccount, error) {
	w := &TemplateAccount{}
	if err := DB.Self.Where("name = ?", name).Preload("Templates").First(&w).Error; err != nil {
		return w, err
	}
	return w, nil
}

func ListTemplateAccounts(offset, limit int) (tas []*TemplateAccount, total int, err error) {
	t := &TemplateAccount{}
	if limit == 0 {
		limit = constvar.DefaultLimit
	}

	if err = DB.Self.Offset(offset).Limit(limit).Find(&tas).Error; err != nil {
		return nil, 0, errors.New("cannot query template from database")
	}

	if err := DB.Self.Model(t).Count(&total).Error; err != nil {
		return tas, 0, errors.New("cannot fetch count of the row")
	}
	return tas, total, nil
}

func DeleteTemplateAccount(id uint64) error {
	t := TemplateAccount{}
	t.ID = id

	DB.Self.Model(&t).Association("Groups").Clear()
	DB.Self.Model(&t).Association("Templates").Clear()

	return DB.Self.Delete(&t).Error
}

//AddGroupUsers :
func AddTemplateAccountRelateGroups(tid uint64, IDList []uint64) (err error) {
	t := &TemplateAccount{}
	t.ID = tid

	tx := DB.Self.Begin()
	var groups []Group
	for _, id := range IDList {
		groups = append(groups, Group{BaseModel: BaseModel{ID: id}})
	}

	tx.Model(&t).Association("Groups").Clear()
	err = tx.Model(&t).Association("Groups").Append(groups).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//AddGroupUsers :
func AddTemplateAccountRelateTemplates(tid uint64, IDList []uint64) (err error) {
	t := &TemplateAccount{}
	t.ID = tid

	tx := DB.Self.Begin()
	var templates []Template
	for _, id := range IDList {
		templates = append(templates, Template{BaseModel: BaseModel{ID: id}})
	}

	tx.Model(&t).Association("Templates").Clear()
	err = tx.Model(&t).Association("Templates").Append(templates).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
