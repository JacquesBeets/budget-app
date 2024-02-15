package controllers

import (
	"budget-app/internal/database"
	"budget-app/internal/models"
)

func GetTransactionsTypes(s database.Service) ([]models.TransactionType, error) {
	transactionTypes := []models.TransactionType{}

	query := `SELECT * FROM transactions_types;`

	rows, err := s.Query(query)
	if err != nil {
		return transactionTypes, err
	}
	defer rows.Close()

	for rows.Next() {
		var transactionType models.TransactionType

		err := rows.Scan(
			&transactionType.ID,
			&transactionType.Title,
			&transactionType.Category,
			&transactionType.CreatedAt,
		)
		if err != nil {
			return transactionTypes, err
		}

		transactionTypes = append(transactionTypes, transactionType)
	}

	return transactionTypes, nil
}

func CreateTransactionType(s database.Service, tnt models.TransactionType) (models.TransactionType, error) {

	query := `INSERT INTO transactions_types (title, category) VALUES (?, ?);`

	result, err := s.Exec(query, tnt.Title, tnt.Category)
	if err != nil {
		return tnt, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return tnt, err
	}
	tnt.ID = int(lastId)

	return tnt, nil
}

func UpdateTransactionType(s database.Service, trtID int) (models.TransactionType, error) {

	transactionType := models.TransactionType{}

	query := `UPDATE transactions_types SET title = ?, category = ? WHERE id = ?;`

	result, err := s.Exec(query, transactionType.Title, transactionType.Category, trtID)
	if err != nil {
		return transactionType, err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return transactionType, err
	}

	return transactionType, nil
}

func DeleteTransactionType(s database.Service, trtID int) error {

	query := `DELETE FROM transactions_types WHERE id = ?;`

	_, err := s.Exec(query, trtID)
	if err != nil {
		return err
	}

	return nil

}
