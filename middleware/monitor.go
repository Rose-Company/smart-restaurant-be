package middleware

import (
	"app-noti/server"

	"github.com/gin-gonic/gin"
)

func Monitoring(sc server.ServerContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		//timer := sc.GetRequestDurationMetric().ObserveDuration(c, c.Request.URL.Path, c.Request.Method)
		//defer func() {
		//	timer.ObserveDuration()
		//}()

		c.Next()

		//go func() {
		//	counter := sc.GetRequestCounterMetric()
		//	counter.Add(c, c.Request.URL.Path, c.Request.Method, fmt.Sprintf("%v", c.Writer.Status()), 1)
		//}()

	}
}
