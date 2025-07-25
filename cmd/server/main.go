package main

import (
	"log"
	"os"

	"crabi-test/internal/infrastructure/database/sqlite"
	"crabi-test/internal/infrastructure/http/routes"
	"crabi-test/pkg/validator"

	_ "crabi-test/docs" // Importar docs generados

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Crabi API
// @version 1.0
// @description API para sistema de PLD (Prevención de Lavado de Dinero) de Crabi
// @termsOfService http://swagger.io/terms/

// @contact.name Crabi API Support
// @contact.url https://github.com/crabi-test
// @contact.email dev@crabi-test.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Cargar variables de entorno desde .env
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}

	// Configurar modo de Gin
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Inicializar base de datos
	db, err := sqlite.InitDB()
	if err != nil {
		log.Fatal("Error inicializando base de datos:", err)
	}
	defer db.Close()

	// Crear router
	r := gin.Default()

	// Middleware de validación personalizada
	r.Use(validator.CustomValidator())

	// Configurar rutas
	routes.SetupRoutes(r, db)

	// Documentación Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Ruta de health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "Crabi API is running",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor iniciando en puerto %s", port)
	log.Printf("Documentación disponible en: http://localhost:%s/swagger/index.html", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Error iniciando servidor:", err)
	}
}
