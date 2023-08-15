package engine

import (
	"bytes"
	"database/sql/driver"
	"errors"
	log "gorm/Log"
	session "gorm/Session"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

// CreateTable, DropTable, HasTable test
func TestEngine0(t *testing.T) {
	type User struct {
		Name string `gorm:"PRIMARY KEY"`
		Age  int
	}
	type Car struct {
		User  string `gorm:"PRIMARY KEY"`
		Brand *string
		Year  int
		Price float64
	}

	cases := []struct {
		name  string
		model interface{}
	}{
		{"User", &User{}},
		{"Car", &Car{Brand: new(string)}},
	}

	engine, _ := NewEngine("sqlite3", "gee.db")
	defer engine.Close()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := engine.NewSession().Model(c.model)
			s.CreateTable()

			log.Infof("create table %v", c.model)
			if !s.HasTable() {
				t.Errorf("Table should exist")
			}

			s.DropTable()
			if s.HasTable() {
				t.Errorf("Table should not exist")
			}
		})
	}
}

// struct without pointer
func TestEngine1(t *testing.T) {
	type User struct {
		Name string `gorm:"PRIMARY KEY"`
		Age  int
	}

	cases := []struct {
		name string
		user User
	}{
		{"Tom&18", User{"Tom", 18}},
		{"Jack&22", User{"Jack", 22}},
		{"Sam&25", User{"Sam", 25}},
		{"Alice&18", User{"Alice", 18}},
		{"Bob&22", User{"Bob", 22}},
		{"Lily&20", User{"Lily", 20}},
	}

	engine, _ := NewEngine("sqlite3", "gee.db")
	defer engine.Close()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := engine.NewSession().Model(&c.user)
			s.CreateTable()

			log.Infof("create table %v", c.user)
			_, _ = s.Insert(&c.user)
			var users []User
			_ = s.Where("Name = ?", c.user.Name).Find(&users)
			if !reflect.DeepEqual(users[0], c.user) {
				t.Errorf("Failed to insert, got %v, expect %v", users[0], c.user)
			}
		})
	}
}

// struct with pointer
func TestEngine2(t *testing.T) {
	type Car struct {
		User  string `gorm:"PRIMARY KEY"`
		Brand *string
		Year  int
		Price float64
	}

	brands := []string{"BMW", "Benz", "Audi", "Lexus", "Toyota", "Honda", "Ford", "Tesla", "Volkswagen", "Nissan"}

	cases := []struct {
		name string
		car  Car
	}{
		{"Tom&18", Car{"Tom", &brands[0], 2001, 10000}},
		{"Jack&22", Car{"Jack", &brands[1], 2008, 20000}},
		{"Sam&25", Car{"Sam", &brands[2], 1997, 14000}},
		{"Alice&18", Car{"Alice", &brands[3], 2010, 30000}},
		{"Bob&22", Car{"Bob", &brands[4], 2021, 532000}},
		{"Lily&20", Car{"Lily", &brands[5], 2017, 100000}},
	}

	engine, _ := NewEngine("sqlite3", "gee.db")
	defer engine.Close()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := engine.NewSession().Model(&c.car)
			s.CreateTable()

			log.Infof("create table %v", c.car)
			_, _ = s.Insert(&c.car)
			var car Car
			_ = s.Where("User = ?", c.car.User).First(&car)
			if !reflect.DeepEqual(car, c.car) {
				t.Errorf("Failed to insert, got %v, expect %v", car, c.car)
			}
		})
	}
}

// update
func TestEngine3(t *testing.T) {
	type User struct {
		Name string `gorm:"PRIMARY KEY"`
		Age  int
	}

	cases := []struct {
		name        string
		user_before User
		update      map[string]any
		user_after  User
	}{
		{"Tom&18", User{"Tom", 18}, map[string]any{"Age": 19}, User{"Tom", 19}},
		{"Jack&22", User{"Jack", 22}, map[string]any{"Age": 23}, User{"Jack", 23}},
		{"Sam&25", User{"Sam", 25}, map[string]any{"Age": 14}, User{"Sam", 14}},
		{"Alice&18", User{"Alice", 18}, map[string]any{"Name": "AliceNew"}, User{"AliceNew", 18}},
		{"Bob&22", User{"Bob", 22}, map[string]any{"Age": 23, "Name": "BobNew"}, User{"BobNew", 23}},
	}

	engine, _ := NewEngine("sqlite3", "gee.db")
	defer engine.Close()
	s := engine.NewSession().Model(&User{})
	s.DropTable()
	s.CreateTable()
	for _, c := range cases {
		s.Insert(&c.user_before)
	}

	for _, c := range cases {
		s.Where("Name = ?", c.user_before.Name).Update(c.update)
	}

	for _, c := range cases {
		var users []User
		_ = s.Where("Name = ?", c.user_after.Name).Find(&users)
		if !reflect.DeepEqual(users[0], c.user_after) {
			t.Errorf("Failed to update, got %v, expect %v", users[0], c.user_after)
		}
	}
}

