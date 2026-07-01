package main

import (
	"fmt"
	"log"

	_ "embed"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vinovest/sqlx"
)

//go:embed mock.sql
var mockSQL string

func main() {
	fmt.Println("Seeding mock academy database...")
	academy()

	fmt.Println("Seeding mock 'external' database...")
	external()
}

func academy() {
	db, err := sqlx.Connect("mysql", "academy:academypass@tcp(localhost:3306)/academy?parseTime=true")
	if err != nil {
		log.Fatalf("could not connect to out database: %v", err)
	}

	tables := []string{"waves", "staff", "threads", "cases", "thread_messages", "case_notes"}
	for _, table := range tables {
		if _, err := db.Exec("TRUNCATE TABLE " + table); err != nil {
			log.Fatalf("could not empty out data: %v", err)
		}
	}

	if _, err := db.Exec("ALTER TABLE waves AUTO_INCREMENT = 1;"); err != nil {
		log.Fatalf("could not reset waves auto increment: %v", err)
	}

	if _, err := db.Exec("ALTER TABLE staff AUTO_INCREMENT = 1;"); err != nil {
		log.Fatalf("could not reset staff auto increment: %v", err)
	}

	if _, err = db.Exec("INSERT INTO waves (id, state, close_at) VALUES (67, 'helper', DATE_ADD(NOW(), INTERVAL 30 DAY))"); err != nil {
		log.Fatalf("could not create wave: %v", err)
	}

	if _, err = db.Exec("INSERT INTO staff (snowflake, username, display_name, wave_id, role) VALUES ('204084691425427466', 'isaac', 'isaac nonlogged', 67, 'admin'), ('164564849915985922', 'dray', 'yarD', 67, 'helper'), ('163008912348413953', 'kieu_', 'Mik', 67, 'trainee')"); err != nil {
		log.Fatalf("could not create staff users: %v", err)
	}
}

func external() {
	db, err := sqlx.Connect("mysql", "academy:academypass@tcp(localhost:3306)/modmail2?parseTime=true&multiStatements=true")
	if err != nil {
		log.Fatalf("could not connect to mock external database: %v\n - does it exist?", err)
	}

	if _, err = db.Exec(mockSQL); err != nil {
		log.Fatal(err)
	}
}
