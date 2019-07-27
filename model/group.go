package model

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"hr-server/pkg/constvar"
	"hr-server/util"
	"os"
	"strconv"
	"strings"
)

// Group :
type Group struct {
	BaseModel
	Name             string            `json:"name" gorm:"column:name;not null"`
	FullName         string            `json:"full_name" gorm:"-"`
	Users            []User            `json:"users" gorm:"many2many:user_groups;"`
	Profiles         []Profile         `json:"profiles" gorm:"many2many:profile_groups;"`
	Tags             []Tag             `json:"tags" gorm:"many2many:tag_groups;"`
	TemplateAccounts []TemplateAccount `gorm:"many2many:templateaccount_groups;"`
	Code             int               `json:"code" gorm:"default:0"` //群组编码
	Parent           uint64            `json:"parent" gorm:"column:parent;"`
	Levels           string            `json:"levels" gorm:"column:levels"`  //保存父子层级关系图,例如 pppid.ppid.pid.id
	Coefficient      float64           `json:"coefficient" gorm:"default:0"` //每个群组有一个系数，可用于计算
	Locked           bool              `json:"locked" `                      //锁定，不可删除，不可修改
	Invalid          bool              `json:"invalid"`                      //作废
	IsDefault        bool              `json:"default"`                      //默认组，不可删除，不可修改，新增员工默认加入该组
}

const GroupTableName = "tb_groups"

// TableName :
func (g *Group) TableName() string {
	return GroupTableName
}

// Create : Create a new Group
func (g *Group) Create() error {
	pm := &Group{}
	g.Levels = "0."
	if g.Parent != 0 {
		pm.BaseModel.ID = g.Parent
		if err := DB.Self.First(&pm).Error; err != nil {
			return errors.New("找不到父目录")
		}
		g.Levels = pm.Levels + util.Uint2Str(g.Parent) + "."
	}

	err := DB.Self.Create(&g).Error
	return err
}

// Update updates an group information.
// only update name and coefficient
func (g *Group) Update() error {
	group, err := GetGroup(g.ID, false)
	if err != nil {
		return err
	}

	if group.Locked {
		return errors.New("group is locked,cannot update ")
	}
	if group.IsDefault {
		return errors.New("group is default,cannot update ")
	}
	tx := DB.Self.Begin()
	if err := tx.Model(&g).Update(map[string]interface{}{"name": g.Name, "coefficient": g.Coefficient, "code": g.Code, "parent": g.Parent}).Error; err != nil {
		tx.Rollback()
		return errors.New("无法更新")
	}
	tx.Commit()
	return nil
}

//MoveGroup :
// move every child group to new parent group when group has moved.
func MoveGroup(id, parent uint64) (err error) {
	g := &Group{}
	g.ID = id
	tx := DB.Self.Begin()
	p := &Group{}
	level := "0."
	if parent != 0 {
		p, err = GetGroup(parent, false)
		if err != nil {
			return errors.New("cannot find parent group")
		}
		level = fmt.Sprintf("%s%d.", p.Levels, parent)
	}

	g, err = GetGroup(id, false)
	if err != nil {
		return errors.New("cannot find target group by group id " + string(id))
	}
	oldLevel := g.Levels + util.Uint2Str(g.ID) + "."
	newLevel := level + util.Uint2Str(g.ID) + "."
	if err := tx.Model(&g).Update(map[string]interface{}{"parent": parent, "levels": level}).Error; err != nil {
		tx.Rollback()
		return errors.New("cannot move group to new group ")
	}
	//change
	// postgresql sql syntax :
	// update  tb_groups set levels  = replace(levels,'0.5.','0.6.') where levels like '%0.5.%'
	sqlStr := "UPDATE  tb_groups SET levels = replace(levels,?,?) WHERE levels Like ?"
	if err := tx.Exec(sqlStr, oldLevel, newLevel, "%"+oldLevel+"%").Error; err != nil {
		tx.Rollback()
		return errors.New("无法删除")
	}
	tx.Commit()
	return nil
}

