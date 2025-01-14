# @Isgpt_bot

Этот проект представляет собой Telegram-бот, который взаимодействует с OpenAI API для обработки запросов пользователей. Проект написан на языке Go и использует фреймворк Fiber для обработки HTTP-запросов, а также GORM для работы с базой данных SQLite.

## Установка

1. *Клонируйте репозиторий:*
   bash
   git clone https://github.com/avmedvedkov/Isgpt_bot 
   

2. *Установите зависимости:*
   Убедитесь, что у вас установлены все необходимые зависимости:
   bash
   go mod tidy
   

3. *Настройте переменные окружения:*
   В файле main.go укажите значения для переменных:
   - AdminChatID — ID администратора чата.
   - APIURL — URL для API.
   - APIKey — ключ для доступа к API.
   - TelegramToken — токен вашего Telegram-бота.

## Использование

Запустите бота командой:
bash
go run main.go


Бот будет принимать запросы от пользователей и отправлять их на OpenAI API, возвращая ответы в чат Telegram.

## Требования

- Go версии 1.16 или выше
- Fiber — фреймворк для создания HTTP API
- GORM — ORM для работы с базой данных SQLite


## Авторы

Автор проекта: @avmedvedkov

## Рабочий бот в тг

@Isgpt_bot

## Лицензия

Этот проект лицензируется на условиях MIT License.


