package modelos

import "go.mongodb.org/mongo-driver/bson/primitive"

// Categoria representa una categoria en la base de datos
type Categoria struct {
	Nombre    string `json:"nombre" bson:"nombre"`
	Slug      string `json:"slug,omitempty" bson:"slug"`
	Timestamp int64  `json:"timestamp,omitempty" bson:"timestamp"`
}

// Producto representa un producto en la base de datos
type Producto struct {
	Nombre      string `json:"nombre" validate:"required,min=2,max=100" bson:"nombre"`
	Precio      int    `json:"precio" validate:"required,gt=0" bson:"precio"`
	Stock       int    `json:"stock" validate:"required,gte=0" bson:"stock"`
	Descripcion string `json:"descripcion" validate:"required,min=10" bson:"descripcion"`
	CategoriaID string `json:"categoria_id" validate:"required,len=24" bson:"categoria_id"`
	Timestamp   int64  `json:"timestamp,omitempty" validate:"omitempty" bson:"timestamp"`
}

type UpdateProducto struct {
	Nombre      string             `json:"nombre"`
	Precio      int                `json:"precio"`
	Stock       int                `json:"stock"`
	Descripcion string             `json:"descripcion"`
	CategoriaID primitive.ObjectID `json:"categoria_id" bson:"categoria_id"`
}
