package main

import (
	"TelegramGameBot/Game/engine"
	"TelegramGameBot/Game/world"
	"fmt"
	"os"
	"reflect"
	"sort"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

// Каждому пользователю, написавшему боту соответствует один игровой
// мир.
var worlds = make(map[int64]*engine.World)

var bot *tgbotapi.BotAPI

func main() {
	var err error

	// Получает token из переменной окружения
	token := os.Getenv("TGBtoken")

	// Инициализирует бота.
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	// Устанавливает время обновления.
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Создаёт канал получения обновлений.
	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {

		// Обрабатывает обновления.
		go handler(update)
	}
}

// Обработчик обновлений.
func handler(update tgbotapi.Update) {

	// Проверяет что сообщение не пустое.
	if update.Message == nil {
		return
	}

	UserName := update.Message.Chat.UserName
	ID := update.Message.Chat.ID
	Request := update.Message.Text

	// Проверяет что от пользователья пришло именно текстовое сообщение.
	if !(reflect.TypeOf(Request).Kind() == reflect.String && Request != "") {
		return
	}

	// Выводит имя пользователя и сообщение в консоль.
	fmt.Printf("%v: %v\n", UserName, Request)

	// Если игрок пишет в первый раз после запуска бота.
	if _, ok := worlds[ID]; !ok || Request == "/start" {
		if Request != "/start" {
			resp := "К сожалению ваш прогрес был утерян в связи с перезапуском игрового сервера. Игра начнётся заново :("
			msg := tgbotapi.NewMessage(ID, resp)
			bot.Send(msg)
		}

		resp := "Игра скоро начнётся..."
		msg := tgbotapi.NewMessage(ID, resp)
		bot.Send(msg)

		worlds[ID] = world.Constructor()
	}

	// Блокировка доступа к игровому миру для других горутин.
	worlds[ID].Mu.Lock()

	for {
		// Получение ответа от главного исполнителя игрового мира.
		resp, remainder := worlds[ID].ActiveHandler.Handle(Request)

		// Отправка ответа пользователю.
		msg := tgbotapi.NewMessage(ID, resp.Msg)
		bot.Send(msg) // Отправка основного сообщения.
		msg = tgbotapi.NewMessage(ID, resp.Status)
		msg.BaseChat.ReplyMarkup = Keyboard(resp.Options)
		bot.Send(msg) // Отправка статуса с клавиатурой.

		if remainder == "" {
			// Если обработан весь запрос, цикл завершается.
			break
		} else {
			// Если нет, обрабатывается оставшаяся часть.
			Request = remainder
		}
	}
	// Возвращение доступа к игровому миру для других горутин.
	// Доступ возвращается после отправки ответа для предотвращения
	// изменения порядка ответов.
	worlds[ID].Mu.Unlock()
}

// Возвращает клавиатуру из доступных опций.
func Keyboard(options [][]string) interface{} {
	l := len(options)

	// Удаление старой клавиатуры при отсутствии опций.
	if l == 0 {
		return tgbotapi.ReplyKeyboardRemove{
			RemoveKeyboard: true,
		}
	}

	str := make([][]tgbotapi.KeyboardButton, 0)
	if l == 1 {
		l := len(options[0])
		if l == 0 {
			return tgbotapi.ReplyKeyboardRemove{
				RemoveKeyboard: true,
			}
		}
		// Создание клавиатуры с двумя кнопками в строке.
		sort.Slice(options[0], func(i, j int) bool { return options[0][i] < options[0][j] })
		for i := 0; i < l; i += 2 {
			pos := []tgbotapi.KeyboardButton{
				{
					Text: options[0][i],
				},
			}
			// В случае нечётного количества опций, последняя строка
			// содержит одну кнопку.
			if i+1 < l {
				pos = append(pos, tgbotapi.KeyboardButton{
					Text: options[0][i+1],
				})
			}
			str = append(str, pos)
		}
	} else {
		// Создание клавиатуры в соответствии с options[строка][столбец].
		for _, opStr := range options {
			pos := make([]tgbotapi.KeyboardButton, 0)
			for _, opBt := range opStr {
				pos = append(pos, tgbotapi.KeyboardButton{
					Text: opBt,
				})
			}
			str = append(str, pos)
		}
	}

	return tgbotapi.ReplyKeyboardMarkup{
		Keyboard:        str,
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
	}
}
