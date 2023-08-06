package engine

import (
	dialect "gorm/Dialect"
	log "gorm/Log"
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
		log.Info(res[i])
		if res[i] != users[i] {
			t.Errorf("Find error")
		}
	}
}

// Limit, Where, OrderBy, Select test
func TestEngine4(t *testing.T) {
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
	_ = s.Limit(2).Where("name = ? AND age = ?", "Jack", 35).OrderBy("age DESC").Find(&res)
	if len(res) != 1 {
		t.Errorf("Find error")
	}
	if res[0] != users[0] {
		t.Errorf("Find error")
	}

	// find 3 users
	var res2 []User
	_ = s.Limit(3).Where("age > ?", 3).OrderBy("age DESC").Find(&res2)
	if len(res2) != 3 {
		t.Errorf("Find error")
	}
}

// Update test
func TestEngine5(t *testing.T) {
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

	// update
	_, _ = s.Where("name = ?", "Jack").Update(map[string]any{"age": 36, "career": "Engineer"})
	var res []User
	_ = s.Find(&res)
	if res[0].Age != 36 || res[0].Career != "Engineer" {
		t.Errorf("Update error")
	}
}

// Delete test
func TestEngine6(t *testing.T) {
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

	// delete
	_, _ = s.Where("name = ?", "Jack").Delete()
	var res []User
	_ = s.Find(&res)
	if len(res) != 3 {
		t.Errorf("Delete error")
	}
}

// first test
func TestEngine7(t *testing.T) {
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

	// first
	var res User
	_ = s.Where("age > ?", 18).OrderBy("age").First(&res)

	if res != users[1] {
		t.Errorf("First error")
	}
}
