package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Store struct {
	id            int64
	name          string
	address       string
	shipping_cost float32
}

func getStoresWithShippingCostLowerThan(conn *sql.DB, price float32) ([]Store, error) {
	var stores []Store

	rows, err := conn.Query("SELECT id, name, address, shipping_cost FROM store WHERE shipping_cost < ?", price)
	if err != nil {
		return nil, fmt.Errorf("Error while querying stores with price: %f: %v\n", price, err)
	}
	defer rows.Close()

	for rows.Next() {
		var st Store
		if err := rows.Scan(&st.id, &st.name, &st.address, &st.shipping_cost); err != nil {
			return nil, fmt.Errorf("Error while iterating over Stores: %v\n", err)
		}
		stores = append(stores, st)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error on rows next: %v\n", err)
	}
	return stores, nil
}

func getStoreLowestShippingPrice(conn *sql.DB) (Store, error) {
	var store Store

	row := conn.QueryRow(`SELECT id, name, address, shipping_cost FROM store
					ORDER BY shipping_cost ASC LIMIT 1;`)

	if err := row.Scan(&store.id, &store.name, &store.address, &store.shipping_cost); err != nil {
		if err == sql.ErrNoRows {
			return store, fmt.Errorf("Cannot find lowest price shipping: %v", err)
		}
		return store, fmt.Errorf("Error while searching for lowest price shipping %v", err)
	}
	return store, nil
}

func addNewStore(conn *sql.DB, new Store) (int64, error) {
	result, err := conn.Exec("INSERT INTO store (name, address, shipping_cost) VALUES (?, ?, ?)",
		new.name, new.address, new.shipping_cost)
	if err != nil {
		return 0, fmt.Errorf("Cannot add new store: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Cannot get id of inserted row: %v", err)
	}
	return id, nil
}

func getRandomProductName() string {
	dummyNames := []string{"STM", "Atmega", "PIC"}
	dummyCodes := []string{"32", "16", "328"}
	return dummyNames[rand.Intn(len(dummyNames))] +
		dummyCodes[rand.Intn(len(dummyCodes))]
}

func getRandomPackage() string {
	packages := []string{"LQFP", "TQFP", "DIP", "QFP"}
	legs := []string{"16", "32", "64", "128"}
	return packages[rand.Intn(len(packages))] + legs[rand.Intn(len(legs))]
}

func addRandomComponents(conn *sql.DB, num int) ([]int64, error) {
	var ids []int64

	trx, err := conn.Begin()
	if err != nil {
		return nil, fmt.Errorf("Cannot create transaction: %v", err)
	}
	defer trx.Rollback()

	for i := 0; i < num; i++ {
		name := getRandomProductName()
		pack := getRandomPackage()
		res, err := trx.Exec(`INSERT INTO component (name, description, package)
					VALUES (?, ?, ?)`, name, "", pack)
		if err != nil {
			return ids, fmt.Errorf("Cannot insert component %q in package %q: %v\n", name, pack, err)
		}

		id, err := res.LastInsertId()
		if err != nil {
			return ids, fmt.Errorf("Cannot get id for %q in package %q: %v\n", name, pack, err)
		}
		ids = append(ids, id)
	}

	err = trx.Commit()
	if err != nil {
		fmt.Errorf("Error while commiting transaction: %v", err)
	}

	return ids, nil
}

func Init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "db",
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatalf("Cannot connect to database: %s", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Cannot ping the database: %s", err)
	}
	fmt.Printf("Connected to %q\n", cfg.DBName)

	// id, err := addNewStore(db, Store{
	// 	name:          "barion",
	// 	address:       "www.barion.com.pl",
	// 	shipping_cost: 8.99,
	// })
	// if err != nil {
	// 	log.Println(err)
	// }

	// if err == nil {
	// 	fmt.Printf("Succesfully added store with ID: %d\n", id)
	// }

	stores, err := getStoresWithShippingCostLowerThan(db, 25.0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(stores)

	store, err := getStoreLowestShippingPrice(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(store)

	ids, err := addRandomComponents(db, 4)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Components added with ids: %v", ids)

}
