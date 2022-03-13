package base

// Интерфейс для манипуляций с именами объектов.
type Namer interface {
	// Возвращает имя, может принимать аргументы.
	// Наличие аргументов предусмотрено для возможности получить
	// имя в нужной форме (падеж, предлог и т.п.).
	Name(...string) string
	// Возвращает строку дополнительной информации по ключу.
	Info(string) string
	// Добавляет форму имени.
	AddName(...Title)
	// Добавляет/обновляет дополнительную информацию по ключу.
	AddInfo(...Title)
	// Возвращает массив имен.
	Names() []string
	// Возвращает строковое представление именованого объекта.
	String() string
}

// Названнние имеет форму (например падеж) и значение.
type Title struct {
	Form  string
	Value string
}

// Родитель именованых объектов.
type StNamer struct {
	names map[string]string
	info  map[string]string
}

// Конструктор.
func (*StNamer) New() *StNamer {
	n := &StNamer{
		make(map[string]string),
		make(map[string]string),
	}
	// Карта всегда должна содержать запись с ключем "". Под ключем ""
	// имя по умолчанию.
	n.names[""] = ""
	return n
}

// Возвращает имя n в первой указанной форме form. Если такая форма,
// отсутствует, то во второй указанной и т.д. Если все указанные формы
// отсутствуют или не указано ни одной формы, то возвращается имя в
// форме по умолчанию form="".
func (n *StNamer) Name(form ...string) string {
	for _, f := range form {
		if name, ok := n.names[f]; ok {
			return name
		}
	}
	return n.names[""]
}

// По умолчанию, Строковое представление именованого объекта - его имя в
// форме по умолчанию.
func (n *StNamer) String() string {
	return "Просто " + n.names[""]
}

// Возвращает информацию (такую как род) по ключу.
func (n *StNamer) Info(key string) string {
	if name, ok := n.info[key]; ok {
		return name
	}
	return ""
}

// Добавляет формы имен.
func (n *StNamer) AddName(names ...Title) {
	for _, name := range names {
		n.names[name.Form] = name.Value
	}
}

// Добавляет информацию (такую как род)
func (n *StNamer) AddInfo(infs ...Title) {
	for _, inf := range infs {
		n.info[inf.Form] = inf.Value
	}
}

// Возвращает массив имен (все формы)
func (n *StNamer) Names() []string {
	mNames := make([]string, 0)
	for _, name := range n.names {
		mNames = append(mNames, name)
	}
	return mNames
}

// Интерфейс для позиционирования объектов.
type Positioner interface {
	Namer
	Position() Conteiner
	Reposition(newPosition Conteiner)
}

// Родитель позиционируемых объектов.
type StPositioner struct {
	position Conteiner
}

// Возвращает позицию объекта.
func (obj *StPositioner) Position() Conteiner {
	return obj.position
}

// Изменяет позицию объекта.
func (obj *StPositioner) Reposition(newPos Conteiner) {
	obj.position = newPos
}

// Расширение интерфейса Positioner с возможностью отказаться от
// перемещения
type Relocater interface {
	//Positioner
	Relocate(newPosition Conteiner) error
}

// Расширение интерфейса Positioner для размерных объектов.
type Sizer interface {
	//Positioner
	Size() int
}

// Родитель размерных объектов.
type StSizer struct {
	size int
}

// Возвращает размер объекта.
func (s *StSizer) Size() int {
	return s.size
}

// Задает размер объекта.
func (s *StSizer) Resize(value int) {
	s.size = value
}

// Интерфейс для манипулирования содержимым объектов.
type Conteiner interface {
	Namer
	Add(obj ...Positioner)
	Remove(obj ...Positioner)
	Content() []Positioner
}

// Родитель объектов, способных быть позицией другим объектам.
type StConteiner struct {
	content []Positioner
}

// Конструктор.
func (*StConteiner) New() *StConteiner {
	c := StConteiner{make([]Positioner, 0)}
	return &c
}

// Добавляет объекты в содержимое.
func (cont *StConteiner) Add(obj ...Positioner) {
	cont.content = append(cont.content, obj...)
}

// Удаляет объекты из содержимого.
func (cont *StConteiner) Remove(obj ...Positioner) {
	ub := len(cont.content) - 1
	for _, del := range obj {
		for i, obj := range cont.content {
			if del == obj {
				cont.content[i] = cont.content[ub]
				cont.content[ub] = nil
				cont.content = cont.content[:ub]
				ub--
				break
			}
		}
	}
}

// Возвращает все содержимое.
func (cont *StConteiner) Content() []Positioner {
	return cont.content
}

// Расширение интерфейса Conteiner с возможностью отказаться от
// взятия на хранение.
type Taker interface {
	//Conteiner
	Take(obj Positioner) error
}

// Расширение интерфейса Conteiner с возможностью отказаться от
// выдачи.
type Giver interface {
	//Conteiner
	Give(obj Positioner, newPosition Conteiner) error
}

// Расширение интерфейса Conteiner для ограничения вместимости.
type Limiter interface {
	//Conteiner
	Capacity() int
	Vacancy() int
}

// РасширениеSt Conteiner с ограниченной вместимостью.
type StLimitedConteiner struct {
	*StConteiner
	capacity int
}

// Конструктор.
func (*StLimitedConteiner) New() *StLimitedConteiner {
	c := StLimitedConteiner{
		StConteiner: (*StConteiner).New(&StConteiner{}),
	}
	return &c
}

// Возвращает вместимость.
func (c *StLimitedConteiner) Capacity() int {
	return c.capacity
}

// Задает вместимость.
func (c *StLimitedConteiner) Recapacity(value int) {
	c.capacity = value
}

// Возвращает оставшуюся вместимость (свободную от содержимого)
func (c *StLimitedConteiner) Vacancy() int {
	var SizeContent int
	for _, obj := range c.Content() {
		if sizeObj, ok := obj.(Sizer); ok {
			SizeContent += sizeObj.Size()
		}
	}
	return c.capacity - SizeContent
}
