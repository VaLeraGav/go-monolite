package helper

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"os"
)

func GetProjectPath() string {
	projectPath := os.Getenv("PROJECT_PATH")
	if projectPath != "" {
		return projectPath
	}

	currentDir, _ := os.Getwd()
	return currentDir
}

// читабельный вывод в консоль
func PrettyPrint(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func Unmarshal(data []byte, v any) (string, error) {
	err := json.Unmarshal(data, v)
	if err != nil {
		if strings.Contains(err.Error(), "invalid UUID length") {
			return "UUID имеет некорректную длину", err
		}
		return "Получен некорректный формат JSON", err
	}
	return "", nil
}

func GenerateCode(length int) string {
	if length <= 0 {
		return ""
	}

	rand.Seed(time.Now().UnixNano())

	code := ""
	for i := 0; i < length; i++ {
		digit := rand.Intn(10) // число от 0 до 9
		code += strconv.Itoa(digit)
	}

	return code
}
