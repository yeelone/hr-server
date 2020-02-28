package salary

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	h "hr-server/handler"
	"hr-server/model"
	"hr-server/pkg/buildinfunc"
	"hr-server/pkg/errno"
	"hr-server/pkg/formula"
	"hr-server/pkg/template"
	"hr-server/util"
	"strconv"
	"time"
)

var DataMap map[uint64]map[string]model.SalaryField
var wModel *model.Salary
var templateAccount *model.TemplateAccount
var profileGroupMap map[string]map[uint64]uint64 //用户与部门、岗位之间进行映射以供调用,比如[部门][用户ID][业务部] [岗位][用户ID][安全岗]
var department = "部门"
var post = "岗位"
var cardIDMap map[string]map[string]map[string]float64 = make(map[string]map[string]map[string]float64) //从上传的文件中取出数据 ，按模板名- 身份证号码 - 字段key- 值 进行存储
var profileCardMap map[uint64]string                                                                    // 身份证跟ID的映射
var profiles []model.Profile
var errMsg = []string{}                           //记录错误信息
var excelUploadedFields = make(map[string]string) // 记录所有的字段，为后面判断从excel上传的字段是否与系统模板配置的字段一致，不一致则警告
var allFieldMap = make(map[string]bool)
var configIdMap = make(map[string]map[uint64]model.SalaryProfileConfig) //记录所有预配置的记录，用field id 作为key ,field id => profile id =>

func Calculate(c *gin.Context) {
	t1 := time.Now() // get current time

	var r CreateRequest
	var err error
	if err = c.Bind(&r); err != nil {
		fmt.Println(err)
		h.SendResponse(c, errno.ErrBind, err.Error())
		return
	}

	// 获取所有特殊用户配置信息，所谓的特殊用户，即有些特殊情况，无法系统性的用规则来计算，可能会存在一些特殊的情况，所以可以在”特殊人员配置“功能里对这些进行配置
	// 比如有 + - * / 四种操作，对计算后的值进行进一步的计算
	configList, err := model.GetSalaryProfileConfig()
	if err != nil {
		configIdMap = nil
	}

	for _, item := range configList {
		if _, ok := configIdMap[item.TemplateFieldId]; !ok {
			configIdMap[item.TemplateFieldId] = make(map[uint64]model.SalaryProfileConfig)
		}
		configIdMap[item.TemplateFieldId][item.ProfileId] = item
	}

	//查询账套信息
	templateAccount, err = model.GetTemplateAccount(r.TemplateAccountID)

	if err != nil {
		h.SendResponse(c, errno.ErrGetTemplateAccount, err.Error())
		return
	}
	if err := model.ClearSalary(r.Year, r.Month, templateAccount.ID); err != nil {
		fmt.Println("清空月度工资出错，出错信息: " + err.Error())
	}

	profiles = getProfiles(templateAccount.Groups)
	if len(r.InitData) > 0 { // 移到下面
		cardIDMap, err = handleUploadExcel(r.InitData)
	}

	//计算模板需要按照顺序来计算
	templateMap := make(map[uint64]model.Template)

	for _, t := range templateAccount.Templates {
		templateMap[t.ID] = t
	}
	errMsg = []string{}
	for _, i := range templateAccount.Order {
		t := templateMap[uint64(i)]
		if len(t.InitData) > 0 { //固定的导入数据
			cardIDMap2, _ := handleUploadExcel(t.InitData)
			//要对两个map进行合并
			for temp, v := range cardIDMap2 {
				if _, ok := cardIDMap[temp]; !ok {
					cardIDMap[temp] = make(map[string]map[string]float64)
				}
				for card, v1 := range v {
					if _, ok := cardIDMap[temp][card]; !ok {
						cardIDMap[temp][card] = make(map[string]float64)
					}
					for field, value := range v1 {
						cardIDMap[temp][card][field] = value
					}
				}
			}
		}
		for _, fields := range DataMap {
			for _, item := range fields {
				allFieldMap[item.Key] = true
			}
		}
		calculateTemplate(c, profiles, t.Name, r.Year, r.Month)
	}

	for field, templateName := range excelUploadedFields {
		if _, ok := allFieldMap[field]; !ok {
			errMsg = append(errMsg, "工资项目不存在，位于模板:["+templateName+"],字段名:["+field+"]")
		}
	}

	elapsed := time.Since(t1)
	fmt.Println("App elapsed: ", elapsed)

	if len(errMsg) > 0 {
		xlsx = excelize.NewFile()
		xlsx.NewSheet("Sheet1")
		// Set value of a cell.
		duplicate := make(map[string]struct{})
		for i, err := range errMsg {
			if _, ok := duplicate[err]; !ok { //去除重复
				xlsx.SetCellValue("Sheet1", "A"+fmt.Sprint(i+1), err)
				duplicate[err] = struct{}{}
			}
		}

		filename := templateAccount.Name + "-计算错误信息表.xlsx"
		err = xlsx.SaveAs("./export/" + filename)
		if err != nil {
			fmt.Println(err)
			h.SendResponse(c, errno.ErrCreateFile, err.Error())
			return
		}

		h.SendResponse(c, errno.ErrSalaryCalculate, CreateResponse{
			//SalaryData: salaryDataMap.Rows,
			ErrorMessageFile: filename,
		})
		return
	}
	h.SendResponse(c, nil, CreateResponse{
		//SalaryData: salaryDataMap.Rows,
	})
	return
}

