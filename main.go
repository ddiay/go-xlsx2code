package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ddiay/go-xlsx2code/export"
	"github.com/tealeg/xlsx"
)

var xlsxPath = flag.String("f", "", "Path to an XLSX file")
var outputPath = flag.String("t", "", "Path to output")
var delimiter = flag.String("d", "\t", "Delimiter to use between fields")

func generateCSVFromXLSXFile(excelFileName string, exporters []export.Exporter) error {
	xlFile, error := xlsx.OpenFile(excelFileName)
	if error != nil {
		return error
	}
	sheetLen := len(xlFile.Sheets)
	switch {
	case sheetLen == 0:
		return errors.New("This XLSX file contains no sheets.")
	}

	var tables []export.Table
	for _, sheet := range xlFile.Sheets {
		// csvContent := ""
		table := export.Table{
			Name: sheet.Name,
		}

		for i, row := range sheet.Rows {
			// var vals []string
			if row != nil {
				r := export.Row{}
				for j, cell := range row.Cells {
					if i == 0 {
						fieldInfo := export.FieldInfo{
							Name: cell.String(),
						}
						table.FieldInfos = append(table.FieldInfos, fieldInfo)
					} else if i == 1 {
						ft := cell.String()
						table.FieldInfos[j].Type = ft
						if strings.HasPrefix(ft, "#") {
							str := strings.TrimLeft(ft, "#")
							table.FieldInfos[j].Type = str
							table.FieldInfos[j].Index = true
						} else if strings.HasPrefix(ft, "<") && strings.HasSuffix(ft, ">") {
							str := strings.TrimPrefix(ft, "<")
							str = strings.TrimSuffix(ft, ">")
							blocks := strings.Split(str, ",")
							switch len(blocks) {
							case 1:
								table.FieldInfos[j].Type = "list"
								table.FieldInfos[j].Value = blocks[0]
							case 2:
								table.FieldInfos[j].Type = "map"
								table.FieldInfos[j].Key = blocks[0]
								table.FieldInfos[j].Key = blocks[1]
							}
						}
					} else if i == 2 {
						table.FieldInfos[j].Desc = cell.String()
					} else {
						field := export.Field{
							Info:  &table.FieldInfos[j],
							Value: cell.String(),
						}
						r.Fields = append(r.Fields, field)
					}
				}
				if i > 2 {
					table.Rows = append(table.Rows, r)
				}
				// csvContent += strings.Join(vals, *delimiter) + "\n"
			}
		}
		// path := filepath.Join(*outputPath, sheet.Name+".txt")
		// err := ioutil.WriteFile(path, []byte(csvContent), 0666)
		// if err != nil {
		// 	return err
		// }
		tables = append(tables, table)
		for _, exporter := range exporters {
			exporter.Save(*outputPath, &table)
		}
		break
	}
	fmt.Println(tables)
	return nil
}

func main() {
	flag.Parse()
	if len(os.Args) < 3 {
		flag.PrintDefaults()
		return
	}
	flag.Parse()
	var exporters []export.Exporter

	csExporter := &export.CSharpExporter{}
	exporters = append(exporters, csExporter)
	if err := generateCSVFromXLSXFile(*xlsxPath, exporters); err != nil {
		fmt.Println(err)
	}
}
