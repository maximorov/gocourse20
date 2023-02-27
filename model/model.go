package model

type Language struct {
	Code        LanguageCode `json:"code" table_name:"languages" schema:"pk" pk:"1"`
	CountryCode CountryCode  `json:"country_code" validation:"required,country-exists"`
	Name        string       `json:"name"`
}

func (l *Language) PrimaryKey() PrimaryKey {
	return NewPrimaryKey(l.Code.String())
}
