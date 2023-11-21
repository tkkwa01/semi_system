package mysql

import (
	"bufio"
	"fmt"
	"gorm.io/gorm"
	"os"
	"semi_systems/attendance/domain"
	"semi_systems/driver"
)

func InitDatabase() {
	db, err := database.NewDB()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}

	loadDataFromFile(db)
}

func loadDataFromFile(db *gorm.DB) {
	file, err := os.Open("names.txt")
	if err != nil {
		fmt.Println("Error opening names.txt:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		name := scanner.Text()

		attendance := &domain.Attendance{Name: name, Status: false}
		result := db.Create(attendance)
		if result.Error != nil {
			fmt.Println("Error writing to database:", result.Error)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading names.txt:", err)
	}
}
