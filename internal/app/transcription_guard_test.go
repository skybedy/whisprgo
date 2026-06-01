package app

import "testing"

func TestShouldRetryParakeetForCzech(t *testing.T) {
	cases := []struct {
		name string
		lang string
		text string
		want bool
	}{
		{name: "czech text", lang: "cs", text: "Ahoj, to je test.", want: false},
		{name: "english text", lang: "cs", text: "Uh I don't test.", want: true},
		{name: "russian text", lang: "cs", text: "Я новым тестую.", want: true},
		{name: "non-cs language", lang: "en", text: "Uh I don't test.", want: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := shouldRetryParakeetForCzech(tc.lang, tc.text)
			if got != tc.want {
				t.Fatalf("unexpected result got=%v want=%v", got, tc.want)
			}
		})
	}
}
