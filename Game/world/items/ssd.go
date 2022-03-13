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
	b.Inf = "На карте памяти имеется текстовый файл. Это последовательность букв которую Вы уже выучили наизусть. С ней проще работать когда она в файле"
	return b
}
