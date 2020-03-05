package app

import (
	"github.com/pestanko/isstat/core"
	"github.com/pestanko/isstat/parsers"
)

// IsStatApp - Is MUNI Statistics application
type IsStatApp struct {
	Client *core.CourseClient
	Parser parsers.Parser
}

// NewBasicApp - create a new basic application
func NewBasicApp(client *core.CourseClient, parser parsers.Parser) IsStatApp {
	return IsStatApp {Client: client, Parser: parser}
}

type IsStatsAppFactory struct {
	Config 
}