package httpx

import (
	"context"
	"errors"
	"github.com/maxwellzp/golang-chat-api/internal/contextkey"
)

func GetUserID(ctx context.Context) (int64, error) {
	val := ctx.Value(contextkey.UserID)
	id, ok := val.(int64)
	if !ok {
		return 0, errors.New("user_id not found in context")
	}
	return id, nil
}
