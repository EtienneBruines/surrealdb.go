package sql

import (
	"context"
	"database/sql/driver"
	"fmt"
	"regexp"
)

var inputRegex = regexp.MustCompile("\\$[a-zA-Z]+")

type Stmt struct {
	conn *Conn

	rawQuery string
	method   string
	thing    string
}

func (s *Stmt) Close() error {
	return nil
}

func (s *Stmt) NumInput() int {
	inputs := inputRegex.FindAllString(s.rawQuery, -1)
	return len(inputs)
}

func (s *Stmt) Exec(args []driver.Value) (driver.Result, error) {
	rows, err := s.conn.Execute(context.Background(), s.rawQuery, args)
	if err != nil {
		return nil, err
	}

	cols := rows.Columns()
	data := make([]driver.Value, len(cols))

	driverResult := Result{}
	if err := rows.Next(data); err != nil {
		driverResult.AffectedRows++
	}
	return driverResult, nil
}

func (s *Stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	named := map[string]any{}
	for _, a := range args {
		named[a.Name] = a.Value
	}

	actual := []driver.Value{
		named,
	}
	rows, err := s.conn.Execute(ctx, s.rawQuery, actual)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (s *Stmt) Query(args []driver.Value) (driver.Rows, error) {
	rows, err := s.conn.Execute(context.Background(), s.rawQuery, args)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

type Result struct {
	AffectedRows int64
}

func (r Result) LastInsertId() (int64, error) {
	return 0, fmt.Errorf("surrealDB does not support numeric/int64 auto-increment ids")
}

func (r Result) RowsAffected() (int64, error) {
	return r.AffectedRows, nil
}
