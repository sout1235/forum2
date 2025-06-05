#!/bin/bash

# Остановка выполнения при ошибке
set -e

# Проверка наличия PostgreSQL
if ! command -v psql &> /dev/null; then
    echo "PostgreSQL не установлен. Пожалуйста, установите PostgreSQL."
    exit 1
fi

# Проверка наличия Go
if ! command -v go &> /dev/null; then
    echo "Go не установлен. Пожалуйста, установите Go."
    exit 1
fi

# Создание тестовой базы данных
echo "Создание тестовой базы данных..."
psql -U postgres -f testdata.sql

# Установка зависимостей
echo "Установка зависимостей..."
go mod tidy

# Запуск тестов
echo "Запуск интеграционных тестов..."
go test -v ./...

# Очистка
echo "Очистка тестовой базы данных..."
psql -U postgres -c "DROP DATABASE IF EXISTS forum_test;"

echo "Тесты завершены." 