package api

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func blockActionMux() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "authorization", "X-API-Key"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

	engine.GET("/health", health)
	// TODO : swagger
	v1Group := engine.Group("/v1")
	{
		v1Group.POST("/signup", signup)
		v1Group.POST("/signin", signin)
		userGroup := v1Group.Group("/user")
		userGroup.GET("/:id", getUser)
	}

	return engine
}

func health(c *gin.Context) {
	c.Status(http.StatusOK)
}

func signin(c *gin.Context) {
	// TODO : signin
}

func signup(c *gin.Context) {
	// TODO : signup
}

func getUser(c *gin.Context) {
	// TODO : get user
}
