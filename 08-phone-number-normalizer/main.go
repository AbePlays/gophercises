package main

import (
	"database/sql"
	"fmt"
	"strings"

	dblocal "github.com/AbePlays/gophercises/08-phone-number-normalizer/db"
	"github.com/AbePlays/gophercises/08-phone-number-normalizer/utils"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "abe"
	password = "joeydoesntsharefood"
	dbname   = "phone_numbers"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", host, port, user, password)
	db, err := sql.Open("postgres", psqlInfo)
	utils.PanicIfError(err)

	err = dblocal.CreateDb(db, dbname)
	utils.PanicIfError(err)
	db.Close()

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	utils.PanicIfError(err)

	defer db.Close()

	err = dblocal.CreatePhoneNumbersTable(db)
	utils.PanicIfError(err)

	err = dblocal.SeedDb(db)
	utils.PanicIfError(err)

	phoneNumbers, err := dblocal.GetAllPhoneNumbers(db)
	utils.PanicIfError(err)

	for _, phoneNumber := range phoneNumbers {
		normalizedPhoneNumber := normalize(phoneNumber)
		if normalizedPhoneNumber == phoneNumber {
			fmt.Printf("%s\n", normalizedPhoneNumber)
		} else {
			fmt.Printf("%s -> %s\n", phoneNumber, normalizedPhoneNumber)
			err = dblocal.UpdatePhoneNumber(db, phoneNumber, normalizedPhoneNumber)
			utils.PanicIfError(err)
		}
	}

	phoneNumbers, err = dblocal.GetAllPhoneNumbers(db)
	utils.PanicIfError(err)

	fmt.Println("Phone numbers:", phoneNumbers)
}

func normalize(phone string) string {
	var newPhone strings.Builder

	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			newPhone.WriteString(string(ch))
		}
	}

	return newPhone.String()
}
