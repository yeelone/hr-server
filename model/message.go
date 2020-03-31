package model

import (
	"hr-server/pkg/constvar"
	"hr-server/util"
	"strings"
)

// SendId：发件人id；  RecId：收件人id,为0时表示所有人；  TextId：消息id；  Status：标识已读1/未读0；
// Message :
type Message struct {
	BaseModel
	TextId uint64 `json:"textId" gorm:"column:text_id;not null"`
	MType  string `json:"messageType" gorm:"column:message_type;not null"`
	RecId  uint64 `json:"recId" gorm:"column:rec_id;not null"`
	Status int    `json:"status" gorm:"column:status;not null"`
}

const MessageTableName = "tb_message"

// TableName :
func (m *Message) TableName() string {
	return MessageTableName
}

func (m *Message) Create() error {
	err := DB.Self.Create(&m).Error
	return err
}

func (m *Message) SetStatus(status int) error {
	tx := DB.Self.Begin()
	if err := tx.Model(&m).Update(map[string]interface{}{"status": status}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func GetMessages(offset, limit int, recId uint64, where string, whereKeyword string) (mts []MessageText, total int , err error) {
	ms := make([]Message,0)

	if limit == 0 {
		limit = constvar.NoLimit
	}

	if len(where) > 0 {
		err = DB.Self.Offset(offset).Limit(limit).Where(where+" = ? and rec_id = ? ", whereKeyword, recId).Find(&ms).Error
		err = DB.Self.Model(Message{}).Where(where+" = ? and rec_id = ? ", whereKeyword, recId).Count(&total).Error
	} else {
		err = DB.Self.Offset(offset).Limit(limit).Find(&ms).Error
		err = DB.Self.Model(Message{}).Count(&total).Error
	}

	ids := make([]string,0)

	// 将text_id 与 message id 进行映射 ，因为最终要的是message id
	idMap := make(map[uint64]uint64)

	for _, m := range ms {
		ids = append(ids,util.Uint2Str(m.TextId))
		idMap[m.TextId] = m.ID
	}

	if len(ids) > 0 {
		sql := `select * from ` + MessageTextTableName + ` as text where text.id in  ( ` + strings.Join(ids, ",") + `)`
		DB.Self.Raw(sql).Scan(&mts)
	}

	for i, mt := range mts {
		mts[i].MessageId =  idMap[mt.ID]
	}

	return mts, total , err
}

func CheckUserMessage(recId uint64, status int) (private, public, global int) {

	// 检查私信
	if err := DB.Self.Model(Message{}).Where(" rec_id = ? AND status = ? AND message_type =  ?", recId, status, "Private").Count(&private).Error; err != nil {
		return 0, 0, 0
	}

	if err := DB.Self.Model(Message{}).Where(" rec_id = ? AND status = ?  AND message_type =  ?", recId, status, "Public").Count(&public).Error; err != nil {
		return 0, 0, 0
	}

	if err := DB.Self.Model(Message{}).Where(" rec_id = ? AND status = ?  AND message_type =  ?", recId, status, "Global").Count(&global).Error; err != nil {
		return 0, 0, 0
	}

	// 检查public，会稍微麻烦一点，需要先查看用户是否属于那个组或者角色，再判断这条信息是否是发给自己，如果是，会写入message表
	// 先获取用户的组和角色

	groupIdList := make([]string, 0)
	groupIds := ""
	roleIdList := make([]string, 0)
	roleIds := ""

	user, err := GetUserWithGroupAndTag(recId)

	if err != nil {
		return 0, 0, 0
	}

	for _, g := range user.Groups {
		groupIdList = append(groupIdList, util.Uint2Str(g.ID))
	}

	where := ""
	if len(groupIdList) > 0 {
		groupIds = strings.Join(groupIdList, ",")
		where = " AND text.group in (" + groupIds + ") "
	}

	for _, r := range user.Roles {
		roleIdList = append(roleIdList, util.Uint2Str(r.ID))
	}

	if len(roleIdList) > 0 {
		roleIds = strings.Join(roleIdList, ",")

		if len(groupIdList) > 0 {
			where = " AND (  text.group in ( " + groupIds + ") OR text.role in ( " + roleIds + ` ) )`
		} else {
			where = " AND text.role in ( " + roleIds + ` ) `
		}
	}

	// 如果是发送给某个组或者角色，同时这条消息又不在于message表中，则进行统计回馈给用户 ，同时也需要将这些消息加入到message表中。
	sql := `select text.id from ` + MessageTextTableName + ` as text where text.message_type='Public' ` +
		where +
		`AND text.id NOT IN (select m.text_id from ` + MessageTableName + ` as m  where m.rec_id=` + util.Uint2Str(recId) + `)`

	rows, _ := DB.Self.Raw(sql).Rows()

	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			return 0, 0, 0
		}
		public += 1
		msg := Message{}
		msg.TextId = id
		msg.RecId = recId
		msg.Status = 0
		msg.MType = "Public"   // 还是要标记一下这条消息来自于组或者角色
		msg.Create()
	}

	// 如果是全部人，属于系统消息，同时这条消息又不在于message表中，则进行统计回馈给用户 ，同时也需要将这些消息加入到message表中。
	sql = `select text.id from ` + MessageTextTableName + ` as text where text.message_type='Global' AND text.id NOT IN (select m.text_id from ` + MessageTableName + ` as m  where m.rec_id=` + util.Uint2Str(recId) + `)`

	rows, _ = DB.Self.Raw(sql).Rows()

	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			return 0, 0, 0
		}
		public += 1
		msg := Message{}
		msg.TextId = id
		msg.RecId = recId
		msg.Status = 0
		msg.MType = "Global"   // 还是要标记一下这条消息来自于组或者角色
		msg.Create()
	}

	return private, public, global
}
