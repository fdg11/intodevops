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
	tmpID = map[int]map[string]int64{}
 	count int = 0
	usGroupId = map[string]int{}
	join bool
)

type Config struct {
	TelegramBotToken string
	ChatIdSite int64  `json:",string"`
	ChatIdOnline int64 `json:",string"`
}

func replyMsg(text string, id int64) tgbotapi.MessageConfig {
	reply := fmt.Sprintf(text)
	msg := tgbotapi.NewMessage(id, reply)
	msg.ParseMode = "markdown"
	return msg
}

/*func removeAtIndex(source []int, index int) []int {
	if index == 0 {
		source = source[1:]
	} else {
		for i, v := range source {
			if v == index {
				copy(source[i:], source[i+1:])
				source[len(source)-1] = 0
				source = source[:len(source)-1]
			}
		}
	}
	return source
}*/

func conf(fileDir string) (string, int64, int64) {
	file, ok := os.Open(fileDir)
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

func chatBot(b *tgbotapi.BotAPI) {
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

		if update.Message.Chat.IsPrivate() && !update.Message.IsCommand() {
			loop:
			for n:=0 ; n < len(tmpID); n++ {
				for _, id := range tmpID[n] {
					if id == update.Message.Chat.ID {
						log.Printf("№ в очереди: [%d]", n)
						reply = fmt.Sprintf(
							"*Переадрисовано от:* _%s_ _%s_\n"+"`в очереди: %d`\n\n"+"%s",
							update.Message.From.FirstName, update.Message.From.LastName, n, update.Message.Text)
						join = true
						log.Println(tmpID)
						log.Println("Cчетчик: ", count)
						log.Println(join)
						break loop
					} else {
						join = false
					}
				}
			}
			if !join {
				tmpID[count] = map[string]int64{}
				tmpID[count][update.Message.From.FirstName+" "+update.Message.From.LastName] = update.Message.Chat.ID
				log.Printf("№ в очереди: [%d]", count)
				log.Println(tmpID)
				log.Println("Cчетчик: ", count)
				log.Println(join)
				reply = fmt.Sprintf(
					"*Переадрисовано от:* _%s_ _%s_\n"+"`в очереди: %d`\n\n"+"%s",
					update.Message.From.FirstName, update.Message.From.LastName, count, update.Message.Text)
				count++
			}
			msg := tgbotapi.NewMessage(ido, reply)
			msg.ParseMode = "markdown"
			b.Send(msg)
			if update.Message.Document != nil {
				b.Send(tgbotapi.NewDocumentShare(ido,update.Message.Document.FileID))
			}
			if update.Message.Sticker != nil {
				b.Send(tgbotapi.NewStickerShare(ido,update.Message.Sticker.FileID))
			}

		} else if update.Message.Chat.IsGroup() {
			switch update.Message.Command() {
			case "start":
				for j := 0; j < len(tmpID); j++ {
					if update.Message.CommandArguments() == strconv.Itoa(j) {
						reply = fmt.Sprintf("*%s*\n"+"*Специалист:* _%s_\n", "Добрый день. С вами общается:",
							update.Message.From.FirstName)
						for _, id := range tmpID[j] {
							msg := tgbotapi.NewMessage(id, reply)
							msg.ParseMode = "markdown"
							b.Send(msg)
						}
						log.Println(usGroupId)
						usGroupId[update.Message.From.FirstName] = j
						log.Println(usGroupId)
						break
					} else if update.Message.CommandArguments() == "" {
						b.Send(replyMsg("`Укажите номер очереди`\n"+"``` /help - для справки```", ido))
						break
					}
					arg, err := strconv.Atoi(update.Message.CommandArguments())
					if err != nil {
						log.Panic(err)
						break
					}
					if arg >= count {
						b.Send(replyMsg("`Очереди не существует`\n"+"``` /help - для справки```", ido))
						break
					}
				}
			case "stop":
				if len(usGroupId) != 0 {
					for key, value := range usGroupId {
						if key == update.Message.From.FirstName {
							reply = fmt.Sprintf("*Чат c*"+" _%s_ "+"*завершен.*", key)
							for _, id := range tmpID[value] {
								msg := tgbotapi.NewMessage(id, reply)
								msg.ParseMode = "markdown"
								b.Send(msg)
							}
							delete(usGroupId, update.Message.From.FirstName)
							log.Println(usGroupId)
							break
						} else {
							continue
 						}
					}
				} else {
					b.Send(replyMsg("`Все очереди свободны либо Вы не заняли очередь`", ido))
				}
			case "list":
				b.Send(replyMsg("`Занятых:`", ido))
				for keyQ, _ := range tmpID {
					for key, value := range usGroupId {
						if keyQ == value {
							reply = fmt.Sprintf("``` %s, занял очередь: %d```", key, value)
							msg := tgbotapi.NewMessage(ido, reply)
							msg.ParseMode = "markdown"
							b.Send(msg)
						}
					}
				}
				b.Send(replyMsg("`Всего очередей:`", ido))
				for keyQ, valueQ :=range tmpID {
					for name, _ := range valueQ {
						reply = fmt.Sprintf("```  Очеред: %d, клиент: %s```", keyQ, name)
						msg := tgbotapi.NewMessage(ido, reply)
						msg.ParseMode = "markdown"
						b.Send(msg)
					}
				}
			case "help":
				b.Send(replyMsg("*Воспользуйтесь следующими командами для общения с клиентом:*\n"+
					"``` /start [номер в очереди...] - команда запустит чат с клиентом по номеру очереди```\n"+
					"``` /stop - прекратит текущий чат с клиентом```\n"+
					"``` /list - выводит список занятых очередей специалистами```", ido))
			}
			if len(usGroupId) != 0 {
				for key, value := range usGroupId {
					if key == update.Message.From.FirstName {
						if update.Message.IsCommand() {
							continue
						}
						for _, id := range tmpID[value] {
							msg := tgbotapi.NewMessage(id, update.Message.Text)
							b.Send(msg)
							if update.Message.Document != nil {
								b.Send(tgbotapi.NewDocumentShare(id, update.Message.Document.FileID))
							}
							if update.Message.Sticker != nil {
								b.Send(tgbotapi.NewStickerShare(id, update.Message.Sticker.FileID))
							}
						}
						log.Println(usGroupId)
//						log.Println(usGroupFlg)
						log.Println(tmpID)
						break
					}
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
	go chatBot(bot)
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