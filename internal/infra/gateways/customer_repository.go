package gateways

import (
	"fmt"

	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/entities"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/drivers/sql"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/gateways/sqlscripts"
)

type CustomerRepositoryGateway interface {
	FindCustomerById(id int) (entities.Customer, error)
	FindCustomerByCPF(cpf string) (entities.Customer, error)
	SaveCustomer(customer entities.Customer) error
}

type customerRepositoryGateway struct {
	sqlClient sql.SQLClient
}

func NewCustomerRepositoryGateway(sqlClient sql.SQLClient) CustomerRepositoryGateway {
	return customerRepositoryGateway{
		sqlClient: sqlClient,
	}
}

func (r customerRepositoryGateway) FindCustomerById(id int) (entities.Customer, error) {
	getCustomerByIdQuery := fmt.Sprintf(sqlscripts.GetCustomerByIdQuery)

	row := r.sqlClient.FindOne(getCustomerByIdQuery, id)

	var customer entities.Customer
	err := row.Scan(&customer.ID, &customer.Name, &customer.Cpf, &customer.Email, &customer.CreatedAt, &customer.UpdatedAt)
	if err != nil {
		return entities.Customer{}, fmt.Errorf("failed to find customer by id [%d], error %v", id, err)
	}

	return customer, nil
}

func (r customerRepositoryGateway) FindCustomerByCPF(cpf string) (entities.Customer, error) {
	getCustomerByIdQuery := fmt.Sprintf(sqlscripts.GetCustomerByCPFQuery)

	row := r.sqlClient.FindOne(getCustomerByIdQuery, cpf)

	var customer entities.Customer
	err := row.Scan(&customer.ID, &customer.Name, &customer.Cpf, &customer.Email, &customer.CreatedAt, &customer.UpdatedAt)
	if err != nil {
		return entities.Customer{}, fmt.Errorf("failed to find customer by cpf [%s], error %v", cpf, err)
	}

	return customer, nil
}

func (r customerRepositoryGateway) SaveCustomer(customer entities.Customer) error {
	insertCustomerCmd := fmt.Sprintf(sqlscripts.InsertCustomer)

	_, err := r.sqlClient.Exec(insertCustomerCmd, customer.Name, customer.Cpf, customer.Email, customer.CreatedAt, customer.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to save customer, error %v", err)
	}

	return nil
}
