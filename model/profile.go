package model

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/lexkong/log"
	"gopkg.in/go-playground/validator.v9"
	"hr-server/pkg/constvar"
	"hr-server/util"
	"os"
	"strconv"
	"strings"
)

//prevent sql injection
var UserMap = map[string]string{
	"id":            "id",
	"email":         "email",
	"name":          "name",
	"job_number":    "job_number",
	"on_board_date": "on_board_date",
	"picture":       "picture",
	"bank_card":     "bank_card",
}

var ProfileI18nMap = map[string]string{
	"姓名":              "name",
	"name":            "姓名",
	"证件类型":            "type_card",
	"type_card":       "证件类型",
	"证件号码":            "id_card",
	"id_card":         "身份证号码",
	"身份证号码":           "id_card",
	"电话号码":            "phone",
	"phone":           "电话号码",
	"银行账号":            "bank_card",
	"bank_card":       "银行账号",
	"生日":              "birth_day",
	"birth_day":       "生日",
	"性别":              "gender",
	"gender":          "性别",
	"在职状态":            "status",
	"status":          "在职状态",
	"毕业院校":            "school",
	"school":          "毕业院校",
	"招聘来源":            "source",
	"source":          "招聘来源",
	"毕业时间":            "graduation_date",
	"graduation_date": "毕业时间",
	"专业":              "specialty",
	"specialty":       "专业",
	"民族":              "nation",
	"nation":          "民族",
	"婚姻状况":            "marital_status",
	"marital_status":  "婚姻状况",
	"入职时间":            "on_board_date",
	"on_board_date":   "入职时间",
}

var ProfileTableName = TableNames["Profile"]

const ProfileAuditObject = "Profile"

type Profile struct {
	BaseModel
	Name            string  `json:"name" gorm:"column:name;not null" binding:"required"`
	JobNumber       string  `json:"job_number" `
	TypeCard        string  `json:"type_card"`
	Phone           string  `json:"phone"`
	IDCard          string  `json:"id_card" gorm:"not null;unique" binding:"required"`
	Gender          string  `json:"gender"`
	BirthDay        string  `json:"birth_day" `
	Source          string  `json:"source" ` //招聘来源
	School          string  `json:"school"`
	GraduationDate  string  `json:"graduation_date"`          //毕业时间
	Specialty       string  `json:"specialty"`                //专业
	LastCompany     string  `json:"last_company"`             //上一家公司
	FirstJobDate    string  `json:"first_job_date"`           //第一分工作时间
	WorkAge         int     `json:"workage" gorm:"default:0"` //工龄
	Nation          string  `json:"nation"`                   //民族
	MaritalStatus   string  `json:"marital_status"`           //婚姻状况
	AccountLocation string  `json:"account_location"`         //户口所在地
	Address         string  `json:"address"`
	BankCard        string  `json:"bank_card"`      //银行卡
	OnBoardDate     string  `json:"on_board_date" ` //入职日期
	Tags            []Tag   `json:"tags" gorm:"many2many:profile_tags;"`
	Groups          []Group `json:"groups" gorm:"many2many:profile_groups;"`
	Freezed         bool    `json:"freezed"`
	AuditState      int     `json:"audit_state" gorm:"audit_state"` //审核结果
}

func SelectUsers(fields []string) (profiles []*Profile, total uint64, err error) {
	var selectStr string
	var keyArr []string
	for _, val := range fields {
		if _, ok := UserMap[val]; ok {
			keyArr = append(keyArr, val)
		}
	}

	selectStr = strings.Join(keyArr, ",")

	if err := DB.Self.Model(&Profile{}).Select(selectStr).Find(&profiles).Error; err != nil {
		return profiles, 0, err
	}

	if err := DB.Self.Model(&Profile{}).Count(&total).Error; err != nil {
		return profiles, 0, err
	}
	return profiles, total, nil

}

// TableName :
func (p *Profile) TableName() string {
	return ProfileTableName
}

// Create :
func (p *Profile) Create() error {
	p.Freezed = false
	return DB.Self.Create(&p).Error
}

// DeleteProfile deletes the user by the user identifier.
func DeleteProfile(id uint64) error {
	profile := Profile{}
	profile.BaseModel.ID = id
	return DB.Self.Delete(&profile).Error
}

// Update updates an profile
func (p *Profile) Update() (err error) {
	return DB.Self.Save(p).Error
}

func (p *Profile) UpdateState(state int) (err error) {
	tx := DB.Self.Begin()
	if err = tx.Model(&p).Update(map[string]interface{}{"audit_state": state}).Error; err != nil {
		tx.Rollback()
		return errors.New("无法更新")
	}

	tx.Commit()

	return err
}

func CountProfile() (count int , err error ){
	err = DB.Self.Model(&Profile{}).Count(&count).Error
	return count , err
}

