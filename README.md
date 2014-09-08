Wee standalone Go utilities, mostly so I can see how my ignorance progresses
over time.

Each one should be standalone and buildable with:

    $ go build some_tidbit.go

For example:

    $ go build mysql_dumper.go

Running a tidbit with no arguments should give you a little more 
information. For example:

```
$ ./mysql_dumper

Usage:
  ./mysql_dumper mysql-db-uri path/to/query-file.sql

Examples:

  ./mysql_dumper root:@/local?loc=UTC&parseTime=true popular_products.sql > popular_products.txt
  ./mysql_dumper cwinters:blahblah@some-replica.prod.example.com/products?loc=UTC&parseTime=true awesome-users.sql > awesome_users.txt

Note that you can accomplish much the same with the mysql client:

  mysql --quick --compress --batch --raw --skip-column_names --user=you --password=secret \
    --host=some-host.example.com --database=mydb < some_query.sql > some_query_results.txt
```
