package response

import (
	"gopokemon/config"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	ResponseDataNotFound              = "data_not_found"
	ResponseValidationFailed          = "validation_failed"
	ResponseInvalidUsernameOrPassword = "invalid_username_or_password"
	ResponseSuccess                   = "success"
	ResponseNotAllowed                = "not_allowed"

	ErrorFormSelectionInvalid = "error_form_selection_invalid"
	ErrorFormDataNotFound     = "error_form_data_not_found"
	ErrorFormLengthTooBig     = "error_form_length_too_big"
	ErrorFormLengthTooShort   = "error_form_length_too_short"
	ErrorFormInvalidEmail     = "error_form_invalid_email"
	ErrorFormFieldRequired    = "error_form_field_required"
	ErrorFormFieldNumeric     = "error_form_field_numeric"
	ErrorFormFixedLength      = "error_form_fix_length"
	ErrorFormAlreadyExists    = "error_form_already_exists"
	ErrorFormFlexibleMsg      = "error_form_flexible_msg"
	ErrorFormPhoto            = "error_form_photo"
	ErrorFormLowercase        = "error_form_lowercase"
	ErrorFormUppercase        = "error_form_uppercase"
	ErrorFormHiragana         = "error_form_hiragana"
	ErrorFormKatakana         = "error_form_katakana"
	ErrorFormKana             = "error_form_kana"
	ErrorFormKanji            = "error_form_kanji"
)

type FieldError struct {
	Field string `json:"field"`
	Msg   string `json:"msg"`
}

func getFieldError(str ...string) interface{} {
	switch str[1] {
	case ErrorFormSelectionInvalid:
		return FieldError{
			Field: str[0],
			Msg:   strings.ReplaceAll(str[2], "_", " "),
		}
	case ErrorFormDataNotFound:
		return FieldError{
			Field: str[0],
			Msg:   strings.ReplaceAll(str[2], "_", " "),
		}
	case ErrorFormLengthTooBig:
		return FieldError{
			Field: str[0],
			Msg:   strings.ReplaceAll(str[2], "_", " ") + " " + str[3],
		}
	case ErrorFormLengthTooShort:
		return FieldError{
			Field: str[0],
			Msg:   strings.ReplaceAll(str[2], "_", " ") + " " + str[3],
		}
	case ErrorFormInvalidEmail:
		return FieldError{
			Field: str[0],
			Msg:   strings.ReplaceAll(str[2], "_", " "),
		}
	case ErrorFormFieldRequired:
		return FieldError{
			Field: str[0],
			Msg:   strings.ReplaceAll(str[2], "_", " "),
		}
	case ErrorFormFieldNumeric:
		return FieldError{
			Field: str[0],
			Msg:   strings.ReplaceAll(str[2], "_", " "),
		}
	case ErrorFormFixedLength:
		return FieldError{
			Field: str[0],
			Msg:   strings.ReplaceAll(str[2], "_", " ") + " " + str[3],
		}
	case ErrorFormAlreadyExists:
		return FieldError{
			Field: str[0],
			Msg:   strings.ReplaceAll(str[2], "_", " "),
		}
	case ErrorFormFlexibleMsg:
		return FieldError{
			Field: str[0],
			Msg:   str[2],
		}
	case ErrorFormPhoto:
		return FieldError{
			Field: str[0],
			Msg:   "Photo maks " + strconv.FormatInt(config.MaxSizeUploadPhotoByte/1000000, 10) + " mb dan ext: jpg, jpeg, png",
		}
	case ErrorFormLowercase:
		return FieldError{
			Field: str[0],
			Msg:   strings.ReplaceAll(str[2], "_", " "),
		}
	case ErrorFormUppercase:
		return FieldError{
			Field: str[0],
			Msg:   strings.ReplaceAll(str[2], "_", " "),
		}
	case ErrorFormHiragana:
		return FieldError{
			Field: str[0],
			Msg:   str[2],
		}
	case ErrorFormKatakana:
		return FieldError{
			Field: str[0],
			Msg:   str[2],
		}
	case ErrorFormKana:
		return FieldError{
			Field: str[0],
			Msg:   str[2],
		}
	case ErrorFormKanji:
		return FieldError{
			Field: str[0],
			Msg:   str[2],
		}
	default:
		return FieldError{
			Field: str[0],
			Msg:   "Error Message",
		}
	}
}

