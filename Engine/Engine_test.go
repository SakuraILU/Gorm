package engine

import (
	"bytes"
	"database/sql/driver"
	"strconv"
	"strings"
	"testing"
)

// type User struct {
// 	Name   string `gorm:"PRIMARY KEY"`
// 	Age    int
// 	Career string
// }

// // sqlite3 table exist test
// func TestEngine1(t *testing.T) {
// 	// create a table
// 	engine, _ := NewEngine("sqlite3", "tmp.db")
// 	defer engine.Close()
// 	s := engine.NewSession()
// 	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
// 	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()

// 	// test table exist
// 	d := &dialect.Sqlite3{}
// 	sql, vals := d.TableExistSQL("User")
// 	if sql != "SELECT name FROM sqlite_master WHERE type='table' and name = ?" {
// 		t.Errorf("TableExistSQL error")
// 	}
// 	if vals != "User" {
// 		t.Errorf("TableExistSQL error")
// 	}

// 	// test whether the table exists
// 	r, _ := s.Raw(sql, vals).Query()
// 	if !r.Next() {
// 		t.Errorf("Table should exist")
// 	}
// }

// // CreateTable, DropTable, HasTable test
// func TestEngine2(t *testing.T) {
// 	// create a table
// 	engine, _ := NewEngine("sqlite3", "tmp.db")
// 	defer engine.Close()
// 	s := engine.NewSession()

// 	s.Model(&User{}).CreateTable()
// 	if !s.HasTable() {
// 		t.Errorf("Table should exist")
// 	}

// 	s.Model(&User{}).DropTable()
// 	if s.HasTable() {
// 		t.Errorf("Table should not exist")
// 	}
// }

// // insert
// func TestEngine3(t *testing.T) {
// 	// create a table
// 	engine, _ := NewEngine("sqlite3", "tmp.db")
// 	defer engine.Close()
// 	s := engine.NewSession()

// 	// define several users
// 	users := []User{
// 		{"Jack", 35, "Worker"},
// 		{"Bob", 21, "Student"},
// 		{"Tom", 25, "Teacher"},
// 		{"Alice", 18, "Student"},
// 	}
// 	// insert
// 	_ = s.Model(&User{}).DropTable()
// 	s.CreateTable()
// 	for _, user := range users {
// 		_, _ = s.Insert(&user)
// 	}

// 	// find
// 	var res []User
// 	_ = s.Find(&res)
// 	if len(res) != len(users) {
// 		t.Errorf("Find error")
// 	}
// 	for i := 0; i < len(res); i++ {
// 		log.Info(res[i])
// 		if res[i] != users[i] {
// 			t.Errorf("Find error")
// 		}
// 	}
// }

// // Limit, Where, OrderBy, Select test
// func TestEngine4(t *testing.T) {
// 	// create a table
// 	engine, _ := NewEngine("sqlite3", "tmp.db")
// 	defer engine.Close()
// 	s := engine.NewSession()

// 	// define several users
// 	users := []User{
// 		{"Jack", 35, "Worker"},
// 		{"Bob", 21, "Student"},
// 		{"Tom", 25, "Teacher"},
// 		{"Alice", 18, "Student"},
// 	}

// 	// insert
// 	_ = s.Model(&User{}).DropTable()
// 	s.CreateTable()
// 	for _, user := range users {
// 		_, _ = s.Insert(&user)
// 	}
// 	// insert again
// 	for _, user := range users {
// 		_, _ = s.Insert(&user)
// 	}
// 	// insert again
// 	for _, user := range users {
// 		_, _ = s.Insert(&user)
// 	}

// 	// find
// 	var res []User
// 	_ = s.Limit(2).Where("name = ? AND age = ?", "Jack", 35).OrderBy("age DESC").Find(&res)
// 	if len(res) != 1 {
// 		t.Errorf("Find error")
// 	}
// 	if res[0] != users[0] {
// 		t.Errorf("Find error")
// 	}

// 	// find 3 users
// 	var res2 []User
// 	_ = s.Limit(3).Where("age > ?", 3).OrderBy("age DESC").Find(&res2)
// 	if len(res2) != 3 {
// 		t.Errorf("Find error")
// 	}
// }

// // Update test
// func TestEngine5(t *testing.T) {
// 	// create a table
// 	engine, _ := NewEngine("sqlite3", "tmp.db")
// 	defer engine.Close()
// 	s := engine.NewSession()

