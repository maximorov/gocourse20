package model

var (
	SupportedLanguages = []LanguageCode{"at", "bf", "bg", "bn", "cf", "ch", "cz", "de",
		"dk", "ee", "en", "es", "fi", "fr", "gb", "gr", "hu", "it", "ld", "lf",
		"lr", "lt", "lv", "nl", "no", "pl", "pt", "ro", "ru", "se", "si", "sk",
	}
)

type LanguageCode string

func (lc LanguageCode) String() string {
	return string(lc)
}

func (lc LanguageCode) IsValid() bool {
	for _, l := range SupportedLanguages {
		if l == lc {
			return true
		}
	}

	return false
}
