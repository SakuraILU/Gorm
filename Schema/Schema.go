package schema

import (
	"fmt"
	dialect "gorm/Dialect"
	log "gorm/Log"
	"reflect"
)

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

type Schema struct {
	model any

	name       string
	namefields map[string]*Field
	fieldnames []string
	fields     []*Field

	methods map[string]*reflect.Value
}

func NewSchema(v any, dial dialect.Dialect) (s *Schema) {
	// it may be a pointer...so indirect it's value if neccessary
	typ := reflect.Indirect(reflect.ValueOf(v)).Type()

	s = &Schema{
		model:      v,
		name:       typ.Name(),
		namefields: make(map[string]*Field),
		fieldnames: make([]string, 0),
		fields:     make([]*Field, 0),
		methods:    make(map[string]*reflect.Value),
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

	ptyp := reflect.ValueOf(v).Type()
	for i := 0; i < ptyp.NumMethod(); i++ {
		m := ptyp.Method(i)
		log.Infof("schema %s method %s", s.name, m.Name)
		s.methods[m.Name] = &m.Func
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
	for _, field := range s.fields {
		log.Info(field.Name, field.Type, field.Tag)
	}
	return s.fields
}

func (s *Schema) GetMdoel() any {
	return s.model
}

func (s *Schema) GetMethod(name string) (fn *reflect.Value, err error) {
	fn, ok := s.methods[name]
	if !ok {
		err = fmt.Errorf("method %s not found", name)
		log.Error(err)
	}
	return
}
