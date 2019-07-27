package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"hr-server/pkg/token"
	"time"
)

//记录用户在系统里的操作记录
type OperateRecord struct {
	ID        uint64    `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:createdAt" `
	Body      string    `json:"body"`
}

const OperateRecordTableName = "tb_operate_record"

// TableName :
func (a *OperateRecord) TableName() string {
	return OperateRecordTableName
}

// Create :
func (a *OperateRecord) Create() error {
	return DB.Self.Debug().Create(&a).Error
}

func CreateOperateRecord(c *gin.Context, body string) error {
	ctx, err := token.ParseRequest(c)
	if err != nil {
		fmt.Println("err", err)
		return err
	}
	str := fmt.Sprintf("ID: %d | username: %s | operate: %s  ", ctx.ID, ctx.Username, body)
	record := &OperateRecord{
		Body: str,
	}
	return record.Create()
}

// GetOperateRecord : 获取指定年月的变动记录
// @params offset :
// @params limit :
// @return : record list , total count ,error information
func GetOperateRecord(offset, limit int) (rs []OperateRecord, total uint64, err error) {
	if limit == 0 {
		limit = 10000
	}

	if err := DB.Self.Order(`"createdAt" desc`).Offset(offset).Limit(limit).Find(&rs).Error; err != nil {
		return rs, 0, err
	}
	if err := DB.Self.Model(&OperateRecord{}).Count(&total).Error; err != nil {
		return rs, 0, err
	}

	return rs, total, err
}
