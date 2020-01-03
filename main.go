package main

import (
	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/ms"
	_ "github.com/heroku/x/hmetrics/onload"
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	ms.Init()

	router := gin.New()
	router.Use(gin.Logger())
	//	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/ping", func(c *gin.Context) {
		c.Writer.Write([]byte("HI!"))
	})

	router.POST("/action", func(c *gin.Context) {
		ms.HandleGameAction(c.Writer, c.Request)
	})

	router.POST("/newGame", func(c *gin.Context) {
		ms.HandleRestartAction(c.Writer, c.Request)
	})

	router.Run(":" + port)
}
