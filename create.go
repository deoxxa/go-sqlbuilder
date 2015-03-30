package sqlbuilder

// CreateIndexStatement represents a "CREATE INDEX" statement.
type CreateIndexStatement struct {
	table       Table
	columns     []Column
	name        string
	ifNotExists bool
}

// CreateTableStatement represents a "CREATE TABLE" statement.
type CreateTableStatement struct {
	table       Table
	ifNotExists bool
}

// CreateTable returns new "CREATE TABLE" statement. The table is Table object to create.
func CreateTable(table Table) *CreateTableStatement {
	return &CreateTableStatement{
		table: table,
	}
}

// IfNotExists sets "IF NOT EXISTS" clause.
func (b *CreateTableStatement) IfNotExists() *CreateTableStatement {
	b.ifNotExists = true
	return b
}

// CreateIndex returns new "CREATE INDEX" statement. The table is Table object to create index.
func CreateIndex(table Table) *CreateIndexStatement {
	return &CreateIndexStatement{
		table: table,
	}
}

// ToSql generates query string, placeholder arguments, and error.
func (b *CreateTableStatement) ToSql() (query string, args []interface{}, err error) {
	bldr := newBuilder()
	defer func() {
		query, args, err = bldr.Query(), bldr.Args(), bldr.Err()
	}()

	if b.table == nil {
		bldr.SetError(newError("table is nil."))
		return
	} else if _, ok := b.table.(*table); !ok {
		bldr.SetError(newError("This table can not create."))
		return
	}

	bldr.Append("CREATE TABLE ")
	if b.ifNotExists {
		bldr.Append("IF NOT EXISTS ")
	}
	bldr.AppendItem(b.table)

	if len(b.table.Columns()) != 0 {
		bldr.Append(" ( ")
		bldr.AppendItem(createTableColumnList(b.table.Columns()))
		bldr.Append(" )")
	} else {
		bldr.SetError(newError("CreateTableStatement needs one or more columns."))
		return
	}

	// table option
	if tabopt, args, err := dialect().TableOptionToString(b.table.Option()); err == nil {
		if len(tabopt) != 0 {
			bldr.Append(" " + tabopt)
			bldr.AppendValuesRaw(args...)
		}
	} else {
		bldr.SetError(err)
	}

	return
}

// IfNotExists sets "IF NOT EXISTS" clause.
func (b *CreateIndexStatement) IfNotExists() *CreateIndexStatement {
	b.ifNotExists = true
	return b
}

// IfNotExists sets "IF NOT EXISTS" clause. If not set this, returns error on ToSql().
func (b *CreateIndexStatement) Columns(columns ...Column) *CreateIndexStatement {
	b.columns = columns
	return b
}

// Name sets name for index.
// If not set this, auto generated name will be used.
func (b *CreateIndexStatement) Name(name string) *CreateIndexStatement {
	b.name = name
	return b
}

// ToSql generates query string, placeholder arguments, and returns err on errors.
func (b *CreateIndexStatement) ToSql() (query string, args []interface{}, err error) {
	bldr := newBuilder()
	defer func() {
		query, args, err = bldr.Query(), bldr.Args(), bldr.Err()
	}()

	bldr.Append("CREATE INDEX ")
	if b.ifNotExists {
		bldr.Append("IF NOT EXISTS ")
	}

	if len(b.name) != 0 {
		bldr.Append(dialect().QuoteField(b.name))
	} else {
		bldr.SetError(newError("name was not setted."))
		return
	}

	bldr.Append(" ON ")
	if b.table != nil {
		bldr.AppendItem(b.table)
	} else {
		bldr.SetError(newError("table is nil."))
		return
	}

	if len(b.columns) != 0 {
		bldr.Append(" ( ")
		bldr.AppendItem(createIndexColumnList(b.columns))
		bldr.Append(" )")
	} else {
		bldr.SetError(newError("columns was not setted."))
		return
	}

	return
}

type createTableColumnList []Column

func (m createTableColumnList) serialize(bldr *builder) {
	first := true
	for _, column := range m {
		if first {
			first = false
		} else {
			bldr.Append(", ")
		}
		cc := column.config()

		// Column name
		bldr.AppendItem(cc)
		bldr.Append(" ")

		// SQL data name
		str, err := dialect().ColumnTypeToString(cc)
		if err != nil {
			bldr.SetError(err)
		}
		bldr.Append(str)

		str, args, err := dialect().ColumnOptionToString(cc.Option())
		if err != nil {
			bldr.SetError(err)
		}
		if len(str) != 0 {
			bldr.Append(" " + str)
			bldr.AppendValuesRaw(args...)
		}
	}
}

type createIndexColumnList []Column

func (m createIndexColumnList) serialize(bldr *builder) {
	first := true
	for _, column := range m {
		if first {
			first = false
		} else {
			bldr.Append(", ")
		}
		cc := column.config()
		bldr.AppendItem(cc)
	}
}
