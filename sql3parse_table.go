// +build cgo

package parse

/*
#include <stdint.h>
#include "sqlite-createtable-parser/sql3parse_table.h"
#include "sqlite-createtable-parser/sql3parse_table.c"
*/
import "C"
import (
	"fmt"
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
	fmt.Println(sql)
	sqlString := C.CString(sql)
	defer C.free(unsafe.Pointer(sqlString))

	sqlStringLength := C.ulong(len(sql))

	var err C.sql3error_code

	tablePointer := C.sql3parse_table(sqlString, sqlStringLength, &err)
	defer C.sql3table_free(tablePointer)

	fmt.Println(tablePointer, err)
	switch err {
	case C.SQL3ERROR_MEMORY:
		return nil, errors.New("sqlite-createtable-parser: sql3parse_table() reported SQL3ERROR_MEMORY")
	case C.SQL3ERROR_SYNTAX:
		return nil, errors.New("sqlite-createtable-parser: sql3parse_table() reported SQL3ERROR_SYNTAX")
	case C.SQL3ERROR_UNSUPPORTEDSQL:
		return nil, errors.New("sqlite-createtable-parser: sql3parse_table() reported SQL3ERROR_UNSUPPORTEDSQL")
	}
	
	// table := &Table{}
	
	schema := sql3stringToGo(C.sql3table_schema(tablePointer))
	name := sql3stringToGo(C.sql3table_name(tablePointer))
	temporary := bool(C.sql3table_is_temporary(tablePointer))
	fmt.Printf("sql3table_schema(): %T %+v\n", schema, schema)
	fmt.Printf("sql3table_name(): %T %+v\n", name, name)
	fmt.Printf("sql3table_is_temporary(): %T %+v\n", temporary, temporary)

	return nil, nil
}

func sql3stringToGo(input *C.sql3string) *string {
	if input == nil || input.ptr == nil {
		return nil
	}
	converted := C.GoStringN(input.ptr, C.int(input.length))
	fmt.Printf("converted: %s\n", converted)
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
	Name               string
	Type               string
	Length             string
	ConstraintName     string
	PrimaryKey         bool
	Autoincrement      bool
	NotNull            bool
	Unique             bool
	PrimaryKeyOrder    OrderClause
	PrimaryKeyConflict ConflictClause
	NotNullConflict    ConflictClause
	UniqueConflict     ConflictClause
	Check              string
	Default            string
	CollateName        string
	ForeignKey         ForeignKey
}

type IndexedColumn struct {
	Name    string
	Collate string
	Order   OrderClause
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
	Table      string
	Columns    []string
	OnDelete   ForeignKeyAction
	OnUpdate   ForeignKeyAction
	Match      string
	Deferrable ForeignKeyDeferrableType
}

type TableConstraintType int

const (
	TableConstraintTypePrimaryKey TableConstraintType = iota
	TableConstraintTypeUnique
	TableConstraintTypeCheck
	TableConstraintTypeForeignKey
)

type TableConstraint struct {
	Name              string
	Type              TableConstraintType
	IndexedColumns    []IndexedColumn
	ConflictClause    ConflictClause
	Check             string
	ForeignKeyColumns []string
	ForeignKey        ForeignKey
}
