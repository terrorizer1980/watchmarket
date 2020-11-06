package endpoint

import (
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/watchmarket/services/controllers"
	"net/http"
)

func SubscribeHandler(controller controllers.AlertsController) func(c *gin.Context) {
	return func(c *gin.Context) {
		var sr controllers.SubscriptionsRequest
		err := c.BindJSON(&sr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, nil)
			return
		}
		err = controller.HandleSubscriptionsRequest(sr, c.Request.Context())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
			return
		}
		c.JSON(http.StatusOK, nil)
	}
}
