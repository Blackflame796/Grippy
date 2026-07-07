package middlewares

import (
	"ToDoApp/internal/core/transport/http/response"
	"ToDoApp/pkg/logger"
	"fmt"
	"net/http"
)

func ZapLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		res := response.HTTPResponseHandler{ResponseWriter: w, Status: http.StatusOK, Size: 0}
		next.ServeHTTP(&res, r)
		logger.Log.Info(fmt.Sprintf("%s %s %d %s", r.Method, r.URL.Path, res.Status, http.StatusText(res.Status)))
	})
}
