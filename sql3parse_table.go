// +build cgo

package parse

/*
#include <stdint.h>
#include "sqlite-createtable-parser/sql3parse_table.h"
#include "sqlite-createtable-parser/sql3parse_table.c"
*/
import "C"
import (
	"unsafe"

	"github.com/pkg/errors"
)

type Table struct {
	Schema           *string
	Name             *string
	Temporary        bool
	IfNotExists      bool
	WithoutRowid     bool
	Columns          []Column
	TableConstraints []TableConstraint
}

type parseError int

const (
	parseErrorNone parseError = iota
	parseErrorMemory
	parseErrorSyntax
	parseErrorUnsupportedSQL
)

func FromString(sql string) (*Table, error) {
	sqlString := C.CString(sql)
	defer C.free(unsafe.Pointer(sqlString))

	sqlStringLength := C.ulong(len(sql))

	var err C.sql3error_code

	tablePtr := C.sql3parse_table(sqlString, sqlStringLength, &err)
	defer C.sql3table_free(tablePtr)

	switch err {
	case C.SQL3ERROR_MEMORY:
		return nil, errors.New("sqlite-createtable-parser: sql3parse_table() reported SQL3ERROR_MEMORY")
	case C.SQL3ERROR_SYNTAX:
		return nil, errors.New("sqlite-createtable-parser: sql3parse_table() reported SQL3ERROR_SYNTAX")
	case C.SQL3ERROR_UNSUPPORTEDSQL:
		return nil, errors.New("sqlite-createtable-parser: sql3parse_table() reported SQL3ERROR_UNSUPPORTEDSQL")
	}

	table := &Table{
		Schema:       sql3stringToGo(C.sql3table_schema(tablePtr)),
		Name:         sql3stringToGo(C.sql3table_name(tablePtr)),
		Temporary:    bool(C.sql3table_is_temporary(tablePtr)),
		IfNotExists:  bool(C.sql3table_is_ifnotexists(tablePtr)),
		WithoutRowid: bool(C.sql3table_is_withoutrowid(tablePtr)),
	}

	for i := C.ulong(0); i < C.sql3table_num_columns(tablePtr); i++ {
		columnPtr := C.sql3table_get_column(tablePtr, i)
		if columnPtr == nil {
			return nil, errors.Errorf("sqlite-createtable-parser: sql3table_get_column() returned NULL (i: %d)", i)
		}
		column, err := fromColumnPtr(columnPtr)
		if err != nil {
			return nil, err
		}
		if column == nil {
			return nil, errors.New("unexpected column == nil (table column)")
		}
		table.Columns = append(table.Columns, *column)
	}

	for i := C.ulong(0); i < C.sql3table_num_constraints(tablePtr); i++ {
		tableConstraintPtr := C.sql3table_get_constraint(tablePtr, i)
		if tableConstraintPtr == nil {
			return nil, errors.Errorf("sqlite-createtable-parser: sql3table_get_constraint() returned NULL (i: %d)", i)
		}
		tableConstraint, err := fromTableConstraintPtr(tableConstraintPtr)
		if err != nil {
			return nil, err
		}
		if tableConstraint == nil {
			return nil, errors.New("unexpected tableConstraint == nil (table constraint)")
		}
		table.TableConstraints = append(table.TableConstraints, *tableConstraint)
	}

	return table, nil
}

func sql3stringToGo(ptr *C.sql3string) *string {
	if ptr == nil || ptr.ptr == nil {
		return nil
	}
	converted := C.GoStringN(ptr.ptr, C.int(ptr.length))
	return &converted
}

type OrderClause int

const (
	OrderNone OrderClause = iota
	OrderAsc
	OrderDesc
)

type ConflictClause int

const (
	ConflictNone ConflictClause = iota
	ConflictRollback
	ConflictAbort
	ConflictFail
	ConflictIgnore
	ConflictReplace
)

type Column struct {
	Name               *string
	Type               *string
	Length             *string
	ConstraintName     *string
	PrimaryKey         bool
	Autoincrement      bool
	NotNull            bool
	Unique             bool
	PrimaryKeyOrder    OrderClause
	PrimaryKeyConflict ConflictClause
	NotNullConflict    ConflictClause
	UniqueConflict     ConflictClause
	Check              *string
	Default            *string
	CollateName        *string
	ForeignKey         *ForeignKey
}

func fromColumnPtr(ptr *C.sql3column) (*Column, error) {
	column := &Column{
		Name:               sql3stringToGo(C.sql3column_name(ptr)),
		Type:               sql3stringToGo(C.sql3column_type(ptr)),
		Length:             sql3stringToGo(C.sql3column_length(ptr)),
		ConstraintName:     sql3stringToGo(C.sql3column_constraint_name(ptr)),
		PrimaryKey:         bool(C.sql3column_is_primarykey(ptr)),
		Autoincrement:      bool(C.sql3column_is_autoincrement(ptr)),
		NotNull:            bool(C.sql3column_is_notnull(ptr)),
		Unique:             bool(C.sql3column_is_unique(ptr)),
		PrimaryKeyOrder:    OrderClause(C.sql3column_pk_order(ptr)),
		PrimaryKeyConflict: ConflictClause(C.sql3column_pk_conflictclause(ptr)),
		NotNullConflict:    ConflictClause(C.sql3column_notnull_conflictclause(ptr)),
		UniqueConflict:     ConflictClause(C.sql3column_unique_conflictclause(ptr)),
		Check:              sql3stringToGo(C.sql3column_name(ptr)),
		Default:            sql3stringToGo(C.sql3column_name(ptr)),
		CollateName:        sql3stringToGo(C.sql3column_name(ptr)),
	}

	if foreignKeyPtr := C.sql3column_foreignkey_clause(ptr); foreignKeyPtr != nil {
		foreignKey, err := fromForeignKeyPtr(foreignKeyPtr)
		if err != nil {
			return nil, err
		}
		if foreignKey == nil {
			return nil, errors.New("unexpected foreignKey == nil (column foreign key)")
		}
		column.ForeignKey = foreignKey
	}

	return column, nil
}

