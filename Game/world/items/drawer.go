package items

import (
	"TelegramGameBot/Game/base"
	"TelegramGameBot/Game/engine"
	"TelegramGameBot/Game/world/buildTools"
)

// Не перемещаемая коробка.
type Drawer struct {
	*base.StNamer
	*base.StPositioner
	*base.StSizer
	*base.StLimitedConteiner
	NoRePlace error
}

// Конструктор.
func (b *Drawer) New() *Drawer {
	b = &Drawer{
		StNamer:            (*base.StNamer).New(&base.StNamer{}),
		StLimitedConteiner: (*base.StLimitedConteiner).New(&base.StLimitedConteiner{}),
		StPositioner:       &base.StPositioner{},
		StSizer:            &base.StSizer{},
	}
	buildTools.SetName(b, "выдвижной ящик")
	buildTools.SetAdditName(b, "ящик", "сокр")
	b.Resize(20)
	b.Recapacity(20)
	return b
}

func (b *Drawer) String() string {
	return engine.StrConteiner(b)
}

// При попытки перемещения возвращается noRebase.Place. Ошибка noRePlase
// может принимать значение nil, тогда nrpb все же можно переместить.
func (nrpb *Drawer) Relocate(newPos base.Conteiner) error {
	return nrpb.NoRePlace
}
