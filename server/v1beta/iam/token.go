package iam_v1beta

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/choral-io/gommerce-server-aio/data/models"
	"github.com/choral-io/gommerce-server-aio/data/repos"
)

func validatePassword(login *models.Login, password string) error {
	if login == nil || !login.Credential.Valid {
		return errors.New("password not set")
	}
	if login.Disabled {
		return errors.New("password disabled")
	}
	if login.ExpiresAt.Valid && login.ExpiresAt.Time.Before(time.Now()) {
		return errors.New("password expired")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(login.Credential.String), []byte(password)); err != nil {
		return errors.New("password not match")
	}
	return nil
}

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
	var login *models.Login
	err := p.drs.RunInTx(ctx, nil, func(ctx context.Context, dr repos.DataRepos) error {
		var err error
		login, err = dr.Logins().FindByIdentifier(ctx, realmId, p.Name(), username, repos.WithRelation("User"))
		if err == sql.ErrNoRows {
			return status.Error(codes.InvalidArgument, "login not found")
		}
		if err != nil {
			return err
		}
		if err := validatePassword(login, password); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return login, nil
}

type smsOTPCodeLoginProvider struct {
	drs repos.DataRepos
}

func NewSMSOTPCodeLoginProvider(drs repos.DataRepos) LoginProvider {
	return &smsOTPCodeLoginProvider{
		drs: drs,
	}
}

func (p *smsOTPCodeLoginProvider) Name() string {
	return LoginProviderSmsOtpCode
}

func (p *smsOTPCodeLoginProvider) Login(ctx context.Context, realmId, username, password, _ string, _ []string) (*models.Login, error) {
	var login *models.Login
	err := p.drs.RunInTx(ctx, nil, func(ctx context.Context, dr repos.DataRepos) error {
		var err error
		login, err = dr.Logins().FindByIdentifier(ctx, realmId, p.Name(), username, repos.WithRelation("User"))
		if err == sql.ErrNoRows {
			return status.Error(codes.InvalidArgument, "login not found")
		}
		if err != nil {
			return err
		}
		if err := validatePassword(login, password); err != nil {
			return err
		}
		if err := dr.Logins().DisableById(ctx, login.Id); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return login, nil
}
