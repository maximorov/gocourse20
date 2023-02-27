package model

type CountryCode string

func (cc CountryCode) String() string {
	return string(cc)
}
