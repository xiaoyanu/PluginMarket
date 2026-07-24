package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

type PageData struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

func OK(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{Code: 200, Msg: msg})
}

func OKData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "成功", Data: data})
}

func OKDataMsg(c *gin.Context, data interface{}, msg string) {
	c.JSON(http.StatusOK, Response{Code: 200, Msg: msg, Data: data})
}

func OKMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{Code: 200, Msg: msg})
}

func OKPage(c *gin.Context, list interface{}, total int64) {
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "成功",
		Data: PageData{List: list, Total: total},
	})
}

func Fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{Code: code, Msg: msg})
}

func BadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{Code: 400, Msg: msg})
}

func Unauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{Code: 401, Msg: msg})
}

func Forbidden(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{Code: 403, Msg: msg})
}

func NotFound(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{Code: 404, Msg: msg})
}

func ServerError(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{Code: 500, Msg: msg})
}
