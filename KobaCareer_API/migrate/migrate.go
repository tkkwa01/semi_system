package migrate

import (
	"fmt"
	"log"
	"semi_systems/KobaCareer_API/db"
	"semi_systems/KobaCareer_API/domain"
)

func init() {
	fmt.Println("Starting migration...")
	dbConn := db.NewDB()
	defer fmt.Println("Successfully Migrated")
	defer db.CloseDB(dbConn)
	err := dbConn.AutoMigrate(&domain.Internships{})
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	} else {
		log.Println("Migration completed successfully")
	}
}
