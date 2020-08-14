package cache

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type Cache struct {
	Key string
	Value string
}

func init() {
	c := &Cache{}
	c.init("./spartan.db")

}

func (c *Cache) init(dbpath string) {
	dir := "cache"
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		fmt.Println("creating cache directory")
		err := os.Mkdir(dir, 0755)
		if err != nil {
			fmt.Printf("os.Mkdir('%s') failed with '%s'\n", dir)
		}
	}

	//dblayout
	db, err := sql.Open("sqlite3",dbpath)
	if err != nil {
		fmt.Printf("error: failed to open db: %v", err)
	}

	defer db.Close()

	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS cache (id TEXT PRIMARY KEY, value TEXT)")
	if err != nil {
		fmt.Printf("error: failed to initialize db: %v", err)
	}

	_, err = statement.Exec()
	if err != nil {
		fmt.Printf("error: failed to initialize statement exec: %v", err)
	}
}

func (c *Cache)Set(key string, value string)  {
	db, err := sql.Open("sqlite3","./spartan.db" )
	if err != nil {
		fmt.Printf("error: failed to open db: %v", err)
	}

	defer db.Close()

	statement, _ := db.Prepare("INSERT INTO cache(id, value) VALUES (?, ?)")
	//if err != nil {
	//	fmt.Printf("error: failed to insert db: %v", err)
	//}

	_, _ = statement.Exec(key, value)
	//if err != nil {
	//	fmt.Printf("error: failed to exec insert db: %v", err)
	//}

}

func (c *Cache)Get(key string) string {
	db, err := sql.Open("sqlite3","./spartan.db" )
	if err != nil {
		fmt.Printf("error: failed to open db: %v", err)
	}
	defer db.Close()

	var value string
	_ = db.QueryRow("SELECT value FROM cache WHERE id =?", key).Scan(&value)
	//if err != nil {
	//	fmt.Printf("error: failed to get cache value db: %v", err)
	//}

	return value
}

func (c *Cache)Delete(key string) {
	db, err := sql.Open("sqlite3","./spartan.db" )
	if err != nil {
		fmt.Printf("error: failed to open db: %v", err)
	}
	defer db.Close()

	statement, err := db.Prepare("DELETE FROM cache WHERE id = ?")
	if err != nil {
		fmt.Printf("error: failed to insert db: %v", err)
	}
	statement.Exec(key)

}

func (c *Cache)ClearAll() {
	db, err := sql.Open("sqlite3","./spartan.db" )
	if err != nil {
		fmt.Printf("error: failed to open db: %v", err)
	}
	defer db.Close()

	statement, err := db.Prepare("DELETE FROM cache")
	if err != nil {
		fmt.Printf("error: failed to remove all cache records: %v", err)
	}
	statement.Exec()
}

func (c *Cache)DisplayAll() {
	db, err := sql.Open("sqlite3","./spartan.db" )
	if err != nil {
		fmt.Printf("error: failed to open db: %v", err)
	}
	defer db.Close()

	row, err := db.Query("SELECT * FROM cache ORDER BY id")
	if err != nil {
		fmt.Printf("error: failed to display all cache records: %v", err)
	}
	defer row.Close()

	for row.Next() { // Iterate and fetch the records from result cursor
		var id string
		var value string
		row.Scan(&id, &value)
		fmt.Println(id + ":" + value)
	}
}