//GetAllProfileWidthGroupAndTag :
func GetAllProfileWidthGroupAndTag() (profiles []*Profile, err error) {
	//if err := DB.Self.Preload("Groups").Preload("Tags").Where("audit_state = ?", 1).Find(&profiles).Error; err != nil {
	if err := DB.Self.Preload("Groups").Preload("Tags").Find(&profiles).Error; err != nil {
		return profiles, err
	}

	return profiles, nil
}

//GetAllProfileWidthDepartment :
func GetAllProfileWidthGroup(name string) (profiles []*Profile, err error) {
	g := &Group{}
	if err := DB.Self.Where("name = ?", name).First(&g).Error; err != nil {
		return nil, errors.New("cannot find group which name '部门'")
	}

	if err := DB.Self.Preload("Groups", "levels LIKE ? OR id = ?", "%."+util.Uint2Str(g.ID)+".%", g.ID).Find(&profiles).Error; err != nil {
		return profiles, err
	}

	return profiles, nil
}

//GetAllProfile :
func GetAllProfile() (profiles []*Profile, err error) {
	if err := DB.Self.Find(&profiles).Error; err != nil {
		return profiles, err
	}
	return profiles, nil
}

//GetProfileRelatedGroup :
// func GetProfileWidthGroupAndTag(id uint64) (profile Profile, err error) {
// 	p := Profile{}
// 	p.ID = id
// 	if err := DB.Self.Preload("Groups").Preload("Tags").First(&p).Error; err != nil {
// 		return profile, err
// 	}

// 	return p, nil
// }

func GetProfileWithGroupAndTag(ids []uint64) (profiles []Profile, err error) {
	// if err := DB.Self.Preload("Groups").Preload("Tags").Where("id in (?) AND audit_state = ?", ids, 1).Find(&profiles).Error; err != nil {
	if err := DB.Self.Debug().Preload("Groups").Preload("Tags").Where("id in (?)", ids).Find(&profiles).Error; err != nil {
		return profiles, err
	}

	return profiles, nil
}

// ListProfile List all profiles
func ListProfile(key, value string, offset, limit int, freezed bool) ([]*Profile, uint64, error) {
	if limit == 0 {
		limit = constvar.DefaultLimit
	}
	freezeStr := ""

	if freezed {
		freezeStr = "true"
	} else {
		freezeStr = "false"
	}

	profiles := make([]*Profile, 0)
	var count uint64

	where := ""
	if len(key) > 0 {
		// where = fmt.Sprintf(key+" like '%%%s%%' AND audit_state = %s", value, 1)
		where = fmt.Sprintf(key+" like '%%%s%%' and freezed="+freezeStr, value)
	} else {
		where = "freezed=" + freezeStr
	}

	if err := DB.Self.Model(&Profile{}).Where(where).Count(&count).Error; err != nil {
		return profiles, count, err
	}

	if err := DB.Self.Where(where).Offset(offset).Limit(limit).Order("id desc").Find(&profiles).Error; err != nil {
		return profiles, count, err
	}

	return profiles, count, nil
}

// GetUser gets an user by the user identifier.
func GetProfile(id uint64) (Profile, error) {
	p := Profile{}
	d := DB.Self.Where("id = ?", id).Preload("Groups").First(&p)

	return p, d.Error
}

// GetUser gets an user by the user identifier.
func GetProfileWithTags(id uint64) (Profile, error) {
	p := Profile{}
	d := DB.Self.Where("id = ?", id).Preload("Tags").First(&p)

	return p, d.Error
}

// GetUser gets an user by the user identifier.
func GetProfileByIDCard(card string) (Profile, error) {
	p := Profile{}
	d := DB.Self.Where("id_card = ?", card).Preload("Groups").First(&p)

	return p, d.Error
}
func GetProfiles(ids []uint64) (ps []Profile, err error) {
	d := DB.Self.Where("id in (?)", ids).Find(&ps)
	return ps, d.Error
}

