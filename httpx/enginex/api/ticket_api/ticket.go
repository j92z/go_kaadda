package ticket_api

import (
	"github/j92z/go_kaadda/httpx/enginex/api"
	"github/j92z/go_kaadda/model/ticket_detail_model"
	"github/j92z/go_kaadda/model/ticket_type_model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddTicket(c *gin.Context) {
	var param struct {
		Submitter  string `form:"submitter" binding:"required"`
		TicketType int    `form:"ticket_type" binding:"required"`
		Introduce  string `form:"introduce" binding:"required"`
		Content    string `form:"content" binding:"required"`
		Contact    string `form:"contact" binding:"required"`
	}
	if err := c.ShouldBind(&param); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	userInfo, err := api.RequestUserInfo(param.Submitter)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ticketType := ticket_type_model.FindById(param.TicketType)

	if ticketType.ID <= 0 {
		c.String(http.StatusBadRequest, "找不到工单类型")
		return
	}
	ticketDetail := &ticket_detail_model.TicketDetail{
		CompanyID:      userInfo[0].Parent.ID,
		CompanyName:    userInfo[0].Parent.Name,
		TicketType:     param.TicketType,
		TicketTypeName: ticketType.Name,
		Submitter:      param.Submitter,
		SubmitterName:  userInfo[0].Name,
		Introduce:      param.Introduce,
		Content:        param.Content,
		Status:         0,
		Contact:        param.Contact,
	}
	if err := ticketDetail.Add(); err != nil || ticketDetail.ID <= 0 {
		c.String(http.StatusBadRequest, "写入数据失败")
		return
	}
	c.JSON(http.StatusOK, ticketDetail)
}

func EditTicket(c *gin.Context) {
	var param struct {
		ID         int    `form:"id" binding:"required"`
		Submitter  string `form:"submitter" binding:"required"`
		TicketType int    `form:"job_type"`
		Introduce  string `form:"introduce"`
		Content    string `form:"content"`
		Contact    string `form:"contact"`
	}
	if err := c.ShouldBind(&param); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	ticketDetail := ticket_detail_model.FindById(param.ID)
	if ticketDetail.ID <= 0 {
		c.String(http.StatusBadRequest, "Can't Find This Ticket")
		return
	}
	userInfo, err := api.RequestUserInfo(param.Submitter)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if userInfo[0].ID != ticketDetail.Submitter {
		c.String(http.StatusBadRequest, "Submitter Can't Edit")
		return
	}

	if param.TicketType > 0 && param.TicketType != ticketDetail.TicketType {
		ticketType := ticket_type_model.FindById(param.TicketType)
		if ticketType.ID <= 0 {
			c.String(http.StatusBadRequest, "找不到工单类型")
			return
		}
		ticketDetail.TicketType = param.TicketType
		ticketDetail.TicketTypeName = ticketType.Name
	}

	if len(param.Introduce) > 0 && param.Introduce != ticketDetail.Introduce {
		ticketDetail.Introduce = param.Introduce
	}

	if len(param.Content) > 0 && param.Content != ticketDetail.Content {
		ticketDetail.Content = param.Content
	}

	if len(param.Contact) > 0 && param.Contact != ticketDetail.Contact {
		ticketDetail.Contact = param.Contact
	}

	if err := ticketDetail.Update(); err != nil {
		c.String(http.StatusBadRequest, "更新工单失败")
		return
	}
	c.JSON(http.StatusOK, ticketDetail)
}

func ChangeStatus(c *gin.Context) {
	var param struct {
		ID     int `form:"id" binding:"required"`
		Status int `form:"status" binding:"gte=0,lte=1"`
	}
	if err := c.ShouldBind(&param); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	ticketDetail := ticket_detail_model.FindById(param.ID)
	if ticketDetail.ID <= 0 {
		c.String(http.StatusBadRequest, "Can't Find This Ticket")
		return
	}
	if param.Status <= 0 {
		ticketDetail.Status = 0
	} else {
		ticketDetail.Status = param.Status
	}

	if err := ticketDetail.Update(); err != nil {
		c.String(http.StatusBadRequest, "更新状态失败")
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "status": ticketDetail.Status})
}

func RemoveTicket(c *gin.Context) {
	var param struct {
		ID int `form:"Id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&param); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	ticketDetail := ticket_detail_model.FindById(param.ID)
	if ticketDetail.ID <= 0 {
		c.String(http.StatusBadRequest, "Can't Find This Ticket")
		return
	}
	if err := ticketDetail.Remove(); err != nil {
		fmt.Println(err)
		c.String(http.StatusBadRequest, "删除失败")
		return
	}
	c.JSON(http.StatusOK, ticketDetail)
}

func GetById(c *gin.Context) {
	var param struct {
		ID int `form:"Id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&param); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	ticketDetail := ticket_detail_model.FindById(param.ID)
	if ticketDetail.ID <= 0 {
		c.String(http.StatusBadRequest, "Can't Find This Ticket")
		return
	}
	c.JSON(http.StatusOK, ticketDetail)
}

func GetListByPage(c *gin.Context) {
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
	list := ticket_detail_model.GetList(param.Page, size, status)
	count := ticket_detail_model.GetListCount(status)
	c.JSON(http.StatusOK, gin.H{
		"page":  param.Page,
		"size":  size,
		"count": count,
		"data":  list,
	})
}
