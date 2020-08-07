package banner_api

import (
	"github/j92z/go_kaadda/model/announcement_model"
	"github/j92z/go_kaadda/model/banner_model"
	"github/j92z/go_kaadda/model/file_model"
	"github/j92z/go_kaadda/setting"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func AddBanner(c *gin.Context) {
	var param struct {
		Fid   int    `form:"fid" binding:"required,gt=0"`
		Aid   int    `form:"aid" binding:"gt=0"`
		APath string `form:"a_path"`
	}
	if err := c.ShouldBind(&param); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	bannerModel := &banner_model.Banner{
		Fid:   param.Fid,
		FPath: setting.EnvSetting.Server.Path + "/File/ById?Id=" + strconv.Itoa(param.Fid),
	}
	if _, ok := c.GetPostForm("aid"); ok {
		if !announcement_model.ExistById(param.Aid) {
			c.String(http.StatusBadRequest, "找不到所添加的公告")
			return
		}
		bannerModel.Aid = param.Aid
	}
	if _, ok := c.GetPostForm("a_path"); ok {
		bannerModel.APath = param.APath
	}
	if err := bannerModel.Add(); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, bannerModel)
}

func RemoveBanner(c *gin.Context) {
	var param struct {
		ID int `form:"Id" binding:"required,gte=0"`
	}
	if err := c.ShouldBindQuery(&param); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	bannerModel := banner_model.FindById(param.ID)
	if bannerModel.ID <= 0 {
		c.String(http.StatusBadRequest, "找不到当前Banner")
		return
	}
	if err := bannerModel.Remove(); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, bannerModel)
}

func EditBanner(c *gin.Context) {
	var param struct {
		ID    int    `form:"id" binding:"required,gt=0"`
		Fid   int    `form:"fid" binding:"gt=0"`
		Aid   int    `form:"aid" binding:"gt=0"`
		APath string `form:"a_path"`
	}
	if err := c.ShouldBind(&param); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	bannerModel := banner_model.FindById(param.ID)
	if bannerModel.ID <= 0 {
		c.String(http.StatusBadRequest, "找不到当前Banner")
		return
	}
	if _, ok := c.GetPostForm("fid"); ok {
		if !file_model.ExistById(param.Fid) {
			c.String(http.StatusBadRequest, "找不到添加的文件")
			return
		}
		bannerModel.Fid = param.Fid
	}
	if _, ok := c.GetPostForm("aid"); ok {
		if !announcement_model.ExistById(param.Aid) {
			c.String(http.StatusBadRequest, "找不到所添加的公告")
			return
		}
		bannerModel.Aid = param.Aid
	}
	if _, ok := c.GetPostForm("a_path"); ok {
		bannerModel.APath = param.APath
	}
	if err := bannerModel.Update(); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, bannerModel)
}

func GetBannerById(c *gin.Context) {
	var param struct {
		ID int `form:"Id" binding:"required,gte=0"`
	}
	if err := c.ShouldBindQuery(&param); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	bannerModel := banner_model.FindById(param.ID)
	if bannerModel.ID <= 0 {
		c.String(http.StatusBadRequest, "找不到当前Banner")
		return
	}
	c.JSON(http.StatusOK, bannerModel)
}

func GetBannerList(c *gin.Context) {
	var param struct {
		Page int `form:"Page" binding:"gte=0"`
		Size int `form:"Size" binding:"gte=0"`
	}
	if err := c.ShouldBindQuery(&param); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	var size = 20
	if param.Size > 0 {
		size = param.Size
	}
	bannerList := banner_model.GetList(param.Page, size)
	bannerCount := banner_model.GetListCount()
	c.JSON(http.StatusOK, gin.H{
		"page":  param.Page,
		"size":  size,
		"count": bannerCount,
		"data":  bannerList,
	})
}

func Get5Banner(c *gin.Context) {
	bannerList := banner_model.GetListOrder(0, 5, "created_at desc")
	c.JSON(http.StatusOK, bannerList)
}
