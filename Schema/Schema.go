package schema

import (
	dialect "gorm/Dialect"
	log "gorm/Log"
	"reflect"
)

type Schema struct {
	name       string
	namefields map[string]*Field
	fieldnames []string
	fields     []*Field
}

func NewSchema(v any, dial dialect.Dialect) (s *Schema) {
	// it may be a pointer...so indirect it's value if neccessary
	typ := reflect.Indirect(reflect.ValueOf(v)).Type()

	s = &Schema{
		name:       typ.Name(),
		namefields: make(map[string]*Field),
		fieldnames: make([]string, 0),
		fields:     make([]*Field, 0),
	}

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		field := &Field{
			Name: f.Name,
			Type: dial.DataTypeOf(reflect.Indirect(reflect.New(f.Type)).Interface()),
			Tag:  f.Tag.Get("gorm"),
		}
		s.namefields[f.Name] = field
		s.fieldnames = append(s.fieldnames, f.Name)
		s.fields = append(s.fields, field)
	}
	return
}

func (s *Schema) GetName() string {
	return s.name
}

func (s *Schema) GetFieldNames() []string {
	return s.fieldnames
}

func (s *Schema) GetField(name string) *Field {
	return s.namefields[name]
}

func (s *Schema) GetFields() []*Field {
	log.Warnf("fields %v", s.fields)
	return s.fields
}
