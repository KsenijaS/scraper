package main

import (
	"database/sql"
	"github.com/KsenijaS/scraper"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/smtp"
	"regexp"
	"strconv"
	"time"
)

func compareAndUpdatePrice(last_price float32, price float32, id int, db *sql.DB) bool {
	sqlStatement := `UPDATE urls SET last_price=? where id=?`
	_, err := db.Exec(sqlStatement, price, id)
	if err != nil {
		log.Fatal(err)
	}

	if price < last_price {
		return true
	}

	return false
}

func send(to string, url string) {
	from := "...@gmail.com"
	pass := "..."

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Hello \n\n" +
		"There is a sale, please visit " + url

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}

func main() {
	for {
		time.Sleep(20 * time.Second)
		db, err := sql.Open("mysql", "ksenija:tajna@/Coupons")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		rows, err := db.Query(`SELECT id, last_price, url, user_id from urls`)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var last_price, new_price float32
			var url, username string
			var user_id, id int
			var notify bool

			if err = rows.Scan(&id, &last_price, &url, &user_id); err != nil {
				log.Fatal(err)
			}

			strPrice, err := scraper.ParseUrl(url)
			if err != nil {
				log.Fatal(err)
			}

			re := regexp.MustCompile("[0-9]+[/.]*[0-9]*")
			strNum := re.FindString(strPrice)
			value, err := strconv.ParseFloat(strNum, 32)
			if err != nil {
				log.Fatal(err)
			}

			new_price = float32(value)

			if new_price != last_price {
				notify = compareAndUpdatePrice(last_price, new_price, id, db)
			}

			if notify {
				err = db.QueryRow(`SELECT username where id=?`, user_id).Scan(&username)
				send(username, url)
			}
		}
	}
}
