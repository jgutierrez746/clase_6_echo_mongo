package rutas

import (
	"clase_6_echo_mongo/database"
	"clase_6_echo_mongo/jwt"
	"clase_6_echo_mongo/modelos"
	"clase_6_echo_mongo/validaciones"
	"context"
	"net/http"
	"time"

	echo "github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func LoginUsuario(mongoClient *database.MongoDBClient, dbName, collectionName string) echo.HandlerFunc {
	return func(c echo.Context) error {
		usuarioLogin := new(modelos.LoginDto)

		// Bindear el JSON
		if err := c.Bind(usuarioLogin); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Error al procesar el JSON: " + err.Error()})
		}

		// Validación de campos
		if err := validaciones.ValidarLogin(*usuarioLogin); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		// INTEGRANDO VALIDACIÓN DE CORREO EN BASE DE DATOS
		filter := bson.M{
			"correo": usuarioLogin.Correo,
		} // Filtro base, ej. de query params

		documento, err := mongoClient.BuscarDocumentoExistente(context.TODO(), dbName, collectionName, filter)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		// Validar si ya existe
		if len(documento) < 1 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Las credenciales ingresadas son inválidas"})
		}

		// Proceso de comparación de password
		passwordBytes := []byte(usuarioLogin.Password)
		passwordEnDocumento := documento[0]["password"]
		passwordBD := []byte(passwordEnDocumento.(string))
		errPassword := bcrypt.CompareHashAndPassword(passwordBD, passwordBytes)

		if errPassword != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Las credenciales ingresadas son inválidas"})
		} else {
			correoEnDocumento := documento[0]["correo"].(string)
			usuarioEnDocumento := documento[0]["nombre"].(string)
			idEnDocumento := documento[0]["_id"].(primitive.ObjectID).Hex()
			jwtKey, err := jwt.GenerarJWT(correoEnDocumento, usuarioEnDocumento, idEnDocumento)

			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Error al intentar generar el token" + err.Error()})
			} else {
				retorno := modelos.LoginRespuestaDto{
					Nombre: usuarioEnDocumento,
					Token:  "Bearer " + jwtKey,
				}
				return c.JSON(http.StatusOK, retorno)
			}
		}
	}
}

func RegistroUsuario(mongoClient *database.MongoDBClient, dbName, collectionName string) echo.HandlerFunc {
	return func(c echo.Context) error {
		usuario := new(modelos.UsuarioDto)

		// Bindear el JSON
		if err := c.Bind(usuario); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Error al procesar el JSON: " + err.Error()})
		}

		// Validación de campos
		if err := validaciones.ValidarUsuario(*usuario); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		// INTEGRANDO VALIDACIÓN DE CORREO EN BASE DE DATOS
		filter := bson.M{
			"correo": usuario.Correo,
		} // Filtro base, ej. de query params

		documento, err := mongoClient.BuscarDocumentoExistente(context.TODO(), dbName, collectionName, filter)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		// Validar si ya existe
		if len(documento) > 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "El correo ya está registrado"})
		}

		// generar Hash con Bcrypt para contraseña
		// wololo.Pass2 contraseña ejemplo
		costo := 8
		bytes, _ := bcrypt.GenerateFromPassword([]byte(usuario.Password), costo)
		password := string(bytes)

		// Creando documento para inserción
		usuario.Timestamp = time.Now().Unix()
		documentoUsuario := bson.M{
			"nombre":    usuario.Nombre,
			"correo":    usuario.Correo,
			"telefono":  usuario.Telefono,
			"timestamp": usuario.Timestamp,
		}

		// Asigno el valor encriptado de la contraseña al map bson.M que se insertará en la base de datos
		documentoUsuario["password"] = password

		// Insertar en MongoDB usando BSON
		err = mongoClient.InsertDocumento(context.TODO(), dbName, collectionName, documentoUsuario)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al guardar en la base de datos: " + err.Error()})
		}

		return c.JSON(http.StatusCreated, map[string]string{
			"mensaje": "Usuario creado correctamente",
			"estado":  "ok",
		})
	}
}
