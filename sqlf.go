package sqlf

import (
	"bytes"
	"fmt"
	"reflect"
	"sync"
)

// Builder builds query and args usable for database/sql's API.
type Builder interface {
	BuildSQL() (query string, args []interface{})
}

// Printf is the main API; It simply creates SQL.
// The first argument format is that of fmt.Printf,
// with one special verb "%_", which tries to expand args
// to an SQL placeholder and args.
func Printf(format string, values ...interface{}) SQL {
	return SQL{
		Format: format,
		Values: values,
	}
}

// SQL is a buildable query object; Call its BuildSQL() method
// to obtain query and args.
type SQL struct {
	Format string
	Values []interface{}
}

// BuildSQL produces query and args usable for database/sql's API.
func (s SQL) BuildSQL() (string, []interface{}) {
	args := []interface{}{}

	var mutex sync.Mutex
	addArgs := func(a ...interface{}) {
		mutex.Lock()
		defer mutex.Unlock()
		args = append(args, a...)
	}

	values := make([]interface{}, 0, len(s.Values))
	for _, v := range s.Values {
		values = append(values, formatter{value: v, addArgs: addArgs})
	}

	// In Sprintf, farg.Format's will be called and args filled.
	query := fmt.Sprintf(s.Format, values...)

	return query, args
}

// formatter implements fmt.Formatter.
// It replaces fmt verb "%_" by SQL placeholders.
type formatter struct {
	// The underlying value
	value interface{}

	// Callback func to add resulting SQL args
	addArgs func(args ...interface{})
}

func (f formatter) Format(s fmt.State, c rune) {
	if c == '_' {
		// "%_", The special SQL placeholder
		if b, ok := f.value.(Builder); ok {
			// Embed already-built SQL
			query, args := b.BuildSQL()
			s.Write([]byte(query))
			f.addArgs(args...)
		} else if vv, ok := f.value.([]interface{}); ok {
			// Expand slices
			// e.g. "IN (%_)", []{1,2,3} -> "IN (?,?,?)", 1, 2, 3
			s.Write(bytes.Repeat([]byte(",?"), len(vv))[1:])
			f.addArgs(vv...)
		} else if rv := reflect.ValueOf(f.value); rv.Type().Kind() == reflect.Slice {
			// Any slices
			l := rv.Len()
			s.Write(bytes.Repeat([]byte(",?"), l)[1:])
			for i := 0; i < l; i++ {
				f.addArgs(rv.Index(i).Interface())
			}
		} else {
			s.Write([]byte{'?'})
			f.addArgs(f.value)
		}
		return
	}

	// Fallback to default fmt implementation.
	fmt.Fprintf(s, "%"+string(c), f.value)
}
