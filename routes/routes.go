package routes

import (
	"github.com/Kdsingh333/miniLink-go-backend/helper"
	"github.com/gin-gonic/gin"
)

func Routers()*gin.Engine{
	r := gin.Default();
	r.GET("/:code",helper.Redirect);
	r.POST("/shorten",helper.Shorten);
	r.POST("/custom",helper.Custom);

	return r ;
}