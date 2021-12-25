package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"strings"
	"webserver/logging"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if isBrokenPipeError(err) {
					handleBrokenPipeError(c, logger, err)
				} else {
					handleRecoveredError(c, logger, err)
				}
			}
		}()

		c.Next()
	}
}

func isBrokenPipeError(err interface{}) bool {
	ne, ok := err.(*net.OpError)
	if !ok {
		return false
	}

	if !errors.As(ne, &os.SyscallError{}) {
		return false
	}

	msg := strings.ToLower(ne.Error())
	if strings.Contains(msg, "broken pipe") || strings.Contains(msg, "connection reset by peer") {
		return true
	}

	return false
}

func handleBrokenPipeError(c *gin.Context, logger *zap.Logger, err interface{}) {
	httpRequest, _ := httputil.DumpRequest(c.Request, false)
	headers := strings.Split(string(httpRequest), "\r\n")
	for idx, header := range headers {
		current := strings.Split(header, ":")
		if current[0] == "Authorization" {
			headers[idx] = current[0] + ": *"
		}
	}

	logger.
		With(logging.ContextField(c)...).
		Error("broken pipe",
			zap.Any("error", err),
			zap.String("method", c.Request.Method),
			zap.String("URL", c.Request.URL.String()),
			zap.Strings("headers", headers),
		)

	// If the connection is dead, we can't write a status to it.
	_ = c.Error(err.(error))
	c.Abort()
}

func handleRecoveredError(c *gin.Context, logger *zap.Logger, err interface{}) {
	stack := stack(3)
	logger.
		With(logging.ContextField(c)...).
		Error("recovered",
			zap.Any("error", err),
			zap.String("method", c.Request.Method),
			zap.String("URL", c.Request.URL.String()),
			zap.String("stack", string(stack)),
		)

	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    -1,
		"message": fmt.Sprintf("%v", err),
	})
}

func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		_, _ = fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		_, _ = fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// source returns a space-trimmed slice of the nth line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptr_method
	// and want
	//	*T.ptr_method
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.ReplaceAll(name, centerDot, dot)
	return name
}
