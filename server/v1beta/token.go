package v1beta

import (
	"context"
	"errors"
	"fmt"

	"github.com/choral-io/gommerce-server-aio/data/models"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

const (
	LOGIN_PROVIDER_FORM_PASSWORD = "FORM_PASSWORD"
	LOGIN_PROVIDER_SMS_OTP_CODE  = "SMS_OTP_CODE"
)

type LoginProvider interface {
	Name() string
	Login(ctx context.Context, realm, username, password, idToken string, scope []string) (*models.Login, error)
}

type FormPasswordLoginProvider struct {
	bdb bun.IDB
}

func NewFormPasswordLoginProvider(bdb bun.IDB) LoginProvider {
	return &FormPasswordLoginProvider{bdb: bdb}
}

func (p *FormPasswordLoginProvider) Name() string {
	return LOGIN_PROVIDER_FORM_PASSWORD
}

func (p *FormPasswordLoginProvider) Login(ctx context.Context, realmId, username, password, idToken string, scope []string) (*models.Login, error) {
	var login models.Login
	if err := p.bdb.NewSelect().Model(&login).
		Where(`"login"."realm_id" = ?`, realmId).
		Where(`"login"."provider" = ?`, LOGIN_PROVIDER_FORM_PASSWORD).
		Where(`"login"."identifier" = ?`, username).
		Relation("User").Scan(ctx); err != nil {
		return nil, err
	}
	if !login.Credential.Valid {
		return nil, errors.New("password not set")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(login.Credential.String), []byte(password)); err != nil {
		return nil, errors.New("password not match")
	}
	return &login, nil
}

type SMSOTPCodeLoginProvider struct{}

func NewSMSOTPCodeLoginProvider() LoginProvider {
	return &SMSOTPCodeLoginProvider{}
}

func (p *SMSOTPCodeLoginProvider) Name() string {
	return LOGIN_PROVIDER_SMS_OTP_CODE
}

func (p *SMSOTPCodeLoginProvider) Login(ctx context.Context, realmId, username, password, idToken string, scope []string) (*models.Login, error) {
	return nil, fmt.Errorf("login provider '%s' not implemented", LOGIN_PROVIDER_SMS_OTP_CODE)
}
