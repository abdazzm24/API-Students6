package main

import (
	"context"
	"log"

	"api-students/app/repository"
	"api-students/app/service"
	"api-students/config"
	"api-students/database"
)

func main() {

	// Load environment variables
	config.LoadEnv()

	// Membuat koneksi database
	pool, err := database.NewPool(
		context.Background(),
	)

	if err != nil {
		log.Fatalf(
			"database: %v",
			err,
		)
	}

	defer pool.Close()

	// Membuat repository
	studentRepository :=
		repository.NewStudentRepository(pool)

	// Membuat service
	studentService :=
		service.NewStudentService(
			studentRepository,
		)

	// Membuat logger
	logger :=
		config.NewLogger()

	// Membuat aplikasi
	app :=
		config.NewApp(
			pool,
			studentService,
			logger,
		)

	// Menjalankan server
	log.Fatal(
		app.Listen(":3000"),
	)
}