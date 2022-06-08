package sqlstorage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"log"
)

type Storage struct {
	db      *sql.DB
	cleaner Cleaner
	buffer  []RecordToDelete
}

func New(dsn string) *Storage {
	s := Storage{}
	s.cleaner = NewCleaner(&s)
	s.cleaner.Run()

	var err error
	s.db, err = sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	_, err = s.db.Exec("CREATE TABLE IF NOT EXISTS links (" +
		"key VARCHAR(50), " +
		"url VARCHAR(500) UNIQUE, " +
		"user_id VARCHAR(50), " +
		"is_deleted BOOLEAN DEFAULT FALSE)" +
		";")
	if err != nil {
		panic(err)
	}
	return &s
}

func (s *Storage) LoadLinksPair(key string) string {
	var output sql.NullString
	var isDeleted sql.NullString
	_ = s.db.QueryRow("select url, is_deleted from links where key=$1", key).Scan(&output, &isDeleted)
	if isDeleted.Valid {
		if isDeleted.String == "true" {
			return "gone"
		}

	}
	if output.Valid {
		return output.String
	} else {
		return ""
	}
}

func (s *Storage) SaveLinksPair(uid string, link string, key string) (string, error) {
	_, err := s.db.Exec("insert into links(key,url,user_id) values ($1,$2,$3);", key, link, uid)
	var pqErr *pq.Error

	if errors.As(err, &pqErr) && pqErr.Code == pgerrcode.UniqueViolation {
		s.db.QueryRow("select key from links where url=$1", link).Scan(&key)
		err = NewAlreadyExistErr(link, err)
		return key, err
	}
	return "", err
}

func (s *Storage) GetAllUserURLs(uid string) map[string]string {
	output := make(map[string]string)
	rows, _ := s.db.Query("select url, key from links where user_id=$1", uid)
	err := rows.Err()

	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var url, key string
		err = rows.Scan(&url, &key)
		if err != nil {
			panic(err)
		}
		output[key] = url
	}

	return output
}

func (s *Storage) GetChannelForDelete() chan map[string]string {
	return s.cleaner.UserDeleteRequests
}
func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) CheckReqToDelete(r RecordToDelete) bool {
	var resp string
	err := s.db.QueryRow("select user_id from links where key=$1", r.urlID).Scan(&resp)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if resp == r.userID {
		return true
	}
	return false

}
func (s *Storage) DeleteRecords(delBatch []RecordToDelete) {
	tx, err := s.db.Begin()
	if err != nil {
		log.Fatalf("unable to open db connection")
	}
	stmt, err := tx.Prepare("update links set is_deleted=TRUE where key=$1;")
	if err != nil {
		log.Fatalf("db error")
	}

	for _, r := range delBatch {
		if _, err = stmt.Exec(r.urlID); err != nil {
			if err = tx.Rollback(); err != nil {
				log.Fatalf("update drivers: unable to rollback: %v", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("update drivers: unable to commit: %v", err)
	}

	s.buffer = s.buffer[:0]

}
