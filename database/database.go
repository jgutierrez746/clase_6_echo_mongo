package database

import (
	"clase_6_echo_mongo/config"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBClient encapsulala conexión a MongoDB
type MongoDBClient struct {
	Client *mongo.Client
}

// Inicializar la conexión a MongoDB y crea las colecciones necesarias
func Connect(uri, dbName string) (*MongoDBClient, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	db := client.Database(dbName)

	// Crear automáticamente todas las colecciones
	if err := createAllCollections(db, config.Collections); err != nil {
		return nil, err
	}

	log.Println("Conectado correctamente a MongoDB y colecciones listas.")
	return &MongoDBClient{Client: client}, nil
}

func createAllCollections(db *mongo.Database, collections map[string]string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	existing, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return err
	}

	existingMap := make(map[string]bool)
	for _, name := range existing {
		existingMap[name] = true
	}

	for _, name := range collections {
		if !existingMap[name] {
			if err := db.CreateCollection(ctx, name); err != nil {
				return fmt.Errorf("error creando colección %s: %w", name, err)
			}
			log.Printf("Colección %s creada", name)
		} else {
			log.Printf("Colección %s ya existe", name)
		}
	}
	return nil
}

// Close cierra la conexión a MongoDB
func (c *MongoDBClient) Close() {
	if err := c.Client.Disconnect(context.TODO()); err != nil {
		log.Printf("Error al cerrar la conexión a MongoDB: %v", err)
	}
}

// GetCollection devuelve una colección especifica de una base de datos
func (c *MongoDBClient) GetCollection(dbName, collectionName string) *mongo.Collection {
	return c.Client.Database(dbName).Collection(collectionName)
}

// InsertDocumento inserta un documento en la colección usando BSON
func (c *MongoDBClient) InsertDocumento(ctx context.Context, dbName, collectionName string, documento interface{}) error {
	collection := c.GetCollection(dbName, collectionName)
	_, err := collection.InsertOne(ctx, documento)
	return err
}

// ListDocumentos lista todos los documentos de una colección
func (c *MongoDBClient) ListDocumentos(ctx context.Context, dbName, collectionName string, pipeline mongo.Pipeline) ([]interface{}, error) {
	collection := c.GetCollection(dbName, collectionName)

	// Ejecutar agregación
	cursor, err := collection.Aggregate(ctx, pipeline) // Se pasa pipeline como parámetro para la busqueda
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var resultados []interface{}
	for cursor.Next(ctx) {
		var documento bson.M
		if err := cursor.Decode(&documento); err != nil {
			return nil, err
		}
		resultados = append(resultados, documento)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return resultados, nil
}

func (c *MongoDBClient) ListDocumentoPorId(ctx context.Context, dbName, collectionName string, pipeline mongo.Pipeline) ([]bson.M, error) {
	collection := c.GetCollection(dbName, collectionName)

	// Ejecutar agregación
	cursor, err := collection.Aggregate(ctx, pipeline) // Se pasa pipeline como parámetro para la busqueda
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var resultado []bson.M
	for cursor.Next(ctx) {
		var documento bson.M
		if err := cursor.Decode(&documento); err != nil {
			return nil, err
		}
		resultado = append(resultado, documento)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(resultado) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return resultado, nil
}

func (c *MongoDBClient) UpdateDocumento(ctx context.Context, dbName, collectionName, id string, updateFields bson.M) (*mongo.UpdateResult, error) {
	collection := c.GetCollection(dbName, collectionName)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Usar $set para actualizar solo los campos proporcionados
	update := bson.D{{Key: "$set", Value: updateFields}}
	resultado, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return nil, err
	}
	if resultado.MatchedCount == 0 {
		return nil, mongo.ErrNoDocuments // No encontrado
	}
	return resultado, nil
}

func (c *MongoDBClient) DeleteDocumento(ctx context.Context, dbName, collectionName, id string) (*mongo.DeleteResult, error) {
	collection := c.GetCollection(dbName, collectionName)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: objID}} // Filtro con bson.D: busca por _id (usa ObjectID)

	// Ejecutar DeleteOne
	resultado, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	if resultado.DeletedCount == 0 {
		return nil, mongo.ErrNoDocuments // No se elimino el registro o no se encontró
	}

	return resultado, nil
}

// BuscarDocumentoExistente ejecuta un pipeline con un filtro // y proyección.
// Devuelve []bson.M si hay resultados, o un error si ocurre un fallo de consulta.
func (c *MongoDBClient) BuscarDocumentoExistente(ctx context.Context, dbName, collectionName string, filter bson.M /*project bson.D*/) ([]bson.M, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		// {{Key: "$project", Value: project}},
	}

	resultado, err := c.ListDocumentoPorId(context.TODO(), dbName, collectionName, pipeline)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Devuelve slice vacío, sin error, para permitir continuar
			return []bson.M{}, nil
		}
		return nil, errors.New("error al consultar la base de datos: " + err.Error())
	}

	return resultado, nil
}