func calculateTemplate(c *gin.Context, profiles []model.Profile, name, year, month string) {

	wModel = &model.Salary{}
	wModel.TemplateAccount = templateAccount.ID
	wModel.Template = name
	wModel.Year = year
	wModel.Month = month

	tp, err := template.ResolveTemplate(wModel.Template)

	if err != nil {
		fmt.Println("模板解析失败，请检查模板格式是否正确", err.Error())
		h.SendResponse(c, errno.ErrTemplateInvalid, nil)
		return
	}

	//initial SalaryData struct for saving the temporary data of salary
	salaryDataMap = &SalaryData{}
	DataMap = make(map[uint64]map[string]model.SalaryField)

	//todo : 放在calculateTemplate 里似乎有重复计算，要优化
	//第一步，查找职工基本信息
	handleProfile(wModel.Template, profiles, tp.Base)
	// 第二步，处理上传文件 handle upload
	handleUpload(wModel.Template, profiles, tp.Upload)

	//第三步，固定输入
	err = handleInput(wModel.Template, tp.Input)
	if err != nil {
		log.Error("handleInput error", err)
		h.SendResponse(c, err, nil)
		return
	}

	//第三步，查找职工所在群组系数，比如岗位，学历
	err = handleCoefficient(wModel.Template, profiles, tp.Coefficient)

	if err != nil {
		log.Error("handleCoefficient error", err)
		h.SendResponse(c, err, nil)
		return
	}

	//第四步，处理关于从其它模板导入的数据
	handleRelateOtherTemplate(wModel.Template, tp.Related)

	//第五步，处理需要计算的字段，比如工龄（使用内建功能函数来计算），还有公式运算，此步会计算计算所依赖的先后顺序 。
	handleCalculate(wModel.Template, tp.NeedCalculate)

	//最后，保存进数据库
	saveSalary(tp.Related)
}

func handleUpload(template string, profiles []model.Profile, fields []template.Field) error {
	//t1 := time.Now()
	for _, profile := range profiles {
		salaryField := model.SalaryField{}
		for _, field := range fields {
			salaryField.ProfileID = profile.ID
			salaryField.Key = field.Key
			salaryField.Name = field.Name
			salaryField.Alias = field.Alias
			salaryField.Value = cardIDMap[template][profile.IDCard][field.Key]
			salaryField.DepartmentGroupID = profileGroupMap["department"][profile.ID]
			salaryField.PostGroupID = profileGroupMap["post"][profile.ID]
			salaryField.Year = wModel.Year
			salaryField.Month = wModel.Month
			salaryField.FitIntoYear, salaryField.FitIntoMonth = resolveFitMonth(field, wModel.Year, wModel.Month)
			DataMap[profile.ID][field.Key] = salaryField
		}
	}
	//

	//elapsed := time.Since(t1)
	//fmt.Println("handleUpload", elapsed)
	return nil
}

