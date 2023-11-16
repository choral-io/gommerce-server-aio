package server

import (
	"context"
	"encoding/base64"
	"strings"
	"time"

	"github.com/choral-io/gommerce-server-aio/data/models"
	"github.com/choral-io/gommerce-server-core/secure"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

type BasicTokenStore struct {
	bdb bun.IDB
}

var _ secure.TokenStore = (*BasicTokenStore)(nil)

func NewBasicTokenStore(bdb bun.IDB) (*BasicTokenStore, error) {
	return &BasicTokenStore{
		bdb: bdb,
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

func (s *BasicTokenStore) Issue(*secure.Token, time.Duration) (string, error) {
	return "", secure.ErrUnsupportedOperation
}

func (s *BasicTokenStore) Renew(string, time.Duration) (string, error) {
	return "", secure.ErrUnsupportedOperation
}

func (s *BasicTokenStore) Verify(value string) (*secure.Token, error) {
	if username, password, err := parseBasicAuth(value); err == nil {
		client := &models.Client{}
		if err := s.bdb.NewSelect().Model(client).
			Relation("Realm", func(sq *bun.SelectQuery) *bun.SelectQuery { return sq.Column("name") }).
			Where(`secret_key = ?`, username).Scan(context.Background()); err != nil {
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
		return secure.NewToken("basic", client.Realm.Name, client.Id, client.Id, []string{}), nil
	}
	return nil, secure.ErrInvalidToken
}

func (s *BasicTokenStore) Revoke(string) (*secure.Token, error) {
	return nil, secure.ErrUnsupportedOperation
}

func NewServerAuthorizer(uts secure.TokenStore, cts *BasicTokenStore) *secure.ServerAuthorizer { // create server authorizer
	return secure.NewServerAuthorizer(map[string]secure.TokenStore{
		secure.AUTH_SCHEMA_BEARER: uts,
		secure.AUTH_SCHEMA_BASIC:  cts,
	})
}
