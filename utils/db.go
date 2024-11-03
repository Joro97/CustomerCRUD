package utils

import "database/sql"

func GetLocalDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./customers.db")
	if err != nil {
		return nil, err
	}

	createTableSQL := `
        CREATE TABLE IF NOT EXISTS customers (
            id UUID PRIMARY KEY,
            first_name TEXT NOT NULL,
            middle_name TEXT,
            last_name TEXT NOT NULL,
            email TEXT NOT NULL UNIQUE,
            phone_number TEXT
        );
        `

	if _, err = db.Exec(createTableSQL); err != nil {
		return nil, err
	}
	return db, nil
}
