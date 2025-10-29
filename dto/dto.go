package dto

type GenericoDto struct {
	Estado  string `json:"nombre" xml:"nombre"`
	Mensaje string `json:"mensaje" xml:"mensaje"`
}

type CategoriaDto struct {
	Nombre string `json:"nombre"`
}

/*type ProductoDto struct {
	Nombre      string `json:"nombre"`
	Precio      int    `json:"precio"`
	Stock       int    `json:"stock"`
	Descripcion string `json:"descripcion"`
	CategoriaID string `json:"categoria_id"`
}

type UsuarioDto struct {
	Nombre   string `json:"nombre"`
	Correo   string `json:"correo"`
	Telefono string `json:"telefono"`
	Password string `json:"password"`
}

type LoginDto struct {
	Correo   string `json:"correo"`
	Password string `json:"password"`
}

type LoginRespuestaDto struct {
	Nombre string `json:"nombre"`
	Token  string `json:"token"`
}*/
