package parsers

import (
	"github.com/pestanko/isstat/core"
)

// KontrFunctionalityParser - parses functionality points
type KontrFunctionalityParser struct {
}


/*
Parse the notepad points

Format:
# zapsáno z Kontru 2020-02-18 08:45, v2.2.1

%%       datum    cas  body
 1  2020-02-18  08:45    *1

# POZOR: Tento blok NEUPRAVUJTE!

# Kontr může veškeré změny kdykoliv přepsat.
# Poznámky k odevzdání a hodnocení čistoty pište
# do bloku určeného pro tyto účely.
*/
func (parser *KontrFunctionalityParser) Parse(content *core.NotepadContent) ([]core.Submission, error) {

}
