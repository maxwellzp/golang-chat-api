package httpx

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func ParseInt64Param(r *http.Request, name string) (int64, error) {
	raw := chi.URLParam(r, name)
	return strconv.ParseInt(raw, 10, 64)
}
