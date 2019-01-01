package main

import (
	"database/sql"
	"encoding/json"
	"github.com/KsenijaS/scraper"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

var url string

type myData struct {
	Url   string
	Email string
}

type myUser struct {
	Email string
}

func placeDataUrls(data myData, price float32) {
	var id int

	db, err := sql.Open("mysql", "ksenija:tajna@/Coupons")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.QueryRow(`SELECT id FROM users WHERE username=?`, data.Email).Scan(&id)
	if err != nil {
		log.Fatal(err)
	}

	sqlStatement := `INSERT INTO urls(last_price, url, user_id) VALUES (?, ?, ?)`
	_, err = db.Exec(sqlStatement, price, data.Url, id)
	if err != nil {
		log.Fatalf("sql %v", err)
	}

}

func placeDataUsers(user myUser) {
	db, err := sql.Open("mysql", "ksenija:tajna@/Coupons")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO users(username) VALUES (?)`
	_, err = db.Exec(sqlStatement, user.Email)
	if err != nil {
		log.Fatalf("sql %v", err)
	}
}

func urls(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		decoder := json.NewDecoder(req.Body)

		var data myData
		err := decoder.Decode(&data)
		if err != nil {
			log.Printf("Error decoding data: %s", err)
			return
		}
		log.Println(data)

		strPrice, err := scraper.ParseUrl(data.Url)
		if err != nil {
			log.Fatal(err)
		}

		re := regexp.MustCompile("[0-9]+[/.]*[0-9]*")
		strNum := re.FindString(strPrice)
		log.Println(strNum)
		value, err := strconv.ParseFloat(strNum, 32)
		if err != nil {
			log.Fatal(err)
		}
		floatPrice := float32(value)

		placeDataUrls(data, floatPrice)
		return
	}
}

func users(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		var user myUser

		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&user)
		if err != nil {
			log.Printf("Error decoding data: %s", err)
			return
		}
		log.Println(user)

		placeDataUsers(user)
		return
	}
}

func main() {
	log.Println("Start...")
	http.HandleFunc("/urls", urls)
	http.HandleFunc("/users", users)
	http.ListenAndServe(":8080", nil)
}
