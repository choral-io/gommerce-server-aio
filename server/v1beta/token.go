package v1beta

import (
	"context"
	"errors"
	"fmt"

	"github.com/choral-io/gommerce-server-aio/data/models"
	"github.com/choral-io/gommerce-server-aio/data/repos"
	"golang.org/x/crypto/bcrypt"
)

const (
	LoginProviderFormPassword = models.LoginProviderFormPassword
	LoginProviderSmsOtpCode   = models.LoginProviderSmsOtpCode
)

type LoginProvider interface {
	Name() string
	Login(ctx context.Context, realm, username, password, idToken string, scope []string) (*models.Login, error)
}

type formPasswordLoginProvider struct {
	drs repos.DataRepos
}

func NewFormPasswordLoginProvider(drs repos.DataRepos) LoginProvider {
	return &formPasswordLoginProvider{drs: drs}
}

func (p *formPasswordLoginProvider) Name() string {
	return LoginProviderFormPassword
}

func (p *formPasswordLoginProvider) Login(ctx context.Context, realmId, username, password, _ string, _ []string) (*models.Login, error) {
	login, err := p.drs.Logins().FindOneByIdentifier(ctx, realmId, p.Name(), username, repos.WithRelation("User"))
	if err != nil {
		return nil, err
	}
	if !login.Credential.Valid {
		return nil, errors.New("password not set")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(login.Credential.String), []byte(password)); err != nil {
		return nil, errors.New("password not match")
	}
	return login, nil
}

type smsOTPCodeLoginProvider struct{}

func NewSMSOTPCodeLoginProvider() LoginProvider {
	return &smsOTPCodeLoginProvider{}
}

func (p *smsOTPCodeLoginProvider) Name() string {
	return LoginProviderSmsOtpCode
}

func (p *smsOTPCodeLoginProvider) Login(context.Context, string, string, string, string, []string) (*models.Login, error) {
	return nil, fmt.Errorf("login provider '%s' not implemented", p.Name())
}
