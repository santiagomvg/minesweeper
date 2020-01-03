package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"minesweeper/ms"
	"net/http"
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
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/static", "static")

	router.GET("/ping", func(c *gin.Context) {
		c.Writer.Write([]byte("HI!"))
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.POST("/action", func(c *gin.Context) {
		ms.HandleGameAction(c.Writer, c.Request)
	})

	router.POST("/newGame", func(c *gin.Context) {
		ms.HandleRestartAction(c.Writer, c.Request)
	})

	router.Run(":" + port)
}
