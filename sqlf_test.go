package sqlf_test

import (
	"fmt"
	"testing"

	"github.com/motemen/go-sqlf"
	"github.com/stretchr/testify/assert"
)

func TestPrintf(t *testing.T) {
	assert := assert.New(t)

	query, args := sqlf.Printf(
		"SELECT %s FROM %s WHERE col1 = %_ AND col2 IN (%_) AND col3 = %_",
		"id",
		"table",
		"value",
		[]interface{}{1, 2, 3},
		[]byte("hello"),
	).BuildSQL()

	assert.Equal(
		"SELECT id FROM table WHERE col1 = ? AND col2 IN (?,?,?) AND col3 = ?",
		query,
	)

	assert.Equal(
		args,
		[]interface{}{
			"value",
			1, 2, 3,
			[]byte("hello"),
		},
	)
}

func TestPrintf_Builder(t *testing.T) {
	assert := assert.New(t)

	wherePart := sqlf.Printf("col1 IN (%_)", []interface{}{"x", "y"})

	query, args := sqlf.Printf(
		"SELECT id FROM table WHERE %_ AND col2 = %_",
		wherePart,
		"z",
	).BuildSQL()

	assert.Equal(
		"SELECT id FROM table WHERE col1 IN (?,?) AND col2 = ?",
		query,
	)

	assert.Equal(
		args,
		[]interface{}{
			"x", "y",
			"z",
		},
	)
}

func ExamplePrintf() {
	query, args := sqlf.Printf(
		"SELECT %s FROM %s WHERE col1 = %_ AND col2 IN (%_)",
		"id",    // SELECT %s
		"table", // FROM %s
		"x",     // col1 = %_
		[]interface{}{1, 2, 3}, // col2 IN (%_)
	).BuildSQL()

	fmt.Println(query)
	fmt.Println(args)
	// Output:
	// SELECT id FROM table WHERE col1 = ? AND col2 IN (?,?,?)
	// [x 1 2 3]
}
