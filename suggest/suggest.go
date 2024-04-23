package suggest

var suggestions []Suggestion = []Suggestion{}

type Suggestion struct {
	Module    string
	From      string
	Suggested string
}

func Register(suggest Suggestion) {
	suggestions = append(suggestions, suggest)
}

func Get() []Suggestion {
	return suggestions
}
