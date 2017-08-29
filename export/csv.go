package export

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type CsvExporter struct {
}

func (c *CsvExporter) makeHeadsStr(t *Table) string {
	var names []string
	var types []string
	var descs []string

	for _, fi := range t.FieldInfos {
		names = append(names, fi.Name)
		types = append(types, fi.Type)
		descs = append(descs, fi.Desc)
	}

	str := strings.Join(names, "\t") + "\n"
	str += strings.Join(types, "\t") + "\n"
	str += strings.Join(descs, "\t") + "\n"

	return str
}

func (c *CsvExporter) makeRowsStr(t *Table) string {
	str := ""
	for _, row := range t.Rows {
		var temp []string
		for _, field := range row.Fields {
			temp = append(temp, field.Value)
		}
		str += strings.Join(temp, "\t") + "\n"
	}

	return str + "\n"
}

func (c *CsvExporter) makeTableStr(t *Table) string {
	str := c.makeHeadsStr(t)
	str += c.makeRowsStr(t)
	return str
}

func (c *CsvExporter) Save(path string, tables []Table) error {
	fullpath := ""
	csvpath := filepath.Join(path, "csv")
	os.MkdirAll(csvpath, 0777)
	for _, t := range tables {
		str := c.makeTableStr(&t)
		fullpath = filepath.Join(csvpath, t.Name+".csv")
		ioutil.WriteFile(fullpath, []byte(str), 0666)
	}
	return nil
}
