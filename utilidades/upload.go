package utilidades

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const rutaLocal = "public/uploads/productos"

func SubirArchivo(file *multipart.FileHeader) (map[string]string, error) {
	// Obtiene el nombre original del archivo
	src, err := file.Open()
	if err != nil {
		return map[string]string{
			"error": "Erro al abrir el archivo",
		}, err
	}
	defer src.Close()

	// Renombramos el archivo
	var extension = strings.Split(file.Filename, ".")[1]
	unixTime := time.Now().Unix()
	nomnbreArchivo := strconv.FormatInt(unixTime, 10) + "." + extension

	// Path destino del archivo
	dstPath := filepath.Join(rutaLocal, nomnbreArchivo)

	dst, err := os.Create(dstPath)
	if err != nil {
		return map[string]string{
			"error": "Erro al crear el archivo",
		}, err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return map[string]string{"error": "Erro al copiar el archivo"}, err
	}

	return map[string]string{
		"nombre": nomnbreArchivo,
	}, nil
}

func EliminarArchivo(nombre string) (map[string]string, error) {
	rutaArchivo := filepath.Join(rutaLocal, nombre)
	if err := os.Remove(rutaArchivo); err != nil {
		return map[string]string{"error": "Fallo eliminar archivo '" + nombre + "' o no existe en el servidor"}, err
	}
	return nil, nil
}
