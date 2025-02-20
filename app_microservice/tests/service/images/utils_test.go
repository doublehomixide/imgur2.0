package service

import (
	"pictureloader/app_microservice/service"
	"strings"
	"testing"
)

func TestGenerateSK(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedBase string
	}{
		{"Обычная строка с пробелами и заглавными буквами", "Test Description", "test_description"},
		{"Строка с апострофами", "User's Picture", "users_picture"},
		{"Строка с несколькими пробелами подряд", "Multiple   Spaces", "multiple___spaces"},
		{"Строка с символами, которые не нужно изменять", "image-123.jpg", "image-123.jpg"},
		{"Пустая строка", "", ""},
		{"Строка с только пробелами", "   ", "___"},
		{"Строка с только апострофами", "'''", ""},
		{"Строка с символами верхнего регистра", "ALLUPPERCASE", "alluppercase"},
		{"Строка с символами нижнего регистра", "already_lowercase", "already_lowercase"},
		{"Строка с пробелами в начале и конце", "  Trim Me  ", "__trim_me__"},
		{"Строка с не-ASCII символами", "Café au Lait", "café_au_lait"},
		{"Строка с цифрами", "Image 123", "image_123"},
		{"Строка с подчеркиваниями", "already_has_underscores", "already_has_underscores"},
		{"Строка с разными символами", "Test!@#Description$%^", "test!@#description$%^"},
		{"Строка с эмодзи", "Test 😊 Description", "test_😊_description"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.GenerateSK(tt.input)

			expectedLength := len(tt.expectedBase) + 8
			if len(result) != expectedLength {
				t.Errorf("Test '%s' failed: expected length %d, got %d", tt.name, expectedLength, len(result))
			}

			if !strings.HasPrefix(result, tt.expectedBase) {
				t.Errorf("Test '%s' failed: expected result to start with '%s', got '%s'", tt.name, tt.expectedBase, result)
			}

			randomPart := result[len(tt.expectedBase):]
			for _, c := range randomPart {
				if !((c >= 'a' && c <= 'f') || (c >= '0' && c <= '9')) {
					t.Errorf("Test '%s' failed: random part contains invalid hex character '%c'", tt.name, c)
				}
			}
		})
	}
}
