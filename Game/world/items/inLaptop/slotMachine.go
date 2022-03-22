package inLaptop

import (
	"TelegramGameBot/Game/engine"
	worldGame "TelegramGameBot/Game/world/items/inLaptop/wordGame"
)

// Машина слотов.
type SlotMakhine struct {
	*worldGame.SlotMakhine
	NextHandler engine.Handler
	W           *engine.World
}

// Возвращает варианты взаимодействия с машиной слотов.
func (sm *SlotMakhine) Options() [][]string {
	options := [][]string{
		{},
		{"<-", "зафиксировать слово", "->"},
		{"заново", "x"},
	}
	for _, mr := range sm.Str {
		R := mr.Name()
		sR := string([]rune{R})
		options[0] = append(options[0], sR)
	}
	return options
}

func (sm *SlotMakhine) Status() string {
	return "[дешифровка]"
}

func (sm *SlotMakhine) Handle(str string) (engine.Response, string) {
	l := len(sm.Words)
	switch str {
	case "<-":
		sm.Pos--
		if sm.Pos < 0 {
			sm.Pos += l
		}
		sm.Update()
	case "->":
		sm.Pos++
		if sm.Pos == l {
			sm.Pos = 0
		}
		sm.Update()
	case "заново":
		sm.Text.Reset()
		sm.Text.Down()
		sm.Update()
	case "зафиксировать слово":
		for key, R := range sm.ReplaceMap() {
			sm.Text.Replace(key, R)
		}
		sm.Update()
	case "x":
		return sm.W.NewActiveHandler(sm.NextHandler), ""
	default:
		for _, R := range str {
			sm.SmartClick(R)
		}
	}
	msg := make([]rune, 0)
	for _, pos := range sm.Str {
		msg = append(msg, pos.Name())
	}
	return engine.Response{
		Msg:     string(msg),
		Status:  "[клавиатура]",
		Options: sm.Options(),
	}, ""
}
