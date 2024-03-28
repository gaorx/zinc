package zinc

type Table interface {
	Name() string
	Columns() []TableColumn
}

type TableColumn interface {
	Name() string
}
