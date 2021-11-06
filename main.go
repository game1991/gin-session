package main

import (
	"demo/router"

	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()

	router.API(g)

	g.Run(":9099")
}
