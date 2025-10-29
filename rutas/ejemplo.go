package rutas

import (
	"clase_6_echo_mongo/dto"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	echo "github.com/labstack/echo/v4"
)

func Ejemplo_get(c echo.Context) error {
	cabecero := c.Request().Header.Get("Authorization")
	resp := &dto.GenericoDto{
		Estado:  "ok",
		Mensaje: "autorización: " + cabecero,
	}
	// c.Response().Header().Set("tamila", "www.tamila.cl")
	return c.JSON(http.StatusOK, resp)
}

func Ejemplo_get_con_parametros(c echo.Context) error {
	id := c.Param("id")
	resp := &dto.GenericoDto{
		Estado:  "ok",
		Mensaje: "Método GET | id = " + id,
	}
	return c.JSON(http.StatusOK, resp)
}

func Ejemplo_post(c echo.Context) error {
	categoria := new(dto.CategoriaDto)

	// Bindear el cuerpo del JSON a la estructura dto
	if err := c.Bind(categoria); err != nil {
		resp := &dto.GenericoDto{
			Estado:  "error",
			Mensaje: "Error al procesar el JSON: " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, resp)
	}

	// Validar que los campos no estén vacíos
	if categoria.Nombre == "" {
		resp := &dto.GenericoDto{
			Estado:  "error",
			Mensaje: "El nombre es un campo obligatorio",
		}
		return c.JSON(http.StatusBadRequest, resp)
	}

	resp := &dto.GenericoDto{
		Estado:  "ok",
		Mensaje: "Nombre: " + categoria.Nombre,
	}
	return c.JSON(http.StatusOK, resp)
}

func Ejemplo_put(c echo.Context) error {
	id := c.Param("id")
	resp := &dto.GenericoDto{
		Estado:  "ok",
		Mensaje: "Método PUT | id = " + id,
	}
	return c.JSON(http.StatusOK, resp)
}

func Ejemplo_delete(c echo.Context) error {
	id := c.Param("id")
	resp := &dto.GenericoDto{
		Estado:  "ok",
		Mensaje: "Método DELETE | id = " + id,
	}
	return c.JSON(http.StatusOK, resp)
}

func Ejemplo_query_string(c echo.Context) error {
	id := c.QueryParam("id")
	slug := c.QueryParam("slug")
	resp := &dto.GenericoDto{
		Estado:  "ok",
		Mensaje: "Query String | id = " + id + " | slug = " + slug,
	}
	return c.JSON(http.StatusOK, resp)
}

func Ejemplo_upload(c echo.Context) error {
	file, err := c.FormFile("foto")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Renombramos el archivo
	var extension = strings.Split(file.Filename, ".")[1]
	unixTime := time.Now().Unix()
	foto := strconv.FormatInt(unixTime, 10) + "." + extension
	var archivo string = "public/uploads/fotos/" + foto

	dst, err := os.Create(archivo)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	respuesta := map[string]string{
		"estado":  "ok",
		"mensaje": "todo bien",
		"foto":    foto,
	}

	return c.JSON(http.StatusOK, respuesta)
}
