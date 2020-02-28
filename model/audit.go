package model

import (
	"errors"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
	"github.com/lib/pq"
	"hr-server/pkg/constvar"
	"hr-server/util"
	"os"
	"strconv"
	"strings"
)

const (
	AUDITCREATEACTION = 1
	AUDITUPDATEACTION = 2
	AUDITDELETEACTION = 3
	AUDITQUERYACTION  = 4
	AUDITMOVEACTION   = 5
)

const (
	AuditStateWaiting = 0
	AuditStatePermit  = 1
	AuditStateDeny    = 2
)

type Audit struct {
	BaseModel
	Operator     User `json:"operator" validate:"-" gorm:"foreignkey:OperatorID"`
	OperatorID   uint64
	Auditor      User `json:"auditor" validate:"-" gorm:"foreignkey:AuditorID"`
	AuditorID    uint64
	Object       string        `json:"object"` //审核对象
	Action       int           `json:"action"` //审核操作,包括增删查改
	Fields       string        `json:"fields" gorm:"fields"`
	OrgObjectID  pq.Int64Array `json:"org_object" gorm:"type:integer[]"`  //审核前数据项的ID
	DestObjectID pq.Int64Array `json:"dest_object" gorm:"type:integer[]"` //审核后数据项的ID
	State        int           `json:"state"`
	Reply        string        `json:"reply"`
	Body         string        `json:"body"`
	Remark       string        `json:"remark"` //备注
}

const AuditTableName = "tb_audit"

// TableName :
func (a *Audit) TableName() string {
	return AuditTableName
}

// Create :
func (a *Audit) Create() error {
	return DB.Self.Create(&a).Error
}

// Update
func (a *Audit) Update() (err error) {
	state := a.State
	tx := DB.Self.Begin()
	if err := tx.Model(&a).Update(map[string]interface{}{"reply": a.Reply, "state": a.State, "auditor": a.AuditorID}).Error; err != nil {
		tx.Rollback()
		log.Info("Update audit state error", lager.Data{"body": err.Error()})
		return errors.New("无法更新")
	}

	//审核更新会比较麻烦，需要更新自身模块的状态之外 ，还要更新被审核模板的状态。
	if err = DB.Self.First(&a).Error; err != nil {
		tx.Rollback()
		return errors.New("无法查询到该项审核")
	}
	table := TableNames[a.Object]
	//sqlStr := "UPDATE " + table + " SET audit_state = " + strconv.Itoa(state) + " WHERE id = " + util.Uint2Str(a.OrgObjectID)
	ids := []string{}
	for _, id := range a.OrgObjectID {
		s := strconv.FormatInt(id, 10)
		ids = append(ids, s)
	}
	sqlStr := "UPDATE " + table + " SET audit_state = " + strconv.Itoa(state) + " WHERE id in (" + strings.Join(ids, ",") + ")"

	if err := tx.Exec(sqlStr).Error; err != nil {
		tx.Rollback()
		return errors.New("无法更新审核状态")
	}

	tx.Commit()

	//如果是profile且是action是创建的话，这时需将用户添加进“部门”组里，不然在系统里是不可见的。
	switch a.Object {
	case "Profile":
		handleProfile(a, state)
	case "Template":
		handleTemplate(a, state)
	}
	return err
}

func CountAudit() (count int, err error) {
	err = DB.Self.Model(&Audit{}).Where("state  = ? ", 0).Count(&count).Error
	return count, err
}

func ListAudit(state int, offset, limit int) (audits []Audit, total uint64, err error) {
	if limit == 0 {
		limit = constvar.DefaultLimit
	}

	if state == -1 {
		if err := DB.Self.Model(&Audit{}).Count(&total).Error; err != nil {
			return audits, 0, err
		}

		if err := DB.Self.Offset(offset).Limit(limit).Order("id desc").Preload("Operator").Find(&audits).Error; err != nil {
			fmt.Println("err", err)
			return audits, 0, err
		}
	} else {
		if err := DB.Self.Model(&Audit{}).Where("state  = ? ", state).Count(&total).Error; err != nil {
			return audits, 0, err
		}

		if err := DB.Self.Where("state  =? ", state).Offset(offset).Limit(limit).Order("id desc").Preload("Operator").Find(&audits).Error; err != nil {
			fmt.Println("err", err)
			return audits, 0, err
		}
	}

	return audits, total, err
}

