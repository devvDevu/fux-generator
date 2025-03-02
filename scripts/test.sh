#!/bin/bash
set -euo pipefail

cleanup() {
    echo "Прерывание выполнения. Папка example_name_of_project сохранена."
    exit 1
}

trap cleanup SIGINT SIGTERM

run_program() {
    if ! go run cmd/main.go; then
        echo "Программа завершилась с ошибкой. Папка example_name_of_project сохранена."
        exit 1
    fi
}

for i in {1..100}; do
    echo "Iteration $i"
    
    # Запуск с проверкой на ошибку
    if ! run_program; then
        exit 1
    fi
    
    sleep 1
    rm -rf example_name_of_project || {
        echo "Ошибка при удалении папки"
        exit 1
    }
done

echo "Все 20 итераций выполнены успешно"
