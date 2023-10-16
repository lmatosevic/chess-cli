package repository

import (
	"database/sql"
	"fmt"
	"github.com/lmatosevic/chess-cli/pkg/utils"
	"strings"
)

// PrepareQueryParams returns where, sort, order, and args parameters used for further building of SQL query.
// If filter is empty the where param will be an empty string. The filter must match following pattern:
// (key1)(operator1)(value1);(and|or|not)?;(key2)(operator2)(value2)...
//
// Supported operators: = (equals), != (not equals), <= (less than or equal), >= (greater than or equal), < (less than),
// > (greater than), ~ (like), !~ (not like), -> (in), !-> (not in)
//
// Example filter: id>=1;or;username=someName;and;startedAt!=null
//
// Returned args will contain at least two values at the end: limit and offset. If where query needs additional
// parameters, they will be added in args slice.
func PrepareQueryParams(filter string, page int, size int, sort string) (string, string, string, []any) {
	var where = ""
	var args []any

	if filter != "" {
		where = "WHERE"
		filterParts := strings.Split(filter, ";")
		needOperator := false
		for _, f := range filterParts {
			if strings.ToLower(f) == "and" || strings.ToLower(f) == "or" || strings.ToLower(f) == "not" {
				where = fmt.Sprintf("%s %s", where, f)
				needOperator = false
				continue
			}

			if needOperator && f != "" {
				where = fmt.Sprintf("%s %s", where, "and")
			}

			needOperator = true

			values := strings.SplitN(f, "<=", 2)
			if len(values) > 1 {
				where = fmt.Sprintf(`%s "%s" <= $%d`, where, values[0], len(args)+1)
				args = append(args, values[1])
				continue
			}

			values = strings.SplitN(f, ">=", 2)
			if len(values) > 1 {
				where = fmt.Sprintf(`%s "%s" >= $%d`, where, values[0], len(args)+1)
				args = append(args, values[1])
				continue
			}

			values = strings.SplitN(f, "!=", 2)
			if len(values) > 1 {
				if strings.ToLower(values[1]) == "null" {
					where = fmt.Sprintf(`%s "%s" IS NOT NULL`, where, values[0])
				} else {
					where = fmt.Sprintf(`%s ("%s" != $%d or "%s" IS NULL)`, where, values[0], len(args)+1, values[0])
					args = append(args, values[1])
				}
				continue
			}

			values = strings.SplitN(f, "=", 2)
			if len(values) > 1 {
				if strings.ToLower(values[1]) == "null" {
					where = fmt.Sprintf(`%s "%s" IS NULL`, where, values[0])
				} else {
					where = fmt.Sprintf(`%s "%s" = $%d`, where, values[0], len(args)+1)
					args = append(args, values[1])
				}
				continue
			}

			values = strings.SplitN(f, "!->", 2)
			if len(values) > 1 {
				var argNums []string
				elems := strings.Split(values[1], ",")
				for i := 0; i < len(elems); i++ {
					argNums = append(argNums, fmt.Sprintf("$%d", len(args)+1))
					args = append(args, elems[i])
				}
				where = fmt.Sprintf(`%s "%s" NOT IN (%s)`, where, values[0], strings.Join(argNums, ","))
				continue
			}

			values = strings.SplitN(f, "->", 2)
			if len(values) > 1 {
				var argNums []string
				elems := strings.Split(values[1], ",")
				for i := 0; i < len(elems); i++ {
					argNums = append(argNums, fmt.Sprintf("$%d", len(args)+1))
					args = append(args, elems[i])
				}
				where = fmt.Sprintf(`%s "%s" IN (%s)`, where, values[0], strings.Join(argNums, ","))
				continue
			}

			values = strings.SplitN(f, ">", 2)
			if len(values) > 1 {
				where = fmt.Sprintf(`%s "%s" > $%d`, where, values[0], len(args)+1)
				args = append(args, values[1])
				continue
			}

			values = strings.SplitN(f, "<", 2)
			if len(values) > 1 {
				where = fmt.Sprintf(`%s "%s" < $%d`, where, values[0], len(args)+1)
				args = append(args, values[1])
				continue
			}

			values = strings.SplitN(f, "!~", 2)
			if len(values) > 1 {
				where = fmt.Sprintf(`%s "%s" NOT LIKE $%d`, where, values[0], len(args)+1)
				args = append(args, values[1])
				continue
			}

			values = strings.SplitN(f, "~", 2)
			if len(values) > 1 {
				where = fmt.Sprintf(`%s "%s" LIKE $%d`, where, values[0], len(args)+1)
				args = append(args, values[1])
				continue
			}

			needOperator = false
		}
	}

	args = append(args, size, (page-1)*size)

	if sort == "" {
		sort = "id"
	}

	order := "ASC"
	if strings.HasPrefix(sort, "-") {
		order = "DESC"
		sort = sort[1:]
	}

	return where, sort, order, args
}

func SqlDateFormat(dt sql.NullTime) interface{} {
	if dt.Valid {
		return utils.ISODate(dt.Time.UTC())
	} else {
		return nil
	}
}
