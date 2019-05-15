package model

import "time"

type Record struct {
	ID        uint64     `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	CreatedAt time.Time  `gorm:"column:createdAt" `
	Body         string `json:"body"`
}

const RecordTableName = "tb_record"

// TableName :
func (a *Record) TableName() string {
	return RecordTableName
}

// Create :
func (a *Record) Create() error {
	return DB.Self.Create(&a).Error
}

// GetRecordListByMonth : 获取指定年月的变动记录
// @params date : "2018-01"  只包括年月，但是 "2018-01-01"也可以
// @params offset :
// @params limit :
// @return : record list , total count ,error information
func GetRecordListByMonth(date string , offset, limit int  )( rs []Record, total uint64 , err error) {
	if limit == 0 {
		limit = 10000
	}

	if err := DB.Self.Where(`"createdAt" >= to_date(?,'YYYY-MM')`, date ).Offset(offset).Limit(limit).Find(&rs).Error; err != nil {
		return rs, 0, err
	}
	if err := DB.Self.Model(&Record{}).Where(`"createdAt" >= to_date(?,'YYYY-MM')`, date ).Count(&total).Error; err != nil {
		return rs, 0, err
	}

	return rs , total ,err
}