package usecases

import (
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/usecases/dto"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/drivers/authorizer"
	log "github.com/sirupsen/logrus"
)

type AuthorizerUsecase interface {
	AuthorizeUser(cpf string) (dto.AuthorizerResponse, error)
}

type authorizerUsecase struct {
	authorizer authorizer.Authorizer
}

func NewAuthorizerUsecase(authorizer authorizer.Authorizer) AuthorizerUsecase {
	return authorizerUsecase{
		authorizer: authorizer,
	}
}

func (u authorizerUsecase) AuthorizeUser(cpf string) (dto.AuthorizerResponse, error) {
	authorizerResponse, err := u.authorizer.AuthorizeUser(cpf)
	if err != nil {
		log.Errorf("failed to authorize user, error: %v", err)
		return dto.AuthorizerResponse{}, err
	}

	return authorizerResponse, nil
}
