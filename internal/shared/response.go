package shared

import "github.com/gin-gonic/gin"

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Error   any    `json:"error"`
}

func OK(c *gin.Context, message string, data any) {
	c.JSON(200, Response{Success: true, Message: message, Data: data, Error: nil})
}

func Created(c *gin.Context, message string, data any) {
	c.JSON(201, Response{Success: true, Message: message, Data: data, Error: nil})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, Response{Success: false, Message: message, Data: nil, Error: message})
}