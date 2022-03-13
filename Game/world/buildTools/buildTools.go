package buildTools

import (
	"TelegramGameBot/Game/base"
	"TelegramGameBot/Game/engine"
	"TelegramGameBot/Game/world/texts"
	"fmt"
)

// Находит формы имени и род в NamesMap по name и присваивает их
// объекту obj.
func SetName(obj base.Namer, name string) {
	m, ok := texts.NamesMap()[name]
	if !ok {
		fmt.Printf("в NamesMap нет \"%s\"", name)
		return
	}
	m[""] = name
	val, ok := m["род"]
	if !ok {
		return
	}
	obj.AddInfo(base.Title{Form: "род", Value: val})
	delete(m, "род")
	mas := make([]base.Title, 0)
	for key, val := range m {
		mas = append(mas, base.Title{Form: key, Value: val})
	}
	obj.AddName(mas...)
}

// Находит формы имени в NamesMap и присваивает их объекту obj в качестве
// дополнительных имен для поиска с помощью краткого названия.
func SetAdditName(obj base.Positioner, name string, prefix string) {
	m, ok := texts.NamesMap()[name]
	if !ok {
		fmt.Printf("в NamesMap нет \"%s\"", name)
		return
	}
	m[""] = name
	delete(m, "род")
	mas := make([][]string, 0)
	for key, val := range m {
		mas = append(mas, []string{key, val})
	}
	engine.NamesAddWithPrefix(obj, mas, prefix)
}
