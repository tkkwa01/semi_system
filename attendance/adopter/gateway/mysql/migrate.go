package mysql

import (
	"bufio"
	"fmt"
	"gorm.io/gorm"
	"os"
	"semi_systems/attendance/domain"
	"semi_systems/driver"
	userdomain "semi_systems/keijiban/domain"
	"semi_systems/keijiban/domain/vobj"
	"strings"
)

func InitDatabase() {
	db, err := driver.NewDB()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}

	loadDataFromFile(db)
	loadUsersFromFile(db)
}

func loadDataFromFile(db *gorm.DB) {
	file, err := os.Open("name.txt")
	if err != nil {
		fmt.Println("Error opening name.txt:", err)
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
		fmt.Println("Error reading name.txt:", err)
	}
}

func loadUsersFromFile(db *gorm.DB) {
	file, err := os.Open("user.txt")
	if err != nil {
		fmt.Println("Error opening user.txt:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			fmt.Printf("Invalid line format: %s\n", line)
			continue
		}

		name, password := parts[0], parts[1]

		// パスワードをハッシュ化
		hashedPassword, err := vobj.NewPassword(password)
		if err != nil {
			fmt.Printf("Error hashing password for user %s: %v\n", name, err)
			continue
		}

		// 新しいユーザーオブジェクトを作成
		user := userdomain.User{
			Name:          name,
			Password:      *hashedPassword,
			RecoveryToken: vobj.NewRecoveryToken(""),
		}

		result := db.Create(&user)
		if result.Error != nil {
			fmt.Println("Error writing user to database:", result.Error)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading user.txt:", err)
	}
}
