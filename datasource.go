package sbms

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic("configuration file load failed")
	}
}

type Recipent struct {
	Address string
	Name    string
}

type Email struct {
	From    string
	To      string
	Subject string
	Body    string
}

func GetFromDB() ([]Recipent, error) {
	var records []Recipent

	db, err := sql.Open("mysql", viper.GetString("mysql.dbname"))
	defer db.Close()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(viper.GetString("mysql.query"))
	if err != nil {
		return nil, err
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		var value string
		var record Recipent
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			switch columns[i] {
			case "email":
				record.Address = value
			case "name":
				record.Name = value
			}
		}
		records = append(records, record)
	}

	return records, nil
}

func LoadTextBody() (Email, error) {
	var email Email

	ans := make([]string, 10)
	filename := viper.GetString("mail.file")
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return email, fmt.Errorf(filename + " can't be opened")
	}
	ans = strings.Split(string(data), "\n")

	email.From = viper.GetString("mail.from")
	email.Subject = ans[0]
	email.Body = strings.Join(ans[1:], "\n")
	return email, nil
}