func ImportProfileFromExcel(filepath string,operatorId uint64) (file string, err error) {
	//分析第一行
	xlsx, err := excelize.OpenFile(filepath)
	if err != nil {
		fmt.Println("OpenFile", err)
		return "", err
	}

	rows,_ := xlsx.GetRows("Sheet1")
	cols := []string{}

	pidIndex := 0
	groupIndex := make(map[int]string, 0)               //记录群组所在的col
	groupMap := make(map[string]map[string][]string, 0) //记录群组所含成员
	for index, colCell := range rows[0] {
		//如果导入表有指定部门、学历等组时
		if colCell == "部门" || colCell == "学历" || colCell == "岗位" || colCell == "职称" {
			groupMap[colCell] = make(map[string][]string, 0)
			groupIndex[index] = colCell
			continue
		}
		if ProfileI18nMap[colCell] == "id_card" {
			pidIndex = index
		}

		cols = append(cols, ProfileI18nMap[colCell])
	}

	sql := "insert into " + ProfileTableName + "(" + strings.Join(cols, ",") + ",audit_state,freezed) values"
	values := []string{}

	for _, row := range rows[1:] {
		str := []string{}
		for i, cell := range row {
			if _, ok := groupIndex[i]; ok {
				if len(cell) > 0 {
					groupMap[groupIndex[i]][cell] = append(groupMap[groupIndex[i]][cell], row[pidIndex])
				}
				continue
			}
			if len(cell) > 0 {
				str = append(str, "'"+cell+"'")
			}
		}

		if len(str) > 0 {
			values = append(values, "("+strings.Join(str, ",")+",0,false)")
		}
	}

	type insertResult struct {
		Id uint64
		Name string
		IDCard string
	}

	newProfile := insertResult{}
	errs := []string{}
	//一行一行的插入，如果遇到错误，则写入到excel中。
	for i, v := range values {
		s := sql + v + "  returning id, name,id_card"
		execrows, _ := DB.Self.Raw(s).Rows()

		for execrows.Next() {
			if err := execrows.Scan(&newProfile.Id,&newProfile.Name,&newProfile.IDCard); err != nil {
				log.Error("批量插入新档案失败：" , err )
			}
		}

		if err != nil {
			errs = append(errs, "第"+strconv.Itoa(i+1)+"行出现错误,提示：" + err.Error() )
			log.Error("批量插入新档案失败：" , err )
			xlsx.SetCellValue("Sheet1", util.ConvertToNumberingScheme(len(rows[0])+1)+strconv.Itoa(i+2), err.Error())
		}else{

			//// 添加审核条目
			////创建的同时需同时创建审核条目
			audit := &Audit{}
			audit.OperatorID = operatorId
			audit.Object = ProfileAuditObject
			audit.Action = AUDITCREATEACTION
			audit.OrgObjectID = []int64{int64(newProfile.Id)}
			audit.State = AuditStateWaiting
			audit.Body = "描述:创建职工档案;" +
				"档案ID:" + util.Uint2Str(newProfile.Id) + "; " +
				"员工姓名:" + newProfile.Name + ";" +
				"身份证号码:" + newProfile.IDCard + ";"

			if err := audit.Create(); err != nil {
				fmt.Println(err)
			}
		}
	}
	// Save xlsx file by the given path.
	exportPath := "/export/"
	newFile := "importResult.xlsx"

	if !util.Exists("export/") {
		os.MkdirAll("export/",os.ModePerm) //创建文件
	}

	err = xlsx.SaveAs("." + exportPath + newFile)
	if err != nil {
		fmt.Println(err)
	}
	importProfileIDToGroup(groupMap)

	if len(errs) > 0 {
		return newFile, errors.New(strings.Join(errs, ";"))
	}

	return newFile, nil
}

func importProfileIDToGroup(groupMap map[string]map[string][]string) {
	// 根据身份证查询ID
	var pids []uint64

	for k, v := range groupMap {
		for k1, cards := range v {
			parent := &Group{}
			g := &Group{}
			DB.Self.Where("name = ?", k).First(&parent)
			DB.Self.Where("name = ? AND parent = ?", k1, parent.ID).First(&g)
			DB.Self.Table(ProfileTableName).Where("id_card in (?)", cards).Pluck("id", &pids)
			AddGroupProfiles(g.ID, pids)
		}
	}
}

// 利用postgresql \copy的功能来进行批量导入。
// \copy tb_profile(name,type_card,id_card,bank_card,status) from 'E:\gopath\src\hr-server\upload\人员表.csv' with csv;
// gorm 似乎无法执行 \copy 这个命令，这个函数暂时放弃。
// func ImportProfileFromCSV(fields, filename string) (err error) {
// 	p := &Profile{}
// 	sql := ` \\copy ` + p.TableName() + `(` + fields + `) from ` + filename + ` with csv`

// 	fmt.Println(sql)
// 	tx := DB.Self.Begin()
// 	err = tx.Exec(sql).Error
// 	fmt.Println(err)
// 	tx.Commit()

// 	return err
// }

func FreezeProfile(pids []uint64) error {
	tx := DB.Self.Begin()
	if err := tx.Model(&Profile{}).Where("id in (?)", pids).Update(map[string]interface{}{"freezed": true, "audit_state": AuditStateWaiting}).Error; err != nil {
		tx.Rollback()
		return errors.New("无法冻结")
	}
	tx.Commit()
	return nil
}

func UnFreezeProfile(pids []uint64) error {
	tx := DB.Self.Begin()
	if err := tx.Model(&Profile{}).Where("id in (?)", pids).Update(map[string]interface{}{"freezed": false, "audit_state": AuditStateWaiting}).Error; err != nil {
		tx.Rollback()
		return errors.New("无法激活")
	}
	tx.Commit()
	return nil
}

// Validate the fields.
func (p *Profile) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
