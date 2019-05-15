package model

import (
	"hrgdrc/util"
	"strconv"
)

type Statistics struct {
	Name       string `json:"name"`
	IDCard     string `json:"id_card"`
	Profile    string `json:"profile"`
	Department string `json:"department"`
	Post       string `json:"post"`
	Field      string `json:"field"`
	Year       string `json:"year"`
	Month      string `json:"month"`
	Value      string `json:"value"`
	Total      string `json:"total"`
	Number     int    `json:"number"`
}

// FetchAnnualIncome : 获取指定年份的总收入
func FetchAnnualIncome(year string, pids []uint64) (sf []Statistics, err error) {
	var result []Statistics
	if len(pids) > 0 {
		//sql := `select name as field,year,month,content,value,profile_id as profile,post_group_id as post,department_group_id as department ` +
		//	`from tb_salary_fields  where (year=? or fit_into_year=?) and profile_id in (?)  and is_income=true and value<>-1 order by profile_id, year,month;`
		sql := `select profile_id as profile,department_group_id as department ,year,sum(value) as total from tb_salary_fields  where (year=? or fit_into_year=?) ` +
			`and profile_id in (?) and is_income=true and value<>-1 group by profile_id,department_group_id,year;`
		err = DB.Self.Raw(sql, year, year, util.ArrayToString(pids, ",")).Scan(&result).Error
	} else {
		//sql := `select name as field,year,month,content,value,profile_id as profile,post_group_id as post,department_group_id as department ` +
		//	`from tb_salary_fields  where (year=? or fit_into_year=?)  and is_income=true and value<>-1 order by profile_id, year,month ;`
		sql := `select profile_id as profile,department_group_id as department ,year,sum(value) as total from tb_salary_fields  where (year=? or fit_into_year=?) and is_income=true and value<>-1 group by profile_id,department_group_id,year;`
		err = DB.Self.Raw(sql, year, year).Scan(&result).Error
	}

	fillProfile(result)

	return result, err
}

// FetchDetail : 获取指定年份的总收入
func FetchDetail(year string, pids []uint64) (sf []Statistics, err error) {
	var result []Statistics
	if len(pids) > 0 {
		//sql := `select name as field,year,month,content,value,profile_id as profile,post_group_id as post,department_group_id as department ` +
		//	`from tb_salary_fields  where (year=? or fit_into_year=?) and profile_id in (?)  and is_income=true and value<>-1 order by profile_id, year,month;`
		sql := `select profile_id as profile,department_group_id as department ,year,sum(value) as total from tb_salary_fields  where (year=? or fit_into_year=?) ` +
			`and profile_id in (?) and is_income=true and value<>-1 group by profile_id,department_group_id,year;`
		err = DB.Self.Raw(sql, year, year, util.ArrayToString(pids, ",")).Scan(&result).Error
	} else {
		//sql := `select name as field,year,month,content,value,profile_id as profile,post_group_id as post,department_group_id as department ` +
		//	`from tb_salary_fields  where (year=? or fit_into_year=?)  and is_income=true and value<>-1 order by profile_id, year,month ;`
		sql := `select profile_id as profile,department_group_id as department ,year,sum(value) as total from tb_salary_fields  where (year=? or fit_into_year=?) and is_income=true and value<>-1 group by profile_id,department_group_id,year;`
		err = DB.Self.Raw(sql, year, year).Scan(&result).Error
	}

	fillProfile(result)

	return result, err
}

func fillProfile(result []Statistics) {
	//取出来的的profile, department,post为id, 需要通过这些ID来取得对应的名字展示给客户端.
	pidMap := make(map[string]string)
	cardMap := make(map[string]string)
	departMap := make(map[string]string)
	postMap := make(map[string]string)

	for i, r := range result {
		if name, ok := pidMap[r.Profile]; ok {
			result[i].Profile = name
			result[i].IDCard = cardMap[r.Profile]
		} else {
			id, err := strconv.ParseUint(r.Profile, 10, 64)
			if err != nil {
				continue
			}
			p, err := GetProfile(id)
			if err != nil {
				continue
			}
			pidMap[r.Profile] = p.Name
			cardMap[r.Profile] = p.IDCard
			result[i].Profile = p.Name
			result[i].IDCard = p.IDCard
		}

		if name, ok := departMap[r.Department]; ok {
			result[i].Department = name
		} else {
			id, err := strconv.ParseUint(r.Department, 10, 64)
			if err != nil {
				continue
			}
			d, err := GetGroup(id, false)
			if err != nil {
				continue
			}
			departMap[r.Department] = d.Name
			result[i].Department = d.Name
		}

		if name, ok := postMap[r.Post]; ok {
			result[i].Post = name
		} else {
			id, err := strconv.ParseUint(r.Post, 10, 64)
			if err != nil {
				continue
			}
			p, err := GetGroup(id, false)
			if err != nil {
				continue
			}
			postMap[r.Post] = p.Name
			result[i].Post = p.Name
		}
	}
}
