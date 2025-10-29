package rutas

import (
	"clase_6_echo_mongo/database"
	"clase_6_echo_mongo/utilidades"
	"context"
	"net/http"
	"time"

	echo "github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UploadFotoProducto(mongoClient *database.MongoDBClient, dbName, collectionName string) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" || !primitive.IsValidObjectID(id) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID inválido o requerido"})
		}

		// Obtener el archivo del formulario multipart
		file, err := c.FormFile("file")
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "No se encontró el archivo"})
		}

		mensaje, err := utilidades.SubirArchivo(file)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, mensaje)
		}

		// Crear Map
		productoID, _ := primitive.ObjectIDFromHex(id)
		timestamp := time.Now().Unix()
		documentoFotoProducto := bson.M{
			"nombre":      mensaje["nombre"],
			"producto_id": productoID,
			"timestamp":   timestamp,
		}

		// Insertar en MongoDB usando BSON
		err = mongoClient.InsertDocumento(context.TODO(), dbName, collectionName, documentoFotoProducto)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al guardar en la base de datos: " + err.Error()})
		}

		return c.JSON(http.StatusCreated, map[string]string{
			"mensaje": "Foto cargada y registrada correctamente",
			"estado":  "ok",
		})
	}
}

func ListarFotosPorIdProducto(mongoClient *database.MongoDBClient, dbName, collectionName string) echo.HandlerFunc {
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
			"producto_id": objID,
		} // Filtro base, ej. de query params

		pipeline := mongo.Pipeline{
			{{Key: "$match", Value: filter}},
			{{Key: "$project", Value: bson.D{ // $project permite omitir campos al asignar el valor a 0
				{Key: "timestamp", Value: 0},
				{Key: "producto_id", Value: 0},
			}}},
			{{Key: "$addFields", Value: bson.M{ // $addFields me permite agregar texto a un campo al devolver su valor desde MongoDB
				"nombre": bson.M{
					"$concat": []interface{}{"http://localhost:8086/imagenes/", "$nombre"},
				},
			}}},
		}

		// Listar documentos de la colección
		documentos, err := mongoClient.ListDocumentos(context.TODO(), dbName, collectionName, pipeline)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al listar categorias: " + err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"mensaje":     "Imágenes encontradas",
			"producto_id": id,
			"imagenes":    documentos,
		})
	}
}

func EliminarFotoProducto(mongoClient *database.MongoDBClient, dbName, collectionName string) echo.HandlerFunc {
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

		documento, err := mongoClient.BuscarDocumentoExistente(context.TODO(), dbName, collectionName, filter)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		if len(documento) == 0 {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Elemento no encontrado"})
		}

		nombreEnDocumento := documento[0]["nombre"]
		nombreArchivo, _ := nombreEnDocumento.(string)

		mensaje, err := utilidades.EliminarArchivo(nombreArchivo)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, mensaje)
		}

		// Eliminar documento de la colección en MongoDB
		resultado, err := mongoClient.DeleteDocumento(context.TODO(), dbName, collectionName, id)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Elemento no encontrado: " + err.Error()})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al eliminar documento: " + err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"mensaje":   "Imágen eliminada correctamente",
			"eliminado": resultado.DeletedCount > 0, // Confirma que se eliminó
			"id":        id,
		})
	}
}
