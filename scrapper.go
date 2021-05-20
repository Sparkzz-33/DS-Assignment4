package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gocolly/colly"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	database, _ := sql.Open("sqlite3", "./live_shares.db")
	for {
		delete_stmt, _ := database.Prepare("DROP TABLE shares")
		delete_stmt.Exec()

		create_stmt, _ := database.Prepare("CREATE TABLE shares (company TEXT, ltp TEXT, change TEXT, volume TEXT, buy_price TEXT, sell_price TEXT, buy_qty TEXT, sell_qty TEXT)")
		create_stmt.Exec()
		c := colly.NewCollector()
		var store [30][8]string
		count := 0
		i := 0
		j := 0
		c.OnHTML("tr", func(e *colly.HTMLElement) {
			j = 0
			e.ForEach("td", func(_ int, e2 *colly.HTMLElement) {
				//fmt.Println([]string{e2.Text}[0])
				if i < 30 && count > 1 {
					i = count - 2
					store[i][j] = []string{e2.Text}[0]
					//fmt.Println(store[i][j])
					//fmt.Println(j)
					j = j + 1
				}
			})
			i = i + 1
			count = count + 1
			//fmt.Println("---------------------------")

		})
		c.Visit("https://www.moneycontrol.com/markets/indian-indices/")

		for itr := 0; itr < 30; itr++ {
			company := store[itr][0]
			ltp := store[itr][1]
			change := store[itr][2]
			volume := store[itr][3]
			buy_price := store[itr][4]
			sell_price := store[itr][5]
			buy_qty := store[itr][6]
			sell_qty := store[itr][7]
			//fmt.Println(store[itr][0])
			insert_stmt, _ := database.Prepare("INSERT INTO shares (company, ltp, change, volume, buy_price, sell_price, buy_qty, sell_qty) VALUES(?,?,?,?,?,?,?,?)")
			insert_stmt.Exec(company, ltp, change, volume, buy_price, sell_price, buy_qty, sell_qty)

		}
		fmt.Println("Scraping Completed")
		time.Sleep(200 * time.Second)
	}

}
