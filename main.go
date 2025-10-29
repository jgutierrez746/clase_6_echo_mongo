package main

import (
	"clase_6_echo_mongo/config"
	"clase_6_echo_mongo/database"
	"clase_6_echo_mongo/middleware_custom"
	"clase_6_echo_mongo/rutas"
	"log"
	"os"

	"github.com/joho/godotenv"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var prefijo string = "/api/v1/"

func main() {
	// Cargar el archivo .env
	errorVariables := godotenv.Load()
	if errorVariables != nil {
		log.Fatal("Error al cargar el archivo .env: ", errorVariables)
	}

	mongoURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DB_NAME")

	variablesVacias := func(valores ...string) bool {
		for _, v := range valores {
			if v == "" {
				return true
			}
		}
		return false
	}(mongoURI, dbName)

	if variablesVacias {
		log.Fatal("Las variables de entorno no están bien definidas!")
	}

	// Conectar a MongoDB y crear colecciones
	mongoClient, err := database.Connect(mongoURI, dbName)
	if err != nil {
		log.Fatal("Error al conectar a MongoDB", err)
	}
	defer mongoClient.Close()

	// Alias local para las colecciones
	cols := config.Collections

	// Instancia de echo framework
	e := echo.New()

	// Middleware
	// e.Use(middleware.Logger())
	e.Use(middleware.BodyLimit("5M"))

	e.Static("/imagenes", "public/uploads/productos")

	// Rutas ejemplo
	ejemploGroup := e.Group(prefijo + "ejemplo")
	ejemploGroup.GET("", rutas.Ejemplo_get)
	ejemploGroup.GET("/:id", rutas.Ejemplo_get_con_parametros)
	ejemploGroup.POST("", rutas.Ejemplo_post)
	ejemploGroup.PUT("/:id", rutas.Ejemplo_put)
	ejemploGroup.DELETE("/:id", rutas.Ejemplo_delete)
	e.GET(prefijo+"query-string", rutas.Ejemplo_query_string)
	e.POST(prefijo+"upload", rutas.Ejemplo_upload)

	// Rutas MongoDB 'Categorias'
	categoriaGroup := e.Group(prefijo + "categorias")
	categoriaGroup.GET("", rutas.ListarCategorias(mongoClient, dbName, cols["categorias"]))
	categoriaGroup.GET("/:id", rutas.ListarCategoriaPorId(mongoClient, dbName, cols["categorias"]))
	categoriaGroup.POST("", rutas.CrearCategoria(mongoClient, dbName, cols["categorias"]))
	categoriaGroup.PUT("/:id", rutas.EditarCategoria(mongoClient, dbName, cols["categorias"]))
	categoriaGroup.DELETE("/:id", rutas.EliminarCategoria(mongoClient, dbName, cols["categorias"]))

	// Rutas MongoDB 'Productos'
	productoGroup := e.Group(prefijo+"productos", middleware_custom.ValidarJWT) // Validación de token para acceder a productos
	productoGroup.GET("", rutas.ListarProductos(mongoClient, dbName, cols["productos"], cols["categorias"]))
	productoGroup.GET("/:id", rutas.ListarProductoPorId(mongoClient, dbName, cols["productos"], cols["categorias"]))
	productoGroup.POST("", rutas.CrearProducto(mongoClient, dbName, cols["productos"]))
	productoGroup.PUT("/:id", rutas.EditarProducto(mongoClient, dbName, cols["productos"]))
	productoGroup.DELETE("/:id", rutas.EliminarProducto(mongoClient, dbName, cols["productos"]))

	// Rutas MongoDB 'Productos-fotos'
	productoFotosGroup := e.Group(prefijo + "productos-fotos")
	productoFotosGroup.GET("/:id", rutas.ListarFotosPorIdProducto(mongoClient, dbName, cols["productos_fotos"]))
	productoFotosGroup.POST("/:id", rutas.UploadFotoProducto(mongoClient, dbName, cols["productos_fotos"]))
	productoFotosGroup.DELETE("/:id", rutas.EliminarFotoProducto(mongoClient, dbName, cols["productos_fotos"]))

	// Ruta 'Seguridad' registro y login, elementos protegidos
	seguridadGroup := e.Group(prefijo + "seguridad")
	{
		seguridadGroup.POST("/registro", rutas.RegistroUsuario(mongoClient, dbName, cols["usuarios"]))
		seguridadGroup.POST("/login", rutas.LoginUsuario(mongoClient, dbName, cols["usuarios"]))
	}

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8086"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
