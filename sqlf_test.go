package sqlf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintf(t *testing.T) {
	assert := assert.New(t)

	query, args := Printf(
		"SELECT %s FROM %s WHERE col1 = %_ AND col2 IN (%_)",
		"id",
		"table",
		"value",
		[]interface{}{1, 2, 3},
	).Build()

	assert.Equal(
		"SELECT id FROM table WHERE col1 = ? AND col2 IN (?,?,?)",
		query,
	)

	assert.Equal(
		args,
		[]interface{}{
			"value",
			1, 2, 3,
		},
	)
}

func TestPrintf_Builder(t *testing.T) {
	assert := assert.New(t)

	wherePart := Printf("col1 IN (%_)", []interface{}{"x", "y"})

	query, args := Printf(
		"SELECT id FROM table WHERE %_ AND col2 = %_",
		wherePart,
		"z",
	).Build()

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
