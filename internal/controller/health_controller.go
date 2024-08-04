package controller

import (
	"gogin-template/baselib/dto"
	"gogin-template/baselib/identifier"
	"gogin-template/bootstrap"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HealthController struct {
	cfg *bootstrap.Container
}

func NewHealthController(server *gin.Engine, cfg *bootstrap.Container) {
	controller := &HealthController{
		cfg: cfg,
	}

	route := server.Group("/health")
	{
		route.GET("", controller.GetHealth)
	}
}

// @Summary      Get Health
// @Description  Get Health
// @Tags         Health
// @Produce      json
// @Success      200  {array}  	dto.Response[string]
// @Failure      500  {object}  dto.Response[any]
// @Router       /health [get]
func (c *HealthController) GetHealth(ctx *gin.Context) {
	resp := &dto.Response[string]{
		ResponseCode:    strconv.Itoa(http.StatusOK),
		ResponseMessage: "Success",
		LogReff:         identifier.GetLogReff(ctx),
		TraceId:         identifier.GetTraceId(ctx),
		Data:            "Hello World",
	}

	ctx.JSON(http.StatusOK, resp)
}
