package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"log"
	"minesweeper/ms"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	//	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/ping", func(c *gin.Context) {
		c.Writer.Write([]byte("HI!"))
	})

	router.GET("/action", func(c *gin.Context) {
		ms.HandleGameAction(c.Writer, c.Request)
	})

	router.POST("/newGame", func(c *gin.Context) {
		ms.HandleRestartAction(c.Writer, c.Request)
	})

	router.Run(":" + port)
}