func getListError(err error) Payload {
	listError := Payload{}
	fieldsError := err.(validator.ValidationErrors)

	// fmt.Println("fieldsError ", fieldsError)

	for _, fieldError := range fieldsError {
		// fmt.Println("-----------------------")
		// fmt.Println("fieldError.ActualTag() ", fieldError.ActualTag())
		// fmt.Println("fieldError.Type() ", fieldError.Type())
		// fmt.Println("fieldError.Kind() ", fieldError.Kind())
		// fmt.Println("fieldError.Namespace() ", fieldError.Namespace())
		// fmt.Println("fieldError.Field() ", fieldError.Field())
		// fmt.Println("fieldError.StructField() ", fieldError.StructField())
		// fmt.Println("fieldError.StructNamespace() ", fieldError.StructNamespace())
		// fmt.Println("fieldError.Tag() ", fieldError.Tag())

		switch fieldError.ActualTag() {
		case "notexists":
			listError[fieldError.Field()] = getFieldError(fieldError.Field(), ErrorFormAlreadyExists, ErrorFormAlreadyExists)
		case "existsdata":
			listError[fieldError.Field()] = getFieldError(fieldError.Field(), ErrorFormDataNotFound, ErrorFormDataNotFound)
		case "no_hp":
			listError[fieldError.Field()] = getFieldError(fieldError.Field(), ErrorFormFlexibleMsg, "Format nomor HP tidak benar")
		case "oneof", "exists", "weekday":
			listError[fieldError.Field()] = getFieldError(fieldError.Field(), ErrorFormSelectionInvalid, ErrorFormSelectionInvalid)
		case "lte", "max":
			listError[fieldError.Field()] = getFieldError(fieldError.Field(), ErrorFormLengthTooBig, ErrorFormLengthTooBig, fieldError.Param())
		case "gte", "min":
			listError[fieldError.Field()] = getFieldError(fieldError.Field(), ErrorFormLengthTooShort, ErrorFormLengthTooShort, fieldError.Param())
		case "email":
			listError[fieldError.Field()] = getFieldError(fieldError.Field(), ErrorFormInvalidEmail, ErrorFormInvalidEmail)
		case "required", "required_with":
			listError[fieldError.Field()] = getFieldError(fieldError.Field(), ErrorFormFieldRequired, ErrorFormFieldRequired)
		case "numeric":
			listError[fieldError.Field()] = getFieldError(fieldError.Field(), ErrorFormFieldNumeric, ErrorFormFieldNumeric)
		case "len":
			listError[fieldError.Field()] = getFieldError(fieldError.Field(), ErrorFormFixedLength, ErrorFormFixedLength, fieldError.Param())
		case "passwdComplex":
			listError[fieldError.Field()] = getFieldError(fieldError.Field(), ErrorFormFlexibleMsg, "Password harus 1 lowercase 1 uppercase 1 numberic")
		case "photo":
			listError[fieldError.Field()] = getFieldError(fieldError.Field(), ErrorFormPhoto)
		case "lowercase":
			listError[fieldError.Field()] = getFieldError(fieldError.Field(), ErrorFormLowercase, ErrorFormLowercase)
		case "uppercase":
			listError[fieldError.Field()] = getFieldError(fieldError.Field(), ErrorFormUppercase, ErrorFormUppercase)
		case "hiragana":
			listError[fieldError.Field()] = getFieldError(fieldError.Field(), ErrorFormHiragana, ErrorFormHiragana)
		case "katakana":
			listError[fieldError.Field()] = getFieldError(fieldError.Field(), ErrorFormKatakana, ErrorFormKatakana)
		case "kana":
			listError[fieldError.Field()] = getFieldError(fieldError.Field(), ErrorFormKana, ErrorFormKana)
		case "kanji":
			listError[fieldError.Field()] = getFieldError(fieldError.Field(), ErrorFormKanji, ErrorFormKanji)
		}
	}

	// fmt.Println("listError ", listError)

	return listError
}
