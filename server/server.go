package server

import (
	"context"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"go.uber.org/fx"

	"github.com/choral-io/gommerce-server-core/server"
)

func NewSelectorMatcher() selector.Matcher {
	return selector.MatchFunc(func(ctx context.Context, callMeta interceptors.CallMeta) bool {
		return strings.HasPrefix(callMeta.Service, "gommerce.")
	})
}

const ServerRegistrationsTag = `group:"server/grpc.registrations"`

var serverRegistrationsAnns = []fx.Annotation{fx.As((*any)(nil)), fx.ResultTags(ServerRegistrationsTag)}

func AnnotateRegistrations(reg ...any) []any {
	res := make([]any, 0, len(reg))
	for _, r := range reg {
		res = append(res, fx.Annotate(r, serverRegistrationsAnns...))
	}
	return res
}

func ProviceRegistrations() []any {
	return AnnotateRegistrations(
		server.NewHealthServiceServer,
	)
}
