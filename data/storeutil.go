package data

import (
	"errors"
	"fmt"
	"github.com/gocraft/dbr/v2"
)

type junction struct {
	table1        string
	table2        string
	junctionTable string
	table1Pk      string
	table2Pk      string
	junctionFk1   string
	junctionFk2   string
}

func (s *Store) selectJunction(db *dbr.Session, lookupId interface{}, j junction) *dbr.SelectStmt {
	if j.table1Pk == "" {
		j.table1Pk = "id"
	}

	if j.table2Pk == "" {
		j.table2Pk = "id"
	}

	// wrapping table names in quotes prevents errors when tables/columns are named after reserved words
	return db.
		Select(quotes(j.table1)+".*").
		From(quotes(j.table1)).
		Join(
			j.junctionTable,
			fmt.Sprintf("%s.%s = %s.%s", quotes(j.table1), j.table1Pk, j.junctionTable, j.junctionFk1),
		).
		Join(
			j.table2,
			fmt.Sprintf("%s.%s = %s.%s", quotes(j.table2), j.table2Pk, j.junctionTable, j.junctionFk2),
		).
		Where(fmt.Sprintf("%s.%s = ?", quotes(j.table2), j.table2Pk), lookupId)
}

func quotes(s string) string {
	return fmt.Sprintf("\"%s\"", s)
}

func (s *Store) create(table string, record interface{}, columns []string) (interface{}, error) {
	var id interface{}

	err := s.db.
		InsertInto(table).
		Columns(columns...).
		Record(record).
		Returning("id").
		Load(&id)

	if err != nil {
		return nil, err
	}

	return id, nil
}

func (s *Store) update(table string, id interface{}, fields []string, updateSets ...set) error {
	setMap := makeSetMap(fields, updateSets...)

	result, err := s.db.
		Update(table).
		SetMap(setMap).
		Where("id = ?", id).
		Exec()

	if err != nil {
		return err
	}

	count, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New(ErrResourceDNE)
	}

	return err
}

func (s *Store) getById(table string, id interface{}, resource interface{}) (interface{}, int, error) {
	count, err := s.db.
		Select("*").
		From(fmt.Sprintf(`"%s"`, table)).
		Where("id = ?", id).
		Load(resource)

	if err != nil {
		return nil, 0, err
	}

	return resource, count, nil
}

func (s *Store) delete(table string, id interface{}) error {
	result, err := s.db.DeleteFrom(table).Where("id = ?", id).Exec()

	if err != nil {
		return err
	}

	count, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New(ErrResourceDNE)
	}

	return err
}

func includes(strings []string, val string) bool {
	for _, v := range strings {
		if v == val {
			return true
		}
	}

	return false
}

type set struct {
	Field string
	Col   string
	Val   interface{}
}

func makeSetMap(fields []string, sets ...set) map[string]interface{} {
	setMap := map[string]interface{}{}

	for _, set := range sets {
		if includes(fields, set.Field) {
			setMap[set.Col] = set.Val
		}
	}

	return setMap
}
