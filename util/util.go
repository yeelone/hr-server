package util

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/teris-io/shortid"
)

func GenShortId() (string, error) {
	return shortid.Generate()
}

func GetReqID(c *gin.Context) string {
	v, ok := c.Get("X-Request-Id")
	if !ok {
		return ""
	}
	if requestId, ok := v.(string); ok {
		return requestId
	}
	return ""
}

func Uint2Str(i uint64) string {
	return fmt.Sprintf("%v", i)
}

const (
	empty = ""
	tab   = "\t"
)

func PrettyJson(data interface{}) string {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent(empty, tab)

	err := encoder.Encode(data)
	if err != nil {
		return empty
	}
	return buffer.String()
}

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func MoveFile(orgFile, desFile string) error {
	err := os.Rename(orgFile, desFile)
	return err
}

//ConvertToNumberingScheme : 将数字转化成ASCII字符，符合excel
func ConvertToNumberingScheme(number int) string {
	baseChar := []rune("A")[0]
	letters := ""

	for number > 0 {
		number -= 1
		letters = string(int(baseChar)+(number%26)) + letters
		number = (number / 26) >> 0 // quick `floor`
	}
	return letters
}

func ExtractFileName(filename string) (name, subffix string) {

	var filenameWithSuffix string
	filenameWithSuffix = path.Base(filename)

	var fileSuffix string
	fileSuffix = path.Ext(filenameWithSuffix) //获取文件后缀

	var filenameOnly string
	filenameOnly = strings.TrimSuffix(filenameWithSuffix, fileSuffix) //获取文件名

	return filenameOnly, fileSuffix

}

//LastMonth 根据传入的日期取得上个月
//https://play.golang.org/p/lnVudDwFYXK
func LastMonth(year, month string) (string, string) {
	y, _ := strconv.Atoi(year)
	m, _ := strconv.Atoi(month)

	if (m - 1) == 0 { //即1月份
		m = 12
		y = y - 1
	} else {
		m -= 1
	}

	y2 := strconv.Itoa(y)
	m2 := ""
	if m < 10 && m > 0 {
		m2 = "0" + strconv.Itoa(m)
	} else {
		m2 = strconv.Itoa(m)
	}

	return y2, m2
}

//Decimal
func Decimal(value float64) float64 {
	//value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	//value = value + 0.0005
	value = math.Round((value)*100) / 100
	return value
}

func ArrayToString(a []uint64, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

func Strip(str string) string {
	//先去除所有的空格
	str = strings.Replace(str, " ", "", -1)
	//同时去除换行符
	str = strings.Replace(str, "\n", "", -1)
	return str
}

// StringSliceEqualBCE 对比两个slice是否相等
func StringSliceEqualBCE(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	b = b[:len(a)]
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

// Read a whole file into the memory and store it as array of lines
func ReadLines(path string) (lines []string, err error) {
	var (
		file   *os.File
		part   []byte
		prefix bool
	)
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

func WriteLines(lines []string, path string) (err error) {
	var (
		file *os.File
	)

	if file, err = os.Create(path); err != nil {
		return
	}
	defer file.Close()

	for _, item := range lines {
		_, err := file.WriteString(strings.TrimSpace(item) + "\n")
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	return
}

func StructToMap(model interface{}) map[string]interface{} {
	var sMap map[string]interface{}
	j, _ := json.Marshal(model)
	json.Unmarshal(j, &sMap)
	return sMap
}

//FindUpdatedField : 通过对比新旧两个model来找出变化的字段
func FindUpdatedField(oldModel interface{}, newModel interface{}) (result map[string]map[string]interface{}) {
	s1 := StructToMap(oldModel)
	s2 := StructToMap(newModel)
	result = make(map[string]map[string]interface{}) // 字段 --》 旧值 新值
	for k := range s1 {
		if s1[k] != s2[k] {
			result[k] = make(map[string]interface{})
			result[k]["from"] = s1[k]
			result[k]["to"] = s2[k]
		}
	}
	return result
}

// RoundUp 四舍五入
func RoundUp(input float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Ceil(digit)
	newVal = round / pow
	return
}

// ZipFiles compresses one or many files into a single zip archive file.
//压缩多个文件到一个文件里面
// Param 1: 输出的zip文件的名字
// Param 2: 需要添加到zip文件里面的文件
//Param 3: 由于file是绝对路径，打包后可能不是想要的目录，oldform就是filename中需要被替换的掉的路径
//Param 4: 要替换成的路径
func ZipFiles(filename string, files []string, oldform, newform string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// 把files添加到zip中
	for _, file := range files {

		zipfile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer zipfile.Close()

		// 获取file的基础信息
		info, err := zipfile.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		//使用上面的FileInforHeader() 就可以把文件保存的路径替换成我们自己想要的了，如下面
		header.Name = filepath.Base(file)
		fmt.Println(header.Name)
		// 优化压缩
		// 更多参考see http://golang.org/pkg/archive/zip/#pkg-constants
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if _, err = io.Copy(writer, zipfile); err != nil {
			return err
		}
	}
	return nil
}

func ZipDir(dir string, files []string, zipFile string) error {

	fz, err := os.Create(zipFile)
	if err != nil {
		log.Fatalf("Create zip file failed: %s\n", err.Error())
		return err
	}
	defer fz.Close()

	w := zip.NewWriter(fz)
	defer w.Close()

	for _, f := range files {
		fDest, err := w.Create(f)
		if err != nil {
			log.Printf("Create failed: %s\n", err.Error())
			return err
		}
		fSrc, err := os.Open(f)
		if err != nil {
			log.Printf("Open failed: %s\n", err.Error())
			return err
		}
		defer fSrc.Close()
		_, err = io.Copy(fDest, fSrc)
		if err != nil {
			log.Printf("Copy failed: %s\n", err.Error())
			return err
		}
	}

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fDest, err := w.Create(path)
			if err != nil {
				log.Printf("Create failed: %s\n", err.Error())
				return nil
			}
			fSrc, err := os.Open(path)
			if err != nil {
				log.Printf("Open failed: %s\n", err.Error())
				return nil
			}
			defer fSrc.Close()
			_, err = io.Copy(fDest, fSrc)
			if err != nil {
				log.Printf("Copy failed: %s\n", err.Error())
				return nil
			}
		}
		return nil
	})

	return nil
}

func UnzipDir(zipFile string, dir string) {

	r, err := zip.OpenReader(zipFile)
	if err != nil {
		log.Fatalf("Open zip file failed: %s\n", err.Error())
	}
	defer r.Close()

	for _, f := range r.File {
		func() {
			path := dir + string(filepath.Separator) + f.Name
			os.MkdirAll(filepath.Dir(path), 0755)
			fDest, err := os.Create(path)
			if err != nil {
				log.Printf("Create failed: %s\n", err.Error())
				return
			}
			defer fDest.Close()

			fSrc, err := f.Open()
			if err != nil {
				log.Printf("Open failed: %s\n", err.Error())
				return
			}
			defer fSrc.Close()

			_, err = io.Copy(fDest, fSrc)
			if err != nil {
				log.Printf("Copy failed: %s\n", err.Error())
				return
			}
		}()
	}
}

func Uint64ArrayToInt64Array(a []uint64) []int64 {
	arr := make([]int64, len(a))

	for i := range a {
		arr[i] = int64(a[i])
	}
	return arr
}