func calculateValueWithConfig(profileId uint64, field template.Field, value float64) (newValue float64) {
	newValue = value

	if configIdMap != nil {
		if _, ok := configIdMap[field.ID][profileId]; ok {
			switch configIdMap[field.ID][profileId].Operate {
			case "+":
				newValue = value + configIdMap[field.ID][profileId].Value
			case "-":
				newValue = value - configIdMap[field.ID][profileId].Value
			case "*":
				newValue = value * configIdMap[field.ID][profileId].Value
			case "/":
				v := value / configIdMap[field.ID][profileId].Value
				newValue = float64(int64(v + 0.5))
			}
		}
	}

	return newValue
}

//getProfiles 根据客户端指定的组ID获取所有用户档案
func getProfiles(groups []model.Group) (profiles []model.Profile) {
	//t1 := time.Now()

	profileGroupMap = make(map[string]map[uint64]uint64)

	departGroup, _ := model.GetGroupByName(department)
	postGroup, _ := model.GetGroupByName(post)

	profileIDs := []uint64{}
	profileGroupMap["department"] = make(map[uint64]uint64)
	profileGroupMap["post"] = make(map[uint64]uint64)
	for _, g := range groups {
		tempProfiles, err := model.GetGroupRelatedAllProfiles(g.ID)
		if err != nil {
			fmt.Println("model.GetGroupRelatedAllProfiles(g.ID) err", util.PrettyJson(g), err)
		}
		for _, profile := range tempProfiles {
			for _, childg := range profile.Groups {
				if childg.Parent == departGroup.ID {
					profileGroupMap["department"][profile.ID] = childg.ID
				}
				if childg.Parent == postGroup.ID {
					profileGroupMap["post"][profile.ID] = childg.ID
				}
			}

			profileIDs = append(profileIDs, profile.ID)
		}
	}
	profiles, _ = model.GetProfileWithGroupAndTag(profileIDs)
	//elapsed := time.Since(t1)
	//fmt.Println("getProfiles", elapsed)
	return profiles
}

func handleRelateOtherTemplate(template string, fields []template.Field) {
	//t1 := time.Now()
	for _, field := range fields {
		if len(field.RelateMonth) < 1 { //如果没有指定年月，只默认为当前计算的月份
			field.RelateMonth = wModel.Month
			field.RelateYear = wModel.Year
		}
		account := templateAccount.ID

		if len(field.RelateTemplateAccount) > 0 {
			a, err := model.GetTemplateAccountByName(field.RelateTemplateAccount)
			if err != nil {
				errMsg = append(errMsg, "找不到关联的账套:["+field.RelateTemplateAccount+"],可能影响到["+field.Alias+"]的计算")
			} else {
				account = a.ID
			}
		}

		//todo : 这里有个问题，就是关联其它账套时，会把账套里的人全部取出来，而自己账套所关联的员工只有几个。比如说A账套关联全行员工，而B账套只关联了一个问题的员工。这就有点浪费
		results := model.GetRelatedTemplateValue(field.RelateYear, field.RelateMonth, field.RelateTemplate, account, []string{"profile_id", field.RelateKey})
		for _, m := range results {
			if _, ok := DataMap[m.ProfileID]; !ok {
				continue
			}
			salaryField := model.SalaryField{}
			salaryField.Key = field.Key
			salaryField.Name = field.Name
			salaryField.Alias = field.Alias
			salaryField.ProfileID = m.ProfileID
			salaryField.DepartmentGroupID = m.DepartmentGroupID

			salaryField.PostGroupID = m.PostGroupID
			// 已经存在，表示在用户上传的模板数据中有存在这一项数据，优先处理上传的数据
			if _, ok := cardIDMap[template][profileCardMap[m.ProfileID]][field.Key]; ok {
				salaryField.Value = cardIDMap[template][profileCardMap[m.ProfileID]][field.Key]
			} else {
				salaryField.Value = m.Value
			}
			salaryField.Value = calculateValueWithConfig(salaryField.ProfileID, field, salaryField.Value)
			salaryField.Year = wModel.Year
			salaryField.Month = wModel.Month
			salaryField.FitIntoYear, salaryField.FitIntoMonth = resolveFitMonth(field, wModel.Year, wModel.Month)

			DataMap[m.ProfileID][field.Key] = salaryField
		}
	}
	//elapsed := time.Since(t1)
	//fmt.Println("handleRelateOtherTemplate", elapsed)
}

