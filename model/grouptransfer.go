package model

import (
	"errors"
	"fmt"
	"github.com/lib/pq"
	"hrgdrc/util"
	"regexp"
	"time"
)

// GroupTransfer :
type GroupTransfer struct {
	ID                  uint64        `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	CreatedAt           time.Time     `gorm:"column:createdAt" json:"created_at"`
	Profile             uint64        `json:"profile" gorm:"not null"`
	OldGroup            uint64        `json:"old_group" gorm:"not null"`
	NewGroup            uint64        `json:"new_group" gorm:"not null"`
	NewGroupCombination pq.Int64Array `json:"-" gorm:"column:new_group_combination;type:integer[]"`
	Description         string        `json:"description" `
	AddedTagsRecord     pq.Int64Array `json:"-" gorm:"column:added_tags_record;type:integer[]"` //记录当调动组时，为profile添加了哪些标签
	ProfileName         string        `json:"profile_name" gorm:"-"`
	OldGroupName        string        `json:"old_group_name" gorm:"-"`
	NewGroupName        string        `json:"new_group_name" gorm:"-"`
}

var TRANSFER_TABLENAME = "tb_group_transfer"

// TableName :
func (t *GroupTransfer) TableName() string {
	return TRANSFER_TABLENAME
}

// Create : Create a new Group Transfer
// old group and new group must in the same parent group
func (t *GroupTransfer) Create() error {

	newg, err := GetGroup(t.NewGroup, false)
	if err != nil {
		return err
	}

	oldg, err := GetGroup(t.OldGroup, false)
	if err != nil {
		oldg = &Group{}
		oldg.ID = 0
		oldg.Parent = newg.Parent
		oldg.Levels = newg.Levels
	}

	// 父群组的格式为  0.1.
	// 0表示最顶层，1是父群组的ID
	reg := regexp.MustCompile(`0\.(\d+)\..*`)
	m1 := reg.FindStringSubmatch(oldg.Levels)
	m2 := reg.FindStringSubmatch(newg.Levels)
	if len(m1) < 1 {
		return errors.New("cannot find group's parent")
	}
	if len(m2) < 1 {
		return errors.New("cannot find group's parent")
	}

	if m1[1] == m2[1] {
		err := MoveProfileToNewGroup(t.Profile, t.OldGroup, t.NewGroup)
		if err != nil {
			return errors.New("cannot move to the new group,err :" + err.Error())
		}
		return DB.Self.Create(&t).Error
	}
	return errors.New("not in the same parent group")

}

// Update updates an group information.
// only update name and coefficient
func (g *GroupTransfer) Update() error {
	tx := DB.Self.Begin()
	if err := tx.Model(&g).Debug().Update(map[string]interface{}{"added_tags_record": g.AddedTagsRecord, "new_group_combination": g.NewGroupCombination}).Error; err != nil {
		tx.Rollback()
		return errors.New("无法更新")
	}
	tx.Commit()
	return nil
}

func GetTransfer(id uint64) (t GroupTransfer, err error) {
	g := GroupTransfer{}
	g.ID = id

	err = DB.Self.Where("id = ?", id).First(&t).Error
	return t, err
}

func GetTransferByNewGroupAndProfile(gid, pid uint64) (t GroupTransfer, err error) {
	err = DB.Self.Where("new_group = ? and profile = ?", gid, pid).First(&t).Error
	return t, err
}

func GetTransferByNewGroupCombination(gid, pid uint64) (t GroupTransfer, err error) {
	sql := `select * from ` + TRANSFER_TABLENAME + ` where profile =` + util.Uint2Str(pid) +
		` and new_group_combination @> array[` + util.Uint2Str(gid) + `]`
	//sql := `select * from ` + TRANSFER_TABLENAME + ` where profile = ? and new_group_combination @> array[?]`

	//err = DB.Self.Debug().Where(" profile = ? and new_group_combination @> array[?]" , pid , gid ).First(&t).Error
	fmt.Println(sql)
	err = DB.Self.Raw(sql).Scan(&t).Error
	return t, err
}

// GetProfileTransfer : 得到指定员工的调动记录
func GetProfileTransfer(pid uint64) (gs []GroupTransfer, err error) {
	sql := `select b.name as profile_name, c.name as old_group_name, d.name as new_group_name,a."createdAt" ,a.description as description  from ` + TRANSFER_TABLENAME + ` as a  ` +
		` left join ` + ProfileTableName + ` as b   on a.profile = b.id ` +
		` left join ` + GroupTableName + ` as c on a.old_group=c.id ` +
		` left join ` + GroupTableName + ` as d on a.new_group=d.id ` +
		` where a.profile = ?  order by a."createdAt" desc`
	err = DB.Self.Raw(sql, pid).Scan(&gs).Error
	return gs, err
}
func GetAllProfileTransfer() (gs []GroupTransfer, err error) {
	sql := `select b.name as profile_name, c.name as old_group_name, d.name as new_group_name,a."createdAt" ,a.description as description  from ` + TRANSFER_TABLENAME + ` as a  ` +
		` left join ` + ProfileTableName + ` as b   on a.profile = b.id ` +
		` left join ` + GroupTableName + ` as c on a.old_group=c.id ` +
		` left join ` + GroupTableName + ` as d on a.new_group=d.id ` +
		` order by a."createdAt" desc`
	err = DB.Self.Raw(sql).Scan(&gs).Error
	return gs, err
}