func handleProfile(a *Audit, state int) error {
	if a.Action == AUDITCREATEACTION {
		if state == AuditStatePermit { //允许创建的话
			if err := AddProfileToDefaultGroup(uint64(a.OrgObjectID[0])); err != nil {
				return errors.New("新增职工档案审核通过，但无法将职工档案添加到部门里，错误信息:" + err.Error())
			}

			profile, err := GetProfile(uint64(a.OrgObjectID[0]))
			if err != nil {
				return errors.New("新增职工档案审核通过，但无法将新增职工信息添加到记录表里，错误信息:" + err.Error())
			}
			record := Record{}
			record.Body = "描述:新增职工信息; 姓名：" + profile.Name + ";身份证号码:" + profile.IDCard + ";"
			if err = record.Create(); err != nil {
				return errors.New("新增职工档案审核通过，但无法将新增职工信息添加到记录表里，错误信息:" + err.Error())
			}
		}
		return nil
	}
	if a.Action == AUDITUPDATEACTION {
		if state == AuditStatePermit { //允许创建的话
			record := Record{}
			record.Body = a.Body
			if err := record.Create(); err != nil {
				return errors.New("更新职工档案审核通过，但无法将职工更新信息添加到记录表里，错误信息:" + err.Error())
			}
		}
		return nil
	}
	if a.Action == AUDITDELETEACTION {
		if state == AuditStatePermit { //允许删除
			if err := DeleteProfile(uint64(a.OrgObjectID[0])); err != nil {
				return errors.New("审核通过，但可能发生一些错误，无法删除.错误信息:" + err.Error())
			}

			if err := ClearProfileForGroups(uint64(a.OrgObjectID[0])); err != nil {
				return errors.New("删除职工档案审核通过，但无法将职工档案从相应的组里清除，错误信息:" + err.Error())
			}

			profile, err := GetProfile(uint64(a.OrgObjectID[0]))
			if err != nil {
				return errors.New("新增职工档案审核通过，但无法将新增职工信息添加到记录表里，错误信息:" + err.Error())
			}

			record := Record{}
			record.Body = "描述:删除职工信息; 姓名：" + profile.Name + ";身份证号码:" + profile.IDCard + ";"
			if err = record.Create(); err != nil {
				return errors.New("删除职工档案审核通过，但无法将新增职工信息添加到记录表里，错误信息:" + err.Error())
			}

		} else {
			p := &Profile{}
			p.ID = uint64(a.OrgObjectID[0])
			p.UpdateState(AuditStateDeny)
		}
		return nil
	}
	if a.Action == AUDITMOVEACTION {
		fmt.Println("move action")
		// 员工移动之后
		// 第一步，查询调动表得到调动信息
		m, err := GetTransfer(uint64(a.DestObjectID[0]))
		if err != nil {
			return errors.New("员工调动审核通过,但数据库发生错误:" + err.Error())
		}
		if state == AuditStatePermit {
			profile, err := GetProfile(m.Profile)
			if err != nil {
				return errors.New("新增职工档案审核通过，但无法将新增职工信息添加到记录表里，错误信息:" + err.Error())
			}

			// 获取旧组信息
			oldGroup, err := GetGroup(m.OldGroup, false)
			if err != nil {
				return errors.New("员工调动审核通过,但数据库发生错误:" + err.Error())
			}
			// 获取新组信息
			newGroup, err := GetGroup(m.NewGroup, false)
			fmt.Println("new group ", util.PrettyJson(newGroup))
			if err != nil {
				return errors.New("员工调动审核通过,但数据库发生错误:" + err.Error())
			}

			//获取所有标签的信息, 主要是为了建立标签ID 与 Name 的映射表
			topTags, _, err := GetAllTags(0, 10000, "parent", "0")
			topTagMap := make(map[uint64]string)

			parentTagMap := make(map[uint64]Tag) // 记录所有父标签

			for _, tag := range topTags {
				topTagMap[tag.ID] = tag.Name

				if tag.Parent == 0 {
					parentTagMap[tag.ID] = *tag
				}

			}

			rules, _ := getRulesFromCSV(newGroup)

			//第一步： 先查找出旧的组所关联的标签，然后根据这些标签，删除掉用户与这些标签的关联.
			// 根据当时
			deleteTagRecord := ""
			deleteIDList := []uint64{}
			if len(oldGroup.Tags) > 0 {
				for _, tag := range oldGroup.Tags {
					deleteIDList = append(deleteIDList, tag.ID)
				}
			}
			//  第二步： 再找出用户当前所关联的标签，如果这些标签与将要新添加的标签有重复的，以新标签为准。
			//  比如当你手工为用户添加了 车改补贴.800 ，此时将用户调动到副经理时，标签应该为车改补贴.400 ,那就应该以车改补贴.400 为准
			profileTags, _ := GetProfileWithTags(profile.ID)

			for _, ptag := range profileTags.Tags {
				associatedGroupIds := make([]int64, 0)

				if ptag.Parent != 0 {
					associatedGroupIds = parentTagMap[ptag.Parent].CommensalismGroupIds
				}

				//标签共生组的情况下，如果目标组有该标签，则为用户加上该标签，如果没有该标签，则将该标签删除，类似共生体
				for _, gid := range associatedGroupIds {
					if uint64(gid) == newGroup.Parent { //如果标签有关联到这个组
						hasTag := false
						//看看这个组是不是有关联该标签
						for _, gtag := range newGroup.Tags {
							if ptag.Parent == gtag.Parent {
								hasTag = true
							}
						}

						if !hasTag {
							deleteIDList = append(deleteIDList, ptag.ID)
						}
					}
				}

				for _, gtag := range newGroup.Tags {
					if ptag.Parent == gtag.Parent { // 同属于车改补贴这个大标签类，
						deleteIDList = append(deleteIDList, ptag.ID)
					}
				}
			}
			//查询调动表，看看当时调动到这个群组时为profile所添加的标签都有哪些，将这些标签都移除.
			oldTransfer, err := GetTransferByNewGroupCombination(oldGroup.ID, profile.ID)
			if err != nil {
				fmt.Println("err", err)
			}

			for _, id := range oldTransfer.AddedTagsRecord {
				deleteIDList = append(deleteIDList, uint64(id))
			}

			tags, err := GetTagsByIDList(deleteIDList)
			for _, tag := range tags {
				deleteTagRecord += "数值名:" + topTagMap[tag.Parent] + "; 数值:" + fmt.Sprint(tag.Coefficient) + ";"
			}

			// 新增标签这里会有个特殊情况，比如如下场景
			// 组与标签会关联一个系数 ，但有一个特殊情况，比如说A员工属于总行营业部，岗位任高级会计主管。但却是派驻白塔支行。
			// 根据单位的车补规则，总行营业部应该是200， 白塔支行应该是300。此时A员工虽然身份属于总行营业部高级会计主管，但实际领的是白塔支行的车补400。
			// 所以这个时候不能再根据岗位来关联车补的标签。而是需要根据多个群组来判断关联哪个标签。
			// 这个规则表会保存在服务器上的配置文件 conf/group_tag_rule.csv 中。
			// 格式如下：
			// 岗位.高级会计主管, 部门.营业部,车补.400
			// 第一项表示 当前的group ,比如说 职称.员级
			// 第二项表示 关联的group, 比如说 部门.桂岭支行
			// 第三项表示 当员工调动之后 ，如果员工同时属于第一项和第二项的话，就为员工添加第三项指定的标签系数
			//  这个功能在前端和 handler已有完成代码，就差这一步把配置文件读出来并分析，以后再实现。

			addTagRecord := ""
			addTagErrRecord := ""
			idList := []uint64{}
			if len(newGroup.Tags) > 0 {
				for _, tag := range newGroup.Tags {
					idList = append(idList, tag.ID)
					addTagRecord += "数值名:" + topTagMap[tag.Parent] + "; 数值:" + fmt.Sprint(tag.Coefficient) + ";"
				}
			}

			// 如果存在规则
			// todo 这段应可以优化，时间不够。
			if len(rules) > 0 {
				for _, rule := range rules {
					s := strings.Split(rule, ",")
					belongGroupOne, err := GetGroupByName(s[0]) //要查看用户此时是否归属于这个群
					belongGroupTwo, err := GetGroupByName(s[1]) //要查看用户此时是否归属于这个群
					if err != nil {
						addTagErrRecord += "发生错误:" + s[1] + "出错，错误信息" + err.Error() + ";"
						continue
					}
					//判断是否同时属于
					isBelongOne := false
					isBelongTwo := false
					belongOneID := int64(0)
					belongTwoID := int64(0)
					for _, g := range profile.Groups {
						if g.ID == belongGroupOne.ID {
							isBelongOne = true
							belongOneID = int64(g.ID)
						}
						if g.ID == belongGroupTwo.ID {
							isBelongTwo = true
							belongTwoID = int64(g.ID)
						}
					}
					if isBelongOne && isBelongTwo {
						//用户在同时相应的组里，此时为用户添加规则指定的标签
						tag, err := GetTagByName(s[2])
						if err != nil {
							addTagErrRecord += "发生错误:" + s[2] + "出错，错误信息" + err.Error() + ";"
							continue
						}

						for _, ptag := range profileTags.Tags {
							fmt.Println("ptag", util.PrettyJson(ptag), util.PrettyJson(tag), ptag.Parent == tag.Parent)
							if ptag.Parent == tag.Parent { // 同属于车改补贴这个大标签类，
								deleteIDList = append(deleteIDList, ptag.ID)
							}
						}

						idList = append(idList, tag.ID)
						m.NewGroupCombination = []int64{belongOneID, belongTwoID} //记录转移时的转移组合，也就是记录下，当转移时，是根据哪些组合来添加标签的.
					}
				}
			}
			err = RemoveProfileTags(m.Profile, deleteIDList)
			if err != nil {
				fmt.Println(" err.Error()", err.Error())
				return errors.New("员工调动审核通过,但数据库发生错误:" + err.Error())
			}

			err = AddProfileTags(m.Profile, idList)
			if err != nil {
				return errors.New("员工调动审核通过,但数据库发生错误:" + err.Error())
			}
			//将这些profile新添加的tag id 记录进 group_transfer表中，以便以后查找
			addedList := []int64{}
			for _, id := range idList {
				addedList = append(addedList, int64(id))
			}
			m.AddedTagsRecord = addedList
			if err := m.Update(); err != nil {
				fmt.Println("transfer update error", err)
			}
			record := Record{}
			record.Body = "profile"
			record.Body = "描述:职工调动; 姓名：" + profile.Name + ";身份证号码:" + profile.IDCard + ";从：" + oldGroup.Name + "; 到 :" + newGroup.Name + ";"
			record.Body += "调动数值变化, 从：" + fmt.Sprint(oldGroup.Coefficient) + ";到:" + fmt.Sprint(newGroup.Coefficient) + ";"
			if len(deleteTagRecord) > 0 {
				record.Body += "数值变化 ，删除了以下数值：" + deleteTagRecord
			}
			if len(addTagRecord) > 0 {
				record.Body += "数值变化 ，新增了以下数值：" + addTagRecord
			}
			if len(addTagErrRecord) > 0 {
				record.Body += addTagErrRecord
			}
			fmt.Println("move action", record.Body)
			if err = record.Create(); err != nil {
				return errors.New("删除职工档案审核通过，但无法将新增职工信息添加到记录表里，错误信息:" + err.Error())
			}

		} else {
			p := &Profile{}
			p.ID = uint64(a.OrgObjectID[0])
			//审核不通过，将用户转移回原先的状态
			MoveProfileToNewGroup(p.ID, m.NewGroup, m.OldGroup)
			p.UpdateState(AuditStateDeny)
		}
	}
	return nil
}

