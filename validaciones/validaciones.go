package validaciones

import (
	"clase_6_echo_mongo/modelos"
	"errors"
	"fmt"
	"strings"
	"unicode"

	validator "github.com/go-playground/validator/v10"
)

func passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	hasMinLen := len(password) >= 8
	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case (char >= 33 && char <= 47) || (char >= 58 && char <= 64) ||
			(char >= 91 && char <= 96) || (char >= 123 && char <= 126):
			hasSpecial = true
		}
	}

	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func ValidarUsuario(dto modelos.UsuarioDto) error {
	validate := validator.New()

	// Registrar la validación personalizada
	if err := validate.RegisterValidation("password", passwordValidator); err != nil {
		return err
	}

	if err := validate.Struct(&dto); err != nil {
		var mensajes []string

		for _, e := range err.(validator.ValidationErrors) {
			campo := e.Field() // Nombre del campo
			tag := e.Tag()     // Regla que falló
			valor := e.Value() // Valor que causó el error

			// Mensaje personalizado
			var msg string
			switch tag {
			case "required":
				msg = fmt.Sprintf("El campo '%s' es requerido", campo)
			case "min":
				msg = fmt.Sprintf("El campo '%s' debe tener al menos %s caracteres", campo, e.Param())
			case "email":
				msg = fmt.Sprintf("El campo '%s' debe contener un correo válido", campo)
			case "password":
				msg = fmt.Sprintf("El campo '%s' debe presentar un formato válido", campo)
			default:
				msg = fmt.Sprintf("Error en '%s': %s no es válido (%s)", campo, valor, tag)
			}
			mensajes = append(mensajes, msg)
		}
		// Unir mensajes en solo uno
		return errors.New(strings.Join(mensajes, "; "))
	}
	return nil
}

func ValidarLogin(dto modelos.LoginDto) error {
	validate := validator.New()

	// Registrar la validación personalizada
	if err := validate.RegisterValidation("password", passwordValidator); err != nil {
		return err
	}

	if err := validate.Struct(&dto); err != nil {
		var mensajes []string

		for _, e := range err.(validator.ValidationErrors) {
			campo := e.Field() // Nombre del campo
			tag := e.Tag()     // Regla que falló
			valor := e.Value() // Valor que causó el error

			// Mensaje personalizado
			var msg string
			switch tag {
			case "required":
				msg = fmt.Sprintf("El campo '%s' es requerido", campo)
			case "email":
				msg = fmt.Sprintf("El campo '%s' debe contener un correo válido", campo)
			case "password":
				msg = fmt.Sprintf("El campo '%s' debe presentar un formato válido", campo)
			default:
				msg = fmt.Sprintf("Error en '%s': %s no es válido (%s)", campo, valor, tag)
			}
			mensajes = append(mensajes, msg)
		}
		// Unir mensajes en solo uno
		return errors.New(strings.Join(mensajes, "; "))
	}
	return nil
}

func ValidarProducto(dto modelos.Producto) error {
	validate := validator.New()
	if err := validate.Struct(&dto); err != nil {
		var mensajes []string

		for _, e := range err.(validator.ValidationErrors) {
			campo := e.Field() // Nombre del campo
			tag := e.Tag()     // Regla que falló
			valor := e.Value() // Valor que causó el error

			// Mensaje personalizado
			var msg string
			switch tag {
			case "required":
				msg = fmt.Sprintf("El campo '%s' es requerido", campo)
			case "min":
				msg = fmt.Sprintf("El campo '%s' debe tener al menos %s caracteres", campo, e.Param())
			case "gt":
				msg = fmt.Sprintf("El campo '%s' debe ser mayor que %s", campo, e.Param())
			case "gte":
				msg = fmt.Sprintf("El campo '%s' debe ser mayor o igual que %s", campo, e.Param())
			default:
				msg = fmt.Sprintf("Error en '%s': %s no es válido (%s)", campo, valor, tag)
			}
			mensajes = append(mensajes, msg)
		}
		// Unir mensajes en solo uno
		return errors.New(strings.Join(mensajes, "; "))
	}
	return nil
}

// NO SE ESTÁ UTILIZANDO
func ValidarUploadFotoProducto(dto modelos.UploadFotoProducto) error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(&dto); err != nil {
		var mensajes []string
		for _, e := range err.(validator.ValidationErrors) {
			campo := e.Field() // Nombre del campo
			tag := e.Tag()     // Regla que falló
			valor := e.Value() // Valor que causó el error

			// Mensaje personalizado
			var msg string
			switch tag {
			case "required":
				msg = fmt.Sprintf("El campo '%s' es requerido", campo)
			case "min":
				msg = fmt.Sprintf("El campo '%s' debe tener al menos %s caracteres", campo, e.Param())
			case "mongodb":
				msg = fmt.Sprintf("El campo '%s' debe ser un objectID válido", campo)
			default:
				msg = fmt.Sprintf("Error en '%s': %s no es válido (%s)", campo, valor, tag)
			}
			mensajes = append(mensajes, msg)
		}
		// Unir mensajes en solo uno
		return errors.New(strings.Join(mensajes, "; "))
	}
	return nil
}
