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
	users map[string]uuid.UUID `json:"users"`
}

// NewStudentsRegister - create a new instance
func NewStudentsRegister() StudentsRegister {
	return StudentsRegister{users: make(map[string]uuid.UUID)}
}

// Register a new parser
func (register *StudentsRegister) Register(uco string, uuid uuid.UUID) {
	register.users[uco] = uuid
}

// GetOrRegister new UUID for the provided uco
func (register *StudentsRegister) GetOrRegister(uco string) uuid.UUID {
	value, ok := register.users[uco]
	if !ok {
		value = uuid.New()
		register.users[uco] = value
	}
	return value
}

// Get a parser instance
func (register *StudentsRegister) Get(uco string) (uuid.UUID, error) {
	value, ok := register.users[uco]
	if !ok {
		return uuid.UUID{}, fmt.Errorf("User with uco not found: %s", uco)
	}
	return value, nil
}

// Export students register to a provided file
func (register *StudentsRegister) Export(file string) {
	content, err := json.MarshalIndent(register.users, "", "  ")
	if err != nil {
		log.WithError(err).Error("Unable to marshall file")
	}

	if err = ioutil.WriteFile(file, content, 0644); err != nil {
		log.WithError(err).WithField("filepath", file).Error("unable to save marshall file")
	}
}

// Import the file content to the refgister
func (register *StudentsRegister) Import(file string) error {

	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.WithField("filepath", file).WithError(err).Error("unable to read a file")
	}

	if err = json.Unmarshal(content, register.users); err != nil {
		log.WithField("filepath", file).WithError(err).Error("unable to unmarshal a file")

	}

}

