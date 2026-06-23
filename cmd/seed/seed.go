package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vinovest/sqlx"
)

func main() {
	fmt.Println("Seeding dev database...")

	db, err := sqlx.Connect("mysql", "academy:academypass@tcp(localhost:3306)/academy?parseTime=true")
	if err != nil {
		log.Fatalf("could not connect to out database: %v", err)
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

	if _, err = db.Exec("INSERT INTO staff (snowflake, username, display_name, wave_id, role) VALUES ('204084691425427466', 'isaac', 'isaac nonlogged', 67, 'admin'), ('164564849915985922', 'dray', 'yarD', 67, 'helper')"); err != nil {
		log.Fatalf("could not create staff users: %v", err)
	}
}
