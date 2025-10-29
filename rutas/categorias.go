package rutas

import (
	"clase_6_echo_mongo/database"
	"clase_6_echo_mongo/modelos"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gosimple/slug"
	echo "github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func ListarCategorias(mongoClient *database.MongoDBClient, dbName, collectionName string) echo.HandlerFunc {
	return func(c echo.Context) error {
		filter := bson.M{} // Filtro base, ej. de query params

		pipeline := mongo.Pipeline{
			{{Key: "$match", Value: filter}},
			{{Key: "$sort", Value: bson.M{
				"_id": -1,
			}}},
		}

		// Listar documentos de la colección
		documentos, err := mongoClient.ListDocumentos(context.TODO(), dbName, collectionName, pipeline)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al listar categorias: " + err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"mensaje": "Categorías listadas correctamente",
			"datos":   documentos,
		})
	}
}

func ListarCategoriaPorId(mongoClient *database.MongoDBClient, dbName, collectionName string) echo.HandlerFunc {
	return func(c echo.Context) error {
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
		}

		// Listar documento de la colección
		documento, err := mongoClient.ListDocumentoPorId(context.TODO(), dbName, collectionName, pipeline)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Elemento no encontrado: " + err.Error()})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al buscar categoria: " + err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"mensaje": "Categoria encontrada",
			"datos":   documento,
		})
	}
}

func CrearCategoria(mongoClient *database.MongoDBClient, dbName, collectionName string) echo.HandlerFunc {
	return func(c echo.Context) error {
		categoria := new(modelos.Categoria)

		// Bindear el JSON
		if err := c.Bind(categoria); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Error al procesar el JSON: " + err.Error()})
		}

		// Validación simple
		if categoria.Nombre == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Nombre es un campo obligatorio"})
		}

		// Agregar slug
		categoria.Slug = slug.Make(categoria.Nombre)

		// Agregar timestamp
		categoria.Timestamp = time.Now().Unix()

		// Insertar en MongoDB usando BSON
		err := mongoClient.InsertDocumento(context.TODO(), dbName, collectionName, categoria)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al guardar en la base de datos: " + err.Error()})
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"mensaje": "Categoria creada correctamente",
			"datos":   categoria,
		})
	}
}

func EditarCategoria(mongoClient *database.MongoDBClient, dbName, collectionName string) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" || !primitive.IsValidObjectID(id) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID inválido o requerido"})
		}

		categoria := new(modelos.Categoria)

		// Bindear el JSON
		if err := c.Bind(categoria); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Error al procesar el JSON: " + err.Error()})
		}

		// Validación simple
		if categoria.Nombre == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Nombre es un campo obligatorio"})
		}

		// Preparar campos para $set (solo los no vacíos)
		updateFields := bson.M{}
		if categoria.Nombre != "" {
			updateFields["nombre"] = strings.TrimSpace(categoria.Nombre)
			updateFields["slug"] = slug.Make(strings.TrimSpace(categoria.Nombre))
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
			"mensaje":    "Categoria actualizada correctamente",
			"modificado": result.MatchedCount > 0, // Indica si se cambió algo
		})
	}
}

func EliminarCategoria(mongoClient *database.MongoDBClient, dbName, collectionName string) echo.HandlerFunc {
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
