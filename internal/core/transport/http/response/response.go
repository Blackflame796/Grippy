package response

import (
	"net/http"
)

type HTTPResponseHandler struct {
	http.ResponseWriter
	Status int
	Size   int
}

func (r *HTTPResponseHandler) WriteHeader(statusCode int) {
	r.Status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *HTTPResponseHandler) Write(data []byte) (int, error) {
	size, err := r.ResponseWriter.Write(data)
	r.Size += size
	return size, err
}
