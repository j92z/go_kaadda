package file_api

import (
	"github/j92z/go_kaadda/model/file_model"
	"github/j92z/go_kaadda/pkg/util/file_util"
	"github/j92z/go_kaadda/setting"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	contentType := file.Header.Get("Content-Type")
	if !checkAllowFileType(contentType) {
		c.String(http.StatusBadRequest, "不允许"+contentType+"类型文件上传")
		return
	}
	if file.Size > setting.EnvSetting.File.MaxSize {
		c.String(http.StatusBadRequest, "文件大小 > 100M, 请上传小于100M的文件")
		return
	}
	dateString := time.Now().Format("20060102")
	filePath := file_util.PathJoin(setting.EnvSetting.File.Path, dateString)
	file_util.CheckDir(filePath)
	uniqueFileName := file_util.GetUniqueFileName(filePath, file.Filename)
	fileFullPath := file_util.PathJoin(filePath, uniqueFileName)
	if err := c.SaveUploadedFile(file, fileFullPath); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	fileModel := &file_model.File{
		Name:     file.Filename,
		Path:     fileFullPath,
		MimeType: contentType,
	}
	if err := fileModel.Add(); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, fileModel)
}

func Resource(c *gin.Context) {
	var param struct {
		ID int `form:"Id"	binding:"required,gt=0"`
	}
	if err := c.ShouldBindQuery(&param); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	fileModel := file_model.FindById(param.ID)
	if fileModel.ID <= 0 || !file_util.CheckFile(fileModel.Path) {
		c.String(http.StatusBadRequest, "找不到当前文件")
		return
	}
	fileContent, err := ioutil.ReadFile(fileModel.Path)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Header("Content-Disposition", "filename=\""+fileModel.Name+"\"")
	c.Data(http.StatusOK, fileModel.MimeType, fileContent)
}

func Download(c *gin.Context) {
	var param struct {
		ID int `form:"Id"	binding:"required,gt=0"`
	}
	if err := c.ShouldBindQuery(&param); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	fileModel := file_model.FindById(param.ID)
	if fileModel.ID <= 0 || !file_util.CheckFile(fileModel.Path) {
		c.String(http.StatusBadRequest, "找不到当前文件")
		return
	}
	fileContent, err := ioutil.ReadFile(fileModel.Path)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Header("Content-Disposition", "attachment;filename=\""+fileModel.Name+"\"")
	c.Data(http.StatusOK, fileModel.MimeType, fileContent)
}

func RemoveFile(c *gin.Context) {
	var param struct {
		ID int `form:"Id"	binding:"required,gt=0"`
	}
	if err := c.ShouldBindQuery(&param); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	fileModel := file_model.FindById(param.ID)
	if fileModel.ID <= 0 {
		c.String(http.StatusBadRequest, "找不到当前文件")
		return
	}
	if err := fileModel.Remove(); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, fileModel)
}

func checkAllowFileType(fileType string) bool {
	for _, v := range setting.AllowFileType {
		if v == fileType {
			return true
		}
	}
	return false
}
