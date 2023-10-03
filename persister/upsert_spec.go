package persister

import (
	"fmt"
)

type UpsertSpec struct {
	Table  string
	Fields Fields
	OnConf Fields
	UpdSet Fields
}

func NewUpsertSpec(tbl string, flds ...string) *UpsertSpec {
	if len(flds) == 0 {
		flds = []string{"id"}
	}
	for i := 0; i < 3; i++ {
		if i < len(flds) {
			continue
		}
		flds = append(flds, flds[i-1])
	}
	return &UpsertSpec{
		Table:  tbl,
		Fields: NewFieldsFromString(flds[0]),
		OnConf: NewFieldsFromString(flds[1]),
		UpdSet: NewFieldsFromString(flds[2]),
	}
}

func (us *UpsertSpec) UpsertSQL() string {
	sql := `INSERT INTO %s (%s) VALUES (%s) ON CONFLICT (%s) DO UPDATE SET %s RETURNING *;`
	return fmt.Sprintf(sql,
		us.Table,
		us.Fields.Names(),
		us.Fields.PlaceHolders(),
		us.OnConf.Names(),
		us.UpdSet.DoUpdateSet(),
	)
}
