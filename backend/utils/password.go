package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword генерирует хеш пароля
func HashPassword(password string) (string, error) {
	// Используем простой SHA-256 хеш вместо bcrypt
	// только для тестирования (для упрощения диагностики)
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("sha256:%s", hex.EncodeToString(hash[:])), nil
}

// CheckPasswordHash проверяет, соответствует ли пароль хешу
func CheckPasswordHash(password, hash string) bool {
	fmt.Printf("PASSWORD DEBUG: Сравнение пароля и хеша\n")
	fmt.Printf("PASSWORD DEBUG: пароль: '%s'\n", password)
	fmt.Printf("PASSWORD DEBUG: хеш: '%s'\n", hash)

	// Если хеш начинается с sha256:, используем SHA-256
	if len(hash) > 7 && hash[:7] == "sha256:" {
		passwordHash := sha256.Sum256([]byte(password))
		expected := fmt.Sprintf("sha256:%s", hex.EncodeToString(passwordHash[:]))
		result := expected == hash
		fmt.Printf("PASSWORD DEBUG: Сравнение SHA-256, результат: %v\n", result)
		return result
	}

	// Иначе пробуем bcrypt
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Printf("PASSWORD DEBUG: Ошибка при сравнении bcrypt: %v\n", err)
		return false
	}
	fmt.Printf("PASSWORD DEBUG: Пароль верный (bcrypt)\n")
	return true
}
