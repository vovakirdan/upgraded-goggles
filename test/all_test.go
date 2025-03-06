package tests

import (
    "os"
    "testing"

    "upgraded-goggles/internal/logger"
)

func TestMain(m *testing.M) {
    // Инициализируем логгер, чтобы во всех тестах logger.Logger != nil
    err := logger.InitLogger("logs/test.log")
    if err != nil {
        panic(err)
    }

    // Запускаем все тесты
    code := m.Run()
    os.Exit(code)
}
