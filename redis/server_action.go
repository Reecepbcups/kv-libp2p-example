package redis

import (
	"fmt"
	"strings"
)

func HandleMsg(msg string, store *Store) []byte {
	// TODO: add protocol version here? (smart but not necessary for example demo)
	// formats:
	// * set;table;key,value
	// * get;table;key
	// * delete;table;key
	// * keys;table
	// * values;table
	// * all

	msg = strings.TrimSuffix(msg, "\n")

	args := strings.Split(msg, ";")
	fmt.Println("(Server) Debugging HandleMsg", args)

	action := args[0]

	var table string
	if len(args) > 1 {
		table = args[1]
	}

	switch action {
	case "set":
		tuple := strings.Split(args[2], ",")
		key, value := tuple[0], tuple[1]

		store.Table(table).Set(key, value)
		return []byte(`{"status":"OK"}`)
	case "get":
		key := args[2]
		res, ok := store.Table(table).Get(key)
		if !ok {
			return []byte(fmt.Sprintf(`{"error":"Key '%s' not found in table '%s'"}`, key, table))
		}
		return []byte(fmt.Sprintf(`{"result":"%s","key":"%s"}`, res, key))
	case "keys":
		keys := store.Table(table).Keys()
		return []byte(fmt.Sprintf(`{"result":["%s"]}`, strings.Join(keys, `", "`)))
	case "values":
		values := store.Table(table).Values()
		return []byte(fmt.Sprintf(`{"result":["%s"]}`, strings.Join(values, `", "`)))
	case "all":
		return []byte(fmt.Sprintf(`{"result":%s}`, store.tables.String()))
	case "delete":
		key := args[2]
		store.Table(table).Delete(key)
		return []byte(`{"status":"OK"}`)

	default:
		return []byte(fmt.Sprintf(`{"error":"Invalid action '%s'"}`, action))
	}
}
