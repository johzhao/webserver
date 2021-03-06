package transport

import (
	"context"
	"net/http"
	"webserver/controller"
	"webserver/router"
	"webserver/router/encoder"
	"webserver/server"

	"go.uber.org/zap"
)

func SetupRouters(logger *zap.Logger, server server.Server, userController controller.User) {
	routers := make([]router.Router, 0)
	routers = append(routers, userRouters(logger, userController)...)

	routers = append(routers, router.NewCustomRouter(http.MethodGet, "/ping", nil, PingHandler, encoder.NewJSONResponseEncoder(logger)))

	for _, serviceRouter := range routers {
		server.HandleRouter(serviceRouter)
	}
}

//goland:noinspection GoUnusedParameter
func PingHandler(ctx context.Context, req interface{}) (interface{}, error) {
	return map[string]string{
		"message": "pong",
	}, nil
}
