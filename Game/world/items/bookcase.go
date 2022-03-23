package items

import (
	"TelegramGameBot/Game/base"
	"TelegramGameBot/Game/engine"
	"TelegramGameBot/Game/world/buildTools"
)

// Коробка - простой предмет, который может содержать другие предметы.
type Bookcase struct {
	*base.StNamer
	*base.StPositioner
	*base.StSizer
	*base.StLimitedConteiner
}

// По умолчанию строковым представлением коробки является перечисление
// содержимого.
func (c *Bookcase) String() string {
	return engine.StrConteiner(c)
}

// Конструктор.
func (b *Bookcase) New() *Bookcase {
	b = &Bookcase{
		StNamer:            (*base.StNamer).New(&base.StNamer{}),
		StLimitedConteiner: (*base.StLimitedConteiner).New(&base.StLimitedConteiner{}),
		StPositioner:       &base.StPositioner{},
		StSizer:            &base.StSizer{},
	}
	buildTools.SetName(b, "книжный шкаф")
	buildTools.SetAdditName(b, "шкаф", "сокр")
	b.Resize(2500)
	b.Recapacity(2000)
	return b
}
