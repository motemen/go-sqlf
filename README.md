# sqlf

[![GoDoc](https://godoc.org/github.com/motemen/go-sqlf?status.svg)](https://godoc.org/github.com/motemen/go-sqlf)
[![Build Status](https://travis-ci.org/motemen/go-sqlf.svg?branch=master)](https://travis-ci.org/motemen/go-sqlf)

Package sqlf provides Printf-like methods to generate SQL queries with placeholders.
It produces query and args which can be passed to database/sql APIs.
It assumes a special format verb "%\_" in addition to those of package fmt,
which expands to SQL placeholders.

For example, see the example for Printf.

## Examples

### Printf

```go
query, args := sqlf.Printf(
    "SELECT %s FROM %s WHERE col1 = %_ AND col2 IN (%_)",
    "id",    // SELECT %s
    "table", // FROM %s
    "x",     // col1 = %_
    []interface{}{1, 2, 3}, // col2 IN (%_)
).BuildSQL()

fmt.Println(query)
fmt.Println(args)
```

Output:

```
SELECT id FROM table WHERE col1 = ? AND col2 IN (?,?,?)
[x 1 2 3]
```

## Author

motemen <https://motemen.github.io/>
