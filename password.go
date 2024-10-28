package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

const (
	LowerLetters = "abcdefghijklmnopqrstuvwxyz"
	UpperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Digits       = "0123456789"
	Symbols      = "!@#$%&*()_+-={}[]:?,."
)

// Генерация паролей путём объединения информации о пользователе и добавления случайных символов
func generateRandomPasswords(targetInfo []string, desiredLength int, passwdLengthMin int, passwdLengthMax int, knownPart string) []string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	passwords := make(map[string]struct{})
	mergedInfo := mergeUserInfo(targetInfo)

	bar := progressbar.Default(int64(desiredLength))

	var known string
	var position int

	if strings.HasPrefix(knownPart, "P") && strings.HasSuffix(knownPart, "P") {
		_, err := fmt.Sscanf(knownPart, "P%dP", &position)
		if err != nil {
			fmt.Println("Ошибка чтения позиции:", err)
			return nil
		}
		known = knownPart[2 : len(knownPart)-1]
	} else {
		known = knownPart
		position = 0
	}

	// Генерация паролей из объединённой информации
	for _, base := range mergedInfo {
		if len(passwords) >= desiredLength {
			break
		}
		addPasswordWithKnown(base, desiredLength, passwdLengthMin, passwdLengthMax, known, position, passwords, bar, r)
	}

	// Если сгенерировано недостаточно паролей, заполняем оставшееся пространство
	for len(passwords) < desiredLength {
		base := mergedInfo[r.Intn(len(mergedInfo))]
		password := generatePassword(base, passwdLengthMin, passwdLengthMax, known, position, r)
		passwords[password] = struct{}{}
		bar.Add(1)
	}

	result := make([]string, 0, len(passwords))
	for password := range passwords {
		result = append(result, password)
	}

	return result
}

func mergeUserInfo(targetInfo []string) []string {
	var merged []string

	// Объединение различных полей
	for i := 0; i < len(targetInfo); i++ {
		for j := i + 1; j < len(targetInfo); j++ {
			merged = append(merged, targetInfo[i]+targetInfo[j]) // Пример: john1990
			merged = append(merged, targetInfo[j]+targetInfo[i]) // Пример: 1990john
		}
	}

	// Также добавляем одиночные поля
	merged = append(merged, targetInfo...)

	return merged
}

// Случайное добавление специальных символов с низкой вероятностью
func addRandomSymbols(base string, r *rand.Rand) string {
	const specialProbability = 0.01 // 1% вероятность добавления символа

	var result strings.Builder
	for _, char := range base {
		result.WriteRune(char)

		// Случайно решаем добавить символ
		if r.Float64() < specialProbability {
			symbol := Symbols[r.Intn(len(Symbols))]
			result.WriteByte(symbol)
		}
	}
	return result.String()
}

func addPasswordWithKnown(base string, desiredLength int, minLength int, maxLength int, known string, position int, passwords map[string]struct{}, bar *progressbar.ProgressBar, r *rand.Rand) {
	for i := 0; i < 3 && len(passwords) < desiredLength; i++ {
		password := generatePassword(base, minLength, maxLength, known, position, r)
		passwords[password] = struct{}{}
		bar.Add(1)
	}
}

func generatePassword(base string, minLength int, maxLength int, known string, position int, r *rand.Rand) string {
	// Корректируем длину, чтобы гарантировать наличие места для известной части
	length := r.Intn(maxLength-minLength+1) + minLength
	if position < 1 {
		position = 1
	}
	if position+len(known)-1 > length {
		length = position + len(known) - 1
	}

	var result strings.Builder

	// Создание префикса перед известной частью
	prefixLength := position - 1
	randomPrefix := addRandomSymbols(base, r)

	if prefixLength > len(randomPrefix) {
		prefixLength = len(randomPrefix)
	}
	result.WriteString(randomPrefix[:prefixLength]) // Добавляем префикс

	// Вставка известной части в правильное положение
	result.WriteString(known)

	// Заполняем оставшуюся часть пароля случайными символами
	suffixLength := length - result.Len()
	if suffixLength > 0 {
		randomSuffix := addRandomSymbols(base, r)
		if suffixLength > len(randomSuffix) {
			suffixLength = len(randomSuffix)
		}
		result.WriteString(randomSuffix[:suffixLength]) // Добавляем суффикс
	}

	// Если результат превышает необходимую длину, обрезаем его
	if result.Len() > length {
		return result.String()[:length]
	}
	return result.String()
}

/* Фильтр */

func filterFromWordlist(targetInfo []string, desiredLength int, minPaswdLength int, maxPaswdLength int, wordlist string) []string {
	passwords := make(map[string]struct{})

	file, err := os.Open(wordlist)
	if err != nil {
		fmt.Println("Ошибка открытия списка слов: ", err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	bar := progressbar.Default(-1)

	for scanner.Scan() {
		word := scanner.Text()

		if len(word) >= minPaswdLength && len(word) <= maxPaswdLength {
			for _, target := range targetInfo {
				if strings.Contains(word, target) {
					passwords[word] = struct{}{}
					break
				}
			}
		}

		bar.Add(1)
	}

	// Проверяем на ошибки, которые могли возникнуть во время сканирования
	if err := scanner.Err(); err != nil { 
		fmt.Println("Ошибка чтения списка слов: ", err)
		return nil
	}

	result := make([]string, 0, len(passwords))
	for password := range passwords {
		result = append(result, password)
	}

	if len(result) > desiredLength {
		result = result[:desiredLength]
	}

	return result
}
