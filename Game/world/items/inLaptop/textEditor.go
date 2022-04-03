package inLaptop

import (
	"TelegramGameBot/Game/engine"
	"TelegramGameBot/Game/world/buildTools"
	worldGame "TelegramGameBot/Game/world/items/inLaptop/wordGame"
	"fmt"
	"strings"
)

// Псевдопредмет: текстовый редактор
type TextEditor struct {
	*engine.TreeHandlers
	Text        worldGame.QwestText
	NextHandler engine.Handler
	W           *engine.World
}

// Конструктор.
func (te *TextEditor) New() *TextEditor {
	te = &TextEditor{
		TreeHandlers: (*engine.TreeHandlers).New(&engine.TreeHandlers{}),
	}
	buildTools.SetName(te, "текстовый редактор")

	// Установка статусной строки при использовании редактора.
	te.Stat = func() string {
		return "[редактор текста]"
	}

	// Замена функции обработки входной команды с целью учитывания
	// регистра символов в коммандах заменить и посчитать.
	te.InputFormat = func(s string) string {
		// Применить стандартную функцию форматирования к первому слову.
		mStr := strings.SplitN(s, " ", 2)
		mStr[0] = engine.InputFormat(mStr[0])

		// В случае если первое слово соответствует команде "заменить"
		// или "посчитать", следующие слова не форматируются
		if mStr[0] == "заменить" || mStr[0] == "посчитать" {
			return strings.Join(mStr, " ")
		}

		// В противном случае ко всей строке применяется стандартная
		// функция форматирования.
		return engine.InputFormat(s)
	}

	// Выходное форматирование не используется, требуется следить за
	// выходными текстами!
	te.OutputFormat = func(s string) string { return s }

	te.Applications["показать"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		return te.StResp(te.Text.Print()), args
	})

	te.Applications["заменить"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		if args == "" {
			resp := te.W.NewActiveHandler(&engine.Complementer{
				W:               te.W,
				LastCommand:     "заменить",
				NextImplementer: te,
			})
			return resp, ""
		}
		ar := strings.SplitN(args, " ", 3)
		if len(ar) < 2 {
			return te.StResp("Команда \"заменить\" принимает два символа, разделенные пробелом: первый - заменяемый, второй - заменяющий"), ""
		}
		endStr := ""
		if len([]rune(ar[0])) != 1 || len([]rune(ar[1])) != 1 {
			return te.StResp("Команда \"заменить\" принимает два символа, разделенные пробелом: первый - заменяемый, второй - заменяющий"), ""
		}
		if len(ar) > 2 {
			endStr = ar[2]
		}
		te.Text.Replace([]rune(ar[0])[0], []rune(ar[1])[0])
		return te.StResp("Символ '" + ar[0] + "' заменен на символ '" + ar[1] + "'"), endStr
	})

	te.Applications["нижний регистр"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		te.Text.Down()
		return te.StResp("текст переведен в нижжний регистр"), args
	})

	te.Applications["верхний регистр"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		te.Text.Up()
		return te.StResp("текст переведен в верхний регистр"), args
	})

	te.Applications["посчитать"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		if args == "" {
			resp := te.W.NewActiveHandler(&engine.Complementer{
				W:               te.W,
				LastCommand:     "посчитать",
				NextImplementer: te,
			})
			return resp, ""
		}
		ar := strings.SplitN(args, " ", 2)
		endStr := ""
		if len([]rune(ar[0])) != 1 {
			return te.StResp("Команда \"посчитать\" принимает один символ"), ""
		}
		if len(ar) > 2 {
			endStr = ar[2]
		}
		s := []rune(ar[0])[0]
		n := te.Text.Count(s)
		all := te.Text.CountAll()
		p := float32(100*n) / float32(all)
		return te.StResp(fmt.Sprintf("Символ '%c' встречается в тексте %v раз, что составляет %.2f%% от общего количества символов в тексте", s, n, p)), endStr
	})

	te.Applications["заново"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		te.Text.Reset()
		return te.StResp("Все измененияя отменены"), args
	})

	te.Applications["х"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		return te.W.NewActiveHandler(te.NextHandler), args
	})

	return te
}

// Формирует стандартный ответ из сообщения.
func (te *TextEditor) StResp(Msg string) engine.Response {
	return engine.Response{
		Msg:     Msg,
		Status:  te.Status(),
		Options: te.Options(),
	}
}
