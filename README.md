# 📦 FUX-Generator: Project Archeticture Generator

![License](https://img.shields.io/badge/License-MIT-blue.svg)
![Go version](https://img.shields.io/badge/Golang-1.23.6-blue)

**FUX-Generator** — мощный инструмент для генерации кодовой базы и структуры папок. Позволяет легко генерировать архитектуру проекта с частичной кодогенерацией для быстрого старта нового проекта.

## 🌟 Особенности

- 📁 Создание родительских и дочерних папок
- 🔄 Генерация custom types, custom errors, domains object, results, config
- ⚙️ Конфигурация через JSON-файлы


## 🚀 Быстрый старт

### Требования
- Golang 1.23.6+

### Установка Linux
```bash
git clone -b ver_1 https://github.com/devvDevu/ca-generator.git
cd ca-generator
go build -o ca-gen cmd/main.go
sudo mv ca-gen /usr/local/bin
```

### Запуск Linux
```bash
ca-gen
```

## Примеры кода:
1. [Конфигурация](https://github.com/devvDevu/ca-generator/blob/ver_1/internal/json_file_gen/example_settings.json)
