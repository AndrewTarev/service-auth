package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"service-auth/internal/configs"
)

// SetupLogger настраивает глобальный логгер logrus в соответствии с конфигурацией.
func SetupLogger(cfg *configs.LoggerConfig) {
	// Устанавливаем уровень логирования
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		// Если уровень не распознан, выводим предупреждение и используем уровень по умолчанию
		logrus.Warnf("Не удалось установить уровень логирования '%s', используется уровень по умолчанию: info", cfg.Level)
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	// Устанавливаем формат вывода
	switch cfg.Format {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	default:
		// По умолчанию вывод в текстовом формате с полными временными метками
		logrus.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.RFC3339,
			FullTimestamp:   true,
		})
	}

	// Устанавливаем вывод логов (файл или консоль)
	if cfg.OutputFile != "" {
		// Пытаемся открыть файл для записи
		file, err := os.OpenFile(cfg.OutputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			// В случае ошибки открытия файла, выводим предупреждение и используем консоль
			logrus.Warnf("Не удалось записать логи в файл '%s', используется вывод в консоль: %v", cfg.OutputFile, err)
			logrus.SetOutput(os.Stdout)
		} else {
			// Если файл открыт успешно, устанавливаем его как вывод для логов
			logrus.SetOutput(file)
			// Закрытие файла нужно организовать в вызывающем коде или через хук
		}
	} else {
		// Если путь к файлу не задан, выводим логи в консоль
		logrus.SetOutput(os.Stdout)
	}
}
