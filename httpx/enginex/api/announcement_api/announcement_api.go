package announcement_api

import (
	"github/j92z/go_kaadda/httpx/enginex/api"
	"github/j92z/go_kaadda/httpx/enginex/service/announcement_service"
	"github/j92z/go_kaadda/httpx/enginex/service/file_service"
	"github/j92z/go_kaadda/model/announcement_model"
	"github/j92z/go_kaadda/model/file_model"
	"github/j92z/go_kaadda/model/star_relation_model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func AddAnnouncement(c *gin.Context) {
	var param struct {
		Title   string `form:"title" binding:"required"`
		Content string `form:"content" binding:"required"`
		Creator string `form:"creator" binding:"required"`
		FIds    string `form:"fids"`
	}
	if err := c.ShouldBind(&param); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := api.RequestUserInfo(param.Creator)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	announcementModel := &announcement_model.Announcement{
		Title:       param.Title,
		Content:     param.Content,
		Creator:     param.Creator,
		CreatorName: userInfo[0].Name,
	}
	FidSlice := strings.Split(param.FIds, ",")
	if len(FidSlice) > 0 {
		if file_model.CountByIds(FidSlice) != len(FidSlice) {
			c.String(http.StatusBadRequest, "找不到文件")
			return
		}
		announcementModel.FIds = param.FIds
	}
	if err := announcementModel.Add(); err != nil {
		c.String(http.StatusInternalServerError, "添加失败")
		return
	}
	c.JSON(http.StatusOK, announcementModel)
}

func Up(c *gin.Context) {
	var param struct {
		ID int `form:"id" binding:"required,gt=0"`
	}
	if err := c.ShouldBind(&param); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	announcementModel := announcement_model.FindById(param.ID)
	if announcementModel.ID <= 0 {
		c.String(http.StatusBadRequest, "找不到当前公告")
		return
	}
	if announcementModel.Status == 1 {
		c.String(http.StatusBadRequest, "当前公告已上架")
		return
	}
	announcementModel.Status = 1
	if err := announcementModel.Update(); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func Down(c *gin.Context) {
	var param struct {
		ID int `form:"id" binding:"required,gt=0"`
	}
	if err := c.ShouldBind(&param); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	announcementModel := announcement_model.FindById(param.ID)
	if announcementModel.ID <= 0 {
		c.String(http.StatusBadRequest, "找不到当前公告")
		return
	}
	if announcementModel.Status == 0 {
		c.String(http.StatusBadRequest, "当前公告已下架")
		return
	}
	announcementModel.Status = 0
	if err := announcementModel.Update(); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func View(c *gin.Context) {
	var param struct {
		ID int `form:"id" binding:"required,gt=0"`
	}
	if err := c.ShouldBind(&param); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	announcementModel := announcement_model.FindById(param.ID)
	if announcementModel.ID <= 0 {
		c.String(http.StatusBadRequest, "找不到当前公告")
		return
	}

	announcementModel.ViewCount = announcementModel.ViewCount + 1
	if err := announcementModel.Update(); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"Success": true, "ViewCount": announcementModel.ViewCount})
}

func Star(c *gin.Context) {
	var param struct {
		ID  int    `form:"id" binding:"required,gt=0"`
		Uid string `form:"uid" binding:"required"`
	}
	if err := c.ShouldBind(&param); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if _, err := api.RequestUserInfo(param.Uid); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if err := announcement_service.StarOrNot(param.Uid, param.ID); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	announcementModel := announcement_model.FindById(param.ID)
	c.JSON(http.StatusOK, gin.H{"Success": true, "StarCount": announcementModel.StarCount})
}