//handleProfile
func handleProfile(template string, profiles []model.Profile, fields []template.Field) error {
	//t1 := time.Now()
	var profileArr []string
	indexMap := make(map[string]int)
	for index, field := range fields {
		if field.From == "profile" && field.Type == "Base" {
			profileArr = append(profileArr, field.Key)
			indexMap[field.Key] = index
		}
	}
	profileCardMap = make(map[uint64]string)

	for _, profile := range profiles {
		profileCardMap[profile.ID] = profile.IDCard
		salaryField := model.SalaryField{}
		pMap := structToMap(profile)
		DataMap[profile.ID] = make(map[string]model.SalaryField)
		for _, item := range profileArr {
			salaryField.Key = item
			salaryField.ProfileID = profile.ID
			salaryField.DepartmentGroupID = profileGroupMap["department"][profile.ID]
			salaryField.PostGroupID = profileGroupMap["post"][profile.ID]
			salaryField.Name = fields[indexMap[item]].Name
			salaryField.Value = -1
			salaryField.Alias = fields[indexMap[item]].Alias
			salaryField.Content = fmt.Sprintf("%v", pMap[item])
			salaryField.Year = wModel.Year
			salaryField.Month = wModel.Month
			DataMap[profile.ID][item] = salaryField
		}
	}

	//elapsed := time.Since(t1)
	//fmt.Println("handleProfile", elapsed)
	return nil
}

//handleInput
func handleInput(template string, fields []template.Field) error {
	//t1 := time.Now()
	for key := range DataMap {
		salaryField := model.SalaryField{}
		for _, field := range fields {
			salaryField.ProfileID = key
			salaryField.Key = field.Key
			salaryField.Name = field.Name
			salaryField.Alias = field.Alias
			if _, ok := cardIDMap[template][profileCardMap[key]][field.Key]; ok {
				salaryField.Value = cardIDMap[template][profileCardMap[key]][field.Key]
			} else {
				salaryField.Value = float64(field.Value.(float64))
			}

			salaryField.Value = calculateValueWithConfig(salaryField.ProfileID, field, salaryField.Value)

			salaryField.DepartmentGroupID = profileGroupMap["department"][key]
			salaryField.PostGroupID = profileGroupMap["post"][key]
			salaryField.Year = wModel.Year
			salaryField.Month = wModel.Month
			salaryField.FitIntoYear, salaryField.FitIntoMonth = resolveFitMonth(field, wModel.Year, wModel.Month)
			DataMap[key][field.Key] = salaryField
		}
	}
	//elapsed := time.Since(t1)
	//fmt.Println("handleInput", elapsed)
	return nil
}

