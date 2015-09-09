package sqlf

import (
	"bytes"
	"fmt"
)

type Builder interface {
	Build() (query string, args []interface{})
}

func Printf(format string, args ...interface{}) SQL {
	return SQL{
		Format: format,
		Args:   args,
	}
}

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
