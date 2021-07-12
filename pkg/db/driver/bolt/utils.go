package bolt

import (
	bh "github.com/timshannon/bolthold"
)

func queryFromFilters(filters map[string]interface{}) *bh.Query {
	var query *bh.Query

	for key, value := range filters {
		if query == nil {
			query = bh.Where(key).Eq(value).Index(key)
		} else {
			query = query.And(key).Eq(value).Index(key)
		}
	}

	return query
}
