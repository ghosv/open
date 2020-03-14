package utils

import (
	"context"

	"github.com/micro/go-micro/v2/server"
)

// MicroHandlerContextWithValue Util
func MicroHandlerContextWithValue(k, v interface{}) server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			ctx = context.WithValue(ctx, k, v)
			return fn(ctx, req, rsp)
		}
	}
}

// MicroSubscriberContextWithValue Util
func MicroSubscriberContextWithValue(k, v interface{}) server.SubscriberWrapper {
	return func(fn server.SubscriberFunc) server.SubscriberFunc {
		return func(ctx context.Context, msg server.Message) error {
			ctx = context.WithValue(ctx, k, v)
			return fn(ctx, msg)
		}
	}
}
