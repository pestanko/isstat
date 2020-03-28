package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// StudentsRegister - container for all of the registered parsers
type StudentsRegister struct {
	Users map[string]uuid.UUID `json:"Users"`
}

// NewStudentsRegister - create a new instance
func NewStudentsRegister() StudentsRegister {
	return StudentsRegister{Users: make(map[string]uuid.UUID)}
}

// Register a new parser
func (register *StudentsRegister) Register(uco string, uuid uuid.UUID) {
	register.Users[uco] = uuid
}

// GetOrRegister new UUID for the provided uco
func (register *StudentsRegister) GetOrRegister(uco string) uuid.UUID {
	value, ok := register.Users[uco]
	if !ok {
		value = uuid.New()
		register.Users[uco] = value
	}
	return value
}

// Get a parser instance
func (register *StudentsRegister) Get(uco string) (uuid.UUID, error) {
	value, ok := register.Users[uco]
	if !ok {
		return uuid.UUID{}, fmt.Errorf("User with uco not found: %s", uco)
	}
	return value, nil
}

// Export students register to a provided file
func (register *StudentsRegister) Export(file string) {
	content, err := json.MarshalIndent(register.Users, "", "  ")
	if err != nil {
		log.WithError(err).Error("Unable to marshall file")
	}

	if err = ioutil.WriteFile(file, content, 0644); err != nil {
		log.WithError(err).WithField("filepath", file).Error("unable to save marshall file")
	}
}

// Import the file content to the register
func (register *StudentsRegister) Import(file string) error {

	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.WithField("filepath", file).WithError(err).Error("unable to read a file")
	}

	if err = json.Unmarshal(content, register.Users); err != nil {
		log.WithField("filepath", file).WithError(err).Error("unable to unmarshal a file")
	}

	return nil
}
