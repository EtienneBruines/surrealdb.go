package sql

import (
	"context"
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"

	"github.com/surrealdb/surrealdb.go"
)

var newLineRegex = regexp.MustCompile(`\r?\n`)
var tabRegex = regexp.MustCompile(`\t`)

type Conn struct {
	*surrealdb.DB
}

func (s *Conn) Prepare(query string) (driver.Stmt, error) {
	//method, thing, err := s.parseMethod(query)
	//if err != nil {
	//	return nil, fmt.Errorf("invalid rawQuery: %w", err)
	//}

	return s.PrepareContext(context.Background(), query)
}

func (s *Conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	return &Stmt{
		conn:     s,
		rawQuery: query,
		//method:   method,
		//thing:    thing,
	}, nil
}

func (s *Conn) Close() error {
	s.DB.Close()
	return nil
}

func (s *Conn) Begin() (driver.Tx, error) {
	return nil, fmt.Errorf("this method is deprecated")
}

func (s *Conn) Ping(ctx context.Context) error {
	// TODO: Is there something more reliable?
	// TODO: How do we utilize context.Context ?
	_, err := s.Select("1")
	return err
}

func (s *Conn) ResetSession(ctx context.Context) error {
	return nil // We can do some cleanup here, once needed
}

func (s *Conn) IsValid() bool {
	return true // Might change once we have something that will invalidate the connection
}

func (s *Conn) parseMethod(query string) (string, string, error) {
	idx := strings.IndexRune(query, ' ')
	if idx <= 0 {
		return query, "", nil
	}

	// TODO: do we validate this?
	return strings.ToLower(query[:idx]), query[idx+1:], nil
}

func (s *Conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	actual := NamedValuesToValues(args)
	return s.Execute(ctx, query, actual)
}

func (s *Conn) ExecerContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	actual := NamedValuesToValues(args)
	rows, err := s.Execute(ctx, query, actual)
	if err != nil {
		return nil, err
	}

	return Result{
		AffectedRows: int64(len(rows.RawData)),
	}, nil
}

func (s *Conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	actual := NamedValuesToValues(args)
	rows, err := s.Execute(ctx, query, actual)
	if err != nil {
		return nil, err
	}

	return Result{
		AffectedRows: int64(len(rows.RawData)),
	}, nil
}

func (s *Conn) Execute(ctx context.Context, query string, args []driver.Value) (*Rows, error) {
	query = PrepareQuery(query)
	argInterfaces := make([]interface{}, len(args)+1)
	argInterfaces[0] = strings.TrimSpace(query)

	for idx, arg := range args {
		argInterfaces[idx+1] = s.convertArgument(arg)
	}

	res, err := s.Send("query", argInterfaces...)
	if err != nil {
		return nil, fmt.Errorf("error during Exec: %w", err)
	}

	arr, ok := res.([]interface{})
	if !ok || len(arr) != 1 {
		// No idea what the result is
		return nil, fmt.Errorf("unknown result")
	}

	lookup, ok := arr[0].(map[string]interface{})
	if !ok {
		// No idea what the result is
		return nil, fmt.Errorf("unknown result, expected map")
	}

	status, _ := lookup["status"]
	//duration, _ := lookup["time"]

	switch status {
	case "ERR":
		detail, _ := lookup["detail"]
		return nil, fmt.Errorf("query error: %s", detail)
	case "OK":
		result, _ := lookup["result"]
		rows, ok := result.([]interface{})
		if !ok {
			return nil, fmt.Errorf("unknown result value")
		}

		return &Rows{RawData: rows}, nil
	default:
		return nil, fmt.Errorf("unknown response status: %s", status)
	}
}

func (s *Conn) convertArgument(val driver.Value) interface{} {
	return val
}

// TODO: Potentially implement: ExecerContext, QueryerContext, ConnPrepareContext, and ConnBeginTx.
