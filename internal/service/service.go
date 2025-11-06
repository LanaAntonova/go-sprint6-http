package service

import (
	"strings"
	"unicode"

	"github.com/LanaAntonova/go-sprint6-http/pkg/morse"
)

// isMorse проверяет, похожа ли строка на код Морзе
func isMorse(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return false
		}
	}
	trimmed := strings.Trim(s, " .-/\t\n\r")
	return trimmed == ""
}

// Convert автоматически определяет тип и конвертирует
func Convert(input string) (string, error) {
	if isMorse(input) {
		result := morse.ToText(input)
		if result == "" {
			return "", &ErrConversion{Input: input, Reason: "no valid Morse code decoded"}
		}
		return result, nil
	}
	result := morse.ToMorse(input)
	return result, nil
}

// ErrConversion — ошибка конвертации
type ErrConversion struct {
	Input  string
	Reason string
}

func (e *ErrConversion) Error() string {
	return "conversion error: " + e.Reason + " (input: " + e.Input + ")"
}
