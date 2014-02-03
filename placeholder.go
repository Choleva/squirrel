package squirrel

import (
	"bytes"
	"fmt"
	"strings"
)

type placeholderFormat string

var (
	Question = placeholderFormat("?")
	Dollar   = placeholderFormat("$?")
)

func (f placeholderFormat) String() string {
	return string(f)
}

func (f placeholderFormat) ReplacePlaceholders(sql string) string {
	switch f {
	case Dollar:
		return replaceDollars(sql)
	default:
		return sql
	}
}

// Replace $? with $1, $2, etc
func replaceDollars(sql string) string {
	buf := &bytes.Buffer{}
	for i := 1;; i++ {
		p := strings.Index(sql, Dollar.String())
		if p == -1 {
			break
		}

		buf.WriteString(sql[:p])
		fmt.Fprintf(buf, "$%d", i)
		sql = sql[p + len(Dollar):]
	}

	buf.WriteString(sql)
	return buf.String()
}
