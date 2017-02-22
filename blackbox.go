package main

import (
	"fmt"

	"github.com/cznic/ql"
)

func init() {

}

func main() {
	db, err := ql.OpenMem()
	if err != nil {
		panic(err)
	}

	rss, _, err := db.Run(ql.NewRWCtx(), `
	BEGIN TRANSACTION;
		CREATE TABLE foo (i int);
		INSERT INTO foo VALUES (10), (20);
		CREATE TABLE bar (fooID int, s string);
		INSERT INTO bar SELECT id(), "ten" FROM foo WHERE i == 10;
		INSERT INTO bar SELECT id(), "twenty" FROM foo WHERE i == 20;
	COMMIT;
	SELECT *
	FROM foo, bar
	WHERE bar.fooID == id(foo)
	ORDER BY id(foo);`,
	)
	if err != nil {
		panic(err)
	}

	for _, rs := range rss {
		if err := rs.Do(false, func(data []interface{}) (bool, error) {
			fmt.Println(data)
			return true, nil
		}); err != nil {
			panic(err)
		}
		fmt.Println("----")
	}
}
