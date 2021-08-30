package mysql

import (
	"database/sql"
	"errors"
	"golangify.com/snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (s *SnippetModel) Insert(title, content, expires string) (int, error) {
	command := "INSERT INTO snippets (title, content, created, expires) VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))"
	result, err := s.DB.Exec(command, title, content, expires)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
func (s *SnippetModel) Get(id int) (*models.Snippet, error){
	snp := &models.Snippet{}
	command := "SELECT id, title, content, created, expires FROM snippets\n WHERE expires > UTC_TIMESTAMP() AND id = ?"
	row := s.DB.QueryRow(command, id)
	err := row.Scan(&snp.ID, &snp.Title, &snp.Content, &snp.Created, &snp.Expires)
	//err := s.DB.QueryRow("SELECT id, title, content, created, expires FROM snippets\n    WHERE expires > UTC_TIMESTAMP() AND id = ?", id).Scan(&snp.ID, &snp.Title, &snp.Content, &snp.Created, &snp.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrorNoRecord
		} else {
			return nil, err
		}
	}
	return snp, nil
}
func (s *SnippetModel) Latest() ([]*models.Snippet, error) {
	command := "SELECT id, title, content, created, expires FROM snippets\n    WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10"
	row, err := s.DB.Query(command)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	var snippets []*models.Snippet
	for row.Next(){
		snp := &models.Snippet{}
		err := row.Scan(&snp.ID, &snp.Title, &snp.Content, &snp.Created, &snp.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, snp)
	}
	if err = row.Err(); err != nil{
		return nil, err
	}
	return snippets, nil
}