//GetAllGroup :
func ListGroup(offset, limit int, where string, whereKeyword string) (gs []*Group, total int, err error) {
	g := &Group{}
	if limit == 0 {
		limit = 10000
	}
	fieldsStr := "id,name,code,coefficient,parent,levels,locked,invalid,is_default"
	if len(where) > 0 { // parent == 0 的时候表示获取所有用户组，0为最顶层
		//var w string
		//if where == "parent"{
		//	w = "levels LIKE '%." + whereKeyword + ".%' OR id = " + whereKeyword
		//}else{
		//	w = where + " = " + whereKeyword
		//}
		//fmt.Println(w)
		if err := DB.Self.Select(fieldsStr).Where(where+" = ?", whereKeyword).Order("code").Offset(offset).Limit(limit).Find(&gs).Error; err != nil {
			//if err := DB.Self.Select(fieldsStr).Where(w).Offset(offset).Limit(limit).Find(&gs).Error; err != nil {
			return gs, 0, errors.New("cannot get group list by where " + where + " and keyword " + whereKeyword)
		}

		if err := DB.Self.Model(g).Where(where+" = ?", whereKeyword).Count(&total).Error; err != nil {
			return gs, 0, errors.New("cannot fetch count of the row")
		}
	} else {
		if err := DB.Self.Select(fieldsStr).Offset(offset).Limit(limit).Find(&gs).Error; err != nil {
			return gs, 0, errors.New("cannot get group list ")
		}
		if err := DB.Self.Model(g).Count(&total).Error; err != nil {
			return gs, 0, errors.New("cannot fetch count of the row")
		}
	}

	return gs, total, nil

}

//GetGroupRelatedUsers :
func GetGroupRelatedUsers(id uint64, offset, limit int) (users []*User, total int, err error) {
	if limit == 0 {
		limit = constvar.DefaultLimit
	}

	g := &Group{}
	g.ID = id

	if err := DB.Self.Model(&g).Offset(offset).Limit(limit).Preload("Profile").Association("Users").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	total = DB.Self.Model(g).Association("Users").Count()

	return users, total, nil
}

//GetGroupRelatedProfiles :
//第一步，先根据组ID，获取其下属所有组的ID
//第二步，根据这里组ID，获取相关联Profile ID 列表
//第三步，根据profile id 列表获取相关profile
func GetGroupRelatedProfiles(id uint64, offset, limit int, freezed bool) (profiles []Profile, total int, err error) {
	if limit == 0 {
		limit = constvar.DefaultLimit
	}
	gs := []Group{}
	if err := DB.Self.Where("levels LIKE ? OR id = ?", "%."+util.Uint2Str(id)+".%", id).Find(&gs).Error; err != nil {
		return nil, 0, err
	}

	gids := make([]string, len(gs))
	for i, g := range gs {
		gids[i] = util.Uint2Str(g.ID)
	}
	profileIDs := []uint64{}

	selectSql := ""
	countSql := ""
	freezeStr := ""

	if freezed {
		freezeStr = "true"
	} else {
		freezeStr = "false"
	}

	if id == 0 {
		selectSql = "SELECT profile_id from profile_groups offset " + strconv.Itoa(offset) + " limit " + strconv.Itoa(limit)
		countSql = "select count(id) from tb_profile  where id in ( select profile_id from profile_groups where freezed=" + freezeStr + ")"
	} else {
		selectSql = "SELECT profile_id from profile_groups where group_id in (" + strings.Join(gids, ",") + ")" + " offset " + strconv.Itoa(offset) + " limit " + strconv.Itoa(limit)
		//countSql = "SELECT  count(profile_id) from profile_groups where group_id in (" + strings.Join(gids, ",") + ")"
		countSql = "select count(id) from tb_profile  where id in ( select profile_id from profile_groups where group_id in (" + strings.Join(gids, ",") + ")) and freezed=" + freezeStr
	}

	rows, _ := DB.Self.Raw(selectSql).Rows() // Note: Ignoring errors for brevity

	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			return nil, 0, err
		}
		profileIDs = append(profileIDs, id)
	}
	if err := DB.Self.Where(" id in (?) AND audit_state = ? AND freezed = ? ", profileIDs, 1, freezed).Find(&profiles).Error; err != nil {
		return profiles, 0, err
	}
	rows, err = DB.Self.Raw(countSql).Rows()
	if err != nil {
		fmt.Println("err", err)
	} else {
		for rows.Next() {
			rows.Scan(&total)
		}
	}

	return profiles, total, nil
}

//GetGroupRelatedAllProfiles : 获取指定群组所关联的所有职员档案
func GetGroupRelatedAllProfiles(id uint64) (profiles []Profile, err error) {
	gs := []Group{}
	if err := DB.Self.Where("levels LIKE ? OR id = ?", "%."+util.Uint2Str(id)+".%", id).Preload("Profiles", "audit_state=1 AND freezed=false").Preload("Profiles.Groups").Find(&gs).Error; err != nil {
		return nil, err
	}
	for _, g := range gs {
		for _, p := range g.Profiles {
			profiles = append(profiles, p)
		}
	}

	return profiles, nil
}