//handleCoefficient
func handleCoefficient(templateName string, profiles []model.Profile, coes []template.Field) error {
	//t1 := time.Now()
	//第三步 ,get group coefficient
	groupMap := make(map[uint64]template.Field)
	tagMap := make(map[uint64]template.Field)
	fieldTagDefaultValueMap := make(map[string]model.Tag)

	//比如想查询 上层分类如岗位，学历，第一步先查出岗位、学历的Group
	for _, field := range coes {
		if field.From == "group" {
			group, err := model.GetGroupByName(field.Name)
			if err != nil {
				return errors.New("cannot find group name " + field.Name)
			}
			field.Value = group.Coefficient
			if group.ID > 0 {
				groupMap[group.ID] = field
			}
		}

		if field.From == "tag" {
			tag, err := model.GetTagByName(field.Name)
			if err != nil {
				return errors.New("cannot find tag by name :" + field.Name)
			}
			field.Value = tag.Coefficient
			if tag.ID > 0 {
				tagMap[tag.ID] = field
				fieldTagDefaultValueMap[field.Key] = *tag
			}

		}

		for _, profile := range profiles { //同样进行遍历所有的用户
			salaryField := model.SalaryField{}
			//默认设置为0
			salaryField.ProfileID = profile.ID
			salaryField.Key = field.Key
			salaryField.Name = field.Name
			salaryField.Alias = field.Alias
			salaryField.Year = wModel.Year
			salaryField.Month = wModel.Month

			if _, ok := cardIDMap[templateName][profileCardMap[profile.ID]][salaryField.Key]; ok {
				salaryField.Value = cardIDMap[templateName][profileCardMap[profile.ID]][salaryField.Key]
			} else {
				if field.From == "tag" {
					salaryField.Value = fieldTagDefaultValueMap[field.Key].Coefficient
				} else {
					salaryField.Value = 0.00
				}
			}

			salaryField.Value = calculateValueWithConfig(salaryField.ProfileID, field, salaryField.Value)

			if _, ok := profileGroupMap["department"][profile.ID]; ok {
				salaryField.DepartmentGroupID = profileGroupMap["department"][profile.ID]
			}
			if _, ok := profileGroupMap["post"][profile.ID]; ok {
				salaryField.PostGroupID = profileGroupMap["post"][profile.ID]
			}

			DataMap[profile.ID][field.Key] = salaryField
		}
	}
	for _, profile := range profiles { //同样进行遍历所有的用户
		salaryField := model.SalaryField{}
		salaryField.ProfileID = profile.ID
		salaryField.Year = wModel.Year
		salaryField.Month = wModel.Month
		salaryField.FitIntoYear = wModel.Year
		salaryField.FitIntoMonth = wModel.Month
		for _, g := range profile.Groups { //取得用户所有的组，进行判断，若其父组为模板中指定的，则将它加入SalaryFields列表中
			if _, ok := groupMap[g.Parent]; ok {
				salaryField.Value = g.Coefficient
				salaryField.Key = groupMap[g.Parent].Key
				salaryField.Name = groupMap[g.Parent].Name
				salaryField.Alias = groupMap[g.Parent].Alias
				salaryField.DepartmentGroupID = profileGroupMap["department"][profile.ID]
				salaryField.PostGroupID = profileGroupMap["post"][profile.ID]

				if _, ok := cardIDMap[templateName][profileCardMap[profile.ID]][salaryField.Key]; ok {
					salaryField.Value = cardIDMap[templateName][profileCardMap[profile.ID]][salaryField.Key]
				}
				salaryField.Value = calculateValueWithConfig(salaryField.ProfileID, groupMap[g.Parent], salaryField.Value)
				DataMap[profile.ID][salaryField.Key] = salaryField
			}

		}

		for _, t := range profile.Tags {
			salaryField.Key = tagMap[t.Parent].Key
			salaryField.Name = tagMap[t.Parent].Name
			salaryField.Alias = tagMap[t.Parent].Alias
			salaryField.DepartmentGroupID = profileGroupMap["department"][profile.ID]
			salaryField.PostGroupID = profileGroupMap["post"][profile.ID]

			//比如该标签: 发放比例， 有下属两个标签: [新员工] = 0.5 , [暂停发放] =0 ，属于新员工和暂停发放的用户则赋予设定的值 ，其它员工则赋予[发放比例]该父标签的默认值
			if _, ok := tagMap[t.Parent]; ok {
				//如果用户属于某个tag，则取于该tag的value
				salaryField.Value = t.Coefficient
			}
			// 如果用户有上传数据 ，最终以上传的为准
			if _, ok := cardIDMap[templateName][profileCardMap[profile.ID]][salaryField.Key]; ok {
				salaryField.Value = cardIDMap[templateName][profileCardMap[profile.ID]][salaryField.Key]
			}

			salaryField.Value = calculateValueWithConfig(salaryField.ProfileID, tagMap[t.Parent], salaryField.Value)
			DataMap[profile.ID][salaryField.Key] = salaryField
		}
	}
	//elapsed := time.Since(t1)
	//fmt.Println("handleCoefficient", elapsed)
	return nil
}

