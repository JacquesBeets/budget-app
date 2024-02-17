package database

// func CreateTransactionTable(s service) {
// 	_, err := s.db.Exec(`
//         CREATE TABLE IF NOT EXISTS transactions (
//             id INTEGER PRIMARY KEY AUTOINCREMENT,
//             transaction_type,
//             transaction_date,
//             transaction_amount FLOAT NOT NULL,
//             transaction_id TEXT,
//             transaction_name TEXT,
//             transaction_memo TEXT,
//             created_at DATETIME NOT NULL,
//             bank_name TEXT NOT NULL,
//             transaction_type_id INTEGER,
//             budget_item_id INTEGER,
//             FOREIGN KEY (transaction_type_id) REFERENCES transaction_types(id),
//             FOREIGN KEY (budget_item_id) REFERENCES budget(id)
//         );
//     `)
// 	if err != nil {
// 		log.Fatalf(fmt.Sprintf("Error creating transactions table: %v", err))
// 	}

// 	// runMigration(s)
// }

// func runMigration(s service) {
// 	_, err := s.db.Exec(`
//         ALTER TABLE transactions
//         ADD COLUMN budget_item_id INTEGER DEFAULT 0 REFERENCES budget(id)
//     `)
// 	if err != nil {
// 		log.Fatalf(fmt.Sprintf("Error creating migrations table: %v", err))
// 	}
// }
