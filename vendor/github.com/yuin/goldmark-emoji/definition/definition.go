package definition

// Emoji is a data structure that holds a single emoji.
type Emoji struct {
	// Name is a name of this emoji.
	Name string

	// ShortNames is a shorter representation of this emoji.
	ShortNames []string

	// Unicode is an unicode representation of this emoji.
	Unicode []rune
}

// NewEmoji returns a new Emoji.
func NewEmoji(name string, unicode []rune, shortNames ...string) Emoji {
	if len(shortNames) == 0 {
		panic("Emoji must have at leat 1 short name.")
	}
	if unicode == nil || len(unicode) == 0 {
		unicode = []rune{0xFFFD}
	}
	return Emoji{
		Name:       name,
		ShortNames: shortNames,
		Unicode:    unicode,
	}
}

// IsUnicode returns true if this emoji is defined in unicode, otherwise false.
func (em *Emoji) IsUnicode() bool {
	return !(len(em.Unicode) == 1 && em.Unicode[0] == 0xFFFD)
}

// Emojis is a collection of emojis.
type Emojis interface {
	// Get returns (*Emoji, true) if found mapping associated with given short name, otherwise (nil, false).
	Get(shortName string) (*Emoji, bool)

	// Add adds new emojis to this collection.
	Add(Emojis)

	// Clone clones this collection.
	Clone() Emojis
}

type emojis struct {
	list     []Emoji
	m        map[string]*Emoji
	children []Emojis
}

// NewEmojis returns a new Emojis.
func NewEmojis(es ...Emoji) Emojis {
	m := &emojis{
		list:     es,
		m:        map[string]*Emoji{},
		children: []Emojis{},
	}
	for i, _ := range es {
		emoji := &m.list[i]
		for _, s := range emoji.ShortNames {
			m.m[s] = emoji
		}
	}
	return m
}

func (m *emojis) Add(emojis Emojis) {
	m.children = append(m.children, emojis)
}

func (m *emojis) Clone() Emojis {
	es := &emojis{
		list:     m.list,
		m:        m.m,
		children: make([]Emojis, len(m.children)),
	}
	copy(es.children, m.children)
	return es
}

func (m *emojis) Get(shortName string) (*Emoji, bool) {
	v, ok := m.m[shortName]
	if ok {
		return v, ok
	}

	for _, es := range m.children {
		v, ok := es.Get(shortName)
		if ok {
			return v, ok
		}
	}
	return nil, false
}

// EmojisOption sets options for Emojis.
type EmojisOption func(Emojis)

// WithEmojis is an EmojisOption that adds emojis to the Emojis.
func WithEmojis(emojis ...Emoji) EmojisOption {
	return func(m Emojis) {
		m.Add(NewEmojis(emojis...))
	}
}
