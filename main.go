package main

import (
	"TelegramGameBot/Game/engine"
	mediafiles "TelegramGameBot/Game/mediaFiles"
	"TelegramGameBot/Game/world"
	"fmt"
	"os"
	"reflect"
	"sort"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

// Каждому пользователю соответствует один игровой мир.
var worlds = make(map[int64]*engine.World)

// Объект API Telegram.
var bot *tgbotapi.BotAPI

func main() {
	// Получает token из переменной окружения
	token := os.Getenv("TGBtoken")

	// Инициализирует бота.
	var err error
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

	ID := update.Message.Chat.ID
	UserName := update.Message.Chat.UserName
	Request := update.Message.Text

	// Обработка паники, что бы фатальный запрос не положил бота.
	defer panicHandler(ID, UserName, Request)

	// Проверяет что от пользователья пришло именно текстовое сообщение.
	if !(reflect.TypeOf(Request).Kind() == reflect.String && Request != "") {
		return
	}

	// Выводит имя пользователя и сообщение в консоль.
	fmt.Printf("%v: %v\n", UserName, Request)

	// Преднамеренная паника для теста обработки паники.
	if Request == "Panic" {
		panic("Игрок паникует!")
	}

	// Прерывает обработку, если мир в статусе создания.
	_, ok := worlds[ID]
	if ok && worlds[ID] == nil {
		return
	}

	var W *engine.World

	// Запрос на создание мира (команда /start или вынужденный запрос - отсутствие в worlds ключа ID)
	if !ok || Request == "/start" {

		// Создает в карте миров ключ с nil значением для того чтобы новые запросы игнорировались до завершения создания (стр.71).
		worlds[ID] = nil

		// Информирует пользователя о потере прогресса в случае вынужденного запроса.
		if Request != "/start" {
			resp := "К сожалению ваш прогрес был утерян в связи с перезапуском игрового сервера. Игра начнётся заново :("
			msg := tgbotapi.NewMessage(ID, resp)
			bot.Send(msg)
		}

		// Информирует пользователя о скором начале игры.
		resp := "Игра скоро начнётся..."
		msg := tgbotapi.NewMessage(ID, resp)
		bot.Send(msg)

		msig := tgbotapi.NewPhotoShare(ID, mediafiles.Image["Juno"])
		bot.Send(msig)

		// Создает новый мир.
		W = world.Constructor()
		worlds[ID] = W
	}

	// Вызов внутриигрового обработчика команд.
	inGameHandler(ID, Request)
}

// Обработчик паникующих запросов пользователя.
func panicHandler(ID int64, UserName string, Request string) {
	if err := recover(); err != nil {
		// Выводит в консоль информацию о паникующем запросе.
		fmt.Printf("fatal error:\n      user:    %s\n   request: %s\n   error:   %s\n", UserName, Request, fmt.Sprint(err))

		// Сообщение пользователю.
		resp := "Возникла критическая ошибка, которая вероятно скоро будет исправлена."

		// Другое сообщение в случае преднамеренной паники пользователя.
		if Request == "Panic" {
			resp = "Не паникуйте!"
		}
		msg := tgbotapi.NewMessage(ID, resp)
		bot.Send(msg)
	}
}

// Внутриигровой обработчик пользовательских запросов.
func inGameHandler(ID int64, Request string) {
	// Блокировка доступа к игровому миру для других горутин.
	worlds[ID].Mu.Lock()

	for {
		// Получение ответа от главного исполнителя игрового мира.
		resp, remainder := worlds[ID].ActiveHandler.Handle(Request)
		// Отправка ответа пользователю.
		msg := tgbotapi.NewMessage(ID, resp.Msg)
		bot.Send(msg) // Отправка основного сообщения.
		msig := tgbotapi.NewPhotoShare(ID, resp.Img)
		bot.Send(msig) // Отправка изображения.
		msg = tgbotapi.NewMessage(ID, resp.Status)
		msg.BaseChat.ReplyMarkup = Keyboard(resp.Options)
		bot.Send(msg) // Отправка статуса с клавиатурой.

		if remainder == "" {
			// Если обработан весь запрос, цикл завершается.
			break
		} else {
			// Если нет, обрабатывается оставшаяся часть запроса.
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
