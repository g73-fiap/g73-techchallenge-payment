package usecases

import (
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/entities"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/usecases/dto"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/gateways"
	"time"

	log "github.com/sirupsen/logrus"
)

type CustomerUsecase interface {
	GetCustomerById(id int) (entities.Customer, error)
	GetCustomerByCPF(cpf string) (entities.Customer, error)
	CreateCustomer(customerDTO dto.CustomerDTO) error
}

type customerUsecase struct {
	customerRepositoryGateway gateways.CustomerRepositoryGateway
}

func NewCustomerUsecase(customerRepository gateways.CustomerRepositoryGateway) CustomerUsecase {
	return customerUsecase{
		customerRepositoryGateway: customerRepository,
	}
}

func (u customerUsecase) CreateCustomer(customerDTO dto.CustomerDTO) error {
	customer := customerDTO.ToCustomer()
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()

	err := u.customerRepositoryGateway.SaveCustomer(customer)
	if err != nil {
		log.Errorf("failed to save customer, error: %v", err)
		return err
	}

	return nil
}

func (u customerUsecase) GetCustomerById(id int) (entities.Customer, error) {
	customer, err := u.customerRepositoryGateway.FindCustomerById(id)
	if err != nil {
		log.Errorf("failed to get customer by id [%d], error: %v", id, err)
		return entities.Customer{}, err
	}

	return customer, nil
}

func (u customerUsecase) GetCustomerByCPF(cpf string) (entities.Customer, error) {
	customer, err := u.customerRepositoryGateway.FindCustomerByCPF(cpf)
	if err != nil {
		log.Errorf("failed to get customer by cpf [%s], error: %v", cpf, err)
		return entities.Customer{}, err
	}

	return customer, nil
}
