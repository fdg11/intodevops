package main

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"os"
	"fmt"
	"net/http"
	"path"
	"html/template"
	"time"
	"encoding/json"
	"strconv"
)

var (
	tmpID[10] int64
	flg bool
	i int = 0
	a int
	usgroup int
)
type Config struct {
	TelegramBotToken string
	ChatIdSite int64  `json:",string"`
	ChatIdOnline int64 `json:",string"`
}

func conf(f string) (string, int64, int64) {
	file, ok := os.Open(f)
	if ok != nil {
		log.Panic(ok)
	}
	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		log.Panic(err)
	}

	return config.TelegramBotToken, config.ChatIdSite, config.ChatIdOnline
}

func index(w http.ResponseWriter, _ *http.Request) {
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

func sendForm(w http.ResponseWriter, r *http.Request) {
	// используя токен создаем новый инстанс бота
	token, ids, _ :=conf("config.json")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	// Парсим содержимое формы
	r.ParseForm()
	// Пишем лог парсинга формы
	log.Printf("Test name: %s", r.PostForm)
	// Проверка на заполнение (Доработать!)
	if r.PostFormValue("name") == "" {
		fmt.Fprintln(w, "Не заполнено поле \"Ваше Имя\"")
	} else if len(r.PostFormValue("phone")) > 12 || r.PostFormValue("phone") == ""  {
		fmt.Fprintln(w, "Не заполнено поле \"Телефон\", либо не коретно введен номер")
	} else if r.PostFormValue("messages") == "" {
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
			r.PostFormValue("name"),
			r.PostFormValue("phone"),
			r.PostFormValue("messages"))
		msg := tgbotapi.NewMessage(ids, text)
		// Парсем в "Markdown"
		msg.ParseMode = "markdown"
		// Отправляем
		bot.Send(msg)
	}
}

func chatbot(b *tgbotapi.BotAPI) {
	_, _, ido := conf("config.json")
	// Создаем конфиг для общения с ботом
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
		reply := ""
		if update.Message == nil {
			continue
		}
		// логируем от кого какое сообщение пришло
		log.Printf("Заголовок чата: [%s] Тип чата: [%s]", update.Message.Chat.Title, update.Message.Chat.Type)
		log.Printf("Имя: [%s]", update.Message.From.FirstName)
		log.Printf("ID отпровителя: [%d]", update.Message.From.ID)
		log.Printf("ID чата:[%d] Сообщение: %s", update.Message.Chat.ID, update.Message.Text)
		if update.Message.Document != nil {
			log.Printf("Имя документа: [%s]", update.Message.Document.FileName)
		}
		if update.Message.Sticker != nil {
			log.Printf("Имя стикера: [%s]", update.Message.Sticker.Emoji)
		}

		if update.Message.Chat.IsPrivate() {
			if tmpID[i] != update.Message.Chat.ID {
				i++
			}
			log.Printf("№ в очереди: [%s]", strconv.Itoa(i))
			tmpID[i] = update.Message.Chat.ID
			reply = fmt.Sprintf(
				"*Переадрисовано от:* _%s_ _%s_\n"+"`в очереди: %d`\n\n"+"%s",
				update.Message.From.FirstName, update.Message.From.LastName, i, update.Message.Text)
			msg := tgbotapi.NewMessage(ido, reply)
			msg.ParseMode = "markdown"
			b.Send(msg)
			if update.Message.Document != nil {
				b.Send(tgbotapi.NewDocumentShare(ido,update.Message.Document.FileID))
			}
			if update.Message.Sticker != nil {
				b.Send(tgbotapi.NewStickerShare(ido,update.Message.Sticker.FileID))
			}

		} else {

			switch update.Message.Command() {
			case "start":
				usgroup = update.Message.From.ID
				for j := 0; j <= len(tmpID); j++ {
					if update.Message.CommandArguments() == strconv.Itoa(j) {
						flg = true
						reply = fmt.Sprintf("*%s*\n"+"*Специалист:* _%s_\n", "Добрый день. С вами общается:",
							update.Message.From.FirstName)
						msg := tgbotapi.NewMessage(tmpID[j], reply)
						msg.ParseMode = "markdown"
						// отправляем
						b.Send(msg)
						a = j
						break
					} else if update.Message.CommandArguments() == "" {
						reply = fmt.Sprintf("`Укажите номер очереди`")
						msg := tgbotapi.NewMessage(ido, reply)
						msg.ParseMode = "markdown"
						// отправляем
						b.Send(msg)
						break
					}
				}
			case "stop":
				if flg {
					flg = false
					reply = fmt.Sprintf("`Чат завершен`")
					msg := tgbotapi.NewMessage(tmpID[i], reply)
					msg.ParseMode = "markdown"
					// отправляем
					b.Send(msg)
				} else {
					reply = fmt.Sprintf("`Вы не начали чат с клиентом`\n"+"``` /help - для спавки```")
					msg := tgbotapi.NewMessage(ido, reply)
					msg.ParseMode = "markdown"
					// отправляем
					b.Send(msg)
				}
			case "help":
				reply = fmt.Sprintf("*Воспользуйтесь следующими командами для общения с клиентом:*\n"+
					"``` /start [номер в очереди...] - команда запустит чат с клиентом по номеру очереди```\n"+
						"``` /stop - прекратит текущий чат с клиентом```")
				msg := tgbotapi.NewMessage(ido, reply)
				msg.ParseMode = "markdown"
				// отправляем
				b.Send(msg)
			}
			if flg && usgroup == update.Message.From.ID {
				if update.Message.Text == "/start "+strconv.Itoa(a) {
					continue
				}
				msg := tgbotapi.NewMessage(tmpID[a], update.Message.Text)
				b.Send(msg)
				if update.Message.Document != nil {
					b.Send(tgbotapi.NewDocumentShare(tmpID[a],update.Message.Document.FileID))
				}
				if update.Message.Sticker != nil {
					b.Send(tgbotapi.NewStickerShare(tmpID[a],update.Message.Sticker.FileID))
				}
			}
		}
	}
}

func main() {
	token, ids, ido := conf("config.json")
	// используя токен создаем новый инстанс бота
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	// горутина
	go chatbot(bot)
	//time log
	loc, _ := time.LoadLocation("Europe/Minsk")
	date := time.Now().In(loc).Format(time.RFC3339)
	// пишем лог авторизации бота
	log.Printf("Authorized on account %s", bot.Self.UserName)
	log.Printf("Token: %s; IdSite %d; IdOnline %d", token, ids, ido)
	// Отправляем сообщение в телеграм, что бот активен
	bot.Send(tgbotapi.NewMessage(ids, fmt.Sprintf("%s: %s", date, "Бот запущен, ожидаются сообщения")))
	// Горутина для общения с ботом
	// Включаем сервер, прописываем роуты, статику
	http.HandleFunc("/", index)
	http.HandleFunc("/process", sendForm)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("/workspace/css"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("/workspace/img"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("/workspace/js"))))
	http.Handle("/vendor/", http.StripPrefix("/vendor/", http.FileServer(http.Dir("/workspace/vendor"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}