//
func handleCalculate(template string, fields []template.Field) {
	//t1 := time.Now()
	for i := len(fields) - 1; i >= 0; i-- {
		field := fields[i]
		switch field.Type {
		case "Buildin":
			handleBuildin(template, field)
		case "Calculate":
			handleFormula(template, field)
		}
	}
	//elapsed := time.Since(t1)
	//fmt.Println("handleCalculate", elapsed)
}

func handleBuildin(template string, field template.Field) error {
	//t1 := time.Now()
	salaryField := model.SalaryField{}
	salaryField.Key = field.Key
	salaryField.Name = field.Name
	salaryField.Alias = field.Alias

	//用于 IFTransferMonthEqualCurrentMonth 函数中，将组名与ID映射起来
	tempGroupIDMap := make(map[string]uint64)

	build := &buildinfunc.BuildinFunc{}
	for profileID, row := range DataMap {
		params := []interface{}{}
		//遍历依赖的字段，然后从已有的datamap里再遍历查询依赖字段的值，这里有两点需要注意，
		//如果 item.Content 有值 ，说明这个字段是个字符串类型的字段 ，Content不为空时Value为空。

		//IF 函数参数需要支持表达式，所以处理方式不一样。
		valueMap := make(map[string]float64)
		if field.Call == "IF" {
			params = append(params, field.Params[0]) // 判断表达式
			params = append(params, field.Params[1]) // 真
			params = append(params, field.Params[2]) // 假
			for _, r := range field.Require {
				if _, ok := row[r]; !ok {
					continue
				}
				valueMap["["+r+"]"] = row[r].Value
			}
			params = append(params, valueMap)
		} else if field.Call == "IFTransferMonthEqualCurrentMonth" {
			params = append(params, profileID) // 组名
			var groupID uint64 = 0
			if _, ok := tempGroupIDMap[field.Params[0]]; ok {
				groupID = tempGroupIDMap[field.Params[0]]
			} else {
				g, err := model.GetGroupByName(field.Params[0])
				if err != nil {
					errMsg = append(errMsg, "查询不到群组名:["+field.Params[0]+"],可能影响到["+field.Alias+"]的计算")
					continue
				}
				tempGroupIDMap[field.Params[0]] = g.ID
				groupID = g.ID
			}

			params = append(params, groupID)                      // 组名
			params = append(params, field.Params[1])              // 真
			params = append(params, field.Params[2])              // 假
			params = append(params, wModel.Year+"-"+wModel.Month) //年月
			for _, r := range field.Require {
				if _, ok := row[r]; !ok {
					errMsg = append(errMsg, "查询不到依赖项:["+r+"],可能影响到["+field.Alias+"]的计算")
					continue
				}
				valueMap["["+r+"]"] = row[r].Value
			}
			params = append(params, valueMap)
		} else if field.Call == "IFGroupEqual" || field.Call == "IFTagEqual" {
			params = append(params, profileID)       // 组名
			params = append(params, field.Params[0]) // 组名
			params = append(params, field.Params[1]) // 真
			params = append(params, field.Params[2]) // 假
			for _, r := range field.Require {
				if _, ok := row[r]; !ok {
					errMsg = append(errMsg, "查询不到依赖项:["+r+"],可能影响到["+field.Alias+"]的计算")
					continue
				}
				valueMap["["+r+"]"] = row[r].Value
			}
			params = append(params, valueMap)
		} else if field.Call == "Addup" { //为了给调用函数加入profile id的参数
			params = append(params, profileID) // 组名
			params = append(params, wModel.Month)
			params = append(params, field.Params[0]) // 组名
			params = append(params, field.Params[1]) // 真
			params = append(params, field.Params[2]) // 假
		} else if field.Call == "LastMonthValue" { //为了给调用函数加入profile id的参数
			params = append(params, profileID) // 组名
			params = append(params, wModel.Month)
			params = append(params, field.Params[0]) // 组名
		} else {
			for _, r := range field.Require {
				for _, item := range row {
					if item.Key == r {
						if len(item.Content) > 0 {
							params = append(params, item.Content)
						} else {
							params = append(params, item.Value)
						}
					}
				}
			}
			for _, param := range field.Params {
				i, _ := strconv.ParseFloat(param, 64)
				params = append(params, i)
			}
		}

		v, err := build.Call(field.Call, params)

		if err != nil {
			fmt.Println("error", err)
			return err
		}
		if len(v) > -1 {
			salaryField.ProfileID = profileID
			salaryField.DepartmentGroupID = profileGroupMap["department"][profileID]
			salaryField.PostGroupID = profileGroupMap["post"][profileID]

			if _, ok := cardIDMap[template][profileCardMap[profileID]][field.Key]; ok {
				salaryField.Value = cardIDMap[template][profileCardMap[profileID]][field.Key]
			} else {
				salaryField.Value = util.Decimal(float64(v[0].(float64)))
				//salaryField.Value = util.RoundUp(v[0].(float64), 2)
			}
			salaryField.Value = calculateValueWithConfig(salaryField.ProfileID, field, salaryField.Value)
			salaryField.Year = wModel.Year
			salaryField.Month = wModel.Month
			salaryField.FitIntoYear, salaryField.FitIntoMonth = resolveFitMonth(field, wModel.Year, wModel.Month)
			DataMap[profileID][salaryField.Key] = salaryField
		}

	}
	//elapsed := time.Since(t1)
	//fmt.Println("handleBuildin", elapsed)
	return nil
}

