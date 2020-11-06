package endpoint

import (
	"github.com/gin-gonic/gin"
	alertscontroller "github.com/trustwallet/watchmarket/services/controllers/alerts"
	"net/http"
)

func SubscribeHandler(controller alertscontroller.Controller) func(c *gin.Context) {
	return func(c *gin.Context) {
		var sr alertscontroller.SubscriptionsRequest
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
