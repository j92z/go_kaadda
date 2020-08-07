package file_service

import (
	"github/j92z/go_kaadda/model/file_model"
	"github/j92z/go_kaadda/setting"
	"strconv"
	"strings"
)

type FileResponseStruct struct {
	ID       uint   `json:"Id"`
	Name     string `json:"Name"`
	Path     string `json:"Path"`
	MimeType string `json:"MimeType"`
}

func GetFileInfoByIds(ids string) *[]FileResponseStruct {
	FidSlice := strings.Split(ids, ",")
	var fileList []FileResponseStruct
	if len(FidSlice) > 0 {
		list := file_model.FindByIds(ids)
		for _, v := range list {
			tmp := FileResponseStruct{
				ID:       v.ID,
				Name:     v.Name,
				Path:     setting.EnvSetting.Server.Path + "/File/ById?Id=" + strconv.Itoa(int(v.ID)),
				MimeType: v.MimeType,
			}
			fileList = append(fileList, tmp)
		}
	}
	return &fileList
}
