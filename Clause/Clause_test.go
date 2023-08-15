package clause

import (
	"reflect"
	"testing"
)

// insert
func TestClause1(t *testing.T) {
	// INSERT INTO User (name, age, career) VALUES (?,?,?), (?,?,?)
	c := NewClause()
	c.Set(INSERT, "User", []string{"name", "age", "career"})
	c.Set(VALUES, []any{"Jack", "Worker", 35}, []any{"Bob", "Student", 21})
	cmd, vals := c.Build(INSERT, VALUES)

	res_cmd := "INSERT INTO User (name, age, career) VALUES (?,?,?), (?,?,?)"
	res_vals := []any{"Jack", "Worker", 35, "Bob", "Student", 21}

	if !reflect.DeepEqual(cmd, res_cmd) {
		t.Errorf("cmd: %s != %s", cmd, res_cmd)
	}
	if !reflect.DeepEqual(vals, res_vals) {
		t.Errorf("vals: %v != %v", vals, res_vals)
	}
}

// select
func TestClause2(t *testing.T) {
	// SELECT name, age, career FROM User WHERE name = ? AND age = ? ORDER BY age LIMIT 1
	c := NewClause()
	c.Set(SELECT, "User", []string{"name", "age", "career"})
	c.Set(WHERE, "name = ? AND age = ?", "Jack", 35)
	c.Set(ORDERBY, "age")
	c.Set(LIMIT, 1)
	cmd, vals := c.Build(SELECT, WHERE, ORDERBY, LIMIT)

	res_cmd := "SELECT name,age,career FROM User WHERE name = ? AND age = ? ORDER BY age LIMIT ?"
	res_vals := []any{"Jack", 35, 1}

	if !reflect.DeepEqual(cmd, res_cmd) {
		t.Errorf("cmd: %s != %s", cmd, res_cmd)
	}
	if !reflect.DeepEqual(vals, res_vals) {
		t.Errorf("vals: %v != %v", vals, res_vals)
	}
}

// order by desc
func TestClause3(t *testing.T) {
	// select and order by desc
	c := NewClause()
	c.Set(SELECT, "User", []string{"name", "age", "career"})
	c.Set(WHERE, "name = ? AND age = ?", "Jack", 35)
	c.Set(ORDERBY, "age DESC")
	c.Set(LIMIT, 1)
	cmd, vals := c.Build(SELECT, WHERE, ORDERBY, LIMIT)

	res_cmd := "SELECT name,age,career FROM User WHERE name = ? AND age = ? ORDER BY age DESC LIMIT ?"
	res_vals := []any{"Jack", 35, 1}

	if !reflect.DeepEqual(cmd, res_cmd) {
		t.Errorf("cmd: %s != %s", cmd, res_cmd)
	}
	if !reflect.DeepEqual(vals, res_vals) {
		t.Errorf("vals: %v != %v", vals, res_vals)
	}
}

// update
func TestClause4(t *testing.T) {
	// UPDATE User SET name = ?, age = ?, career = ? WHERE name = ? AND age = ?
	c := NewClause()
	c.Set(UPDATE, "User", map[string]any{"age": 35, "career": "Worker"})
	c.Set(WHERE, "name = ? AND age = ? AND career = ?", "Jack", 24, "Student")
	cmd, vals := c.Build(UPDATE, WHERE)

	res_cmd := "UPDATE User SET age = ?, career = ? WHERE name = ? AND age = ? AND career = ?"
	res_vals := []any{35, "Worker", "Jack", 24, "Student"}

	if !reflect.DeepEqual(cmd, res_cmd) {
		t.Errorf("cmd: %s != %s", cmd, res_cmd)
	}
	if !reflect.DeepEqual(vals, res_vals) {
		t.Errorf("vals: %v != %v", vals, res_vals)
	}
}
