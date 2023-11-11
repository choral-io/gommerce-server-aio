package server

import (
	"context"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
)

func NewSelectorMatcher() selector.Matcher {
	return selector.MatchFunc(func(ctx context.Context, callMeta interceptors.CallMeta) bool {
		return strings.HasPrefix(callMeta.Service, "gommerce.")
	})
}
