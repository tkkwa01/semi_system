// main.go

package main

import (
	"bufio"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
	"semi_systems/kobasemi_attendance_API/api/handler"
	"semi_systems/kobasemi_attendance_API/database"
	"semi_systems/kobasemi_attendance_API/domain"
	"semi_systems/kobasemi_attendance_API/repository"
	"semi_systems/kobasemi_attendance_API/usecase"
)

func main() {
	db, err := database.NewDB()
	if err != nil {
		panic(err)
	}
	repo := repository.NewAttendanceRepository(db)
	uc := usecase.NewAttendanceUsecase(repo)
	handler := handler.NewAttendanceHandler(uc)

	// テキストファイルを開く
	file, err := os.Open("names.txt")
	if err != nil {
		fmt.Println("Error opening names.txt:", err)
		return
	}
	defer file.Close()

	// ファイルから行を一つずつ読み込む
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		name := scanner.Text()

		// データベースに書き込む
		attendance := &domain.Attendance{Name: name, Status: false}
		result := db.Create(attendance)
		if result.Error != nil {
			fmt.Println("Error writing to database:", result.Error)
		}
	}

	if scanner.Err() != nil {
		fmt.Println("Error reading names.txt:", scanner.Err())
		return
	}

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://192.168.50.88"}
	config.AllowMethods = []string{"GET", "POST"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "ngrok-skip-browser-warning"} // この行を追加

	r.Use(cors.New(config))

	r.POST("/attendance/register", handler.UpdateAttendance)
	r.GET("/attendance/watch", handler.GetAllAttendances)
	r.Run(":8085")
}
