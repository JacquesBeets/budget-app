package database

// func CreateBudgetTable(s service) {
// 	_, err := s.db.Exec(`
//         CREATE TABLE IF NOT EXISTS budget (
//             id INTEGER PRIMARY KEY AUTOINCREMENT,
//             name TEXT NOT NULL,
//             amount FLOAT NOT NULL,
//             created_at DATETIME NOT NULL,
//             transaction_type_id INTEGER
//         );
//     `)
// 	if err != nil {
// 		log.Fatalf(fmt.Sprintf("Error creating transactions table: %v", err))
// 	}
// }

// func (s *service) SaveBudgetItem(budget models.Budget) error {
// 	_, err := s.db.Exec(`
//         INSERT INTO budget (name, amount, created_at, transaction_type_id)
//         VALUES (?, ?, ?, ?);
//     `, budget.Name, budget.Amount, time.Now(), budget.TransactionTypeID)
// 	if err != nil {
// 		log.Fatalf(fmt.Sprintf("Error saving budget item: %v", err))
// 		return err
// 	}
// 	return nil
// }
