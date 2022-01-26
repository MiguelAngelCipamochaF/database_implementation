package transactions

import (
	"context"
	"database/sql"
	"log"

	"github.com/MiguelAngelCipamochaF/database_implementation/db"
	"github.com/MiguelAngelCipamochaF/database_implementation/internal/transactions/models"
)

type Repository interface {
	Store(transaction models.Transaction) (models.Transaction, error)
	GetOne(id int) models.Transaction
	Update(transaction models.Transaction) (models.Transaction, error)
	GetAll() ([]models.Transaction, error)
	Delete(id int) error
	GetByCode(code string) (models.Transaction, error)
	GetFullData(id int) models.Transaction
	GetOneWithcontext(ctx context.Context, id int) (models.Transaction, error)
	UpdateWithContext(ctx context.Context, transaction models.Transaction) (models.Transaction, error)
}

type repository struct{}

func NewRepo() Repository {
	return &repository{}
}

func (r *repository) Store(transaction models.Transaction) (models.Transaction, error) {
	db := db.StorageDB
	storeTransaction := "INSERT INTO transactions(code, currency, amount, transmitter, receptor, date) VALUES( ?, ?, ?, ?, ?, ? )"
	stmt, err := db.Prepare(storeTransaction) // se prepara el SQL
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close() // se cierra la sentencia al terminar. Si quedan abiertas se genera consumos de memoria
	var result sql.Result
	result, err = stmt.Exec(transaction.Codigo, transaction.Moneda, transaction.Monto, transaction.Emisor, transaction.Receptor, transaction.Fecha) // retorna un sql.Result y un error
	if err != nil {
		return models.Transaction{}, err
	}
	insertedId, _ := result.LastInsertId() // del sql.Resul devuelto en la ejecución obtenemos el Id insertado
	transaction.ID = int(insertedId)
	return transaction, nil
}

func (r *repository) GetAll() ([]models.Transaction, error) {
	var transactions []models.Transaction
	db := db.StorageDB
	getAll := "SELECT t.id, t.code, t.currency, t.amount, t.transmitter, t.receptor, t.date FROM transactions t"
	rows, err := db.Query(getAll)
	if err != nil {
		return transactions, err
	}

	for rows.Next() {
		var transaction models.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.Codigo, &transaction.Moneda, &transaction.Monto,
			&transaction.Emisor, &transaction.Receptor, &transaction.Fecha); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (r *repository) GetOne(id int) models.Transaction {
	var transaction models.Transaction
	db := db.StorageDB
	getOne := "SELECT t.id, t.code, t.currency, t.amount, t.transmitter, t.receptor, t.date FROM transactions t where id = ?"
	rows, err := db.Query(getOne, id)
	if err != nil {
		log.Println(err)
		return transaction
	}
	for rows.Next() {
		if err := rows.Scan(&transaction.ID, &transaction.Codigo, &transaction.Moneda, &transaction.Monto, &transaction.Emisor,
			&transaction.Receptor, &transaction.Fecha); err != nil {
			log.Println(err.Error())
			return transaction
		}
	}
	return transaction
}

func (r *repository) Update(transaction models.Transaction) (models.Transaction, error) {
	db := db.StorageDB // se inicializa la base
	updateTransaction := "UPDATE transactions SET code = ?, currency = ?, amount = ?, transmitter = ?, receptor = ?, date = ? WHERE id = ?"
	stmt, err := db.Prepare(updateTransaction)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close() // se cierra la sentencia al terminar. Si quedan abiertas se genera consumos de memoria
	_, err = stmt.Exec(transaction.ID, transaction.Codigo, transaction.Moneda, transaction.Monto, transaction.Emisor, transaction.Receptor, transaction.Fecha)
	if err != nil {
		return models.Transaction{}, err
	}
	return transaction, nil
}

func (r *repository) Delete(id int) error {
	db := db.StorageDB
	deleteTransaction := "DELETE FROM transactions where id = ?"
	stmt, err := db.Prepare(deleteTransaction)
	if err != nil {
		return err
	}
	defer stmt.Close() // se cierra la sentencia al terminar.

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetByCode(code string) (models.Transaction, error) {
	var transaction models.Transaction
	db := db.StorageDB
	getByCode := "SELECT t.id, t.code, t.currency, t.amount, t.transmitter, t.receptor, t.date FROM transactions t WHERE code = ?"
	rows, err := db.Query(getByCode, code)
	if err != nil {
		return transaction, err
	}
	for rows.Next() {
		if err := rows.Scan(&transaction.ID, &transaction.Codigo, &transaction.Moneda, &transaction.Monto, &transaction.Emisor,
			&transaction.Receptor, &transaction.Fecha); err != nil {
			log.Println(err.Error())
			return transaction, err
		}
	}
	return transaction, nil
}

func (r *repository) GetFullData(id int) models.Transaction {
	var transaction models.Transaction
	db := db.StorageDB
	innerJoin := "SELECT transactions.id, transactions.code, transactions.currency, transactions.amount, transactions.transmitter, transactions.receptor, transactions.date, warehouses.name, warehouses.adress " +
		"FROM transactions INNER JOIN warehouses ON transactions.id_warehouse = warehouses.id " +
		"WHERE transactions.id = ?"
	rows, err := db.Query(innerJoin, id)
	if err != nil {
		log.Println(err)
		return transaction
	}
	for rows.Next() {
		if err := rows.Scan(&transaction.ID, &transaction.Codigo, &transaction.Moneda, &transaction.Monto, &transaction.Emisor, &transaction.Receptor,
			&transaction.Fecha, &transaction.Warehouse, &transaction.WarehouseAdress); err != nil {
			log.Fatal(err)
			return transaction
		}
	}
	return transaction
}

func (r *repository) GetOneWithcontext(ctx context.Context, id int) (models.Transaction, error) {
	var transaction models.Transaction
	db := db.StorageDB
	// se especifican estrictamente los campos necesarios en la query
	getQuery := "SELECT t.id, t.code, t.currency, t.amount, t.transmitter, t.receptor, t.date FROM transactions t WHERE t.id = ?"
	// ya no se usa db.Query sino db.QueryContext
	rows, err := db.QueryContext(ctx, getQuery, id)
	if err != nil {
		log.Println(err)
		return transaction, err
	}
	for rows.Next() {
		if err := rows.Scan(&transaction.ID, &transaction.Codigo, &transaction.Moneda, &transaction.Monto, &transaction.Emisor, &transaction.Receptor,
			&transaction.Fecha); err != nil {
			log.Fatal(err)
			return transaction, err
		}
	}
	return transaction, nil
}

func (r *repository) UpdateWithContext(ctx context.Context, transaction models.Transaction) (models.Transaction, error) {
	db := db.StorageDB
	updateTransaction := "UPDATE transactions SET code = ?, currency = ?, amount = ?, transmitter = ?, receptor = ?, date = ? WHERE id = ?"
	stmt, err := db.Prepare(updateTransaction)
	if err != nil {
		return models.Transaction{}, err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, transaction.ID, transaction.Codigo, transaction.Moneda, transaction.Monto, transaction.Emisor, transaction.Receptor, transaction.Fecha)
	if err != nil {
		return models.Transaction{}, err
	}
	return transaction, nil
}
