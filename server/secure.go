package server

import (
	"context"
	"encoding/base64"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/choral-io/gommerce-server-core/secure"

	"github.com/choral-io/gommerce-server-aio/data/repos"
)

type BasicTokenStore struct {
	drs repos.DataRepos
}

var _ secure.TokenStore = (*BasicTokenStore)(nil)

func NewBasicTokenStore(drs repos.DataRepos) (*BasicTokenStore, error) {
	return &BasicTokenStore{
		drs: drs,
	}, nil
}

func parseBasicAuth(value string) (string, string, error) {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", "", secure.ErrInvalidToken
	}
	if splits := strings.SplitN(string(decoded), ":", 2); len(splits) != 2 {
		return "", "", secure.ErrInvalidToken
	} else {
		return splits[0], splits[1], nil
	}
}

func (s *BasicTokenStore) Issue(context.Context, *secure.Token, time.Duration) (string, error) {
	return "", secure.ErrUnsupportedOperation
}

func (s *BasicTokenStore) Renew(context.Context, string, time.Duration) (string, error) {
	return "", secure.ErrUnsupportedOperation
}

func (s *BasicTokenStore) Verify(ctx context.Context, value string) (*secure.Token, error) {
	if username, password, err := parseBasicAuth(value); err == nil {
		client, err := s.drs.Clients().FindOneBySecretKey(ctx, username)
		if err != nil {
			return nil, secure.ErrInvalidToken
		}
		if client.Disabled {
			return nil, secure.ErrInvalidToken
		}
		if client.ExpiresAt.Valid && !client.ExpiresAt.Time.After(time.Now()) {
			return nil, secure.ErrInvalidToken
		}
		if !client.SecretCode.Valid {
			return nil, secure.ErrInvalidToken
		}
		if err := bcrypt.CompareHashAndPassword([]byte(client.SecretCode.String), []byte(password)); err != nil {
			return nil, secure.ErrInvalidToken
		}
		return secure.NewToken("basic", "server", client.Id, client.Id, []string{}), nil
	}
	return nil, secure.ErrInvalidToken
}

func (s *BasicTokenStore) Revoke(context.Context, string) (*secure.Token, error) {
	return nil, secure.ErrUnsupportedOperation
}

func NewServerAuthorizer(uts secure.TokenStore, cts *BasicTokenStore) *secure.ServerAuthorizer { // create server authorizer
	return secure.NewServerAuthorizer(map[string]secure.TokenStore{
		secure.AuthSchemaBearer: uts,
		secure.AuthSchemaBasic:  cts,
	})
}
