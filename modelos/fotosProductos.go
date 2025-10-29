package modelos

import "go.mongodb.org/mongo-driver/bson/primitive"

// UploadFotoProducto representa una imagen de producto en la base de datos
// NO SE ESTÁ UTILIZANDO
type UploadFotoProducto struct {
	Nombre     string `json:"nombre" validate:"required,min=2" bson:"nombre"`
	ProductoID string `json:"producto_id" validate:"required,mongodb" bson:"producto_id"`
	Timestamp  int64  `json:"timestamp,omitempty" validate:"omitempty" bson:"timestamp"`
}

type EliminarFotoProducto struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	nombre string             `bson:"nombre"`
}
