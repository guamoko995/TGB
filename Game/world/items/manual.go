package items

import (
	"TelegramGameBot/Game/base"
	"TelegramGameBot/Game/world/buildTools"
)

// Предмет: Конверт
type Manual struct {
	*base.StNamer
	*base.StPositioner
	*base.StSizer
}

func (b *Manual) New() *Manual {
	b = &Manual{
		StNamer:      (*base.StNamer).New(&base.StNamer{}),
		StPositioner: &base.StPositioner{},
		StSizer:      &base.StSizer{},
	}
	buildTools.SetName(b, "инструкция")
	b.Resize(5)
	return b
}

func (b *Manual) String() string {
	return "инструкция"
}