// first and order
func TestEngine4(t *testing.T) {
	type User struct {
		Name string `gorm:"PRIMARY KEY"`
		Age  int
	}

	cases := []struct {
		name string
		user User
	}{
		{"Tom&12", User{"Tom", 12}},
		{"Jack&22", User{"Jack", 22}},
		{"Sam&25", User{"Sam", 25}},
		{"Alice&18", User{"Alice", 18}},
		{"Bob&22", User{"Bob", 22}},
		{"Lily&20", User{"Lily", 30}},
	}

	engine, _ := NewEngine("sqlite3", "gee.db")
	defer engine.Close()
	s := engine.NewSession().Model(&User{})
	s.DropTable()
	s.CreateTable()
	for _, c := range cases {
		s.Insert(&c.user)
	}

	var user User
	_ = s.OrderBy("Age").First(&user)
	if !reflect.DeepEqual(user, cases[0].user) {
		t.Errorf("Failed to first, got %v, expect %v", user, cases[0].user)
	}

	_ = s.OrderBy("Age DESC").First(&user)
	if !reflect.DeepEqual(user, cases[len(cases)-1].user) {
		t.Errorf("Failed to first, got %v, expect %v", user, cases[len(cases)-1].user)
	}
}

// clauses: limit, orderby
func TestEngine5(t *testing.T) {
	type User struct {
		Name string `gorm:"PRIMARY KEY"`
		Age  int
	}

	cases := []struct {
		name string
		user User
	}{
		{"Tom&12", User{"Tom", 12}},
		{"Jack&22", User{"Jack", 22}},
		{"Sam&25", User{"Sam", 25}},
		{"Alice&18", User{"Alice", 18}},
		{"Bob&22", User{"Bob", 22}},
	}

	engine, _ := NewEngine("sqlite3", "gee.db")
	defer engine.Close()
	s := engine.NewSession().Model(&User{})
	s.DropTable()
	s.CreateTable()

	for _, c := range cases {
		s.Insert(&c.user)
	}

	var users []User
	_ = s.OrderBy("Age").Limit(3).Find(&users)
	if len(users) != 3 {
		t.Errorf("Failed to limit, got %v, expect %v", len(users), 3)
	}
	if !reflect.DeepEqual(users[0], cases[0].user) {
		t.Errorf("Failed to limit, got %v, expect %v", users[0], cases[0].user)
	}
	if !reflect.DeepEqual(users[1], cases[3].user) {
		t.Errorf("Failed to limit, got %v, expect %v", users[1], cases[3].user)
	}
	if !reflect.DeepEqual(users[2], cases[1].user) {
		t.Errorf("Failed to limit, got %v, expect %v", users[2], cases[1].user)
	}

	// Desc
	users = []User{}
	_ = s.OrderBy("Age DESC").Limit(3).Find(&users)
	if len(users) != 3 {
		t.Errorf("Failed to limit, got %v, expect %v", len(users), 3)
	}
	if !reflect.DeepEqual(users[0], cases[2].user) {
		t.Errorf("Failed to limit, got %v, expect %v", users[0], cases[2].user)
	}
	if !reflect.DeepEqual(users[1], cases[1].user) {
		t.Errorf("Failed to limit, got %v, expect %v", users[1], cases[1].user)
	}
	if !reflect.DeepEqual(users[2], cases[4].user) {
		t.Errorf("Failed to limit, got %v, expect %v", users[2], cases[4].user)
	}

}

// delete
func TestEngine6(t *testing.T) {
	type User struct {
		Name string `gorm:"PRIMARY KEY"`
		Age  int
	}

	cases_insert := []struct {
		name string
		user User
	}{
		{"Tom&12", User{"Tom", 12}},
		{"Jack&22", User{"Jack", 22}},
		{"Sam&25", User{"Sam", 25}},
		{"Alice&18", User{"Alice", 18}},
		{"Bob&22", User{"Bob", 22}},
	}

	cases_delete := []struct {
		name string
		user User
	}{}

	cases_left := []struct {
		name string
		user User
	}{}

	for idx, c := range cases_insert {
		if idx%2 == 0 {
			cases_delete = append(cases_delete, c)
		} else {
			cases_left = append(cases_left, c)
		}
	}

	engine, _ := NewEngine("sqlite3", "gee.db")
	defer engine.Close()
	s := engine.NewSession().Model(&User{})
	s.DropTable()
	s.CreateTable()

	for _, c := range cases_insert {
		s.Insert(&c.user)
	}

	for _, c := range cases_delete {
		s.Where("Name = ?", c.user.Name).Delete()
	}

	var users []User
	_ = s.Find(&users)

	if len(users) != len(cases_left) {
		t.Errorf("Failed to delete, got %v, expect %v", len(users), len(cases_left))
	}
	for idx, c := range cases_left {
		if !reflect.DeepEqual(users[idx], c.user) {
			t.Errorf("Failed to delete, got %v, expect %v", users[idx], c.user)
		}
	}
}

