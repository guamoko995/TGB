package engine

import (
	"TelegramGameBot/Game/base"
	"fmt"
	"strings"
	"sync"
)

// Игровой мир.
var Worlds = map[int64]*World{}

type Player interface {
	Handler
	base.Positioner
}

type World struct {
	Mu            sync.Mutex
	ActiveHandler Handler
	Pl            Player
	Nr            *Narrator
	*base.StNamer
	*base.StConteiner
}

// Конструктор
func (w *World) New() *World {
	w = &World{
		StNamer:     (*base.StNamer).New(&base.StNamer{}),
		StConteiner: (*base.StConteiner).New(&base.StConteiner{}),
		//Pl:          (*Player).New(&Player{}),
		Nr: &Narrator{},
	}
	w.Nr.W = w
	return w
}

// Добавляет новую локации.
func (w *World) AddLocation() *Location {
	loc := (*Location).New(&Location{})
	base.Place(loc, w)
	return loc
}

// Возвращает указатель на локацию по имени.
func (w *World) Locations(name string) *Location {
	obj, _ := ConsumePositionerFoundByName(name, base.FindPosition{Where: w, Deep: 0, IncludWhere: false})
	if loc, ok := obj.(*Location); ok {
		return loc
	}
	return nil
}

// Задает нового активного исполнителя.
func (w *World) NewActiveHandler(c Handler) Response {
	w.ActiveHandler = c
	return Response{
		Msg:     "",
		Status:  c.Status(),
		Options: c.Options(),
	}
}

// Ищет предмет по заданным параметрам по началу строки, вторым
// аргументом возвращает конец строки (без считанного предмета).
func (W *World) ConsumePositionerFoundByName(name string, fp ...base.FindPosition) (base.Positioner, string, error) {
	if name == "" {
		return nil, "", ErrEmptyStr
	}
	var obj base.Positioner
	var endStr string
	whereList := make([]string, 0)
	for _, p := range fp {
		obj, endStr = ConsumePositionerFoundByName(name, p)
		if obj != nil {
			return obj, endStr, nil
		}
		whereList = append(whereList, p.Where.Name("где"))
	}
	obj, _ = ConsumePositionerFoundByName(name, base.FindPosition{Where: W, Deep: -1, IncludWhere: false})
	if obj != nil {
		return nil, "", fmt.Errorf(List(whereList...) + " нет " + obj.Name("Р"))
	}
	words := strings.Split(name, " ")
	var s string
	for i, w := range words {
		s += w
		if !FindObjByPartName(s, base.FindPosition{Where: W, Deep: -1, IncludWhere: false}) {
			var postfix string
			if i == 0 {
				postfix = "a"
			}
			return nil, "", fmt.Errorf("нет предмета \"" + s + "\" или начинающегося со слов" + postfix + " \"" + s + "\"")
		}
		s += " "
	}
	return nil, "", fmt.Errorf("вероятно вы что-то недоговариваете")
}

// Строковым представлением мира является список строковых
// представлений всех локаций.
func (w *World) String() string {
	locList := make([]string, 0)
	for _, loc := range w.Content() {
		locList = append(locList, loc.String())
	}
	return strings.Join(locList, "\n")
}

// Игровая локация.
type Location struct {
	*base.StNamer
	*base.StPositioner
	*base.StConteiner              // Содержит любые объекты.
	Bridge            *EnConteiner // Проходы на другие локации.
}

type EnConteiner struct {
	*base.StNamer
	*base.StConteiner
}

// Сущность-контейнер.
func (*EnConteiner) New() *EnConteiner {
	c := EnConteiner{
		StNamer:     (*base.StNamer).New(&base.StNamer{}),
		StConteiner: (*base.StConteiner).New(&base.StConteiner{}),
	}
	return &c
}

// Конструктор.
func (*Location) New() *Location {
	loc := &Location{
		StNamer:      (*base.StNamer).New(&base.StNamer{}),
		StPositioner: &base.StPositioner{},
		StConteiner:  (*base.StConteiner).New(&base.StConteiner{}),
		Bridge:       (*EnConteiner).New(&EnConteiner{}),
	}
	//loc.Bridge.Reposition(loc)
	return loc
}

// Строковым представлением локации по умолчанию является ее имя +
// + содержимое + проходы.
func (loc *Location) String() string {
	W := RootConteiner(loc)
	s := loc.Name()
	s += ". здесь "
	ms := make([]string, 0)
	for _, obj := range loc.Content() {
		if W.Pl != obj {
			ms = append(ms, obj.Name("В"))
		}
	}
	if len(ms) == 0 {
		s += ", кажется, пусто"
	} else {
		s += "вы видите " + List(ms...)
	}
	s += ". отсюда "
	ms = make([]string, 0)
	for _, bridg := range loc.Bridge.Content() {
		//fmt.Println(bridg.Name())
		ms = append(ms, bridg.Name("куда"))
	}
	if len(ms) == 0 {
		s += ",кажется, нет выхода"
	} else {
		s += "можно идти " + List(ms...)
	}
	return s
}
