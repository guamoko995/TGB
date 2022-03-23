package items

import (
	"TelegramGameBot/Game/base"
	"TelegramGameBot/Game/engine"
	mediafiles "TelegramGameBot/Game/mediaFiles"
	"TelegramGameBot/Game/world/buildTools"
)

// Предмет: Конверт
type Manual struct {
	//*base.StNamer
	*base.StPositioner
	*base.StSizer
	*engine.TreeHandlers
}

func (b *Manual) New() *Manual {
	b = &Manual{
		//StNamer:      (*base.StNamer).New(&base.StNamer{}),
		StPositioner: &base.StPositioner{},
		StSizer:      &base.StSizer{},
		TreeHandlers: (*engine.TreeHandlers).New(&engine.TreeHandlers{}),
	}

	b.Stat = func() string {
		return "[использование конспекта]"
	}

	buildTools.SetName(b, "конспект по криптоанализу")
	buildTools.SetAdditName(b, "конспект", "сокр")
	b.Resize(5)

	b.Applications["ознакомиться"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		W := engine.RootConteiner(b)
		resp := W.NewActiveHandler(W.Pl)
		resp.Doc = mediafiles.Doc["конспект"]
		return resp, args
	})
	b.Applications["х"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		W := engine.RootConteiner(b)
		resp := W.NewActiveHandler(W.Pl)
		return resp, args
	})
	return b
}

func (b *Manual) String() string {
	return "здесь есть разбор взлома шифра простой замены. Может пригодиться."
}