func handleFormula(template string, field template.Field) error {
	//t1 := time.Now()
	salaryField := model.SalaryField{}
	salaryField.Key = field.Key
	salaryField.Name = field.Name
	salaryField.Alias = field.Alias
	salaryField.Year = wModel.Year
	salaryField.Month = wModel.Month
	for profileID, row := range DataMap {
		values := make(map[string]float64, 0)
		for _, r := range field.Require {
			if _, ok := row[r]; !ok {
				errMsg = append(errMsg, "查询不到依赖项:["+r+"],可能影响到["+field.Key+"]的计算")
				continue
			}
			values["["+r+"]"] = row[r].Value
		}

		if field.Key == "绩效、职级薪酬小计（高管专项职级薪酬）" && profileID == 675 {
			fmt.Println("field.key 绩效、职级薪酬小计（高管专项职级薪酬）", util.PrettyJson(values))
		}

		salaryField.ProfileID = profileID
		salaryField.DepartmentGroupID = profileGroupMap["department"][profileID]
		salaryField.PostGroupID = profileGroupMap["post"][profileID]

		if _, ok := cardIDMap[template][profileCardMap[profileID]][field.Key]; ok {
			salaryField.Value = cardIDMap[template][profileCardMap[profileID]][field.Key]
		} else {
			salaryField.Value = util.Decimal(formula.Resolve(field.Formula, values))
		}
		if field.MustRounding { //四舍五入
			salaryField.Value = float64(int64(salaryField.Value + 0.5))
		}

		salaryField.Value = calculateValueWithConfig(salaryField.ProfileID, field, salaryField.Value)

		salaryField.FitIntoYear, salaryField.FitIntoMonth = resolveFitMonth(field, wModel.Year, wModel.Month)
		DataMap[profileID][salaryField.Key] = salaryField
	}
	//elapsed := time.Since(t1)
	//fmt.Println("handleFormula", elapsed)
	return nil
}