// 	// define several users
// 	users := []User{
// 		{"Jack", 35, "Worker"},
// 		{"Bob", 21, "Student"},
// 		{"Tom", 25, "Teacher"},
// 		{"Alice", 18, "Student"},
// 	}

// 	// insert
// 	_ = s.Model(&User{}).DropTable()
// 	s.CreateTable()
// 	for _, user := range users {
// 		_, _ = s.Insert(&user)
// 	}

// 	// update
// 	_, _ = s.Where("name = ?", "Jack").Update(map[string]any{"age": 36, "career": "Engineer"})
// 	var res []User
// 	_ = s.Find(&res)
// 	if res[0].Age != 36 || res[0].Career != "Engineer" {
// 		t.Errorf("Update error")
// 	}
// }

// // Delete test
// func TestEngine6(t *testing.T) {
// 	// create a table
// 	engine, _ := NewEngine("sqlite3", "tmp.db")
// 	defer engine.Close()
// 	s := engine.NewSession()

// 	// define several users
// 	users := []User{
// 		{"Jack", 35, "Worker"},
// 		{"Bob", 21, "Student"},
// 		{"Tom", 25, "Teacher"},
// 		{"Alice", 18, "Student"},
// 	}

// 	// insert
// 	_ = s.Model(&User{}).DropTable()
// 	s.CreateTable()
// 	for _, user := range users {
// 		_, _ = s.Insert(&user)
// 	}

// 	// delete
// 	_, _ = s.Where("name = ?", "Jack").Delete()
// 	var res []User
// 	_ = s.Find(&res)
// 	if len(res) != 3 {
// 		t.Errorf("Delete error")
// 	}
// }

// // first test
// func TestEngine7(t *testing.T) {
// 	// create a table
// 	engine, _ := NewEngine("sqlite3", "tmp.db")
// 	defer engine.Close()
// 	s := engine.NewSession()

// 	// define several users
// 	users := []User{
// 		{"Jack", 35, "Worker"},
// 		{"Bob", 21, "Student"},
// 		{"Tom", 25, "Teacher"},
// 		{"Alice", 18, "Student"},
// 	}

// 	// insert
// 	_ = s.Model(&User{}).DropTable()
// 	s.CreateTable()
// 	for _, user := range users {
// 		_, _ = s.Insert(&user)
// 	}

// 	// first
// 	var res User
// 	_ = s.Where("age > ?", 18).OrderBy("age").First(&res)

// 	if res != users[1] {
// 		t.Errorf("First error")
// 	}
// }

// type UserHook struct {
// 	Name   string `gorm:"PRIMARY KEY"`
// 	Age    int
// 	Career string
// }

// // test hooks
// func (u *UserHook) BeforeQuery(s *session.Session) error {
// 	log.Infof("BeforeQuery name %v", reflect.TypeOf(u).Elem().Name())
// 	return nil
// }

// func (u *UserHook) AfterQuery(s *session.Session) error {
// 	u.Career = "****"
// 	log.Info("AfterQuery")
// 	return nil
// }

// func (u *UserHook) BeforeInsert(s *session.Session) error {
// 	log.Info("BeforeInsert")
// 	u.Age += 100
// 	return nil
// }

// func (u *UserHook) AfterInsert(s *session.Session) error {
// 	log.Info("AfterInsert")
// 	return nil
// }

// // test hooks
// func TestEngine8(t *testing.T) {
// 	// create a table
// 	engine, _ := NewEngine("sqlite3", "tmp.db")
// 	defer engine.Close()
// 	s := engine.NewSession()

// 	// define several users
// 	users := []UserHook{
// 		{"Jack", 35, "Worker"},
// 		{"Bob", 21, "Student"},
// 		{"Tom", 25, "Teacher"},
// 		{"Alice", 18, "Student"},
// 	}

// 	// insert
// 	_ = s.Model(&UserHook{}).DropTable()
// 	s.CreateTable()
// 	for _, user := range users {
// 		_, _ = s.Insert(&user)
// 	}

// 	// first
// 	var res UserHook
// 	_ = s.Where("age > ?", 18).OrderBy("age").First(&res) // every age add 100...so no one is filtered, will found Alice
// 	t.Log(res)
// 	// check hook
// 	if res.Career != "****" {
// 		t.Errorf("AfterQuery error")
// 	}
// 	if res.Age != 118 {
// 		t.Errorf("BeforeInsert error")
// 	}
// }

// // test commit
// func TestEngine9(t *testing.T) {
// 	// create a table
// 	engine, _ := NewEngine("sqlite3", "tmp.db")
// 	defer engine.Close()

