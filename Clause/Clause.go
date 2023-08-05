package clause

import "strings"

type Clause struct {
	typcmds map[Type]string
	typvals map[Type][]any
}

func NewClause() *Clause {
	return &Clause{
		typcmds: make(map[Type]string),
		typvals: make(map[Type][]any),
	}
}

func (c *Clause) Set(typ Type, vars ...any) {
	cmd, val := generators[typ](vars...)
	c.typcmds[typ] = cmd
	c.typvals[typ] = val
}

func (c *Clause) Build(typs ...Type) (string, []any) {
	cmds := make([]string, 0)
	vals := make([]any, 0)

	for _, typ := range typs {
		cmds = append(cmds, c.typcmds[typ])
		vals = append(vals, c.typvals[typ]...)
	}
	return strings.Join(cmds, " "), vals
}
