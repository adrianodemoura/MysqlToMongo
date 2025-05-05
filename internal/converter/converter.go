package converter

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// cleanSpecialChars remove caracteres especiais como \r e \n
func cleanSpecialChars(str string) string {
	// Remove \r e \n
	str = strings.ReplaceAll(str, "\r", "")
	str = strings.ReplaceAll(str, "\n", "")
	// Remove espaços extras no início e fim
	return strings.TrimSpace(str)
}

// ConvertBinaryToString converte dados binários para string e garante UTF-8
func ConvertBinaryToString(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	// Se for um slice de bytes, tenta converter de base64
	if bytes, ok := value.([]byte); ok {
		// Tenta decodificar de base64
		decoded, err := base64.StdEncoding.DecodeString(string(bytes))
		if err == nil {
			// Verifica se a string é UTF-8 válida
			if utf8.Valid(decoded) {
				return cleanSpecialChars(string(decoded))
			}
			// Se não for UTF-8 válida, tenta converter para UTF-8
			return cleanSpecialChars(string(bytes))
		}
		// Se não for base64, verifica se é UTF-8 válida
		if utf8.Valid(bytes) {
			return cleanSpecialChars(string(bytes))
		}
		// Se não for UTF-8 válida, retorna string vazia
		return ""
	}

	// Se for string, verifica se é UTF-8 válida
	if str, ok := value.(string); ok {
		if utf8.ValidString(str) {
			return cleanSpecialChars(str)
		}
		return ""
	}

	return value
}

// ConvertToTimePtr converte string para *time.Time
func ConvertToTimePtr(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	// Se já for *time.Time, retorna direto
	if t, ok := value.(*time.Time); ok {
		// Converte para o fuso horário de São Paulo
		loc, _ := time.LoadLocation("America/Sao_Paulo")
		tSP := t.In(loc)
		return &tSP
	}

	str := ConvertBinaryToString(value)
	if str == nil {
		return nil
	}

	// Tenta diferentes formatos de data
	formats := []string{
		"20060102", // YYYYMMDD
		"2006-01-02 15:04:05",
		"2006-01-02",
		"02/01/2006",
		"02/01/2006 15:04:05",
	}

	if strValue, ok := str.(string); ok {
		// Se a string estiver vazia, retorna nil
		if strValue == "" {
			return nil
		}
		for _, format := range formats {
			if t, err := time.Parse(format, strValue); err == nil {
				// Converte para o fuso horário de São Paulo
				loc, _ := time.LoadLocation("America/Sao_Paulo")
				tSP := t.In(loc)
				return &tSP
			}
		}
	}

	return nil
}

// ConvertToDatePtr converte string para *time.Time (formato YYYYMMDD)
func ConvertToDatePtr(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	// Se já for *time.Time, retorna direto
	if t, ok := value.(*time.Time); ok {
		return t
	}

	str := ConvertBinaryToString(value)
	if str == nil {
		return nil
	}

	if strValue, ok := str.(string); ok {
		// Se a string estiver vazia, retorna nil
		if strValue == "" {
			return nil
		}
		if t, err := time.Parse("20060102", strValue); err == nil {
			return &t
		}
	}

	return nil
}

// ConvertToDecimal converte para Decimal128
func ConvertToDecimal(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	// Se já for Decimal128, retorna direto
	if d, ok := value.(primitive.Decimal128); ok {
		return d
	}

	// Se for string, tenta converter
	if str, ok := value.(string); ok {
		if str == "" {
			return nil
		}
		if d, err := primitive.ParseDecimal128(str); err == nil {
			return d
		}
	}

	// Se for []byte, tenta converter
	if bytes, ok := value.([]byte); ok {
		if d, err := primitive.ParseDecimal128(string(bytes)); err == nil {
			return d
		}
	}

	// Se for float64, converte para string e depois para Decimal128
	if f, ok := value.(float64); ok {
		str := fmt.Sprintf("%.2f", f)
		if d, err := primitive.ParseDecimal128(str); err == nil {
			return d
		}
	}

	return nil
}

// ConvertOptionalField trata campos opcionais
func ConvertOptionalField(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	// Se for string, verifica se está vazia ou é "0"
	if str, ok := value.(string); ok {
		if str == "" || str == "0" {
			return nil
		}
		return str
	}

	// Se for []byte, converte para string e verifica
	if bytes, ok := value.([]byte); ok {
		str := string(bytes)
		if str == "" || str == "0" {
			return nil
		}
		return str
	}

	// Se for número, verifica se é zero
	if num, ok := value.(float64); ok {
		if num == 0 {
			return nil
		}
		return num
	}

	// Se for int, verifica se é zero
	if num, ok := value.(int64); ok {
		if num == 0 {
			return nil
		}
		return num
	}

	return value
}
