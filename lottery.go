package main

import (
	"./api"
	"./sockects"
	"github.com/gin-gonic/gin"
	"net/http"
)


func main() {

	router := gin.Default()

	//api路由
	router.Group("/api/count").GET("/", api.GetCount)
	router.Group("/api/getAwards").GET("/", api.GetAwards)
	router.Group("/api/init").GET("/", api.InitData)
	router.Group("/api/getNextAction").GET("/", api.GetNextAction)
	router.Group("/api/ndraw").GET("/", api.NDraw)
	router.Group("/api/exdraw").GET("/", api.ExDraw)
	router.Group("/api/pooldraw").GET("/", api.PoolDraw)
	router.Group("/api/addMoney").GET("/", api.AddPoolMoney)
	//web socket 路由
	router.GET("/ws", func(c *gin.Context) { sockects.WebSocketHandler(c.Writer, c.Request) })

	//html页面路由
	router.LoadHTMLGlob("views/*")
	router.Group("/view/").GET("/:name", func(c *gin.Context) { c.HTML(http.StatusOK, c.Param("name")+".html", gin.H{}) })

	////静态文件路由
	router.Static("/wwwroot", "./wwwroot")

	router.Run(":12345")
}
