package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// SetupLogger настраивает глобальный логгер logrus.
func SetupLogger(level string, format string, outputFile string) {
	// Устанавливаем уровень логирования
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		// Если уровень не распознан, выводим предупреждение и используем уровень по умолчанию
		logrus.Warnf("Не удалось установить уровень логирования '%s', используется уровень по умолчанию: info", level)
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)

	// Устанавливаем формат вывода
	switch format {
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
	if outputFile != "" {
		// Пытаемся открыть файл для записи
		file, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			// В случае ошибки открытия файла, выводим предупреждение и используем консоль
			logrus.Warnf("Не удалось записать логи в файл '%s', используется вывод в консоль: %v", outputFile, err)
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
