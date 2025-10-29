package modelos

type UsuarioDto struct {
	Nombre    string `json:"nombre" validate:"required" bson:"nombre"`
	Correo    string `json:"correo" validate:"required,email" bson:"correo"`
	Telefono  string `json:"telefono" validate:"required,numeric" bson:"telefono"`
	Password  string `json:"password" validate:"required,password" bson:"password"`
	Timestamp int64  `json:"timestamp,omitempty" validate:"omitempty" bson:"timestamp"`
}

type LoginDto struct {
	Correo   string `json:"correo" validate:"required,email" bson:"correo"`
	Password string `json:"password" validate:"required,password" bson:"password"`
}

type LoginRespuestaDto struct {
	Nombre string `json:"nombre"`
	Token  string `json:"token"`
}
