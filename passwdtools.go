package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
    "log"
    "flag"
)

func main() {
    
    banner, err := os.ReadFile("banner.txt")
    if err != nil {
        log.Fatalf("Ошибка при чтении файла баннера: %v", err)
        os.Exit(1)
    }

    fmt.Println(string(banner))
        
    reader := bufio.NewReader(os.Stdin)
    mode := flag.String("mode", "random", "Режим для угадывания паролей.")
    knownPart := flag.String("known", "", "Если вы знаете часть пароля, добавьте ее с помощью -known!")
    minLength := flag.Int("min", 4, "Минимальная длина пароля")
    maxLength := flag.Int("max", 12, "Максимальная длина пароля")
    wordlist := flag.String("w", "rockyou.txt", "Словарь для режима фильтрации")

    flag.Parse()
    // Получить информацию о цели от ввода пользователя
    targetInfo := getTargetInfo(reader)

    // Спросить пользователя, сколько вариантов паролей он хочет сгенерировать
    fmt.Print("Введите количество вариантов паролей, которые вы хотите сгенерировать: ")
    lengthInput, _ := reader.ReadString('\n')
    lengthInput = strings.TrimSpace(lengthInput)
    desiredLength, err := strconv.Atoi(lengthInput)
    if err != nil || desiredLength <= 0 {
        fmt.Println("Неверное число. Используется длина по умолчанию - 10000 паролей.")
        desiredLength = 10000
    }
    
    var passwordList []string

    switch *mode {
    case "random":
        passwordList = generateRandomPasswords(targetInfo, desiredLength, *minLength, *maxLength, *knownPart)
    case "smart":
        fmt.Printf("В разработке")
    case "filter":
        passwordList = filterFromWordlist(targetInfo, desiredLength, *minLength, *maxLength, *wordlist)
    }

    // Сохранить в файл
    file, err := os.Create("pontiff.txt")
    if err != nil {
        fmt.Println("Ошибка при создании файла:", err)
        return
    }
    defer file.Close()

    for _, password := range passwordList {
        file.WriteString(password + "\n")
    }

    fmt.Println("Список паролей сгенерирован и сохранен в 'pontiff.txt'")
}
