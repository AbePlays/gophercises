package db

import (
	"database/sql"
	"fmt"
)

func CreateDb(db *sql.DB, name string) error {
	_, err := db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", name))
	if err != nil {
		return err
	}

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", name))
	return err
}

func CreatePhoneNumbersTable(db *sql.DB) error {
	statement := `
		CREATE TABLE IF NOT EXISTS phone_numbers (
			id SERIAL,
			value VARCHAR(255)
		)
		`
	_, err := db.Exec(statement)
	return err
}

func InsertPhoneNumber(db *sql.DB, value string) (int, error) {
	statement := `INSERT INTO phone_numbers (value) VALUES ($1) RETURNING id`
	var id int
	err := db.QueryRow(statement, value).Scan(&id)
	return id, err
}

func GetPhoneNumber(db *sql.DB, id int) (string, error) {
	statement := `SELECT value FROM phone_numbers WHERE id = $1`
	var value string
	err := db.QueryRow(statement, id).Scan(&value)
	return value, err
}

func GetAllPhoneNumbers(db *sql.DB) ([]string, error) {
	statement := `SELECT id, value FROM phone_numbers`
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var phoneNumbers []string
	for rows.Next() {
		var id int
		var value string
		if err := rows.Scan(&id, &value); err != nil {
			return nil, err
		}
		phoneNumbers = append(phoneNumbers, value)
	}
	return phoneNumbers, nil
}

func DeletePhoneNumber(db *sql.DB, id int) error {
	statement := `DELETE FROM phone_numbers WHERE id = $1`
	_, err := db.Exec(statement, id)
	return err
}

func UpdatePhoneNumber(db *sql.DB, oldValue, newValue string) error {
	statement := `UPDATE phone_numbers SET value = $1 WHERE value = $2`
	_, err := db.Exec(statement, newValue, oldValue)
	return err
}

func SeedDb(db *sql.DB) error {
	phoneNumbers := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}

	for _, number := range phoneNumbers {
		_, err := InsertPhoneNumber(db, number)
		if err != nil {
			return err
		}
	}
	return nil
}
