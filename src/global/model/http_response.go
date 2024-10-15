package model

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HttpResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (response HttpResponse) Ok(Data interface{}) HttpResponse {
	response.Code = http.StatusOK
	response.Msg = "ok"
	response.Data = Data
	return response
}

func (response HttpResponse) OkGin(c *gin.Context, Data interface{}) {
	c.JSON(http.StatusOK, response.Ok(Data))
}

func (response HttpResponse) Fail(Data interface{}) HttpResponse {
	response.Code = http.StatusInternalServerError
	response.Msg = "error"
	response.Data = Data
	return response
}

func (response HttpResponse) FailGin(c *gin.Context, Data interface{}) {
	c.JSON(http.StatusInternalServerError, response.Fail(Data))
}
