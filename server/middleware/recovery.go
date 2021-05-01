package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"os"
	"strings"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				//if logger != nil {
				//	stack := stack(3)
				//	httpRequest, _ := httputil.DumpRequest(c.Request, false)
				//	headers := strings.Split(string(httpRequest), "\r\n")
				//	for idx, header := range headers {
				//		current := strings.Split(header, ":")
				//		if current[0] == "Authorization" {
				//			headers[idx] = current[0] + ": *"
				//		}
				//	}
				//	headersToStr := strings.Join(headers, "\r\n")
				//	if brokenPipe {
				//		logger.Printf("%s\n%s%s", err, headersToStr, reset)
				//	} else if IsDebugging() {
				//		logger.Printf("[Recovery] %s panic recovered:\n%s\n%s\n%s%s",
				//			timeFormat(time.Now()), headersToStr, err, stack, reset)
				//	} else {
				//		logger.Printf("[Recovery] %s panic recovered:\n%s\n%s%s",
				//			timeFormat(time.Now()), err, stack, reset)
				//	}
				//}
				if brokenPipe {
					// If the connection is dead, we can't write a status to it.
					_ = c.Error(err.(error))
					c.Abort()
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{
						"code":    -1,
						"message": fmt.Sprintf("%v", err),
						"data":    nil,
					})
				}
			}
		}()
		c.Next()
	}
}
