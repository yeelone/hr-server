package tool

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	h "hr-server/handler"
	"hr-server/util"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

func Backup(c *gin.Context) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	backupFile := dir + "/backup/sql/backup-" + time.Now().Format("2006_01_02_15_04") + ".sql"

	switch runtime.GOOS {
	case "windows":
		cmdStr := fmt.Sprintf(`"host=%s hostaddr=%s port=%s user=%s dbname=%s" `,
			viper.GetString("db.addr"),
			viper.GetString("db.addr"),
			viper.GetString("db.port"),
			viper.GetString("db.username"),
			//viper.GetString("db.password"),
			viper.GetString("db.name"),
		)
		fmt.Println("cmd", cmdStr)
		cmd := exec.Command("cmd", "/C", "SET", "PGPASSWORD="+viper.GetString("db.password"))
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + string(output))
			//h.SendResponse(c, err, nil)
			//return
		}
		// 在windows上会失败
		//cmd = exec.Command("cmd", "/C", "pg_dump", "-Fc", "-f", backupFile, cmdStr)
		//output, err = cmd.CombinedOutput()
		//if err != nil {
		//	fmt.Println(fmt.Sprint(err) + ": " + string(output))
		//	h.SendResponse(c, err, nil)
		//	return
		//}
		//if err != nil {
		//	fmt.Println("err", err, result, cmdStr, backupFile)
		//	backupFile = ""
		//	h.SendResponse(c, err, nil)
		//	return
		//}
	case "linux":
		linuxShellStr := fmt.Sprintf(`pg_dump -Fc "host=%s hostaddr=%s port=%s user=%s password=%s dbname=%s" > %s`,
			viper.GetString("db.addr"),
			viper.GetString("db.addr"),
			viper.GetString("db.port"),
			viper.GetString("db.username"),
			viper.GetString("db.password"),
			viper.GetString("db.name"),
			backupFile,
		)
		_, err = exec.Command("/bin/sh", "-c", linuxShellStr).Output()
	}

	files := []string{backupFile}
	if err != nil {
		h.SendResponse(c, err, nil)
		return
	}

	zipFileName := dir + "/backup/backup_" + time.Now().Format("2006_01_02_15_04") + ".zip"
	//保留原来文件的结构
	err = util.ZipDir("./conf", files, zipFileName)
	if err != nil {
		fmt.Println("zip error", err)
		h.SendResponse(c, err, nil)
		return
	}

	h.SendResponse(c, nil, CreateResponse{File: zipFileName})
	return
}

func ListBackupFiles(c *gin.Context) {
	var files []string
	var exactFilePath []string

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	root := dir + "/backup"
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		files = append(files, "/backup/"+info.Name())
		exactFilePath = append(exactFilePath, path)
		return nil
	})

	if err != nil {
		panic(err)
	}
	h.SendResponse(c, nil, CreateResponse{Files: files})
	return
}
