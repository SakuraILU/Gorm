package engine

import (
	dialect "gorm/Dialect"
	"log"
	"testing"
)

type User struct {
	Name   string `gorm:"PRIMARY KEY"`
	Age    int
	Career string
}

// sqlite3 table exist test
func TestEngine1(t *testing.T) {
	// create a table
	engine, _ := NewEngine("sqlite3", "tmp.db")
	defer engine.Close()
	s := engine.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()

	// test table exist
	d := &dialect.Sqlite3{}
	sql, vals := d.TableExistSQL("User")
	if sql != "SELECT name FROM sqlite_master WHERE type='table' and name = ?" {
		t.Errorf("TableExistSQL error")
	}
	if vals != "User" {
		t.Errorf("TableExistSQL error")
	}

	// test whether the table exists
	r, _ := s.Raw(sql, vals).Query()
	if !r.Next() {
		t.Errorf("Table should exist")
	}
}

// CreateTable, DropTable, HasTable test
func TestEngine2(t *testing.T) {
	// create a table
	engine, _ := NewEngine("sqlite3", "tmp.db")
	defer engine.Close()
	s := engine.NewSession()

	s.Model(&User{}).CreateTable()
	if !s.HasTable() {
		t.Errorf("Table should exist")
	}

	s.Model(&User{}).DropTable()
	if s.HasTable() {
		t.Errorf("Table should not exist")
	}
}

// insert
func TestEngine3(t *testing.T) {
	// create a table
	engine, _ := NewEngine("sqlite3", "tmp.db")
	defer engine.Close()
	s := engine.NewSession()

	// define several users
	users := []User{
		{"Jack", 35, "Worker"},
		{"Bob", 21, "Student"},
		{"Tom", 25, "Teacher"},
		{"Alice", 18, "Student"},
	}
	// insert
	_ = s.Model(&User{}).DropTable()
	s.CreateTable()
	for _, user := range users {
		_, _ = s.Insert(&user)
	}

	// find
	var res []User
	_ = s.Find(&res)
	if len(res) != len(users) {
		t.Errorf("Find error")
	}
	for i := 0; i < len(res); i++ {
		log.Println(res[i])
		if res[i] != users[i] {
			t.Errorf("Find error")
		}
	}
}
