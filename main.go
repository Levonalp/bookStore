package main

import (
	"log"
	"net/http"
	"os"
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/levonalp/go-fiber-postgres/models"
	"github.com/levonalp/go-fiber-postgres/storage"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}

	err := context.BodyParser(&book)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create book"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book has been added"})
	return nil
}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	bookModel := models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.Delete(bookModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete book",
		})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book delete successfully",
	})
	return nil
}

func (r *Repository) GetBooks(context *fiber.Ctx) error {
	bookModels := &[]models.Books{}

	err := r.DB.Find(bookModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get books"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "books fetched successfully",
		"data":    bookModels,
	})
	return nil
}

func (r *Repository) GetBookByID(context *fiber.Ctx) error {
	id := context.Params("id")
	bookModel := &models.Books{}
	if id == "" {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the book"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book retrieved successfully",
		"data":    bookModel,
	})
	return nil
}

func (r *Repository) UpdateBook(context *fiber.Ctx) error {
    book := Book{}
    id := context.Params("id")
    if id == "" {
        context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
            "message": "id cannot be empty",
        })
        return nil
    }

    err := context.BodyParser(&book)
    if err != nil {
        context.Status(http.StatusUnprocessableEntity).JSON(
            &fiber.Map{"message": "request failed"})
        return err
    }

    err = r.DB.Model(&book).Where("id = ?", id).Updates(book).Error
    if err != nil {
        context.Status(http.StatusBadRequest).JSON(
            &fiber.Map{"message": "could not update book"})
        return err
    }

    context.Status(http.StatusOK).JSON(&fiber.Map{
        "message": "book has been updated"})
    return nil
}

	func (r *Repository) SetupRoutes(app *fiber.App) {
		api := app.Group("/api")
		api.Post("/books", r.CreateBook)
		api.Delete("books/:id", r.DeleteBook)
		api.Get("/books/:id", r.GetBookByID)
		api.Get("/books", r.GetBooks)
		api.Put("/books/:id", r.UpdateBook)
		}
		
		func main() {
			// Load the environment variables
			err := godotenv.Load(".env")
			if err != nil {
				log.Fatal(err)
			}
		
		// Create a new database connection
config := &storage.Config{
    Host:     os.Getenv("DB_HOST"),
    Port:     os.Getenv("DB_PORT"),
    Password: os.Getenv("DB_PASS"),
    User:     os.Getenv("DB_USER"),
    SSLMode:  os.Getenv("DB_SSLMODE"),
    DBName:   os.Getenv("DB_NAME"),
}
db, err := storage.NewConnection(config)
if err != nil {
    log.Fatal("could not load the database")
}

// Automigrate
err = models.MigrateBooks(db)
if err != nil {
    log.Fatal(err)
}

// Re-run Automigrate if you want to update your schema
err = models.MigrateBooks(db)
if err != nil {
    log.Fatal(err)
}

// Ensure that the connection is closed when the function returns
 // Automigrate
 err = models.MigrateBooks(db)
 if err != nil {
	 log.Fatal(err)
 }

 // Re-run Automigrate if you want to update your schema
 err = models.MigrateBooks(db)
 if err != nil {
	 log.Fatal(err)
 }

		
			// Migrate the database
			err = models.MigrateBooks(db)
			if err != nil {
				log.Fatal("could not migrate db")
			}
			
		
			// Create a new repository
			r := Repository{
				DB: db,
			}
		
			// Create a new fiber app
			app := fiber.New()
		
			// Enable CORS
			app.Use(cors.New(cors.Config{
				AllowOrigins:     "*",
				AllowMethods:     "GET, POST, PUT, DELETE",
				AllowHeaders:     "*",
				ExposeHeaders:    "Content-Length",
				AllowCredentials: true,
				MaxAge:           3600,
			}))
		
			// Set up the routes
			r.SetupRoutes(app)
		
			// Start the server
			app.Listen(":8080")
		}
		