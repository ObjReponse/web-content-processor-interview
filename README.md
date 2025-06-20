# 🚀 Live Coding Interview: Web Content Processor

## Добро пожаловать!

Это интервью на позицию Go разработчика. Вам предстоит реализовать систему обработки веб-контента.

## 📋 Задание

Необходимо реализовать **package** для обработки веб-контента, который получает данные из канала и выполняет анализ веб-страниц.

### Входные данные

Система должна слушать канал со структурами:

```go
type Entity struct {
    Link  string
    Title string
}
```

### Требуемая функциональность

1. **HTTP запрос**
    - Выполнить GET запрос по ссылке
    - Сохранить статус ответа
    - Сохранить исходный код страницы

2. **Парсинг контента**
    - Извлечь все ссылки из HTML
    - Найти ссылки на .m3u8 файлы

3. **Классификация**
    - Проверить наличие ключевых слов в title и body
    - Проверить наличие m3u8 ссылок
    - На основе критериев выше принять решение о дальнейшей обработке

4. **Результат**
    - Вернуть структурированный результат обработки

### 🎯 Бонусы

- **Graceful shutdown** с context
- **Thread safety** с mutex
- **Error handling** для различных сценариев
- **Настраиваемые таймауты**

## 🛠 Быстрый старт

### Вилка, но не зарплатная

1. Создайте fork репозитория

### Локальная разработка

```bash
# Клонируйте репозиторий
git clone <your-fork-url>
cd web-content-processor-interview

# Запустите тесты
make test

# Запустите пример
make run

# Проверьте код
make lint
```

### VS Code Dev Container

1. Откройте проект в VS Code
2. Нажмите "Reopen in Container" когда появится уведомление
3. Среда готова к работе

## 📁 Структура проекта

```
processor/
├── processor.go      # ← ОСНОВНОЙ ФАЙЛ ДЛЯ РЕАЛИЗАЦИИ
├── types.go         # Типы данных (можно расширить)
└── processor_test.go # Тесты (запускайте для проверки)
```

## 🧪 Тестирование

```bash
# Запуск всех тестов
go test ./...

# Запуск с подробным выводом
go test -v ./processor

# Тесты с coverage
go test -cover ./processor
```

## 💡 Подсказки

- Начните с реализации базовой структуры в `processor/processor.go`
- Используйте стандартные библиотеки: `net/http`, `html`, `strings`
- Не забывайте про обработку ошибок
- Тесты помогут понять ожидаемое поведение
- Можно использовать регулярные выражения для парсинга

## 📝 Критерии оценки

1. **Корректность** - работает ли код
2. **Go идиомы** - каналы, горутины, error handling
3. **Архитектура** - читаемость и структура кода
4. **Тесты** - проходят ли существующие тесты

## ⏰ Время: 45-60 минут

---

**Удачи! Если есть вопросы - не стесняйтесь спрашивать.**