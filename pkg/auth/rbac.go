package auth

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/casbin/casbin"
)

var Enforcer *casbin.Enforcer

func GetEnforcer(model, policy string) *casbin.Enforcer {
	Enforcer = casbin.NewEnforcer(model, policy)
	return Enforcer
}

func RefreshEnforcer() {
	merge("./conf/permission")
	Enforcer.LoadPolicy()
	// GetEnforcer("conf/authz_model.conf", "conf/authz_policy.csv")
}

func getFilelist(path string) (files []string) {
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}

	return files
}

func merge(rootPath string) {
	outFileName := "./conf/authz_policy.csv"
	outFile, openErr := os.OpenFile(outFileName, os.O_WRONLY|os.O_TRUNC, 0600)
	if openErr != nil {
		fmt.Printf("Can not open file %s", outFileName)
	}
	bWriter := bufio.NewWriter(outFile)
	bWriter.Write(make([]byte, 0))
	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		fmt.Println("Processing:", path)
		//这里是文件过滤器，表示我仅仅处理csv文件
		if strings.HasSuffix(path, ".csv") {
			fp, fpOpenErr := os.Open(path)
			if fpOpenErr != nil {
				fmt.Printf("Can not open file %v", fpOpenErr)
				return fpOpenErr
			}
			bReader := bufio.NewReader(fp)
			for {
				buffer := make([]byte, 1024)
				readCount, readErr := bReader.Read(buffer)
				if readErr == io.EOF {
					break
				} else {
					bWriter.Write(buffer[:readCount])
				}
			}
		}
		return err
	})
	bWriter.Flush()
}
