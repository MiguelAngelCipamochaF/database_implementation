package transactions

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/MiguelAngelCipamochaF/database_implementation/internal/transactions/models"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	transaction := models.Transaction{
		Codigo: "test",
	}
	myRepo := NewRepo()
	transactionResult, err := myRepo.Store(transaction)
	if err != nil {
		log.Println(err)
	}
	assert.Equal(t, transaction.Codigo, transactionResult.Codigo)
}

func TestGetByName(t *testing.T) {
	product := []models.Transaction{
		{
			ID:       4,
			Codigo:   "001",
			Moneda:   "USD",
			Monto:    4500,
			Emisor:   "Miguel Cipamocha",
			Receptor: "Miles Morales",
			Fecha:    "26-01-2022",
		},
		{
			ID:       2,
			Codigo:   "002",
			Moneda:   "USD",
			Monto:    4500,
			Emisor:   "Miles Morales",
			Receptor: "Miguel Cipamocha",
			Fecha:    "12-01-2022",
		},
	}
	myRepo := NewRepo()
	productResult, err := myRepo.GetByCode("001")
	if err != nil {
		log.Println(err)
	}
	assert.Equal(t, product[0], productResult)
}

func TestGetOneWithContext(t *testing.T) {
	// usamos un Id que exista en la DB
	id := 6
	// definimos un Transaction cuyo nombre sea igual al registro de la DB
	transaction := models.Transaction{
		Codigo: "test",
	}
	myRepo := NewRepo()
	// se define un context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	transactionResult, err := myRepo.GetOneWithcontext(ctx, id)
	fmt.Println(err)
	assert.Equal(t, transaction.Codigo, transactionResult.Codigo)
}

func TestGetAll(t *testing.T) {
	testTransactions := []models.Transaction{
		{
			ID:              4,
			Codigo:          "001",
			Moneda:          "USD",
			Monto:           4500,
			Emisor:          "Miguel Cipamocha",
			Receptor:        "Miles Morales",
			Fecha:           "26-01-2022",
			Id_warehouse:    0,
			Warehouse:       "",
			WarehouseAdress: "",
		},
		{
			ID:              6,
			Codigo:          "test",
			Moneda:          "",
			Monto:           0,
			Emisor:          "",
			Receptor:        "",
			Fecha:           "",
			Id_warehouse:    0,
			Warehouse:       "",
			WarehouseAdress: "",
		},
	}
	myRepo := NewRepo()
	transactions, err := myRepo.GetAll()
	if err != nil {
		log.Println(err)
	}

	assert.Equal(t, testTransactions, transactions)
}