func EditAnnouncement(c *gin.Context) {
	var param struct {
		ID      int    `form:"id" binding:"required,gt=0"`
		Title   string `form:"title"`
		Content string `form:"content"`
		FIds    string `form:"fids"`
	}
	if err := c.ShouldBind(&param); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	announcementModel := announcement_model.FindById(param.ID)
	if announcementModel.ID <= 0 {
		c.String(http.StatusBadRequest, "找不到当前公告")
		return
	}
	FidSlice := strings.Split(param.FIds, ",")
	if len(FidSlice) > 0 {
		if file_model.CountByIds(FidSlice) != len(FidSlice) {
			c.String(http.StatusBadRequest, "找不到文件")
			return
		}
		announcementModel.FIds = param.FIds
	}
	if _, ok := c.GetPostForm("title"); ok {
		announcementModel.Title = param.Title
	}
	if _, ok := c.GetPostForm("content"); ok {
		announcementModel.Content = param.Content
	}
	if err := announcementModel.Update(); err != nil {
		c.String(http.StatusInternalServerError, "编辑失败")
		return
	}
	c.JSON(http.StatusOK, announcementModel)
}

func RemoveAnnouncement(c *gin.Context) {
	var param struct {
		ID int `form:"Id" binding:"required,gt=0"`
	}
	if err := c.ShouldBindQuery(&param); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	announcementModel := announcement_model.FindById(param.ID)
	if announcementModel.ID <= 0 {
		c.String(http.StatusBadRequest, "找不到当前公告")
		return
	}
	if err := announcementModel.Remove(); err != nil {
		c.String(http.StatusInternalServerError, "移除公告失败")
		return
	}
	c.JSON(http.StatusOK, announcementModel)
}

type AnnouncementResponseStruct struct {
	ID          uint                               `json:"Id"`
	Title       string                             `json:"Title"`
	Content     string                             `json:"Content"`
	Creator     string                             `json:"Creator"`
	CreatorName string                             `json:"CreatorName"`
	CreatedAt   time.Time                          `json:"CreatedAt"`
	Status      int                                `json:"Status"`
	File        *[]file_service.FileResponseStruct `json:"File"`
	StarList    *[]string                          `json:"StarList"`
}

func GetAnnouncementList(c *gin.Context) {
	var param struct {
		Page int `form:"Page" binding:"gte=0"`
		Size int `form:"Size" binding:"gte=0"`
		Type int `form:"Type" binding:"gte=0,lte=1"`
	}
	if err := c.ShouldBindQuery(&param); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	var size = 20
	var status = -1
	if param.Size > 0 {
		size = param.Size
	}
	if _, ok := c.GetQuery("Type"); ok {
		status = param.Type
	}
	fmt.Println(status)
	list := announcement_model.GetList(param.Page, size, status)
	count := announcement_model.GetListCount(status)
	var returnList []AnnouncementResponseStruct
	for _, v := range list {
		tmp := AnnouncementResponseStruct{
			ID:          v.ID,
			Title:       v.Title,
			Content:     v.Content,
			Creator:     v.Creator,
			CreatorName: v.CreatorName,
			CreatedAt:   v.CreatedAt,
			Status:      v.Status,
			File:        file_service.GetFileInfoByIds(v.FIds),
		}
		returnList = append(returnList, tmp)
	}
	c.JSON(http.StatusOK, gin.H{
		"page":  param.Page,
		"size":  size,
		"count": count,
		"data":  returnList,
	})
}

func GetAnnouncementById(c *gin.Context) {
	var param struct {
		ID int `form:"Id" binding:"required,gte=0"`
	}
	if err := c.ShouldBindQuery(&param); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	announcementModel := announcement_model.FindById(param.ID)
	if announcementModel.ID <= 0 {
		c.String(http.StatusBadRequest, "找不到当前公告")
		return
	}
	info := AnnouncementResponseStruct{
		ID:          announcementModel.ID,
		Title:       announcementModel.Title,
		Content:     announcementModel.Content,
		Creator:     announcementModel.Creator,
		CreatorName: announcementModel.CreatorName,
		CreatedAt:   announcementModel.CreatedAt,
		Status:      announcementModel.Status,
		File:        file_service.GetFileInfoByIds(announcementModel.FIds),
		StarList:    star_relation_model.FindUidByAid(param.ID),
	}
	c.JSON(http.StatusOK, info)
}
