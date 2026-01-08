package authcontext

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid"
)

type key string

const (
	requestIDKey key = "x-request-id"
)

func New(request *http.Request) context.Context {
	requestID := request.Header.Get(string(requestIDKey))
	if len(requestID) == 0 {
		requestID = newRequestID()
	}

	return context.WithValue(request.Context(), requestIDKey, requestID)
}

func NewDetachedWithRequestID(ctx context.Context) context.Context {
	requestID := ctx.Value(requestIDKey)
	if requestID == nil {
		return context.Background()
	}
	return context.WithValue(context.Background(), requestIDKey, requestID)
}

func NewTestContext() context.Context {
	return context.WithValue(context.Background(), requestIDKey, newRequestID())
}

func newRequestID() string {
	var id string
	logID, err := uuid.NewV4()
	if err == nil {
		id = logID.String()
	}

	return id
}
