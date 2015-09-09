package sqlf

import (
	"bytes"
	"fmt"
)

// Builder builds query and args useable for database/sql's API.
type Builder interface {
	Build() (query string, args []interface{})
}

// Printf is the main API; It simply creates SQL.
// The first argument format is that of fmt.Printf,
// with one special verb "%_", which tries to expand args
// to an SQL placeholder and args.
func Printf(format string, args ...interface{}) SQL {
	return SQL{
		Format: format,
		Args:   args, // TODO renaming to values?
	}
}

// SQL is a buildable query object; Call its Build() method
// to obtain query and args.
type SQL struct {
	Format string
	Args   []interface{}
}

func (s SQL) Build() (string, []interface{}) {
	args := []interface{}{}

	addArgs := func(a ...interface{}) {
		args = append(args, a...)
	}

	fargs := make([]interface{}, 0, len(s.Args))
	for _, arg := range s.Args {
		fargs = append(fargs, formatter{value: arg, addArgs: addArgs})
	}

	// In Sprintf, farg.Format's will be called and args filled.
	query := fmt.Sprintf(s.Format, fargs...)

	return query, args
}

// formatter implements fmt.Formatter and replaces %_ for SQL placeholders.
type formatter struct {
	// The underlying value
	value interface{}

	// SQL args for output
	addArgs func(args ...interface{})
}

func (f formatter) Format(s fmt.State, c rune) {
	if c == '_' {
		// The special SQL placeholder
		if b, ok := f.value.(Builder); ok {
			query, args := b.Build()
			s.Write([]byte(query))
			f.addArgs(args...)
		} else if vv, ok := f.value.([]interface{}); ok {
			s.Write(bytes.Repeat([]byte(",?"), len(vv))[1:])
			f.addArgs(vv...)
		} else {
			s.Write([]byte{'?'})
			f.addArgs(f.value)
		}
		return
	}

	// Fallback to default fmt implementation.
	fmt.Fprintf(s, "%"+string(c), f.value)
}
