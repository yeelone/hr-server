package model

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"hrgdrc/pkg/constvar"
	"hrgdrc/util"
	"strconv"
	"strings"
)

//Tag :
type Tag struct {
	ID          uint64    `json:"id" gorm:"primary_key"`
	Name        string    `json:"name" gorm:"not null;"`
	Coefficient float64   `json:"coefficient" gorm:"default:0"` //每个群组有一个系数，可用于计算
	Parent      uint64    `json:"parent" gorm:"column:parent;default:0"`
	Users       []User    `json:"users" gorm:"many2many:user_tags;"`
	Profiles    []Profile `json:"profiles" gorm:"many2many:profile_tags;"`
	Groups      []Group   `json:"groups" gorm:"many2many:group_tags;"`
}

const TagTableName = "tb_tags"

// TableName :
func (t *Tag) TableName() string {
	return TagTableName
}

//GetAllTags :
func GetAllTags(offset, limit int, where string, whereKeyword string) (ts []*Tag, total uint64, err error) {
	t := &Tag{}
	if limit == 0 {
		limit = constvar.NoLimit
	}
	fieldsStr := "ID,name,coefficient,parent"
	err = DB.Self.Select(fieldsStr).Offset(offset).Limit(limit).Find(&ts).Error
	err = DB.Self.Model(t).Count(&total).Error

	if len(where) > 0 {
		err = DB.Self.Select(fieldsStr).Offset(offset).Limit(limit).Where(where+" = ?", whereKeyword).Find(&ts).Error
		err = DB.Self.Model(t).Where(where+" = ?", whereKeyword).Count(&total).Error
	} else {
		err = DB.Self.Select(fieldsStr).Offset(offset).Limit(limit).Find(&ts).Error
		err = DB.Self.Model(t).Count(&total).Error
	}

	return ts, total, err
}

// Create : Create a new Tag
func (t *Tag) Create() (err error) {

	temp := &Tag{}
	if err = DB.Self.Where("name = ?", t.Name).First(&temp).Error; err != nil {
		//不存在则可以新建
		return DB.Self.Create(&t).Error
	}

	//如果已找到，则要对比是否在同一个标签组里，同组里不可同名
	if temp.Parent != t.Parent {
		return DB.Self.Create(&t).Error
	}
	err = errors.New("conflict name in same tag group ")

	return err
}

//Update :
func (t Tag) Update() (err error) {
	tx := DB.Self.Begin()

	err = tx.Save(&t).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func GetTagsByIDList(ids []uint64) (tags []Tag, err error) {
	err = DB.Self.Where("id IN (?)", ids).Find(&tags).Error
	return tags, err
}

// relateUsers
func RelatedTagUsers(tagID uint64, keys []uint64) (err error) {
	t := &Tag{ID: tagID}
	tx := DB.Self.Begin()
	var users []*User
	for _, id := range keys {
		user := &User{}
		user.ID = id
		users = append(users, user)
	}

	tx.Model(&t).Association("Users").Clear()
	err = tx.Model(&t).Association("Users").Append(users).Error

	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//GetTagRelatedProfiles :
func GetTagRelatedProfiles(id uint64, offset, limit int) (profiles []Profile, total uint64, err error) {
	if limit == 0 {
		limit = constvar.DefaultLimit
	}
	selectSql := ""
	countSql := ""

	if id == 0 {
		selectSql = "SELECT profile_id from profile_tags offset " + strconv.Itoa(offset) + " limit " + strconv.Itoa(limit)
		countSql = "SELECT count(profile_id) from profile_tags "
	} else {
		selectSql = "SELECT profile_id from profile_tags where tag_id=" + util.Uint2Str(id) + " offset " + strconv.Itoa(offset) + " limit " + strconv.Itoa(limit)
		countSql = "SELECT count(profile_id) from profile_tags where tag_id=" + util.Uint2Str(id) + ""
	}
	rows, _ := DB.Self.Raw(selectSql).Rows() // Note: Ignoring errors for brevity

	profileIDs := []uint64{}
	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			return nil, 0, err
		}
		profileIDs = append(profileIDs, id)
	}

	if err := DB.Self.Where(" id in (?)", profileIDs).Find(&profiles).Error; err != nil {
		return profiles, 0, err
	}

	if id == 0 {
		DB.Self.Model(Profile{}).Count(&total)
	} else {
		rows, _ := DB.Self.Raw(countSql).Rows()
		for rows.Next() {
			rows.Scan(&total)
		}
	}

	return profiles, total, nil
}

