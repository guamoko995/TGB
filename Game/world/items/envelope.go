package items

import (
	"TelegramGameBot/Game/base"
	"TelegramGameBot/Game/engine"
	"TelegramGameBot/Game/world/buildTools"
)

// Предмет: Конверт
type Envelope struct {
	*base.StPositioner
	*base.StSizer
	*base.StLimitedConteiner
	*engine.TreeHandlers
	about  func() string
	opened bool
	label  string
}

func (b *Envelope) Status() string {
	return "[действия с конвертом]"
}

func (b *Envelope) New() *Envelope {
	b = &Envelope{
		StPositioner:       &base.StPositioner{},
		StSizer:            &base.StSizer{},
		StLimitedConteiner: (*base.StLimitedConteiner).New(&base.StLimitedConteiner{}),
		TreeHandlers:       (*engine.TreeHandlers).New(&engine.TreeHandlers{}),
	}
	buildTools.SetName(b, "конверт")
	b.Resize(5)
	b.Recapacity(4)
	b.label = "Прямиком из памяти"
	b.about = func() string {
		return "Конверт запечатан. Подпись: \"" + b.label + "\""
	}

	apdate := func() {
		delete(b.Applications, "вскрыть")
		delete(b.Applications, "х")
		b.about = func() string {
			list := make([]string, 0)
			for _, obj := range b.Content() {
				list = append(list, obj.Name())
			}
			return "Конверт вскрыт. Подпись: \"" + b.label + "\"" + ". В конверте " + engine.List(list...)
		}
		b.opened = true
	}

	b.Applications["вскрыть"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		W := engine.RootConteiner(b)
		resp := W.NewActiveHandler(W.Pl)
		apdate()
		resp.Msg = "конверт вскрыт"
		return resp, args
	})
	b.Applications["х"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		W := engine.RootConteiner(b)
		return W.NewActiveHandler(W.Pl), args
	})
	return b
}

func (b *Envelope) String() string {
	return b.about()
}

func (b *Envelope) Content() []base.Positioner {
	if b.opened {
		return b.StLimitedConteiner.Content()
	}
	return []base.Positioner{}
}
