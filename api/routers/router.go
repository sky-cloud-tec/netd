package routers

import (
	"net/http"

	"github.com/sky-cloud-tec/netd/api/controllers"
	//swaggerFiles "github.com/swaggo/files"
	//"github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

// SetupRouter create gin router and return
func SetupRouter(addr string) *gin.Engine {
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong!")
	})
	r.Use(cors())

	r.POST("/api/operator/hotfix", controllers.OperatorHotfix)
	r.POST("/api/operator/dump", controllers.OperatorDump)

	return r
}

// Cors 跨域设置
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "false")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
