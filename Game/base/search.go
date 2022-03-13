package base

// Переменная, служащая сигналом для передачи через Chan.
// Не занимает память. Полезной информацией является факт  передачи.
var S0 struct{}

// Рекурсивный параллельный поиск объекта по критерию, определяемому
// функцией criterion в поисковой позиции fp до первого совпадения.
func Find(criterion func(Positioner) bool, fp ...FindPosition) Positioner {
	var Obj Positioner
	args := len(fp)
	inChanStop := make(chan struct{}, args)
	outChanEnd := make(chan struct{}, args)
	outChanObj := make(chan Positioner, 1)

	for _, p := range fp {
		if p.IncludWhere {
			if obj, ok := p.Where.(Positioner); ok {
				if criterion(obj) {
					return obj
				}
			}
		}
		go search(criterion, p.Where, p.Deep, outChanObj, outChanEnd,
			inChanStop)
	}
	for range fp {
		select {
		case Obj = <-outChanObj:
			inChanStop <- S0
			break
		case <-outChanEnd:
		}
	}
	select {
	case Obj = <-outChanObj:
	default:
	}
	return Obj
}

// Поисковая позиция.
type FindPosition struct {
	Where       Conteiner // Где искать.
	Deep        int       // Глубина поиска.
	IncludWhere bool      // Место поиска может являться результатом.
}

// Рекурсивный параллельный поиск предмета по критерию criterionс, в
// контейнере where и вложенных в него на глубину deep. Как только
// проверено все содержимое, функция сообщает о завершении себя и всех
// дочерних копий сигналом по каналу outChanEnd. Если объект найден, он
// помещается в канал outChanObj. Чтобы завершить поиск досрочно,
// например в случае когда объект уже найден, необходимо отправить
// сигнал в канал inChanStop. Функцию удобно запускать в go-рутине.
func search(criterion func(Positioner) bool, where Conteiner, deep int,
	outChanObj chan Positioner, outChanEnd chan struct{},
	inChanStop chan struct{}) {

	var ch int // Счетчик запущенных дочерних рекурсивных вызовов.
	l := len(where.Content())
	childEnd := make(chan struct{}, l)
	childStop := make(chan struct{}, l)
	End := func() {
		// Ожидает завершения всех дочерних вызовов.
		for ; ch > 0; ch-- {
			<-childEnd
		}
		// Собщает о завершении вызвавшему родителю.
		outChanEnd <- S0
	}
	// Отправляет сигнал завершения всем дочерним вызовам.
	Stop := func() {
		for i := 0; i < ch; i++ {
			childStop <- S0
		}
		End()
	}
	for _, obj := range where.Content() {
		select {
		case <-inChanStop:
			Stop()
			return
		default:
			if criterion(obj) {
				outChanObj <- obj
				Stop()
				return
			}
			if deep != 0 {
				if k, ok := obj.(Conteiner); ok {
					go search(criterion, k, deep-1, outChanObj, childEnd,
						childStop)
					ch++
				}
			}
		}
	}
	End()
}

// Возвращает корневой контейнер.
func RootConteiner(d Positioner) Conteiner {
	if d.Position() == nil {
		if root, ok := d.(Conteiner); ok {
			return root
		}
		panic("StPositioner нигде")
	}
	pos := d.Position()
	if obj, ok := pos.(Positioner); ok {
		return RootConteiner(obj)
	}
	return pos
}
