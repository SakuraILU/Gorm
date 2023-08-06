package engine

import (
	"database/sql"
	dialect "gorm/Dialect"
	log "gorm/Log"
	session "gorm/Session"

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
