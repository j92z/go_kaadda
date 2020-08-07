package ticket_type_api

import (
	"github/j92z/go_kaadda/model/ticket_type_model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetList(c *gin.Context) {
	list := ticket_type_model.GetAll()
	c.JSON(http.StatusOK, list)
}
