package dialect

import (
	engine "gorm/Engine"
	"log"
	"testing"
	"time"
)

// sqlite3 data type test
func TestDialect1(t *testing.T) {
	// generate a sqlite3 dialect
	d := &Sqlite3{}

	// test several data types
	var v int64 = 1
	if d.DataTypeOf(v) != "bigint" {
		t.Errorf("int64 type should be bigint")
	}
	var v2 int = 1
	if d.DataTypeOf(v2) != "integer" {
		t.Errorf("int type should be integer")
	}
	var v3 string = "hello"
	if d.DataTypeOf(v3) != "text" {
		t.Errorf("string type should be text")
	}
	var v4 bool = true
	if d.DataTypeOf(v4) != "bool" {
		t.Errorf("bool type should be bool")
	}
	var v5 float32 = 1.0
	if d.DataTypeOf(v5) != "real" {
		t.Errorf("float32 type should be real")
	}
	var v6 []byte = []byte("hello")
	if d.DataTypeOf(v6) != "blob" {
		t.Errorf("[]byte type should be blob")
	}
	var v7 time.Time = time.Now()
	if d.DataTypeOf(v7) != "datetime" {
		t.Errorf("time.Time type should be datetime")
	}
}

// sqlite3 table exist test
func TestDialect2(t *testing.T) {
	// generate a sqlite3 dialect
	d := &Sqlite3{}

	// create a table
	engine, _ := engine.NewEngine("sqlite3", "tmp.db")
	defer engine.Close()
	s := engine.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()

	// test table exist
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
	log.Fatal()
}