type IndexedColumn struct {
	Name    *string
	Collate *string
	Order   OrderClause
}

func fromIndexedColumnPtr(ptr *C.sql3idxcolumn) (*IndexedColumn, error) {
	indexedColumn := &IndexedColumn{
		Name:    sql3stringToGo(C.sql3idxcolumn_name(ptr)),
		Collate: sql3stringToGo(C.sql3idxcolumn_collate(ptr)),
		Order:   OrderClause(C.sql3idxcolumn_order(ptr)),
	}

	return indexedColumn, nil
}

type ForeignKeyAction int

const (
	ForeignKeyActionNone ForeignKeyAction = iota
	ForeignKeyActionSetNull
	ForeignKeyActionSetDefault
	ForeignKeyActionCascade
	ForeignKeyActionRestrict
	ForeignKeyActionNoAction
)

type ForeignKeyDeferrableType int

const (
	ForeignKeyDeferrableTypeNone ForeignKeyDeferrableType = iota
	ForeignKeyDeferrableTypeDeferrable
	ForeignKeyDeferrableTypeDeferrableInitiallyDeferred
	ForeignKeyDeferrableTypeDeferrableInitiallyImmediate
	ForeignKeyDeferrableTypeNotDeferrable
	ForeignKeyDeferrableTypeNotDeferrableInitiallyDeferred
	ForeignKeyDeferrableTypeNotDeferrableInitiallyImmediate
)

type ForeignKey struct {
	Table      *string
	Columns    []string
	OnDelete   ForeignKeyAction
	OnUpdate   ForeignKeyAction
	Match      *string
	Deferrable ForeignKeyDeferrableType
}

func fromForeignKeyPtr(ptr *C.sql3foreignkey) (*ForeignKey, error) {
	foreignKey := &ForeignKey{
		Table:      sql3stringToGo(C.sql3foreignkey_table(ptr)),
		OnDelete:   ForeignKeyAction(C.sql3foreignkey_ondelete_action(ptr)),
		OnUpdate:   ForeignKeyAction(C.sql3foreignkey_onupdate_action(ptr)),
		Match:      sql3stringToGo(C.sql3foreignkey_match(ptr)),
		Deferrable: ForeignKeyDeferrableType(C.sql3foreignkey_deferrable(ptr)),
	}

	for i := C.ulong(0); i < C.sql3foreignkey_num_columns(ptr); i++ {
		column := sql3stringToGo(C.sql3foreignkey_get_column(ptr, i))
		if column == nil {
			return nil, errors.Errorf("sqlite-createtable-parser: sql3foreignkey_get_column() returned NULL (i: %d)", i)
		}

		foreignKey.Columns = append(foreignKey.Columns, *column)
	}

	return foreignKey, nil
}

type TableConstraintType int

const (
	TableConstraintTypePrimaryKey TableConstraintType = iota
	TableConstraintTypeUnique
	TableConstraintTypeCheck
	TableConstraintTypeForeignKey
)

type TableConstraint struct {
	Name              *string
	Type              TableConstraintType
	IndexedColumns    []IndexedColumn
	ConflictClause    ConflictClause
	Check             *string
	ForeignKeyColumns []string
	ForeignKey        *ForeignKey
}

func fromTableConstraintPtr(ptr *C.sql3tableconstraint) (*TableConstraint, error) {
	tableConstraint := &TableConstraint{
		Name:           sql3stringToGo(C.sql3table_constraint_name(ptr)),
		Type:           TableConstraintType(C.sql3table_constraint_type(ptr)),
		ConflictClause: ConflictClause(C.sql3table_constraint_conflict_clause(ptr)),
		Check:          sql3stringToGo(C.sql3table_constraint_check_expr(ptr)),
	}

	for i := C.ulong(0); i < C.sql3table_constraint_num_idxcolumns(ptr); i++ {
		indexedColumnPtr := C.sql3table_constraint_get_idxcolumn(ptr, i)
		if indexedColumnPtr == nil {
			return nil, errors.Errorf("sqlite-createtable-parser: sql3foreignkey_get_column() returned NULL (i: %d)", i)
		}
		indexedColumn, err := fromIndexedColumnPtr(indexedColumnPtr)
		if err != nil {
			return nil, err
		}
		if indexedColumn == nil {
			return nil, errors.New("unexpected indexedColumn == nil (table constraint indexed column)")
		}

		tableConstraint.IndexedColumns = append(tableConstraint.IndexedColumns, *indexedColumn)
	}

	for i := C.ulong(0); i < C.sql3table_constraint_num_fkcolumns(ptr); i++ {
		column := sql3stringToGo(C.sql3table_constraint_get_fkcolumn(ptr, i))
		if column == nil {
			return nil, errors.Errorf("sqlite-createtable-parser: sql3table_constraint_get_fkcolumn() returned NULL (i: %d)", i)
		}

		tableConstraint.ForeignKeyColumns = append(tableConstraint.ForeignKeyColumns, *column)
	}

	if foreignKeyPtr := C.sql3table_constraint_foreignkey_clause(ptr); foreignKeyPtr != nil {
		foreignKey, err := fromForeignKeyPtr(foreignKeyPtr)
		if err != nil {
			return nil, err
		}
		if foreignKey == nil {
			return nil, errors.New("unexpected foreignKey == nil (table constraint foreign key)")
		}
		tableConstraint.ForeignKey = foreignKey
	}

	return tableConstraint, nil
}
