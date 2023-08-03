package session

import (
	"database/sql"
	log "gorm/Log"
	"strings"
)

type Session struct {
	db      *sql.DB
	sqlcmds strings.Builder
	sqlvals []any
}

func NewSession(db *sql.DB) (s *Session) {
	return &Session{
		db:      db,
		sqlcmds: strings.Builder{},
		sqlvals: make([]any, 0),
	}
}

func (s *Session) DB() *sql.DB {
	return s.db
}

func (s *Session) Raw(sqlcmd string, sqlvals ...any) *Session {
	s.sqlcmds.WriteString(sqlcmd)
	s.sqlvals = append(s.sqlvals, sqlvals...)
	return s
}

func (s *Session) Clear() {
	s.sqlcmds.Reset()
	s.sqlvals = make([]any, 0)
}

func (s *Session) Exec() (r sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sqlcmds.String(), s.sqlvals)
	r, err = s.db.Exec(s.sqlcmds.String(), s.sqlvals...)
	if err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) Query() (rs *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sqlcmds.String(), s.sqlvals)
	rs, err = s.db.Query(s.sqlcmds.String(), s.sqlvals...)
	if err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) QueryRow() (r *sql.Row) {
	defer s.Clear()
	log.Info(s.sqlcmds.String(), s.sqlvals)
	// the error of QueryRow will be deferred until the Scan() method
	return s.db.QueryRow(s.sqlcmds.String(), s.sqlvals...)
}
