package file_util

import (
	"github.com/rs/xid"
	"os"
	"path/filepath"
	"strings"
)

func PathJoin(path ...string) string {
	newPath := ""
	for k, v := range path {
		if k == 0 {
			newPath += v
		} else {
			if strings.HasSuffix(newPath, "/") || strings.HasPrefix(v, "/") {
				newPath += v
			} else {
				newPath += "/" + v
			}
		}
	}
	return newPath
}

func CheckFile(path string) bool {
	s, err := os.Stat(path) //os.Stat获取文件信息
	return (err == nil && !s.IsDir()) || os.IsExist(err)
}

func CheckDir(path string) {
	s, err := os.Stat(path) //os.Stat获取文件信息
	if !(err == nil && s.IsDir()) {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			panic(err)
		}
	}
}

func GetUniqueFileName(path, fileName string) string {
	uniqueName := genUniqueName(fileName)
	filePath := PathJoin(path, uniqueName)
	if CheckFile(filePath) {
		return GetUniqueFileName(path, fileName)
	}
	return uniqueName
}

func genUniqueName(fileName string) string {
	uniqueId := xid.New().String()
	return uniqueId + "_" + fileName
}

func GetAbsPath(path string) string {
	dir, _ := os.Executable()
	exPath := filepath.Dir(dir)
	return PathJoin(exPath, path)
}
