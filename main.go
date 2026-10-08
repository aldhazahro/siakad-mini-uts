package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"siakad-mini/app/repository"
	"siakad-mini/app/service"
	"siakad-mini/config"
	"siakad-mini/database"
	"siakad-mini/helper"
	"siakad-mini/middleware"
	"siakad-mini/route"
)

func main() {
	cfg := config.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(pool)
	studentRepo := repository.NewStudentRepository(pool)
	courseRepo := repository.NewCourseRepository(pool)
	enrollmentRepo := repository.NewEnrollmentRepository(pool)
	jwt := helper.NewJWTManager(cfg.JWTSecret, cfg.JWTIssuer, cfg.JWTAccessTTL)

	authService := service.NewAuthService(userRepo, studentRepo, jwt)
	studentService := service.NewStudentService(pool, userRepo, studentRepo, enrollmentRepo)
	courseService := service.NewCourseService(courseRepo)
	enrollmentService := service.NewEnrollmentService(pool, studentRepo, courseRepo, enrollmentRepo)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			status := fiber.StatusInternalServerError
			var fe *fiber.Error
			if errors.As(err, &fe) {
				status = fe.Code
			}
			return c.Status(status).JSON(fiber.Map{"success": false, "message": "terjadi kesalahan pada server", "request_id": c.Locals("request_id")})
		},
	})
	app.Use(middleware.RequestID)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "message": "SIAKAD Mini API aktif"})
	})
	route.Register(app, route.Dependencies{Auth: authService, Students: studentService, Courses: courseService, Enrollments: enrollmentService, JWT: jwt})

	log.Printf("SIAKAD Mini API berjalan di http://localhost:%s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
	_ = pgx.ErrNoRows
}
