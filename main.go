package main
// @Isgpt_bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/gofiber/fiber/v2"
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Константы
const (
	AdminChatID   = 0
	APIURL        = "1"
	APIKey        = "1"
	TelegramToken = "1"
)

// Модель пользователя
type User struct {
	gorm.Model
	ChatID int64 `gorm:"uniqueIndex"`
}

// Структура для запроса к OpenAI API
type GPTRequest struct {
	Model       string       `json:"model"`
	Messages    []GPTMessage `json:"messages"`
}

// Структура для сообщения OpenAI
type GPTMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Инициализация базы данных SQLite
func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	db.AutoMigrate(&User{})
  
	return db
}

// Функция для отправки запроса к OpenAI API
func sendGPT(text string) (string, error) {
	client := &http.Client{}

	// Формирование тела запроса
	requestBody := GPTRequest{
		Model: "gpt-4o",
		Messages: []GPTMessage{
			{Role: "user", Content: text},
		},
	}

	// Преобразование структуры в JSON
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", APIURL, bytes.NewBuffer(requestJSON))
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer " + APIKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("API error: %s", body)
	}

	// Чтение ответа
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	// Извлечение ответа из API
	if choices, ok := response["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				return message["content"].(string), nil
			}
		}
	}

	return "", fmt.Errorf("unexpected API response format")
}

// Основная функция
func main() {
	// Инициализация базы данных
	db := initDB()

	// Инициализация Telegram бота
	bot, err := tb.NewBot(tb.Settings{
		Token:  TelegramToken,
		Poller: &tb.LongPoller{Timeout: 10},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Обработчик команды /start
	bot.Handle("/start", func(m *tb.Message) {
		bot.Send(m.Sender, "Send me a message to get a response from GPT-4.")
	})

	// Обработчик команды /count_users
	bot.Handle("/count_users", func(m *tb.Message) {
		if m.Sender.ID == AdminChatID {
			var count int64
			db.Model(&User{}).Count(&count)
			bot.Send(m.Sender, fmt.Sprintf("Count users: %d", count))
		} else {
			bot.Send(m.Sender, "You are not authorized to use this command.")
		}
	})

	// Обработчик текстовых сообщений
	bot.Handle(tb.OnText, func(m *tb.Message) {
		bot.Send(m.Sender, tb.Typing)

		chatID := m.Sender.ID
		// Добавление пользователя в базу данных
		var user User
		if err := db.FirstOrCreate(&user, User{ChatID: chatID}).Error; err != nil {
			log.Println("Failed to create user: ", err)
		}

		// Ответ от GPT
		response, err := sendGPT(m.Text)
		if err != nil {
			bot.Send(m.Sender, "Error: " + err.Error())
			return
		}
		bot.Send(m.Sender, response)
	})

	// Запуск бота в отдельном потоке
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		bot.Start()
		wg.Done()
	}()
	// Инициализация FastAPI сервера
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Listen(":8041")
	wg.Wait()
}