// +build cgo

package parse

/*
#include <stdint.h>
#include "sqlite-createtable-parser/sql3parse_table.h"
#include "sqlite-createtable-parser/sql3parse_table.c"
*/
import "C"

type Table struct {
	Schema           string
	Name             string
	Temporary        bool
	IfNotExists      bool
	WithoutRowid     bool
	Columns          []Column
	TableConstraints []TableConstraint
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
