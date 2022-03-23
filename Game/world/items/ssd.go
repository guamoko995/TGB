package items

import (
	"TelegramGameBot/Game/base"
	"TelegramGameBot/Game/engine"
	"TelegramGameBot/Game/world/buildTools"
)

// Предмет: карта памяти
type Ssd struct {
	*base.StSizer
	*base.StPositioner
	*engine.TreeHandlers
	Inf string
}

func (b *Ssd) New() *Ssd {
	b = &Ssd{
		StPositioner: &base.StPositioner{},
		StSizer:      &base.StSizer{},
		TreeHandlers: (*engine.TreeHandlers).New(&engine.TreeHandlers{}),
	}
	buildTools.SetName(b, "карта памяти")
	b.Resize(1)
	b.Inf = "К ноутбуку подключена карта памяти, на которой записан текстовый файл. " +
		"Текст состоит из последовательности букв, которую Вы уже выучили наизусть. " +
		"С шифром проще работать когда он сохранен в файле"
	return b
}