//MoveUserToNewGroup
func MoveUserToNewGroup(userID, oldGroupID, newGroupID uint64) error {
	user := &User{}
	user.ID = userID

	oldGroup := &Group{}
	oldGroup.ID = oldGroupID
	newGroup := &Group{}
	newGroup.ID = newGroupID
	DB.Self.Model(&user).Association("Groups").Delete(oldGroup)
	return DB.Self.Model(&user).Association("Groups").Append(newGroup).Error
	// return DB.Self.Model(&user).Association("Groups").Replace(oldGroup, newGroup).Error
}

func AddProfileToDefaultGroup(pid uint64) (err error) {
	gs := []Group{}
	//找出所有默认的组

	if err := DB.Self.Where("is_default=?", true).Find(&gs).Error; err != nil {
		fmt.Println(err)
		return err
	}
	for _, g := range gs {
		p := []uint64{pid}
		err = AddGroupProfiles(g.ID, p)
	}

	return err
}

//MoveProfileToNewGroup
func MoveProfileToNewGroup(profileID, oldGroupID, newGroupID uint64) error {
	profile := &Profile{}
	profile.ID = profileID

	oldGroup := &Group{}
	oldGroup.ID = oldGroupID
	newGroup := &Group{}
	newGroup.ID = newGroupID

	if oldGroupID != 0 {
		DB.Self.Model(&profile).Association("Groups").Delete(oldGroup)
	}
	DB.Self.Model(&profile).Association("Groups").Delete(oldGroup)
	return DB.Self.Model(&profile).Association("Groups").Append(newGroup).Error
}

