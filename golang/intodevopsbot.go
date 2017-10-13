package main

import (
	"flag"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"os"
	"fmt"
	"net/http"
	"path"
	"html/template"
)

var (
	telegramBotToken string // токен
	chatID int64 // чат id
	sliceMsg map[string]string = make(map[string]string) // карта полей формы
)

func init() {
	// принимаем на входе флаги -telegrambottoken и chatid
	flag.StringVar(&telegramBotToken, "telegrambottoken", "", "Telegram Bot Token")
	flag.Int64Var(&chatID, "chatid", 0, "chatId to send messages")
	flag.Parse()

	// без токена не запускаемся
	if telegramBotToken == "" {
		log.Print("-telegrambottoken is required")
		os.Exit(1)
	}

	// без чат id не запускаемся
	if chatID == 0 {
		log.Print("-chatid is required")
		os.Exit(1)
	}

}

func client(w http.ResponseWriter, r *http.Request) {
	// создаем шаблон
	fp := path.Join("", "/workspace/client.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// отвечаем по шаблону
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func process(w http.ResponseWriter, r *http.Request) {
	// используя токен создаем новый инстанс бота
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	// Парсим содержимое формы
	r.ParseForm()
	// Пишем лог парсинга формы
	log.Printf("Test name: %s", r.PostForm)

	// Получаем данные из формы
	sliceMsg["name"] = r.PostFormValue("name")
	sliceMsg["phone"] = r.PostFormValue("phone")
	sliceMsg["message"] = r.PostFormValue("messages")

	// Проверка на заполнение (Доработать!)
	if sliceMsg["name"] == "" {
		fmt.Fprintln(w, "Не заполнено поле \"Ваше Имя\"")
	} else if len(sliceMsg["phone"]) > 12 || sliceMsg["phone"] == ""  {
		fmt.Fprintln(w, "Не заполнено поле \"Телефон\", либо не коретно введен номер")
	} else if sliceMsg["message"] == "" {
		fmt.Fprintln(w, "Не заполнено поле \"Сообщение\"")
	} else {
		fmt.Fprintln(w, "<span class=\"_success\">Сообщение отправлено!</span>")
		// Формируем текс сообщения в телеграм
		text := fmt.Sprintf(
			"`%s`\n" +
				"*Имя отправителя:* _%s_\n" +
				"*Телефон:* _%s_\n" +
				"*Сообщение:* %s\n",
				"Сообщение с сайта",
				sliceMsg["name"],
				sliceMsg["phone"],
				sliceMsg["message"])
		msg := tgbotapi.NewMessage(chatID, text)
		// Парсем в "Markdown"
		msg.ParseMode = "markdown"
		// Отправляем
		bot.Send(msg)
	}
}

func main() {
	// используя токен создаем новый инстанс бота
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	// пишем лог авторизации бота
	log.Printf("Authorized on account %s", bot.Self.UserName)
	// Отправляем сообщение в телеграм, что бот активен
	bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("%s", "Бот запущен, ожидает сообщений с сайт")))

	// Структура с конфигом для получения апдейтов
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Включаем сервер, прописываем роуты

	http.HandleFunc("/", client)
	http.HandleFunc("/process", process)
	http.Handle("/modalform/", http.StripPrefix("/modalform/", http.FileServer(http.Dir("/workspace/modalform"))))
	log.Fatal(http.ListenAndServe(":8080", nil))

	// используя конфиг u создаем канал в который будут прилетать новые сообщения
	updates, err := bot.GetUpdatesChan(u)
	// в канал updates прилетают структуры типа Update
	// вычитываем их и обрабатываем
	for update := range updates {
		// универсальный ответ на любое сообщение (Доработать!)
		reply := "Не знаю что сказать"
		if update.Message == nil {
			continue
		}
		// логируем от кого какое сообщение пришло
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		log.Printf("[%s] %s", update.Message.Chat.ID, update.Message.Text)
		// свитч на обработку комманд (Доработать!)
		// комманда - сообщение, начинающееся с "/"
		switch update.Message.Command() {
		case "start":
			reply = "Привет. Я телеграм-бот"
		case "hello":
			reply = "world"
		}
		// создаем ответное сообщение
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		//msg.ReplyToMessageID = update.Message.MessageID
		// отправляем
		bot.Send(msg)
	}

}
