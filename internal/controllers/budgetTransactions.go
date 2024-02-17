package controllers

// func LinkTransactionToBudget(s database.Service, transactionID string, budgetID string) error {
// 	db := s.GetDBPool()

// 	query := `
//         INSERT INTO budget_transactions (transaction_id, budget_id) VALUES (?, ?)
//     `

// 	_, err := db.Exec(query, transactionID, budgetID)
// 	if err != nil {
// 		fmt.Println("Error linking transaction to budget:", err)
// 		return err
// 	}

// 	return nil
// }