func getRulesFromCSV(group *Group) (rules []string, err error) {
	parent, _ := GetGroup(group.Parent, false)
	name := parent.Name + "." + group.Name
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
	lines, err := util.ReadLines(filename)
	if err != nil {
		fmt.Println("open csv file failed, " + err.Error())
	}
	for _, line := range lines {
		index := strings.Index(line, name)
		if index > -1 {
			rules = append(rules, line)
		}
	}

	return rules, err
}

func handleTemplate(a *Audit, state int) error {
	oldT, _ := GetTemplate(uint64(a.OrgObjectID[0]))
	newT, _ := GetTemplate(uint64(a.DestObjectID[0]))

	if a.Action == AUDITCREATEACTION {
		if state == AuditStatePermit { //允许创建或者更新的话e
			filename := "conf/templates/" + newT.Name + "-" + util.Uint2Str(newT.ID) + ".yaml"
			if util.Exists(filename) {
				err := util.MoveFile(filename, "conf/templates/"+newT.Name+".yaml")
				if err != nil {
					fmt.Println("cannot move file to new directory" + err.Error())
				}
			}
			newT.UpdateState(AuditStatePermit)
			return nil
		} else {
			newT.UpdateState(AuditStateDeny)
			return nil
		}
	}

	if a.Action == AUDITUPDATEACTION {
		if state == AuditStatePermit { //允许创建或者更新的话
			//因为账套关联模板，所以当审核通过时，直接新模板替换到老模板
			tempNewID := newT.ID
			//newT.ID = oldT.ID
			//newT.AuditState = AuditStatePermit
			//fmt.Println("new template")
			//fmt.Println(util.PrettyJson(newT))

			if err := newT.Replace(oldT.ID, newT.ID); err != nil {
				return err
			}
			////把原本新的模板删除掉.
			//if err := DeleteTemplate(tempNewID); err != nil {
			//	return err
			//}
			if !util.Exists("conf/templates/old/") {
				os.MkdirAll("conf/templates/old/", os.ModePerm) //创建文件
			}

			//如果有更新了模板名的话，这里把旧模板先移到old文件夹
			filename := "conf/templates/" + oldT.Name + ".yaml"
			if util.Exists(filename) {
				err := util.MoveFile(filename, "conf/templates/old/"+oldT.Name+".yaml")
				if err != nil {
					fmt.Println("cannot move file to new directory" + err.Error())
				}
			}
			//再对新模板改名
			filename = "conf/templates/" + newT.Name + "-" + util.Uint2Str(tempNewID) + ".yaml"
			if util.Exists(filename) {
				err := util.MoveFile(filename, "conf/templates/"+newT.Name+".yaml")
				if err != nil {
					fmt.Println("cannot move file to new directory" + err.Error())
				}
			}

			record := Record{}
			record.Body = a.Body
			if err := record.Create(); err != nil {
				return errors.New("新增职工档案审核通过，但无法将新增职工信息添加到记录表里，错误信息:" + err.Error())
			}

		} else {
			//恢复旧模板的状态
			oldT.UpdateState(AuditStatePermit)
			newT.UpdateState(AuditStateDeny)

			//再对新模板改名
			filename := "conf/templates/" + newT.Name + "-" + util.Uint2Str(newT.ID) + ".yaml"
			if util.Exists(filename) {
				err := util.MoveFile(filename, "conf/templates/old/"+newT.Name+"-"+util.Uint2Str(newT.ID)+".yaml")
				if err != nil {
					fmt.Println("cannot move file to new directory" + err.Error())
				}
			}
		}
	}

	return nil
}
