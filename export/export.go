package export

const (
	ExportTypeUnknow = 0
	ExportTypeNumber = 1
	ExportTypeFloat  = 2
	ExportTypeString = 3
	ExportTypeTable  = 4
	ExportTypeList   = 5
	ExportTypeMap    = 6
)

type FieldInfo struct {
	Name  string
	Type  string
	Desc  string
	Key   string
	Value string
	Index bool
}

type Field struct {
	Info  *FieldInfo
	Value string
}

type Row struct {
	Fields []Field
}

type Table struct {
	Name       string
	FieldInfos []FieldInfo
	Rows       []Row
}

type Database struct {
	Tables []Table
}

type Exporter interface {
	Save(path string, table *Table) error
}

func (d *Database) FindIndexType(tableName string, fieldName string) string {

}
