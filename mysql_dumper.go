package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println(usage())
		os.Exit(1)
	}
	mysqlURI := os.Args[1]
	queryPath := os.Args[2]

	db, err := sql.Open("mysql", mysqlURI)
	if err != nil {
		panic("DB connect fail: " + err.Error())
	}
	if err = db.Ping(); err != nil {
		panic("DB connect fail: " + err.Error())
	}

	queryBytes, err := ioutil.ReadFile(queryPath)
	if err != nil {
		panic(fmt.Sprintf("Cannot read query from %s: %s", queryPath, err.Error()))
	}
	rows, err := db.Query(string(queryBytes))
	if err != nil {
		panic(fmt.Sprintf("Failed to execute query: %s", err.Error()))
	}
	dump(rows)
	rows.Close()
	db.Close()
}

func dump(rows *sql.Rows) {
	columns, _ := rows.Columns()

	// Make a slice for the values
	values := make([]interface{}, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// just so we know when to stop printing tabs
	lastColumn := len(columns) - 1
	rowCount := 0
	rowChunk := 25000
	start := time.Now()
	lap := start
	var elapsed time.Duration

	var displayers []func(interface{})
	for rows.Next() {
		rowCount += 1
		rows.Scan(scanArgs...)

		// first row creates the displayers...
		if displayers == nil {
			displayers = createDisplayers(values)
		}

		// then ever row thereafter executes them
		for idx, value := range values {
			displayers[idx](value)
			if lastColumn != idx {
				fmt.Print("\t")
			}
		}
		fmt.Print("\n")
		if rowCount%rowChunk == 0 {
			elapsed = time.Since(lap)
			os.Stderr.Write([]byte(fmt.Sprintf("%d: %5.2f sec @ %5.2f rows/sec\n",
				rowCount, elapsed.Seconds(), float64(rowChunk)/elapsed.Seconds())))
			lap = time.Now()
		}
	}

	elapsed = time.Since(start)
	os.Stderr.Write([]byte(fmt.Sprintf("DONE: %d in %5.2f sec @ %5.2f rows/sec\n",
		rowCount, elapsed.Seconds(), float64(rowCount)/elapsed.Seconds())))
}

func createDisplayers(values []interface{}) []func(interface{}) {
	displayers := make([]func(interface{}), len(values))
	for idx, value := range values {
		switch value.(type) {
		case nil:
			displayers[idx] = func(v interface{}) { fmt.Print("NULL") }
		case []byte:
			displayers[idx] = func(v interface{}) { fmt.Print(string(v.([]byte))) }
		case bool:
			displayers[idx] = func(v interface{}) { fmt.Printf("%t", value) }
		default:
			displayers[idx] = func(v interface{}) { fmt.Print(value) }
		}
	}
	return displayers
}
func usage() string {
	return `
Usage:
  ./mysql_dumper mysql-db-uri path/to/query-file.sql

Examples:

  ./mysql_dumper root:@/local?loc=UTC&parseTime=true popular_products.sql > popular_products.txt
  ./mysql_dumper cwinters:blahblah@some-replica.prod.example.com/products?loc=UTC&parseTime=true awesome-users.sql > awesome_users.txt

Note that you can accomplish much the same with the mysql client:

  mysql --quick --compress --batch --raw --skip-column_names --user=you --password=secret \
    --host=some-host.example.com --database=mydb < some_query.sql > some_query_results.txt
	`
}
