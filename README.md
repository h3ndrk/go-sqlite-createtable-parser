# go-sqlite-createtable-parser

This repository contains a Golang binding for https://github.com/marcobambini/sqlite-createtable-parser. The underlying sqlite-createtable-parser is a parser for SQLite `CREATE TABLE` statements (see https://www.sqlite.org/lang_createtable.html).

## Installation

```bash
go get github.com/h3ndrk/go-sqlite-createtable-parser
```

## How to use

```go
stmt, err := parse.FromString("CREATE TABLE main.tbl (a INTEGER, b TEXT, FOREIGN KEY(a) REFERENCES othertbl(id))")
if err != nil {
    // handle error
    panic(err)
}

fmt.Println(stmt)
// &parse.Table{
//     Schema:           &"main",
//     Name:             &"tbl",
//     Temporary:        false,
//     IfNotExists:      false,
//     WithoutRowid:     false,
//     Columns:          []parse.Column{
//         parse.Column{
//             Name:               &"a",
//             Type:               &"INTEGER",
//             Length:             nil,
//             ConstraintName:     nil,
//             PrimaryKey:         false,
//             Autoincrement:      false,
//             NotNull:            false,
//             Unique:             false,
//             PrimaryKeyOrder:    parse.OrderNone,
//             PrimaryKeyConflict: parse.ConflictNone,
//             NotNullConflict:    parse.ConflictNone,
//             UniqueConflict:     parse.ConflictNone,
//             Check:              nil,
//             Default:            nil,
//             CollateName:        nil,
//             ForeignKey:         nil,
//         },
//         parse.Column{
//             Name:               &"b",
//             Type:               &"TEXT",
//             Length:             nil,
//             ConstraintName:     nil,
//             PrimaryKey:         false,
//             Autoincrement:      false,
//             NotNull:            false,
//             Unique:             false,
//             PrimaryKeyOrder:    parse.OrderNone,
//             PrimaryKeyConflict: parse.ConflictNone,
//             NotNullConflict:    parse.ConflictNone,
//             UniqueConflict:     parse.ConflictNone,
//             Check:              nil,
//             Default:            nil,
//             CollateName:        nil,
//             ForeignKey:         nil,
//         },
//     },
//     TableConstraints: []parse.TableConstraint{
//         parse.TableConstraint{
//             Name:              nil,
//             Type:              parse.TableConstraintTypeForeignKey,
//             IndexedColumns:    []parse.IndexedColumn{},
//             ConflictClause:    parse.ConflictNone,
//             Check:             nil,
//             ForeignKeyColumns: []string{
//                 "a",
//             },
//             ForeignKey:        &parse.ForeignKey{
//                 Table:      &"othertbl",
//                 Columns:    []string{
//                     "id",
//                 },
//                 OnDelete:   parse.ForeignKeyActionNone,
//                 OnUpdate:   parse.ForeignKeyActionNone,
//                 Match:      nil,
//                 Deferrable: parse.ForeignKeyDeferrableTypeNone,
//             },
//         },
//     },
// }
```

## `CHECK` table constraints not supported

As of the creation of this repository, [sqlite-createtable-parser](https://github.com/marcobambini/sqlite-createtable-parser) does not support parsing `CHECK` table constraints. For example:

```sql
CREATE TABLE a (b INTEGER, CHECK (b >= 42));
```

will return a `SQL3ERROR_SYNTAX` error.


## License

MIT, see `LICENSE` file.
