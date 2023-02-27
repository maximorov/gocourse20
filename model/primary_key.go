package model

import (
	"encoding/json"
	"fmt"
)

type PKType string

const (
	PkTypeID      PKType = "ID"
	PkTypeComplex PKType = "Complex"
	KeyID         string = "id"
)

type Scanable interface {
	Scan(dest ...interface{}) error
}

type PrimaryKeyable interface {
	PrimaryKey() PrimaryKey
}

type PrimaryKey interface {
	IsComplex() bool
	Type() PKType
	GetPartOfComplexInt(k string) (int, error)
	GetPartOfComplexString(k string) (string, error)
	IsID() bool
	String() string
	ToCond() map[string]any
	Fields() []string
	GetInt() (int, error)
	Scan(Scanable) error
}

type PrimaryKeySrc interface {
	int | int64 | string | map[string]any
}

func NewPrimaryKey[S PrimaryKeySrc](key S) PrimaryKey {
	switch val := any(key).(type) {
	case int, int64, string:
		return &primaryKey{
			tp:      PkTypeID,
			key:     map[string]any{KeyID: val},
			_fields: []string{KeyID},
			_values: []any{&val},
		}
	case map[string]any:
		pk := &primaryKey{
			tp:      PkTypeComplex,
			key:     val,
			_fields: make([]string, 0, len(val)),
			_values: make([]any, 0, len(val)),
		}

		for field, value := range val {
			pk._fields = append(pk._fields, field)
			pk._values = append(pk._values, &value)
		}

		return pk
	}

	return nil
}

type primaryKey struct {
	tp      PKType
	key     map[string]any
	_fields []string
	_values []any
}

func (s *primaryKey) Type() PKType {
	return s.tp
}

func (s *primaryKey) GetInt() (int, error) {
	if !s.IsID() {
		return 0, fmt.Errorf("not ID primary key")
	}

	id, ok := s.key[KeyID].(int)
	if !ok {
		return 0, fmt.Errorf("problem with ID from primary key")
	}

	return id, nil
}

func (s *primaryKey) GetPartOfComplexInt(k string) (int, error) {
	if !s.IsComplex() {
		return 0, fmt.Errorf("Not complex primary key")
	}

	if val, ok := s.key[k]; !ok {
		return 0, fmt.Errorf("Part not found")
	} else {
		if intVal, intOk := val.(int); !intOk {
			return 0, fmt.Errorf("Part is not int")
		} else {
			return intVal, nil
		}
	}
}

func (s *primaryKey) GetPartOfComplexString(k string) (string, error) {
	if !s.IsComplex() {
		return "", fmt.Errorf("Not complex primary key")
	}

	if val, ok := s.key[k]; !ok {
		return "", fmt.Errorf("Part not found")
	} else {
		if strVal, intOk := val.(string); !intOk {
			return "", fmt.Errorf("Part is not int")
		} else {
			return strVal, nil
		}
	}
}

func (s *primaryKey) IsComplex() bool {
	return s.Type() == PkTypeComplex
}

func (s *primaryKey) IsID() bool {
	return s.Type() == PkTypeID
}

func (s *primaryKey) String() string {
	return fmt.Sprintf("%v", s.key)
}

func (s *primaryKey) OnlyEq() map[string]any {
	return s.key
}

func (s *primaryKey) IsEmpty() bool {
	switch s.tp {
	case PkTypeID:
		switch val := s.key[KeyID].(type) {
		case int, int64:
			if val == 0 {
				return true
			}
		case string:
			if val == "" {
				return true
			}
		default:
			return true
		}
	case PkTypeComplex:
		for _, part := range s.key {
			switch val := part.(type) {
			case int, int64:
				if val == 0 {
					return true
				}
			case string:
				if val == "" {
					return true
				}
			case *string:
				if val == nil || *val == "" {
					return true
				}
			default:
				return true
			}
		}
	default:
		return true
	}

	return false
}

func (s *primaryKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.key)
}

func (s *primaryKey) Fields() []string {
	return s._fields
}

func (s *primaryKey) ToCond() map[string]any {
	res := make(map[string]any, len(s.key))

	for field, value := range s.key {
		res[field+"_eq"] = value
	}

	return res
}

func (s *primaryKey) Scan(d Scanable) error {
	if err := d.Scan(s._values...); err != nil {
		return err
	}

	for i := range s._fields {
		s.key[s._fields[i]] = s._values[i]
	}

	return nil
}