// 	// define several users
// 	users := []User{
// 		{"Jack", 35, "Worker"},
// 		{"Bob", 21, "Student"},
// 		{"Tom", 25, "Teacher"},
// 	}

// 	// drop table before transaction
// 	// otherwise, if drop in a transaction, it will be rollback
// 	// and the table will not be dropped, some remained data will cause error
// 	s := engine.NewSession()
// 	_ = s.Model(&User{}).DropTable()

// 	// transcation
// 	_, err := engine.Transaction(func(s *session.Session) (any, error) {
// 		_ = s.Model(&User{}).DropTable()
// 		s.CreateTable()
// 		for _, user := range users {
// 			_, _ = s.Insert(&user)
// 		}
// 		return nil, nil
// 	})

// 	// check
// 	if err != nil {
// 		t.Errorf("Transaction error")
// 	}
// 	var res []User
// 	_ = s.Find(&res)
// 	if len(res) != len(users) {
// 		t.Errorf("Transaction error")
// 	}
// 	for i := 0; i < len(res); i++ {
// 		if res[i] != users[i] {
// 			t.Errorf("Transaction error")
// 		}
// 	}
// }

// // rollback test
// func TestEngine10(t *testing.T) {
// 	// create a table
// 	engine, _ := NewEngine("sqlite3", "tmp.db")
// 	defer engine.Close()

// 	// define several users
// 	users := []User{
// 		{"Jack", 35, "Worker"},
// 		{"Bob", 21, "Student"},
// 		{"Tom", 25, "Teacher"},
// 		{"Alice", 18, "Student"},
// 	}

// 	// drop table before transaction
// 	// otherwise, if drop in a transaction, it will be rollback
// 	// and the table will not be dropped, some remained data will cause error
// 	s := engine.NewSession()
// 	_ = s.Model(&User{}).DropTable()

// 	// transcation
// 	// rollback
// 	_, err := engine.Transaction(func(s *session.Session) (any, error) {
// 		_ = s.Model(&User{}).DropTable()
// 		s.CreateTable()
// 		for _, user := range users {
// 			_, _ = s.Insert(&user)
// 		}
// 		return nil, errors.New("rollback")
// 	})
// 	if err == nil {
// 		t.Errorf("Transaction error")
// 	}

// 	var res []User
// 	_ = s.Find(&res)
// 	if len(res) != 0 {
// 		t.Log(res)
// 		t.Errorf("Transaction error")
// 	}
// }

// test migrate
// because two different user defined...,
// if db has no User table, you need to comment below and uncomment above,
// run some tests to create user table and insert some entries (TestEngine4 for example)
// then comment above and uncomment below, test the migration function
//
// however, if db has a User table, then just run below tests
type User struct {
	Name    string `gorm:"PRIMARY KEY"`
	Age     int
	ID      int
	History ints
}

type ints []int

func (i *ints) Scan(value interface{}) error {
	buf := bytes.NewBuffer(value.([]byte))
	for {
		v, err := buf.ReadString(',')
		if err != nil {
			break
		}
		v = strings.Trim(v, ",")
		vi, _ := strconv.Atoi(v)
		*i = append(*i, vi)
	}
	return nil
}

func (i ints) Value() (driver.Value, error) {
	buf := bytes.NewBuffer([]byte{})
	for idx, v := range i {
		if idx != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(strconv.Itoa(v))
	}
	return buf.Bytes(), nil
}

func TestEngine12(t *testing.T) {
	engine, _ := NewEngine("sqlite3", "tmp.db")
	defer engine.Close()
	engine.Migrate(&User{})

	s := engine.NewSession()
	_ = s.Model(&User{})

	// define several users
	users := []User{
		{"Jack", 35, 1245, []int{1, 2, 3}},
		{"Bob", 21, 1246, []int{4, 5, 6}},
		{"Tom", 25, 1247, []int{7, 8, 9}},
		{"Alice", 18, 1248, []int{10, 11, 12}},
	}

	// update id
	for _, user := range users {
		_, _ = s.Where("name = ?", user.Name).Update("id", user.ID, "history", user.History)
	}

	// find
	var res2 []User
	_ = s.Find(&res2)
	if len(res2) != len(users) {
		t.Errorf("Find error")
	}
	for i := 0; i < len(res2); i++ {
		if res2[i].Age != users[i].Age {
			t.Errorf("Find error")
		}
	}
}
