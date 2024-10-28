package main

import (
    "bufio"
    "fmt"
    "strings"
)

func getTargetInfo(reader *bufio.Reader) []string {
        var targetInfo []string

        fmt.Println("Введите следующую целевую информацию (нажимайте Enter после каждой строки):")

        fmt.Print("Имя: ")
        firstName, _ := reader.ReadString('\n')
        targetInfo = append(targetInfo, strings.TrimSpace(firstName))

        fmt.Print("Фамилия: ")
        lastName, _ := reader.ReadString('\n')
        targetInfo = append(targetInfo, strings.TrimSpace(lastName))

        fmt.Print("Имя питомца: ")
        petName, _ := reader.ReadString('\n')
        targetInfo = append(targetInfo, strings.TrimSpace(petName))

        fmt.Print("Спорт: ")
        sport, _ := reader.ReadString('\n')
        targetInfo = append(targetInfo, strings.TrimSpace(sport))

        fmt.Print("День рождения (День): ")
        birthDay, _ := reader.ReadString('\n')
        targetInfo = append(targetInfo, strings.TrimSpace(birthDay))

        fmt.Print("День рождения (Месяц): ")
        birthMonth, _ := reader.ReadString('\n')
        targetInfo = append(targetInfo, strings.TrimSpace(birthMonth))

        fmt.Print("День рождения (Год): ")
        birthYear, _ := reader.ReadString('\n')
        targetInfo = append(targetInfo, strings.TrimSpace(birthYear))
        
        fmt.Print("Имя пользователя: ")
        userName, _ := reader.ReadString('\n')
        targetInfo = append(targetInfo, strings.TrimSpace(userName))

        fmt.Print("Известный пароль: ")
        knownPasswd, _ := reader.ReadString('\n')
        targetInfo = append(targetInfo, strings.TrimSpace(knownPasswd))

        cleanTargetInfo := cleanArray(targetInfo)
        return cleanTargetInfo
}

func cleanArray(array []string) []string {
        var clean []string
        for _, str := range array {
                if str != "" {
                        clean = append(clean, str)
                }
        }
        return clean
}
