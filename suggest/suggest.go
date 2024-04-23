package suggest

import "fmt"

var suggestions []Suggestion = []Suggestion{}

type Suggestion struct {
	Module    string
	From      string
	Suggested string
}

func Register(suggest Suggestion) {
	suggestions = append(suggestions, suggest)
}

func Get() Suggestions {
	return suggestions
}

func (s Suggestion) String() string {
	return fmt.Sprintf("%s suggested from %s coming from %s.",
		s.Suggested,
		s.From,
		s.Module)
}

type Suggestions []Suggestion

func (suggestions Suggestions) String() (output string) {
	for _, s := range suggestions {
		output += fmt.Sprintf("\t- %s\n", s.String())
	}
	return
}

func HasSuggestion() bool { return len(suggestions) > 0 }
