package items

import (
	"TelegramGameBot/Game/base"
	"TelegramGameBot/Game/world/buildTools"
	//"TelegramGameBot/Game/engine"
)

type Armchair struct {
	*base.StNamer
	*base.StPositioner
	*base.StSizer
}

func (b *Armchair) New() *Armchair {
	b = &Armchair{
		StNamer:      (*base.StNamer).New(&base.StNamer{}),
		StPositioner: &base.StPositioner{},
		StSizer:      &base.StSizer{},
	}
	buildTools.SetName(b, "кресло")
	b.Resize(1001)
	return b
}
