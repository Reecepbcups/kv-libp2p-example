package redis

// create a kvstore which has tables (keys in the map) and key value pairs within that table

type KVPairs map[string]string

type Store struct {
	dbName string
	tables map[string]KVPairs
}

func NewStore(name string) *Store {
	return &Store{
		dbName: name,
		tables: make(map[string]KVPairs),
	}
}

func (s *Store) DBName() string {
	return s.dbName
}

// Table returns a table, which is created if it does not already exist.
func (s *Store) Table(name string) KVPairs {
	table, ok := s.tables[name]
	if !ok {
		s.tables[name] = make(KVPairs)
		return s.tables[name]
	}
	return table
}

func (kv KVPairs) Get(key string) (string, bool) {
	value, ok := kv[key]
	return value, ok
}

func (kv KVPairs) Set(key string, value string) {
	kv[key] = value
}