//AddGroupUsers :
func AddGroupProfiles(gid uint64, IDList []uint64) (err error) {
	g := &Group{}
	if g, err = GetGroup(gid, false); err != nil {
		return errors.New("Group is not existed!")
	}

	tx := DB.Self.Begin()
	var profiles []Profile
	for _, id := range IDList {
		profiles = append(profiles, Profile{BaseModel: BaseModel{ID: id}})
	}

	// tx.Model(&g).Association("Profiles").Clear()
	err = tx.Model(&g).Association("Profiles").Append(profiles).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
func ClearProfileForGroups(pid uint64) (err error) {
	tx := DB.Self.Begin()

	err = tx.Model(&Group{}).Exec(" delete from profile_groups where profile_id=" + util.Uint2Str(pid)).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

//RemoveGroupProfiles :
func RemoveGroupProfiles(gid uint64, IDList []uint64) (err error) {
	g := &Group{}
	if g, err = GetGroup(gid, false); err != nil {
		return errors.New("Group is not existed!")
	}

	tx := DB.Self.Begin()

	profileIDs := make([]string, len(IDList))

	for i, id := range IDList {
		profileIDs[i] = util.Uint2Str(id)
	}
	err = tx.Model(&g).Exec(" delete from profile_groups where profile_id in (" + strings.Join(profileIDs, ",") + ") and group_id = " + util.Uint2Str(gid) + " ;").Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// GetGroup :
func GetGroup(id uint64, withUsers bool) (result *Group, err error) {
	g := &Group{}
	if id == 0 {
		return result, errors.New("cannot find group by id " + util.Uint2Str(id))
	}
	err = DB.Self.Model(&g).Preload("Tags").First(&g, id).Error
	if withUsers {
		DB.Self.Model(&result).Select("id").Association("Users").Find(&g.Users)
	}

	gs := []Group{}
	err = DB.Self.Where("id IN (?)", strings.Split(g.Levels[:len(g.Levels)-1], ".")).Find(&gs).Error
	if err != nil {
		fmt.Println("get full name error", err.Error())
	} else {
		fullname := []string{}
		for _, temp := range gs {
			fullname = append(fullname, temp.Name)
		}

		g.FullName = strings.Join(fullname, ".") + "." + g.Name
	}

	return g, err
}

// GetGroupWithProfile :
func GetGroupWithProfile(id uint64, withProfiles bool) (result Group, err error) {
	g := Group{}
	if id == 0 {
		return result, errors.New("cannot find group by id " + util.Uint2Str(id))
	}
	err = DB.Self.Select("id,name,code,coefficient,parent,levels").First(&g, id).Error
	// err = DB.Self.Model(&Group{BaseModel: BaseModel{ID: id}}).First(&result).Error
	if withProfiles {
		DB.Self.Debug().Model(&result).Select("id").Preload("Profiles").Find(&g.Profiles)
	}
	return g, err
}

// GetGroupByName :
// @params name
// name 支持两种格式
// 第一种， 直接输入群组名字，如办事员
// 第二种， 体现父子关系，如 岗位.办事员
func GetGroupByName(name string) (result *Group, err error) {
	hasLevel := strings.Index(name, ".")
	if hasLevel > 0 {
		groups := strings.Split(name, ".")

		var parentID uint64 = 0
		for i, group := range groups {
			gTemp := &Group{}
			err = DB.Self.Select("id,name,code,coefficient,parent,levels").Where("name = ? and parent = ?", util.Strip(group), parentID).First(&gTemp).Error
			if i == (len(groups) - 1) {
				result = gTemp
			} else {
				parentID = gTemp.ID
			}
		}
	} else {
		gTemp := &Group{}
		err = DB.Self.Select("id,name,code,coefficient,parent,levels").Where("name = ?", name).First(&gTemp).Error
		if err == nil {
			result = gTemp
		}
	}

	return result, err
}

//
func IfGroupHaveProfile(group string, profile uint64) (result bool) {
	g, err := GetGroupByName(group)
	if err != nil {
		return false
	}
	sql := `select  * from profile_groups where profile_id=? and group_id=?`
	rows, _ := DB.Self.Raw(sql, profile, g.ID).Rows()

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

// GetGroupAndAllChildren :
func GetGroupWithAllChildren(name string) (result []Group, err error) {
	g := &Group{}

	err = DB.Self.Where("name = ?", name).First(&g).Error
	if err != nil {
		fmt.Println("get group error ", err)
		return nil, err
	}
	err = DB.Self.Where("levels LIKE ? OR id = ?", "%."+util.Uint2Str(g.ID)+".%", g.ID).Find(&result).Error
	if err != nil {
		fmt.Println("get group children err  ", err)
		return nil, err
	}

	return result, err
}

// DeleteGroup : delete children group when parent had deleted
func DeleteGroup(id uint64) error {

	group, err := GetGroup(id, false)
	if err != nil {
		return err
	}

	if group.Locked {
		return errors.New("group is locked,cannot delete ")
	}

	cat := &Group{}
	tx := DB.Self.Begin()
	if err := tx.Where("levels LIKE ? OR id = ?", "%."+util.Uint2Str(id)+".%", id).Delete(&cat).Error; err != nil {
		tx.Rollback()
		return errors.New("无法删除")
	}
	tx.Commit()

	return nil
}

func InvalidGroupOrNot(id uint64, invalid bool) error {
	g := &Group{}
	g.ID = id
	tx := DB.Self.Begin()
	if err := tx.Model(&g).Update(map[string]interface{}{"invalid": invalid}).Error; err != nil {
		tx.Rollback()
		return errors.New("cannot update group invalid status")
	}
	tx.Commit()
	return nil
}

func LockGroupOrNot(id uint64, locked bool) error {
	g := &Group{}
	g.ID = id
	tx := DB.Self.Begin()
	if err := tx.Model(&g).Update(map[string]interface{}{"locked": locked}).Error; err != nil {
		tx.Rollback()
		return errors.New("cannot update group lock status")
	}
	tx.Commit()
	return nil
}

func ImportGroupTagRelationshipFromExcel(file string) (errs []string, err error) {
	//分析第一行
	xlsx, err := excelize.OpenFile(file)
	if err != nil {
		fmt.Println("OpenFile", err)
		return nil, err
	}

	rows := xlsx.GetRows("Sheet1")
	deleteValues := []string{}
	insertValues := []string{}
	rules := []string{}
	for _, row := range rows[1:] {
		g := Group{}
		// 看看能群组直接关联标签，还是有规则存在，规则配置格式如下：
		// 部门.龙尾支行,岗位.总行机关办事员 ， 即中间有 , 相连接
		if strings.Index(row[0], ",") > 0 {
			//存在规则
			rules = append(rules, row[0]+","+row[1])
			continue
		}
		gArr := strings.Split(row[0], ".")
		if len(gArr) == 2 {
			sql := `select id from tb_groups where name=? and parent in (select id from tb_groups where name=?);`
			err = DB.Self.Raw(sql, gArr[1], gArr[0]).Scan(&g).Error
			if err != nil {
				fmt.Println("err", err)
				errs = append(errs, "组"+row[0]+"找不到.")
				continue
			}

		}
		tArr := strings.Split(row[1], ".")
		if len(tArr) == 2 {
			sql := `select id from tb_tags where name=? and parent in (select id from tb_tags where name=?) ;`
			t := Tag{}
			err = DB.Self.Raw(sql, tArr[1], tArr[0]).Scan(&t).Error
			if err != nil {
				errs = append(errs, "系数"+row[0]+"找不到.")
				continue
			}
			insertValues = append(insertValues, "("+util.Uint2Str(g.ID)+","+util.Uint2Str(t.ID)+")")
			deleteValues = append(deleteValues, "(group_id = "+util.Uint2Str(g.ID)+" and tag_id="+util.Uint2Str(t.ID)+")")
		}
	}
	createCSV(rules)
	//先删除原有的关联
	deleteSql := `delete from tag_groups where ` + strings.Join(deleteValues, " or ")
	err = DB.Self.Exec(deleteSql).Error
	if err != nil {
		return errs, err
	}
	insertSql := "insert into tag_groups(group_id,tag_id) values" + strings.Join(insertValues, ",")
	err = DB.Self.Exec(insertSql).Error
	return errs, err
}

//todo : 代码跟 handler/group/relate.go 重复,虽略有修改，再想办法解决冗余。
func createCSV(rules []string) (err error) {

	filename := "conf/group_tag_rule.csv"
	var f *os.File
	if !util.Exists(filename) {
		f, err = os.Create(filename) //创建文件
		if err != nil {
			panic(err)
		}
		defer f.Close()
	}

	// csv的格式类似：职称.员级, 部门.新亨支行,独生子女费.10
	// 当 职称.员级, 部门.新亨支行 存在时，更新 "独生子女费.10"
	// 不存在时则创建
	newLines := []string{}
	lines, err := util.ReadLines(filename)
	uniqueMap := map[string]struct{}{} //去除重复
	for _, line := range lines {
		newLine := util.Strip(line)
		uniqueMap[newLine] = struct{}{}
		newLines = append(newLines, newLine)
	}

	for _, r := range rules {
		newLine := util.Strip(r)
		if _, ok := uniqueMap[newLine]; !ok {
			newLines = append(newLines, newLine)
		}
	}

	util.WriteLines(newLines, filename)

	return err
}

func ImportGroupFromExcel(filepath string) (errs []string, err error) {
	//分析第一行
	xlsx, err := excelize.OpenFile(filepath)
	if err != nil {
		fmt.Println("OpenFile", err)
		return nil, err
	}

	rows := xlsx.GetRows("Sheet1")
	gMap := make(map[string]map[string][]string, 0)  //记录标签所含成员
	coeMap := make(map[string]map[string]float64, 0) //记录系数

	for _, row := range rows[1:] {
		if len(row[2]) > 0 {
			tagArr := strings.Split(row[0], ".")
			if _, ok := gMap[tagArr[0]]; !ok {
				gMap[tagArr[0]] = make(map[string][]string)
				gMap[tagArr[0]][tagArr[1]] = make([]string, 0)

				coeMap[tagArr[0]] = make(map[string]float64)
				coeMap[tagArr[0]][tagArr[1]] = 0.00

			}
			gMap[tagArr[0]][tagArr[1]] = append(gMap[tagArr[0]][tagArr[1]], row[3])
			value, err := strconv.ParseFloat(row[1], 64)
			if err != nil {
				coeMap[tagArr[0]][tagArr[1]] = 0.00
			}
			coeMap[tagArr[0]][tagArr[1]] = value
		}
	}
	return importProfileToGroup(gMap, coeMap), err
}

func importProfileToGroup(gMap map[string]map[string][]string, coeMap map[string]map[string]float64) (errArr []string) {
	for k, v := range gMap {
		parent := &Group{}
		if err := DB.Self.Where("name = ?", k).First(&parent).Error; err != nil {
			errArr = append(errArr, "分类名:"+k+"不存在，系统将为您创建")
			parent.Name = k
			parent.Coefficient = 0.00
			if err := parent.Create(); err != nil {
				errArr = append(errArr, "分类名:"+k+"创建失败，失败原因："+err.Error())
				continue
			}
		}
		for k1, cards := range v {
			t := &Group{}
			if err := DB.Self.Where("name = ? AND parent = ?", k1, parent.ID).First(&t).Error; err != nil {
				errArr = append(errArr, "分类名:"+k+"."+k1+"不存在，系统将为您创建")
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
				AddGroupProfiles(t.ID, []uint64{p.ID})
			}

		}
	}
	return errArr
}

// InitDefault :
func (g Group) InitDefault() error {
	tx := DB.Self.Begin()
	// Create
	g.Name = "normal"
	if err := tx.Create(&g).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
