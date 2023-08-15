package schema

import (
	"fmt"
	dialect "gorm/Dialect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchema(t *testing.T) {
	//a struct
	type User struct {
		Name string `gorm:"PRIMARY KEY"`
		Age  int
	}

	type Car struct {
		User  []byte  `gorm:"PRIMARY KEY"`
		Brand *string `gorm:"PRIMARY KEY"`
		Year  int
		Price float64
	}

	cases := []struct {
		name                 string
		model                any
		expected_name        string
		expected_field_names []string
		expected_field_types []string
		expected_field_tags  []string
	}{
		{
			name:                 "User",
			model:                &User{},
			expected_name:        "User",
			expected_field_names: []string{"Name", "Age"},
			expected_field_types: []string{"text", "integer"},
			expected_field_tags:  []string{"PRIMARY KEY", ""},
		},
		{
			name: "Car",
			model: &Car{
				Brand: new(string),
			},
			expected_name:        "Car",
			expected_field_names: []string{"User", "Brand", "Year", "Price"},
			expected_field_types: []string{"blob", "text", "integer", "real"},
			expected_field_tags:  []string{"PRIMARY KEY", "PRIMARY KEY", "", ""},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := &dialect.Sqlite3{}
			schema := NewSchema(c.model, d)
			assert.Equal(t, c.expected_name, NewSchema(c.model, d).GetName(), fmt.Sprintf("schema name should be %s", c.expected_name))
			for idx, f := range schema.GetFields() {
				assert.Equal(t, c.expected_field_names[idx], f.Name, fmt.Sprintf("schema field name should be %s", c.expected_field_names[idx]))
				assert.Equal(t, c.expected_field_types[idx], f.Type, fmt.Sprintf("schema field type should be %s", c.expected_field_types[idx]))
				assert.Equal(t, c.expected_field_tags[idx], f.Tag, fmt.Sprintf("schema field tag should be %s", c.expected_field_tags[idx]))
			}
		})
	}
}