// hooks
type User struct {
	Name   string `gorm:"PRIMARY KEY"`
	Age    int
	Career string
}

func (u *User) BeforeQuery(s *session.Session) error {
	log.Infof("BeforeQuery name %v", reflect.TypeOf(u).Elem().Name())
	return nil
}

func (u *User) AfterQuery(s *session.Session) error {
	u.Career = "****"
	log.Info("AfterQuery")
	return nil
}

func (u *User) BeforeInsert(s *session.Session) error {
	log.Info("BeforeInsert")
	u.Age += 100
	return nil
}

func (u *User) AfterInsert(s *session.Session) error {
	log.Info("AfterInsert")
	return nil
}

func TestEngine7(t *testing.T) {
	// create a table
	engine, _ := NewEngine("sqlite3", "gee.db")
	defer engine.Close()
	s := engine.NewSession()

	// define several users
	users := []User{
		{"Jack", 35, "Worker"},
		{"Bob", 21, "Student"},
		{"Tom", 25, "Teacher"},
		{"Alice", 18, "Student"},
	}
	expected := User{"Alice", 118, "****"}

	// insert
	_ = s.Model(&User{}).DropTable()
	s.CreateTable()
	for _, user := range users {
		_, _ = s.Insert(&user)
	}

	// first
	var res User
	_ = s.Where("age > ?", 100).OrderBy("Age").First(&res) // every age add 100...so no one is filtered, will found Alice

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("AfterQuery error, got %v, expect %v", res, expected)
	}
}

// test commit
func TestEngine9(t *testing.T) {
	type User struct {
		Name   string `gorm:"PRIMARY KEY"`
		Age    int
		Career string
	}
	// create a table
	engine, _ := NewEngine("sqlite3", "gee.db")
	defer engine.Close()

	// define several users
	users := []User{
		{"Jack", 35, "Worker"},
		{"Bob", 21, "Student"},
		{"Tom", 25, "Teacher"},
	}

	// transcation
	_, err := engine.Transaction(func(s *session.Session) (any, error) {
		_ = s.Model(&User{}).DropTable()
		s.CreateTable()
		for _, user := range users {
			_, _ = s.Insert(&user)
		}
		return nil, nil
	})

	s := engine.NewSession()
	s.Model(&User{})
	if err != nil {
		t.Errorf("Transaction error")
	}
	var res []User
	_ = s.Find(&res)
	if len(res) != len(users) {
		t.Errorf("Transaction error, expect %v, got %v", len(users), len(res))
	}
	for i := 0; i < len(res); i++ {
		if res[i] != users[i] {
			t.Errorf("Transaction error")
		}
	}
}

// rollback test
func TestEngine10(t *testing.T) {
	// create a table
	engine, _ := NewEngine("sqlite3", "gee.db")
	defer engine.Close()

	// define several users
	users := []User{
		{"Jack", 35, "Worker"},
		{"Bob", 21, "Student"},
		{"Tom", 25, "Teacher"},
		{"Alice", 18, "Student"},
	}

	// drop table before transaction
	// otherwise, if drop in a transaction, it will be rollback
	// and the table will not be dropped, some remained data may cause error
	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_ = s.CreateTable()

	// transcation
	// rollback
	_, err := engine.Transaction(func(s *session.Session) (any, error) {
		_ = s.Model(&User{})
		for _, user := range users {
			_, _ = s.Insert(&user)
		}
		return nil, errors.New("rollback")
	})
	if err == nil {
		t.Errorf("Transaction error")
	}

	var res []User
	_ = s.Find(&res)
	if len(res) != 0 {
		t.Log(res)
		t.Errorf("Transaction error")
	}
}

type Car struct {
	Name         string `gorm:"PRIMARY KEY"`
	Brand        string
	HistoryPrice ints
}

type ints []int

func (i *ints) Scan(value interface{}) error {
	buf := bytes.NewBuffer(value.([]byte))
	for {
		v, err := buf.ReadString(',')
		log.Warnf("Scan %v", v)
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
	for _, v := range i {
		buf.WriteString(strconv.Itoa(v))
		buf.WriteString(",")
	}
	return buf.Bytes(), nil
}

func TestEngine12(t *testing.T) {
	engine, _ := NewEngine("sqlite3", "gee.db")
	defer engine.Close()
	engine.Migrate(&Car{})

	s := engine.NewSession()
	_ = s.Model(&Car{})
	_ = s.DropTable()
	_ = s.CreateTable()

	// define several users
	cars := []Car{
		{"Jack", "BMW", ints{10000, 20000, 30000}},
		{"Bob", "Benz", ints{20000, 30000, 40000}},
		{"Tom", "Audi", ints{30000, 40000, 50000}},
		{"Alice", "Lexus", ints{40000, 50000, 60000}},
	}

	// insert
	for _, car := range cars {
		_, _ = s.Insert(&car)
	}

	// find
	var res2 []Car
	_ = s.Find(&res2)
	if !reflect.DeepEqual(res2, cars) {
		t.Errorf("Failed to find, got %v, expect %v", res2, cars)
	}
}
