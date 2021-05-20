package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type brand struct {
	Company    string
	Ltp        string
	Change     string
	Volume     string
	Buy_price  string
	Sell_price string
	Buy_qty    string
	Sell_qty   string
}

type Var struct {
	temp        string
	Result_list []brand
}

var store brand

func renderTemplate(w http.ResponseWriter, tmpl string) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, 1)
}

func query(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "query")
}
func homehandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home")
}

func jobHandler(w http.ResponseWriter, r *http.Request) {
	job := r.FormValue("fill")
	if job == "query" {
		http.Redirect(w, r, "/query/", http.StatusFound)
	} else {
		http.Redirect(w, r, "/output/", http.StatusFound)
	}
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	company_name := r.FormValue("name_action")
	database, _ := sql.Open("sqlite3", "./live_shares.db")
	query := "SELECT * FROM shares WHERE company = '" + company_name + "'"
	rows, _ := database.Query(query)
	//fmt.Println(company_name)
	//fmt.Println(rows)
	for rows.Next() {
		rows.Scan(&store.Company, &store.Ltp, &store.Change, &store.Volume, &store.Buy_price, &store.Sell_price, &store.Buy_qty, &store.Sell_qty)
	}
	fmt.Println(store)
	http.Redirect(w, r, "/queryOutput/", http.StatusFound)
}

func queryOutputHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("queryOutput.html")
	fmt.Println("Hellllloooooo")
	fmt.Println(store)
	t.Execute(w, store)
}

func outputhandler(w http.ResponseWriter, r *http.Request) {
	database, _ := sql.Open("sqlite3", "./live_shares.db")
	rows, _ := database.Query("SELECT * FROM shares")
	var data []brand
	for rows.Next() {
		var temp brand
		rows.Scan(&temp.Company, &temp.Ltp, &temp.Change, &temp.Volume, &temp.Buy_price, &temp.Sell_price, &temp.Buy_qty, &temp.Sell_qty)
		data = append(data, temp)
		fmt.Println(temp)
	}
	fmt.Println(data[1])
	var send_var Var
	send_var.Result_list = data
	t, _ := template.ParseFiles("output.html")
	t.Execute(w, send_var)
}

func main() {

	http.HandleFunc("/output/", outputhandler)
	http.HandleFunc("/home/", homehandler)
	http.HandleFunc("/jobHandler/", jobHandler)
	http.HandleFunc("/query/", query)
	http.HandleFunc("/queryHandler/", queryHandler)
	http.HandleFunc("/queryOutput/", queryOutputHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
