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
	"time"
)

var (
	telegramBotToken string // токен
	chatID int64 // чат id
	mapMsg map[string]string = make(map[string]string) // карта полей формы
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

func client(w http.ResponseWriter, _ *http.Request) {
	// создаем шаблон
	fp := path.Join("", "/workspace/index.html")
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
	mapMsg["name"] = r.PostFormValue("name")
	mapMsg["phone"] = r.PostFormValue("phone")
	mapMsg["message"] = r.PostFormValue("messages")
	// Проверка на заполнение (Доработать!)
	if mapMsg["name"] == "" {
		fmt.Fprintln(w, "Не заполнено поле \"Ваше Имя\"")
	} else if len(mapMsg["phone"]) > 12 || mapMsg["phone"] == ""  {
		fmt.Fprintln(w, "Не заполнено поле \"Телефон\", либо не коретно введен номер")
	} else if mapMsg["message"] == "" {
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
				mapMsg["name"],
				mapMsg["phone"],
				mapMsg["message"])
		msg := tgbotapi.NewMessage(chatID, text)
		// Парсем в "Markdown"
		msg.ParseMode = "markdown"
		// Отправляем
		bot.Send(msg)
	}
}

func chatbot(b tgbotapi.BotAPI) {
	// Создаем конфиг для общения с бото
	u := tgbotapi.NewUpdate( 0)
	u.Timeout = 60
	// используя конфиг u создаем канал в который будут прилетать новые сообщения
	updates, err := b.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}
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
		b.Send(msg)
	}
}

func main() {
	// используя токен создаем новый инстанс бота
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Panic(err)
	}
	loc, _ := time.LoadLocation("Europe/Minsk")
	date := time.Now().In(loc).Format(time.RFC3339)
	// пишем лог авторизации бота
	log.Printf("Authorized on account %s", bot.Self.UserName)
	// Отправляем сообщение в телеграм, что бот активен
	bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("%s: %s", date, "Бот запущен, ожидаются сообщения")))
	// Горутина для общения с ботом
	go chatbot(*bot)
	// Включаем сервер, прописываем роуты, статику
	http.HandleFunc("/", client)
	http.HandleFunc("/process", process)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("/workspace/css"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("/workspace/img"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("/workspace/js"))))
	http.Handle("/vendor/", http.StripPrefix("/vendor/", http.FileServer(http.Dir("/workspace/vendor"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