func handleUploadExcel(filepath string) (data map[string]map[string]map[string]float64, err error) {
	//t1 := time.Now()
	xlsx, err := excelize.OpenFile(filepath)
	if err != nil {
		fmt.Println("handleUploadData open file error ", err)
		return data, err
	}
	//客户传上来的数据可能会有sheet名不小心有空格的情况。所以这里要处理这个问题
	nameMap := make(map[string]string) //去掉空格的和加有空格的进行映射
	data = make(map[string]map[string]map[string]float64)
	for _, name := range xlsx.GetSheetMap() {
		stripName := util.Strip(name)
		nameMap[stripName] = name
		data[stripName] = make(map[string]map[string]float64)
	}

	for sheet := range data {
		rows, _ := xlsx.GetRows(nameMap[sheet]) //nameMap 就是为了这里
		nameRow := []string{}

		//记录身份证号码所在col
		idCardIndex := 0
		//第一步，先取出字段名
		for index, colCell := range rows[0] {
			cellName := util.Strip(colCell)
			if _, ok := model.ProfileI18nMap[colCell]; ok {
				cellName = model.ProfileI18nMap[colCell]
			}
			if cellName == "id_card" {
				idCardIndex = index
			}
			nameRow = append(nameRow, cellName)
			excelUploadedFields[cellName] = sheet
		}

		for _, row := range rows[1:] {
			data[sheet][row[idCardIndex]] = make(map[string]float64)
			for index, cellValue := range row {
				if index != idCardIndex {
					if len(cellValue) > 0 {
						value, _ := strconv.ParseFloat(cellValue, 64)
						data[sheet][row[idCardIndex]][util.Strip(nameRow[index])] = float64(value)
					}
				}
			}
		}
	}

	//elapsed := time.Since(t1)
	//fmt.Println("handleUploadExcel", elapsed)
	return data, nil
}

func saveSalary(relatedFields []template.Field) {
	// 从其它模板关联的数据 不需要再保存进数据库

	//tempMap := make(map[string]struct{})
	//
	//for _, f := range relatedFields {
	//	tempMap[f.Key] = struct{}{}
	//}

	allFields := make([]model.SalaryField, 0)
	for _, fields := range DataMap {
		for _, item := range fields {
			allFields = append(allFields, item)
			allFieldMap[item.Key] = true
			//if _, ok := tempMap[item.Key]; !ok {
			//	allFields = append(allFields, item)
			//}
		}
	}

	//fmt.Println("allFieldMap", util.PrettyJson(allFieldMap))
	//for field, templateName := range excelUploadedFields {
	//	fmt.Println("fields", )
	//	if _, ok := allFieldMap[field]; !ok {
	//		errMsg = append(errMsg, "工资项目不存在，位于模板:["+templateName+"],字段名:["+field+"]")
	//	}
	//}
	// 检查上传的数据中的身份证号码是否正确
	uploadedCards := make(map[string]string) // 身份证号码 为key , 值为模板名
	for templateName, profileData := range cardIDMap {
		for card := range profileData {
			uploadedCards[card] = templateName
		}
	}

	profileCard := make(map[string]struct{})
	for _, p := range profiles {
		profileCard[p.IDCard] = struct{}{}
	}

	for card, templateName := range uploadedCards {
		if _, ok := profileCard[card]; !ok {
			errMsg = append(errMsg, "身份证号码不存在或者发生错误，位于模板:["+templateName+"],身份证号码:["+card+"]")
		}
	}

	wModel.Fields = allFields
	err := wModel.Create()
	if err != nil {
		log.Error("saveSalary function called.", err)
	}
	DataMap = nil
}

//解析归属年月份
//例如加班工资发放，在当月发放，但应归属于上一月份的收入
func resolveFitMonth(field template.Field, year, month string) (string, string) {

	if field.FitIntoMonth == "LASTMONTH" {
		return util.LastMonth(year, month)
	}
	if field.FitIntoMonth == "LASTYEAR" {
		y, _ := strconv.Atoi(year)
		year = fmt.Sprint(y - 1)
		month = "12"
		return year, month
	}
	return year, month
}

func structToMap(model interface{}) map[string]interface{} {
	var sMap map[string]interface{}
	j, _ := json.Marshal(model)
	json.Unmarshal(j, &sMap)
	return sMap
}
