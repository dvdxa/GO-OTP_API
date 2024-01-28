package controller

import (
	"fmt"
	"otp_api/internal/app/middlewares"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Setup() *gin.Engine {
	app := gin.New()

	// Logging to a file.
	// f, _ := os.Create("log/api.log")
	// gin.DisableConsoleColor()
	// gin.DefaultWriter = io.MultiWriter(f)

	// Middlewares
	app.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		return fmt.Sprintf("%s - - [%s] \"%s %s %s %d %s \" \" %s\" \" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("02/Jan/2006:15:04:05"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	app.Use(gin.Recovery())
	app.Use(middlewares.CORS())
	app.NoRoute(middlewares.NoRouteHandler())

	h.log.Infof("[HOST] %s:%s", h.cfg.Sever.Host, h.cfg.Sever.Port)

	//Routes
	app.GET("/health", h.HealthCheck)
	app.POST("/otp", h.CreateOTP)
	app.GET("/otp/:id", h.GetOTP)
	app.PATCH("/otp/:id", h.ValidateOTP)
	return app
}
