package main

// create a kvstore which has tables (keys in the map) and key value pairs within that table

type KVPairs map[string]string

type Store struct {
	tables map[string]KVPairs
}

func NewStore() *Store {
	return &Store{
		tables: make(map[string]KVPairs),
	}
}

func (s *Store) Table(name string) KVPairs {
	table, ok := s.tables[name]
	if !ok {
		table = make(KVPairs)
		s.tables[name] = table
	}
	return table
}

func (kv KVPairs) Get(key string) (any, bool) {
	value, ok := kv[key]
	return value, ok
}

func (kv KVPairs) Set(key string, value string) {
	kv[key] = value
}
