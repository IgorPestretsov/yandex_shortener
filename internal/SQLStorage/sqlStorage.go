package SQLStorage

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(dsn string) *Storage {
	s := Storage{}
	var err error
	s.db, err = sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	_, err = s.db.Exec("CREATE TABLE IF NOT EXISTS links (key VARCHAR(50), url VARCHAR(500) UNIQUE, user_id VARCHAR(50));")
	if err != nil {
		panic(err)
	}
	return &s
}

func (s *Storage) LoadLinksPair(key string) string {
	var output sql.NullString
	_ = s.db.QueryRow("select url from links where key=$1", key).Scan(&output)

	if output.Valid {
		return output.String
	} else {
		return ""
	}
}

func (s *Storage) SaveLinksPair(uid string, link string, key string) (string, error) {
	_, err := s.db.Exec("insert into links(key,url,user_id) values ($1,$2,$3);", key, link, uid)
	if err != nil {
		s.db.QueryRow("select key from links where url=$1", link).Scan(&key)
		return key, err
	} else {
		return "", nil
	}

}

func (s *Storage) GetAllUserURLs(uid string) map[string]string {
	output := make(map[string]string)
	rows, err := s.db.Query("select url, key from links where user_id=$1", uid)
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

func (s Storage) Close() {
	s.db.Close()
}
