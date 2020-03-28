package parsers

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

var parserRegister *Register

// Register - container for all of the registered parsers
type Register struct {
	Parsers map[string]NotepadContentParser
}

// NewRegister - create a new instance
func NewRegister() *Register {
	return &Register{Parsers: make(map[string]NotepadContentParser)}
}

// Register a new parser
func (register *Register) Register(name string, parser NotepadContentParser) {
	register.Parsers[name] = parser
}

// Get a parser instance
func (register *Register) Get(name string) (NotepadContentParser, error) {
	value, ok := register.Parsers[name]
	if !ok {
		return nil, fmt.Errorf("parser not found '%s'", name)
	}
	return value, nil
}


func (register *Register) GetOrDefault(name string) NotepadContentParser {
	parser, err := register.Get(name)

	if err != nil {
		log.WithField("name", name).WithError(err).Debug("No parser found")
		parser, err = register.Get("default")
		if err != nil {
			log.Error("No default parser registered - something is definitely wrong!")
		}
	}

	return parser
}

func GetParserRegister() *Register {
	if parserRegister == nil {
		parserRegister = NewRegister()
	}
	return parserRegister
}