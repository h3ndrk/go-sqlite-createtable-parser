# go-sqlite-createtable-parser

This repository contains a Golang binding for [sqlite-createtable-parser](https://github.com/marcobambini/sqlite-createtable-parser) (MIT license). The underlying sqlite-createtable-parser is a parser for SQLite `CREATE TABLE` statements (see [syntax](https://www.sqlite.org/lang_createtable.html)).

From [sqlite-createtable-parser/README.md](https://github.com/marcobambini/sqlite-createtable-parser/blob/master/README.md):

> ## Motivation
> 
> [SQLite](https://www.sqlite.org/) is a very powerful software but it lacks an easy way to extract complete information about tables and columns constraints. This drawback in addition to the lack of full ALTER TABLE support makes alterring a table a very hard task. The built-in sqlite pragmas provide incomplete information and a manual parsing is required in order to extract all the metadata from a table.

## Installation

A C99 compiler is required.

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

## Implementation status

As of the creation of this repository, [sqlite-createtable-parser](https://github.com/marcobambini/sqlite-createtable-parser) does **not support** parsing `CHECK` **table** constraints. For example

```sql
CREATE TABLE a (b INTEGER, CHECK (b >= 42));
```

will return a `SQL3ERROR_SYNTAX` error. However, `CHECK` **column** constraints are on the other hand **supported** (notice the removed comma):

```sql
CREATE TABLE a (b INTEGER CHECK (b >= 42));
```

All other syntax of `CREATE TABLE` is **supported**. This repository adds many tests to ensure that the parser works.

## License

MIT, see `LICENSE` file.
