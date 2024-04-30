package main

import (
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/api"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/controllers"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// config := configs.NewConfig()
	// appConfig, err := config.ReadConfig()
	// if err != nil {
	// 	panic(err)
	// }

	paymentController := controllers.NewPaymentController()

	api := api.NewApi(paymentController)
	api.Run(":8081")
}
