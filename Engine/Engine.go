package engine

import (
	"database/sql"
	log "gorm/Log"
	session "gorm/Session"

	_ "github.com/mattn/go-sqlite3"
)

type Engine struct {
	db *sql.DB
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
	e = &Engine{
		db: db,
	}
	return
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error(err)
	}

	log.Info("Close database success")
}

func (s *Engine) NewSession() *session.Session {
	return session.NewSession(s.db)
}
