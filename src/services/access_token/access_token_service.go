package access_token

import (
	"github.com/harshasavanth/bookstore_oauth-api/src/domain/access_token"
	"github.com/harshasavanth/bookstore_oauth-api/src/repository/db"
	"github.com/harshasavanth/bookstore_oauth-api/src/repository/rest"
	"github.com/harshasavanth/bookstore_oauth-api/src/utils/errors"
	"strings"
)

type Service interface {
	GetById(string) (*access_token.AccessToken, *errors.RestErr)
	Create(access_token.AccessTokenRequest) (*access_token.AccessToken, *errors.RestErr)
	UpdateExpirationTime(access_token.AccessToken) *errors.RestErr
}

type service struct {
	restUserRepo rest.RestUsersRepository
	dbRepo       db.DBRepository
}

func NewService(userRepo rest.RestUsersRepository, dbRepo db.DBRepository) Service {
	return &service{
		restUserRepo: userRepo,
		dbRepo:       dbRepo,
	}
}

func (s *service) GetById(accessTokenId string) (*access_token.AccessToken, *errors.RestErr) {
	accessTokenId = strings.TrimSpace(accessTokenId)
	if len(accessTokenId) == 0 {
		return nil, errors.NewBadRequestError("invalid access token id")
	}
	accessToken, err := s.dbRepo.GetById(accessTokenId)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

func (s *service) Create(request access_token.AccessTokenRequest) (*access_token.AccessToken, *errors.RestErr) {
	user, err := s.restUserRepo.LoginUser(request.Email, request.Password)
	if err != nil {
		return nil, err
	}
	at := access_token.GetNewAccessToken(user.Id)
	if !at.IsExpired() {
		return nil, errors.NewBadRequestError("token is still not expired")
	}
	at.Generate()
	if err := s.dbRepo.Create(at); err != nil {
		return nil, err
	}
	return &at, nil
}

func (s *service) UpdateExpirationTime(at access_token.AccessToken) *errors.RestErr {
	if err := at.Validate(); err != nil {
		return err
	}
	return s.dbRepo.UpdateExpirationTime(at)
}
