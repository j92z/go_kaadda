package enginex

import (
	"github/j92z/go_kaadda/httpx/enginex/api/announcement_api"
	"github/j92z/go_kaadda/httpx/enginex/api/banner_api"
	"github/j92z/go_kaadda/httpx/enginex/api/file_api"
	"github/j92z/go_kaadda/httpx/enginex/api/ticket_api"
	"github/j92z/go_kaadda/httpx/enginex/api/ticket_type_api"
	"github/j92z/go_kaadda/setting"
	"github.com/gin-gonic/gin"
)

func InitEngineX() *gin.Engine {
	gin.SetMode(setting.EnvSetting.Server.RunMode)
	engine := gin.New()
	//r.Use(TlsHandler())
	InitRouter(engine)

	return engine
}

func InitRouter(r *gin.Engine) {
	ticketApi := r.Group("Ticket")
	{
		ticketApi.PUT("", ticket_api.AddTicket)
		ticketApi.PATCH("", ticket_api.EditTicket)
		ticketApi.PATCH("ChangeStatus", ticket_api.ChangeStatus)
		ticketApi.DELETE("", ticket_api.RemoveTicket)
		ticketApi.GET("ById", ticket_api.GetById)
		ticketApi.GET("GetList", ticket_api.GetListByPage)
	}
	r.GET("TicketType", ticket_type_api.GetList)
	fileApi := r.Group("File")
	{
		fileApi.PUT("", file_api.Upload)
		fileApi.GET("ById", file_api.Resource)
		fileApi.GET("Download", file_api.Download)
		fileApi.DELETE("", file_api.RemoveFile)
	}
	announcementApi := r.Group("Announcement")
	{
		announcementApi.PUT("", announcement_api.AddAnnouncement)
		announcementApi.PATCH("Up", announcement_api.Up)
		announcementApi.PATCH("Down", announcement_api.Down)
		announcementApi.PATCH("View", announcement_api.View)
		announcementApi.PATCH("Star", announcement_api.Star)
		announcementApi.PATCH("", announcement_api.EditAnnouncement)
		announcementApi.DELETE("", announcement_api.RemoveAnnouncement)
		announcementApi.GET("ById", announcement_api.GetAnnouncementById)
		announcementApi.GET("GetList", announcement_api.GetAnnouncementList)
	}
	bannerApi := r.Group("Banner")
	{
		bannerApi.PUT("", banner_api.AddBanner)
		bannerApi.PATCH("", banner_api.EditBanner)
		bannerApi.DELETE("", banner_api.RemoveBanner)
		bannerApi.GET("ById", banner_api.GetBannerById)
		bannerApi.GET("GetList", banner_api.GetBannerList)
		bannerApi.GET("Get5", banner_api.Get5Banner)
	}
}
