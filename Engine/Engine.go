package engine

import (
	"database/sql"
	"fmt"
	dialect "gorm/Dialect"
	log "gorm/Log"
	schema "gorm/Schema"
	session "gorm/Session"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Engine struct {
	db   *sql.DB
	dial dialect.Dialect
}

func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	log.Infof("Connect to the database %v(driver: %v)", source, driver)

	dial, err := dialect.GetDialect(driver)
	if err != nil {
		return
	}

	e = &Engine{
		db:   db,
		dial: dial,
	}
	return
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error(err)
	}

	log.Info("Close database success")
}

func (e *Engine) NewSession() *session.Session {
	return session.NewSession(e.db, e.dial)
}

func (e *Engine) Transaction(fn func(s *session.Session) (any, error)) (value any, err error) {
	s := e.NewSession()

	if err = s.Begin(); err != nil {
		return
	}

	defer func() {
		if err != nil {
			s.Rollback()
		} else if err := recover(); err != nil {
			s.Rollback()
			panic(err)
		} else {
			s.Commit()
		}
	}()

	return fn(s)
}

func (e *Engine) Migrate(value any) (err error) {
	_, err = e.Transaction(func(s *session.Session) (res any, err error) {
		s.Model(value)
		rows, err := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", s.RefTable().GetName())).Query()
		if err != nil {
			log.Error(err)
			return
		}
		oldfnames, err := rows.Columns()
		if err != nil {
			log.Error(err)
			return
		}

		tablename := s.RefTable().GetName()
		newfs := s.RefTable().GetFields()

		addfs := newFields(newfs, oldfnames)
		for _, f := range addfs {
			if _, err = s.Raw(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", tablename, f.Name, f.Type)).Exec(); err != nil {
				log.Error(err)
				return
			}
		}

		tmpname := tablename + "_tmp2e23rwq23"
		if _, err = s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s", tablename, tmpname)).Exec(); err != nil {
			log.Error(err)
			return
		}

		if err = s.CreateTable(); err != nil {
			log.Error(err)
			return
		}

		fieldnames := s.RefTable().GetFieldNames()
		cols := strings.Join(fieldnames, ",")
		if _, err = s.Raw(fmt.Sprintf("INSERT INTO %s SELECT %s FROM %s", tablename, cols, tmpname)).Exec(); err != nil {
			log.Error(err)
			return
		}

		if _, err = s.Raw(fmt.Sprintf("DROP TABLE %s", tmpname)).Exec(); err != nil {
			log.Error(err)
			return
		}

		return
	})
	return
}

func newFields(fs []*schema.Field, oldnames []string) (res []*schema.Field) {
	oldnamemap := make(map[string]any)
	for _, oldname := range oldnames {
		oldnamemap[oldname] = struct{}{}
	}

	for _, f := range fs {
		if _, ok := oldnamemap[f.Name]; !ok {
			res = append(res, f)
		}
	}
	return
}
