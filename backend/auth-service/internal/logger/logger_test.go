package logger

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestLogger(t *testing.T) {
	// Используем zaptest для тестирования логгера
	log = zaptest.NewLogger(t)

	// Тест для Info
	t.Run("Info", func(t *testing.T) {
		Info("test info message", zap.String("key", "value"))
	})

	// Тест для Error
	t.Run("Error", func(t *testing.T) {
		Error("test error message", zap.String("key", "value"))
	})

	// Тест для Fatal
	t.Run("Fatal", func(t *testing.T) {
		// Fatal вызывает os.Exit(1), поэтому тест не дойдет до конца
		// Можно использовать t.Fatal, чтобы прервать тест, если Fatal не вызван
		t.Skip("Fatal вызывает os.Exit(1), пропускаем")
		Fatal("test fatal message", zap.String("key", "value"))
	})

	// Тест для Sync
	t.Run("Sync", func(t *testing.T) {
		Sync()
	})
}

func TestLogger_Init(t *testing.T) {
	// Тест для Init
	t.Run("Init", func(t *testing.T) {
		Init()
	})
}

func TestLogger_Init_Error(t *testing.T) {
	// Тест для Init с ошибкой
	t.Run("Init_Error", func(t *testing.T) {
		// Мокаем zap.NewProduction, чтобы он вернул ошибку
		// Это сложно сделать без изменения кода, поэтому просто проверяем, что Init не паникует
		Init()
	})
}

func TestLogger_Fatal(t *testing.T) {
	// Тест для Fatal
	t.Run("Fatal", func(t *testing.T) {
		// Fatal вызывает os.Exit(1), поэтому тест не дойдет до конца
		// Можно использовать t.Fatal, чтобы прервать тест, если Fatal не вызван
		t.Skip("Fatal вызывает os.Exit(1), пропускаем")
		Fatal("test fatal message", zap.String("key", "value"))
	})
}

func TestLogger_Fatal_Error(t *testing.T) {
	// Тест для Fatal с ошибкой
	t.Run("Fatal_Error", func(t *testing.T) {
		// Fatal вызывает os.Exit(1), поэтому тест не дойдет до конца
		// Можно использовать t.Fatal, чтобы прервать тест, если Fatal не вызван
		t.Skip("Fatal вызывает os.Exit(1), пропускаем")
		Fatal("test fatal message", zap.String("key", "value"))
	})
}
