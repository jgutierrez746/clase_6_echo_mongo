package rutas

import (
	"clase_6_echo_mongo/database"
	"clase_6_echo_mongo/modelos"
	"clase_6_echo_mongo/validaciones"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	echo "github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func ListarProductos(mongoClient *database.MongoDBClient, dbName, productosCollection, categoriasCollection string) echo.HandlerFunc {
	return func(c echo.Context) error {
		filter := bson.M{} // Filtro base, ej. de query params

		pipeline := mongo.Pipeline{
			{{Key: "$match", Value: filter}},
			{{Key: "$lookup", Value: bson.M{
				"from":         categoriasCollection,
				"localField":   "categoria_id",
				"foreignField": "_id",
				"as":           "categoria", // Nombre de la relación
			}}},
			{{Key: "$project", Value: bson.D{
				{Key: "categoria_id", Value: 0},
			}}},
			{{Key: "$sort", Value: bson.M{
				"_id": -1,
			}}},
			// {{Key: "$unwind", Value: "$categoria"}},
		}

		// Listar documentos de la colección
		documentos, err := mongoClient.ListDocumentos(context.TODO(), dbName, productosCollection, pipeline)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al listar categorias: " + err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"mensaje": "Productos listados correctamente",
			"datos":   documentos,
		})
	}
}

func ListarProductoPorId(mongoClient *database.MongoDBClient, dbName, productosCollection, categoriasCollection string) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Con esto obtenemos los valores entregados por el token validado, como el nombre de usuario y otras cosas
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)

		idUsuario := claims["id"].(string)
		nombreUsuario := claims["nombre"].(string)
		// fin de la obtención de valores, esto lo usaré en el return de producto al final de la función

		id := c.Param("id") // Obtener ID de la URL (:id)
		if id == "" || !primitive.IsValidObjectID(id) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID inválido o requerido"})
		}

		objID, err := primitive.ObjectIDFromHex(id) // Convertir el string a ObjectIds
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID no es válido: " + err.Error()})
		}

		filter := bson.M{
			"_id": objID,
		} // Filtro base, ej. de query params

		pipeline := mongo.Pipeline{
			{{Key: "$match", Value: filter}},
			{{Key: "$lookup", Value: bson.M{
				"from":         categoriasCollection,
				"localField":   "categoria_id",
				"foreignField": "_id",
				"as":           "categoria", // Nombre de la relación
			}}},
			{{Key: "$project", Value: bson.D{
				{Key: "categoria_id", Value: 0},
			}}},
			// {{Key: "$unwind", Value: "$categoria"}}, // Separa un array en diferentes bloques individuales
		}

		// Listar documentos de la colección
		documentos, err := mongoClient.ListDocumentoPorId(context.TODO(), dbName, productosCollection, pipeline)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Elemento no encontrado: " + err.Error()})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al listar categorias: " + err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"mensaje":   "Producto encontrado",
			"datos":     documentos,
			"usuario":   "Hola " + nombreUsuario,
			"idUsuario": idUsuario,
		})
	}
}

func CrearProducto(mongoClient *database.MongoDBClient, dbName, collectionName string) echo.HandlerFunc {
	return func(c echo.Context) error {
		producto := new(modelos.Producto)

		// Bindear el JSON
		if err := c.Bind(producto); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Error al procesar el JSON: " + err.Error()})
		}

		// Validación de campos
		if err := validaciones.ValidarProducto(*producto); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		// Crear Map
		categoriaID, _ := primitive.ObjectIDFromHex(producto.CategoriaID)
		producto.Timestamp = time.Now().Unix()
		documentoProducto := bson.M{
			"nombre":       producto.Nombre,
			"precio":       producto.Precio,
			"stock":        producto.Stock,
			"descripcion":  producto.Descripcion,
			"categoria_id": categoriaID,
			"timestamp":    producto.Timestamp,
		}

		// Insertar en MongoDB usando BSON
		err := mongoClient.InsertDocumento(context.TODO(), dbName, collectionName, documentoProducto)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al guardar en la base de datos: " + err.Error()})
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"mensaje": "Producto creado correctamente",
			"datos":   producto,
		})
	}
}

func EditarProducto(mongoClient *database.MongoDBClient, dbName, collectionName string) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" || !primitive.IsValidObjectID(id) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID inválido o requerido"})
		}

		producto := new(modelos.UpdateProducto)

		// Bindear el JSON
		if err := c.Bind(producto); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Error al procesar el JSON: " + err.Error()})
		}

		// Validación de al menos un campo
		if producto.Nombre == "" && producto.Precio == 0 && producto.Stock == 0 && producto.Descripcion == "" && producto.CategoriaID.IsZero() {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Debe proporcionar al menos un campo para actualizar"})
		}

		// Preparar campos para $set (solo los no vacíos)
		updateFields := bson.M{}
		if producto.Nombre != "" {
			updateFields["nombre"] = strings.TrimSpace(producto.Nombre)
		}
		if producto.Precio > 0 {
			updateFields["precio"] = producto.Precio
		}
		if producto.Stock >= 0 {
			updateFields["stock"] = producto.Stock
		}
		if producto.Descripcion != "" {
			updateFields["descripcion"] = strings.TrimSpace(producto.Descripcion)
		}
		if !producto.CategoriaID.IsZero() {
			updateFields["categoria_id"] = producto.CategoriaID
		}

		// Actualizar en MongoDB
		result, err := mongoClient.UpdateDocumento(context.TODO(), dbName, collectionName, id, updateFields)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Elemento no encontrado: " + err.Error()})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al actualizar categoria: " + err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"mensaje":    "Producto actualizado correctamente",
			"modificado": result.MatchedCount > 0, // Indica si se cambió algo
		})
	}
}

func EliminarProducto(mongoClient *database.MongoDBClient, dbName, collectionName string) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" || !primitive.IsValidObjectID(id) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID requerido o inválido"})
		}

		// Eliminar documento de la colección en MongoDB
		resultado, err := mongoClient.DeleteDocumento(context.TODO(), dbName, collectionName, id)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Elemento no encontrado: " + err.Error()})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al eliminar categoria: " + err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"mensaje":   "Categoria eliminada correctamente",
			"eliminado": resultado.DeletedCount > 0, // Confirma que se eliminó
			"id":        id,
		})
	}
}
