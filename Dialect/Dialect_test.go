package dialect

import (
	"testing"
	"time"
)

// sqlite3 data type test
func TestDialect1(t *testing.T) {
	// generate a sqlite3 dialect
	d := &Sqlite3{}

	cases := []struct {
		name   string
		in     interface{}
		expect string
	}{
		{"int", 1, "integer"},
		{"int64", int64(1), "bigint"},
		{"uint8", uint8(1), "unsigned tinyint"},
		{"string", "hello", "text"},
		{"bool", true, "bool"},
		{"float32", float32(1.0), "real"},
		{"[]byte", []byte("hello"), "blob"},
		{"time.Time", time.Now(), "datetime"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if d.DataTypeOf(c.in) != c.expect {
				t.Errorf("expect %s, but got %s", c.expect, d.DataTypeOf(c.in))
			}
		})
	}
}
