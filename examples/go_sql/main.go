package main

import (
	"database/sql"
	"fmt"
	"time"

	sql2 "github.com/surrealdb/surrealdb.go/sql"

	_ "github.com/surrealdb/surrealdb.go/sql"
)

type user struct {
	sql2.BaseEntity
	Name string `json:"name"`
	Age  int    `json:"age"`
}
type post struct {
	sql2.BaseEntity
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

type posted struct {
	sql2.BaseRelationship
	When time.Time `json:"when"`
}

type queryResp struct {
	user
	Posts  sql2.Many[post]   `json:"posts"`
	Posted sql2.Many[posted] `json:"posted"`
}

func main() {

	// Connect the way you would usually
	db, err := sql.Open("surrealdb", "ws://root:root@localhost:8000/?ns=test&db=test")
	if err != nil {
		panic(err)
	}

	// clear tables
	func() {
		db.Exec("DELETE user")
		db.Exec("DELETE post")
		db.Exec("DELETE posted")
	}()

	// Make sure we can ping to it
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// Create some value
	_, err = db.Exec("CREATE user:mark SET name = 'mark', age = 9999")
	if err != nil {
		panic(err)
	}

	for i := 0; i < 20; i++ {
		_, err = db.Exec(fmt.Sprintf("CREATE post:%v SET title = '%v', slug='%v'", i, i, i))
		if err != nil {
			panic(err)
		}

		_, err = db.Exec(fmt.Sprintf("RELATE user:mark->posted->post:%v SET when=time::now()", i))
		if err != nil {
			panic(err)
		}
	}

	// Read it back
	rows, err := db.Query("SELECT id, name, ->posted as posted, ->posted->post AS posts FROM user:mark fetch posts, posted")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		resp := queryResp{}

		// Note: the columns are always sorted alphabetically, regardless of the order in your query
		err := rows.Scan(&resp.ID, &resp.Name, &resp.Posted, &resp.Posts)
		if err != nil {
			panic(err)
		}

		fmt.Println("row", resp)
	}

	// get records by named params
	stmt, err := db.Prepare("SELECT id, name, ->posted as posted, ->posted->post AS posts FROM user where name = $user_name fetch posts, posted")
	if err != nil {
		panic(err)
	}
	rows, err = stmt.Query(sql.Named("user_name", "mark"))
	for rows.Next() {
		resp := queryResp{}
		err := rows.Scan(&resp.ID, &resp.Name, &resp.Posted, &resp.Posts)
		if err != nil {
			panic(err)
		}

		fmt.Println("row", resp)
	}

	// Cleanup
	err = db.Close()
	if err != nil {
		panic(err)
	}
}
