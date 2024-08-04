package controller

import (
	"gogin-template/baselib/dto"
	"gogin-template/baselib/exception"
	"gogin-template/baselib/identifier"
	"gogin-template/bootstrap"
	"gogin-template/internal/service"
	"gogin-template/internal/viewmodel"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SampleController struct {
	service service.SampleService
	cfg     *bootstrap.Container
}

func NewSampleController(service service.SampleService, server *gin.Engine, cfg *bootstrap.Container) {
	controller := &SampleController{
		service: service,
		cfg:     cfg,
	}

	routes := server.Group("/sample")
	{
		routes.GET("", controller.GetSamples)
		routes.GET("/:sample-id", controller.GetSample)
		routes.GET("/:sample-id/version", controller.GetSampleVersions)
		routes.GET("/:sample-id/version/:version-number", controller.GetSampleVersion)

		routes.POST("", controller.SetSampleInsert)
		routes.POST("/versions", controller.SetSampleVersions)
		routes.POST("/version", controller.SetSampleVersionInsert)

		routes.PUT("", controller.SetSampleUpsert)
		routes.PUT("/version", controller.SetSampleVersionUpdate)

		routes.DELETE("/:sample-id", controller.GetSample)
		routes.DELETE("/:sample-id/version/:version-number", controller.SetSampleVersionDelete)
	}
}

// @Summary 	Get Samples
// @Description Get Samples
// @Tags 		Sample
// @Produce  	json
// @Param       sampleId		query  	string  false	"Sample ID"
// @Param       sampleType		query  	string  false	"Search Type"
// @Param 		search			query	string	false	"Search Query"
// @Param       page			query	int		false	"Page Index"
// @Param       pageSize		query	int		false	"Page Size"
// @Param       sortBy			query	string	false	"Sort By"
// @Param       sortDirection	query	string	false	"Sort Direction"
// @Success 	200	{object} 	dto.ApiResponse[*[]viewmodel.SampleRsViewModel]
// @Failure 	500	{object} 	dto.ApiResponse[any]
// @Router 		/sample 	[get]
func (c *SampleController) GetSamples(ctx *gin.Context) {
	var request viewmodel.SampleRqViewModel
	err := ctx.ShouldBindQuery(&request)
	if err != nil {
		ctx.Error(exception.ValidationException(strconv.Itoa(http.StatusBadRequest), err.Error()))
		return
	}

	pagination := dto.PageRequest{}
	err = ctx.ShouldBindQuery(&pagination)
	if err != nil {
		ctx.Error(exception.ValidationException(strconv.Itoa(http.StatusBadRequest), err.Error()))
		return
	}

	response, pageInfo, err := c.service.GetSamples(ctx.Request.Context(), &request, pagination)
	if err != nil {
		ctx.Error(err)
		return
	}

	resp := &dto.Response[*[]viewmodel.SampleRsViewModel]{
		ResponseCode:    strconv.Itoa(http.StatusOK),
		ResponseMessage: "Success",
		LogReff:         identifier.GetLogReff(ctx),
		TraceId:         identifier.GetTraceId(ctx),
		Data:            response,
		PageInfo:        pageInfo,
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	Get Sample
// @Description Get Sample
// @Tags 		Sample
// @Produce  	json
// @Param       sample-id		path  	string  true	"Sample ID"
// @Success 	200	{object} 	dto.ApiResponse[*viewmodel.SampleRsViewModel]
// @Failure 	500	{object} 	dto.ApiResponse[any]
// @Router 		/sample/{sample-id} 	[get]
func (c *SampleController) GetSample(ctx *gin.Context) {
	sampleId := ctx.Param("sample-id")

	request := &viewmodel.SampleRqViewModel{SampleId: sampleId}

	response, err := c.service.GetSample(ctx, request)
	if err != nil {
		ctx.Error(err)
		return
	}

	resp := &dto.Response[*viewmodel.SampleRsViewModel]{
		ResponseCode:    strconv.Itoa(http.StatusOK),
		ResponseMessage: "Success",
		LogReff:         identifier.GetLogReff(ctx),
		TraceId:         identifier.GetTraceId(ctx),
		Data:            response,
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	Get Sample Versions
// @Description Get Sample Versions
// @Tags 		Sample
// @Produce  	json
// @Param       sample-id		path  	string  true	"Sample ID"
// @Success 	200	{object} 	dto.ApiResponse[*[]viewmodel.SampleVersionRsViewModel]
// @Failure 	500	{object} 	dto.ApiResponse[any]
// @Router 		/sample/{sample-id}/version 	[get]
func (c *SampleController) GetSampleVersions(ctx *gin.Context) {
	sampleId := ctx.Param("sample-id")

	request := &viewmodel.SampleVersionRqViewModel{SampleId: sampleId}

	response, err := c.service.GetSampleVersions(ctx, request)
	if err != nil {
		ctx.Error(err)
		return
	}

	resp := &dto.Response[*[]viewmodel.SampleVersionRsViewModel]{
		ResponseCode:    strconv.Itoa(http.StatusOK),
		ResponseMessage: "Success",
		LogReff:         identifier.GetLogReff(ctx),
		TraceId:         identifier.GetTraceId(ctx),
		Data:            response,
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	Get Sample Version
// @Description Get Sample Version
// @Tags 		Sample
// @Produce  	json
// @Param       sample-id			path  	string  true	"Sample ID"
// @Param       version-number		path  	string  true	"Sample Version"
// @Success 	200	{object} 	dto.ApiResponse[*viewmodel.SampleVersionRsViewModel]
// @Failure 	500	{object} 	dto.ApiResponse[any]
// @Router 		/sample/{sample-id}/version/{version-number} 	[get]
func (c *SampleController) GetSampleVersion(ctx *gin.Context) {
	sampleId := ctx.Param("sample-id")
	versionNumber := ctx.Param("version-number")

	request := &viewmodel.SampleVersionRqViewModel{SampleId: sampleId, VersionNumber: versionNumber}

	response, err := c.service.GetSampleVersion(ctx, request)
	if err != nil {
		ctx.Error(err)
		return
	}

	resp := &dto.Response[*viewmodel.SampleVersionRsViewModel]{
		ResponseCode:    strconv.Itoa(http.StatusOK),
		ResponseMessage: "Success",
		LogReff:         identifier.GetLogReff(ctx),
		TraceId:         identifier.GetTraceId(ctx),
		Data:            response,
	}

	ctx.JSON(http.StatusOK, resp)

}

// @Summary 	Set Sample Insert
// @Description Set Sample Insert
// @Tags 		Sample
// @Accept  	json
// @Produce  	json
// @Param       request			body 	viewmodel.SampleRqViewModel  true  "Sample"
// @Success 	200	{object} 	dto.ApiResponse[*any]
// @Failure 	500	{object} 	dto.ApiResponse[any]
// @Router 		/sample [post]
func (c *SampleController) SetSampleInsert(ctx *gin.Context) {
	var request viewmodel.SampleRqViewModel
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.Error(exception.ValidationException(strconv.Itoa(http.StatusBadRequest), err.Error()))
		return
	}

	err = c.service.SetSample(ctx, "Insert", &request)
	if err != nil {
		ctx.Error(err)
		return
	}

	resp := &dto.Response[*any]{
		ResponseCode:    strconv.Itoa(http.StatusOK),
		ResponseMessage: "Success",
		LogReff:         identifier.GetLogReff(ctx),
		TraceId:         identifier.GetTraceId(ctx),
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	Set Sample Upsert
// @Description Set Sample Upsert
// @Tags 		Sample
// @Accept  	json
// @Produce  	json
// @Param       request body 	viewmodel.SampleRqViewModel  true  "Sample"
// @Success 	200	{object} 	dto.ApiResponse[*any]
// @Failure 	500	{object} 	dto.ApiResponse[any]
// @Router 		/sample [put]
func (c *SampleController) SetSampleUpsert(ctx *gin.Context) {
	var request viewmodel.SampleRqViewModel
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.Error(exception.ValidationException(strconv.Itoa(http.StatusBadRequest), err.Error()))
		return
	}

	err = c.service.SetSample(ctx, "Upsert", &request)
	if err != nil {
		ctx.Error(err)
		return
	}

	resp := &dto.Response[*any]{
		ResponseCode:    strconv.Itoa(http.StatusOK),
		ResponseMessage: "Success",
		LogReff:         identifier.GetLogReff(ctx),
		TraceId:         identifier.GetTraceId(ctx),
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	Set Sample Delete
// @Description Set Sample Delete
// @Tags 		Sample
// @Accept  	json
// @Produce  	json
// @Param       sample-id		path  	string	true	"Sample ID"
// @Success 	200	{object} 	dto.ApiResponse[*any]
// @Failure 	500	{object} 	dto.ApiResponse[any]
// @Router 		/sample/{sample-id} [delete]
func (c *SampleController) SetSampleDelete(ctx *gin.Context) {
	sampleId := ctx.Param("sample-id")

	request := &viewmodel.SampleRqViewModel{SampleId: sampleId}

	err := c.service.SetSample(ctx, "Delete", request)
	if err != nil {
		ctx.Error(err)
		return
	}

	resp := &dto.Response[*any]{
		ResponseCode:    strconv.Itoa(http.StatusOK),
		ResponseMessage: "Success",
		LogReff:         identifier.GetLogReff(ctx),
		TraceId:         identifier.GetTraceId(ctx),
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	Set Sample Versions
// @Description Set Sample Versions
// @Tags 		Sample
// @Accept  	json
// @Produce  	json
// @Param       request			body 	[]viewmodel.SampleVersionRqViewModel  true  "Sample Version"
// @Success 	200	{object} 	dto.ApiResponse[*any]
// @Failure 	500	{object} 	dto.ApiResponse[any]
// @Router 		/sample/versions [post]
func (c *SampleController) SetSampleVersions(ctx *gin.Context) {
	var request []viewmodel.SampleVersionRqViewModel
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.Error(exception.ValidationException(strconv.Itoa(http.StatusBadRequest), err.Error()))
		return
	}

	err = c.service.SetSampleVersions(ctx, &request)
	if err != nil {
		ctx.Error(err)
		return
	}

	resp := &dto.Response[*any]{
		ResponseCode:    strconv.Itoa(http.StatusOK),
		ResponseMessage: "Success",
		LogReff:         identifier.GetLogReff(ctx),
		TraceId:         identifier.GetTraceId(ctx),
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	Set Sample Version Insert
// @Description Set Sample Version Insert
// @Tags 		Sample
// @Accept  	json
// @Produce  	json
// @Param       request			body 	viewmodel.SampleVersionRqViewModel  true  "Sample Version"
// @Success 	200	{object} 	dto.ApiResponse[*any]
// @Failure 	500	{object} 	dto.ApiResponse[any]
// @Router 		/sample/version [post]
func (c *SampleController) SetSampleVersionInsert(ctx *gin.Context) {
	var request viewmodel.SampleVersionRqViewModel
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.Error(exception.ValidationException(strconv.Itoa(http.StatusBadRequest), err.Error()))
		return
	}

	err = c.service.SetSampleVersion(ctx, "Insert", &request)
	if err != nil {
		ctx.Error(err)
		return
	}

	resp := &dto.Response[*any]{
		ResponseCode:    strconv.Itoa(http.StatusOK),
		ResponseMessage: "Success",
		LogReff:         identifier.GetLogReff(ctx),
		TraceId:         identifier.GetTraceId(ctx),
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	Set Sample Version Update
// @Description Set Sample Version Update
// @Tags 		Sample
// @Accept  	json
// @Produce  	json
// @Param       request			body 	viewmodel.SampleVersionRqViewModel  true  "Sample Version"
// @Success 	200	{object} 	dto.ApiResponse[*any]
// @Failure 	500	{object} 	dto.ApiResponse[any]
// @Router 		/sample/version [put]
func (c *SampleController) SetSampleVersionUpdate(ctx *gin.Context) {
	var request viewmodel.SampleVersionRqViewModel
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.Error(exception.ValidationException(strconv.Itoa(http.StatusBadRequest), err.Error()))
		return
	}

	err = c.service.SetSampleVersion(ctx, "Update", &request)
	if err != nil {
		ctx.Error(err)
		return
	}

	resp := &dto.Response[*any]{
		ResponseCode:    strconv.Itoa(http.StatusOK),
		ResponseMessage: "Success",
		LogReff:         identifier.GetLogReff(ctx),
		TraceId:         identifier.GetTraceId(ctx),
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	Set Sample Version Delete
// @Description Get Sample Version Delete
// @Tags 		Sample
// @Produce  	json
// @Param       sample-id			path  	string  true	"Sample ID"
// @Param       version-number		path  	string  true	"Sample Version"
// @Success 	200	{object} 	dto.ApiResponse[*any]
// @Failure 	500	{object} 	dto.ApiResponse[any]
// @Router 		/sample/{sample-id}/version/{version-number} 	[delete]
func (c *SampleController) SetSampleVersionDelete(ctx *gin.Context) {
	sampleId := ctx.Param("sample-id")
	versionNumber := ctx.Param("version-number")

	request := &viewmodel.SampleVersionRqViewModel{SampleId: sampleId, VersionNumber: versionNumber}

	err := c.service.SetSampleVersion(ctx, "Delete", request)
	if err != nil {
		ctx.Error(err)
		return
	}

	resp := &dto.Response[*any]{
		ResponseCode:    strconv.Itoa(http.StatusOK),
		ResponseMessage: "Success",
		LogReff:         identifier.GetLogReff(ctx),
		TraceId:         identifier.GetTraceId(ctx),
	}

	ctx.JSON(http.StatusOK, resp)
}