// AddTagProfiles
func AddTagProfiles(tid uint64, IDList []uint64) (err error) {
	t := &Tag{ID: tid}
	tx := DB.Self.Begin()
	var profiles []*Profile
	for _, id := range IDList {
		profile := &Profile{}
		profile.BaseModel.ID = id
		profiles = append(profiles, profile)
	}

	// tx.Model(&t).Association("Profiles").Clear()
	err = tx.Model(&t).Association("Profiles").Append(profiles).Error

	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func IfTagHaveProfile(tag string, profile uint64) (result bool) {
	t, err := GetTagByName(tag)
	if err != nil {
		return false
	}
	sql := `select  * from profile_tags where profile_id=? and tag_id=?`
	rows, _ := DB.Self.Raw(sql, profile, t.ID).Rows()

	for rows.Next() {
		var profile_id uint64
		if err := rows.Scan(&profile_id); err != nil {
			return true
		} else {
			return false
		}
	}

	return false
}

// ClearThenAddProfileTags :替换档案所关联的所有标签，删除旧的关联，添加新的关联
func ClearThenAddProfileTags(pid uint64, tids []uint64) (err error) {
	tx := DB.Self.Begin()
	err = ClearProfileTags(pid)
	if err != nil {
		return err
	}

	insertStr := []string{}
	for _, id := range tids {
		insertStr = append(insertStr, "("+util.Uint2Str(pid)+", "+util.Uint2Str(id)+")")
	}
	err = tx.Exec(" delete from profile_tags where profile_id = " + util.Uint2Str(pid) + " ;").Error

	err = tx.Exec(" insert into profile_tags(profile_id,tag_id) values" + strings.Join(insertStr, ",") + ";").Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

// AddProfileTags : 删除旧的关联，添加新的关联
func AddProfileTags(pid uint64, tids []uint64) (err error) {
	tx := DB.Self.Begin()
	deleteStr := []string{}
	insertStr := []string{}
	for _, id := range tids {
		deleteStr = append(deleteStr, "(profile_id="+util.Uint2Str(pid)+" and tag_id="+util.Uint2Str(id)+")")
		insertStr = append(insertStr, "("+util.Uint2Str(pid)+", "+util.Uint2Str(id)+")")
	}
	err = tx.Debug().Exec(" delete from profile_tags where " + strings.Join(deleteStr, " OR ") + " ;").Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Debug().Exec(" insert into profile_tags(profile_id,tag_id) values" + strings.Join(insertStr, ",") + ";").Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

// ClearThenAddGroupTags :替换群组所关联的所有标签，删除旧的关联，添加新的关联
func ClearThenAddGroupTags(gid uint64, tids []uint64) (err error) {
	tx := DB.Self.Begin()

	insertStr := []string{}
	for _, id := range tids {
		insertStr = append(insertStr, "("+util.Uint2Str(gid)+", "+util.Uint2Str(id)+")")
	}
	err = tx.Exec(" delete from tag_groups where group_id = " + util.Uint2Str(gid) + " ;").Error
	fmt.Println("delete err", err)
	if err != nil {
		tx.Rollback()
		return err
	}
	if len(insertStr) > 0 {
		err = tx.Exec(" insert into tag_groups(group_id,tag_id) values" + strings.Join(insertStr, ",") + ";").Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return err
}

//RemoveTagProfiles :
func RemoveTagProfiles(tid uint64, IDList []uint64) (err error) {
	t := &Tag{}
	if t, err = GetTag(tid, false); err != nil {
		return errors.New("Tag is not existed!")
	}

	tx := DB.Self.Begin()

	profileIDs := make([]string, len(IDList))

	for i, id := range IDList {
		profileIDs[i] = util.Uint2Str(id)
	}
	err = tx.Model(&t).Exec(" delete from profile_tags where profile_id in (" + strings.Join(profileIDs, ",") + ") and tag_id = " + util.Uint2Str(tid) + " ;").Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//ClearProfileTags : 清除用户所有的标签
func ClearProfileTags(pid uint64) (err error) {
	tx := DB.Self.Begin()
	err = tx.Exec(" delete from profile_tags where profile_id = " + util.Uint2Str(pid) + " ;").Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//ClearGroupTags : 清除群组所有的标签
func ClearGroupTags(gid uint64) (err error) {
	tx := DB.Self.Begin()
	err = tx.Exec(" delete from group_tags where group_id = " + util.Uint2Str(gid) + " ;").Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//ClearProfileTags : 清除用户所有的标签
func RemoveProfileTags(pid uint64, tids []uint64) (err error) {
	tidStr := []string{}
	if len(tids) < 1 {
		return nil
	}
	for _, id := range tids {
		tidStr = append(tidStr, fmt.Sprint(id))
	}
	tx := DB.Self.Begin()
	err = tx.Debug().Exec(" delete from profile_tags where profile_id = " + util.Uint2Str(pid) + " AND tag_id IN (" + strings.Join(tidStr, ",") + ")").Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//GetTagParent : 根据子ID来找出其父ID
func GetTagParent(tagID uint64) (result Tag, err error) {
	t := Tag{}
	t.ID = tagID
	child, err := GetTag(tagID, false)
	if err != nil {
		return result, errors.New("cannot find any tag with id " + util.Uint2Str(tagID) + " error:" + err.Error())
	}
	if err = DB.Self.Where("id = ?", child.Parent).First(&result).Error; err != nil {
		return result, errors.New("cannot find parent tag with child id " + util.Uint2Str(tagID) + " error:" + err.Error())
	}
	return result, err
}

//GetTag :
func GetTag(tagID uint64, withUsers bool) (result *Tag, err error) {
	tag := &Tag{}
	if tagID == 0 {
		return result, errors.New("id do not existed")
	}

	if err = DB.Self.Select("id,name,coefficient,parent").First(&tag, tagID).Error; err != nil {
		return result, err
	}

	if withUsers {
		if err = DB.Self.Model(&tag).Preload("Profile").Association("Users").Find(&tag.Users).Error; err != nil {
			return result, err
		}
	}
	return tag, nil

}

//GetTagWithProfile :
func GetTagWithProfile(tagID uint64, withProfiles bool) (result *Tag, err error) {
	tag := &Tag{}
	if tagID == 0 {
		return result, errors.New("id do not existed")
	}

	if err = DB.Self.Select("id,name,coefficient").First(&tag, tagID).Error; err != nil {
		return result, err
	}

	if withProfiles {
		if err = DB.Self.Model(&tag).Association("Profiles").Find(&tag.Profiles).Error; err != nil {
			return result, err
		}
	}
	return tag, nil

}

// GetTagByName :
// @params name
// name 支持两种格式
// 第一种， 直接输入群组名字，如办事员
// 第二种， 体现父子关系，如 岗位.办事员
func GetTagByName(name string) (result *Tag, err error) {
	hasLevel := strings.Index(name, ".")
	if hasLevel > 0 {
		tags := strings.Split(name, ".")

		var parentID uint64 = 0
		for i, tag := range tags {
			temp := &Tag{}
			err = DB.Self.Select("id,name,coefficient,parent").Where("name = ? and parent = ?", util.Strip(tag), parentID).First(&temp).Error
			if i == (len(tags) - 1) {
				result = temp
			} else {
				parentID = temp.ID
			}
		}
	} else {
		temp := &Tag{}
		err = DB.Self.Select("id,name,coefficient,parent").Where("name = ?", name).First(&temp).Error
		if err == nil {
			result = temp
		}
	}

	return result, err
}

// GetTagByName :
func GetSubTag(id uint64) (result []*Tag, err error) {
	if err = DB.Self.Select("id,name,parent,coefficient").Where("parent = ?", id).Find(&result).Error; err != nil {
		return result, err
	}

	return result, nil
}

// DeleteTag :
// @params : ids , tag id list
func DeleteTags(ids []uint64) (err error) {
	tx := DB.Self.Begin()
	if err = tx.Exec("DELETE FROM tb_tags WHERE id IN (?) ", ids).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Exec("DELETE FROM user_tags  WHERE tag_id IN (?) ", ids).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// DeleteTag :
// @params : id , tag id list
func DeleteTag(id uint64) (err error) {
	tx := DB.Self.Begin()

	tag := &Tag{}
	if err := tx.Where(" id = ? or parent = ? ", id, id).Delete(&tag).Error; err != nil {
		tx.Rollback()
		return errors.New("无法删除")
	}

	tx.Commit()
	return nil
}

func ImportTagFromExcel(filepath string) (errs []string, err error) {
	//分析第一行
	xlsx, err := excelize.OpenFile(filepath)
	if err != nil {
		fmt.Println("OpenFile", err)
		return nil, err
	}

	rows := xlsx.GetRows("Sheet1")
	tagMap := make(map[string]map[string][]string, 0) //记录标签所含成员
	coeMap := make(map[string]map[string]float64, 0)  //记录系数

	for _, row := range rows[1:] {
		if len(row[2]) > 0 {
			tagArr := strings.Split(row[0], ".")
			if _, ok := tagMap[tagArr[0]]; !ok {
				tagMap[tagArr[0]] = make(map[string][]string)
				tagMap[tagArr[0]][tagArr[1]] = make([]string, 0)

				coeMap[tagArr[0]] = make(map[string]float64)
				coeMap[tagArr[0]][tagArr[1]] = 0.00

			}
			tagMap[tagArr[0]][tagArr[1]] = append(tagMap[tagArr[0]][tagArr[1]], row[3])

			value, err := strconv.ParseFloat(row[1], 64)
			if err != nil {
				coeMap[tagArr[0]][tagArr[1]] = 0.00
			}
			coeMap[tagArr[0]][tagArr[1]] = value

		}
	}
	return importProfileIDToTag(tagMap, coeMap), err
}

func importProfileIDToTag(tagMap map[string]map[string][]string, coeMap map[string]map[string]float64) (errArr []string) {
	for k, v := range tagMap {
		parent := &Tag{}
		if err := DB.Self.Where("name = ?", k).First(&parent).Error; err != nil {
			errArr = append(errArr, "系数名:"+k+"不存在，系统将为您创建")
			parent.Name = k
			parent.Coefficient = 0.00
			if err := parent.Create(); err != nil {
				errArr = append(errArr, "系数名:"+k+"创建失败，失败原因："+err.Error())
				continue
			}
		}
		for k1, cards := range v {
			t := &Tag{}
			if err := DB.Self.Where("name = ? AND parent = ?", k1, parent.ID).First(&t).Error; err != nil {
				errArr = append(errArr, "系数名:"+k+"."+k1+"不存在，系统将为您创建")
				t.Name = k1
				t.Coefficient = coeMap[k][k1]
				t.Parent = parent.ID
				if err := t.Create(); err != nil {
					errArr = append(errArr, "分类名:"+k+"."+k1+"创建失败，失败原因："+err.Error())
					continue
				}
			}
			for _, c := range cards {
				p, err := GetProfileByIDCard(c)
				if err != nil {
					errArr = append(errArr, "身份证号码 :"+c+"不存在")
					continue
				}
				AddTagProfiles(t.ID, []uint64{p.ID})
			}

		}
	}
	return errArr
}
