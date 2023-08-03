package schema

import (
	dialect "gorm/Dialect"
	"testing"
)

func TestNewSchema1(t *testing.T) {
	// generate a sqlite3 dialect
	d := &dialect.Sqlite3{}

	//a struct
	type User struct {
		Name string `gorm:"PRIMARY KEY"`
		Age  int
	}

	// create a schema
	s := NewSchema(&User{}, d)

	// test schema
	if s.GetName() != "User" {
		t.Errorf("schema name should be User")
	}
	if len(s.GetFieldNames()) != 2 {
		t.Errorf("schema field names length should be 2")
	}
	if s.GetField("Name").Name != "Name" {
		t.Errorf("schema field name should be Name")
	}
	if s.GetField("Age").Type != "integer" {
		t.Errorf("schema field type should be integer")
	}
}

// has primary key and pointer
func TestNewSchema2(t *testing.T) {
	// generate a sqlite3 dialect
	d := &dialect.Sqlite3{}

	//a struct
	type User struct {
		Name *string `gorm:"PRIMARY KEY"`
		Age  int
	}

	// create a schema
	s := NewSchema(&User{}, d)

	// test schema
	if s.GetName() != "User" {
		t.Errorf("schema name should be User")
	}
	if len(s.GetFieldNames()) != 2 {
		t.Errorf("schema field names length should be 2")
	}
	if s.GetField("Name").Name != "Name" {
		t.Errorf("schema field name should be Name")
	}
	if s.GetField("Age").Type != "integer" {
		t.Errorf("schema field type should be integer")
	}
	if s.GetField("Name").Tag != "PRIMARY KEY" {
		t.Errorf("schema field tag should be PRIMARY KEY")
	}
}
