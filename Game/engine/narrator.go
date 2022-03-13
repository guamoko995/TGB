package engine

import (
	"fmt"
	//"TelegramGameBot/Game/world/texts"
)

type Narrator struct {
	Texts           []string
	NumberText      int
	NextImplementer Handler
	W               *World
}

func (n *Narrator) Status() string {
	return fmt.Sprintf("[ %v/%v ]", n.NumberText, len(n.Texts))
}

func (*Narrator) Options() [][]string {
	return [][]string{{"->"}}
}

func (n *Narrator) Handle(str string) (Response, string) {
	l := len(n.Texts) - 1
	if n.NumberText < l {
		t := n.Texts[n.NumberText]
		n.NumberText++
		return Response{
			Msg:     t,
			Status:  n.Status(),
			Options: n.Options(),
		}, ""
	} else {
		n.W.ActiveHandler = n.NextImplementer
		return Response{
			Msg:     n.Texts[n.NumberText],
			Status:  n.W.ActiveHandler.Status(),
			Options: n.W.ActiveHandler.Options(),
		}, ""
	}
}

/*/ Функция которая делает исполнителем пользовательских комманд некий
// помощник - повествователь.
func (n *Narrator) New() *Narrator {
	n = &Narrator{
		Texts: strings.Split(texts.GameText("введение"), "\n"),
	}
	return n
}*/
