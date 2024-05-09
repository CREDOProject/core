package suggest

import "testing"

func Test_Suggest(t *testing.T) {
	suggestion := Suggestion{
		Module:    "Module",
		From:      "From",
		Suggested: "Suggested",
	}
	Register(suggestion)
	if !HasSuggestion() {
		t.Error("Suggestion not registered.")
	}

	suggestions := Get()
	t.Log(suggestions.String())
}
