package squirrel

import (
	"fmt"
	"reflect"
	"strings"
)

// Eq is syntactic sugar for use with the Where/Having methods.
// Ex:
//     .Where(Eq{"id": 1})
type Eq map[string]interface{}

type wherePart struct {
	pred interface{}
	args []interface{}
}

func newWherePart(pred interface{}, args ...interface{}) wherePart {
	return wherePart{pred: pred, args: args}
}

func wherePartsToSql(parts []wherePart, placeholder string) (string, []interface{}, error) {
	sqls := make([]string, 0, len(parts))
	var args []interface{}
	for _, part := range parts {
		partSql, partArgs, err := part.ToSql(placeholder)
		if err != nil {
			return "", []interface{}{}, err
		}
		if len(partSql) > 0 {
			sqls = append(sqls, partSql)
			args = append(args, partArgs...)
		}
	}
	return strings.Join(sqls, " AND "), args, nil
}

func (p wherePart) ToSql(placeholder string) (sql string, args []interface{}, err error) {
	switch pred := p.pred.(type) {
	case nil:
		// no-op
	case Eq:
		return whereEqMap(map[string]interface{}(pred), placeholder)
	case map[string]interface{}:
		return whereEqMap(pred, placeholder)
	case string:
		sql = pred
		args = p.args
	default:
		err = fmt.Errorf("expected string-keyed map or string, not %T", pred)
	}
	return
}

func whereEqMap(m map[string]interface{}, placeholder string) (sql string, args []interface{}, err error) {
	var exprs []string
	for key, val := range m {
		expr := ""
		if val == nil {
			expr = fmt.Sprintf("%s IS NULL", key)
		} else {
			valVal := reflect.ValueOf(val)
			if valVal.Kind() == reflect.Array || valVal.Kind() == reflect.Slice {
				placeholders := make([]string, valVal.Len())
				for i := 0; i < valVal.Len(); i++ {
					placeholders[i] = placeholder
					args = append(args, valVal.Index(i).Interface())
				}
				placeholdersStr := strings.Join(placeholders, ",")
				expr = fmt.Sprintf("%s IN (%s)", key, placeholdersStr)
			} else {
				expr = fmt.Sprintf("%s = %s", key, placeholder)
				args = append(args, val)
			}
		}
		exprs = append(exprs, expr)
	}
	sql = strings.Join(exprs, " AND ")
	return
}
