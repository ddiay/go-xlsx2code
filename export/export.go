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
	Name              string
	Type              string
	Desc              string
	Key               string
	Value             string
	RefClassFieldName string
	Index             bool
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

type Exporter interface {
	Save(path string, table *Table) error
}

var indexTypeMap map[string]string

func AddIndexType(key string, typeStr string) {
	if indexTypeMap == nil {
		indexTypeMap = make(map[string]string)
	}
	indexTypeMap[key] = typeStr
}

func FindIndexType(str string) string {
	v, ok := indexTypeMap[str]
	if ok {
		return v
	}
	return ""
}
