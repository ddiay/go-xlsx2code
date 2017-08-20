package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/tealeg/xlsx"
)

var xlsxPath = flag.String("f", "", "Path to an XLSX file")
var outputPath = flag.String("t", "", "Path to output")
var delimiter = flag.String("d", "\t", "Delimiter to use between fields")

func generateCSVFromXLSXFile(excelFileName string) error {
	xlFile, error := xlsx.OpenFile(excelFileName)
	if error != nil {
		return error
	}
	sheetLen := len(xlFile.Sheets)
	switch {
	case sheetLen == 0:
		return errors.New("This XLSX file contains no sheets.")
	}
	for _, sheet := range xlFile.Sheets {
		csvContent := ""
		for _, row := range sheet.Rows {
			var vals []string
			if row != nil {
				for _, cell := range row.Cells {
					str, err := cell.FormattedValue()
					if err != nil {
						vals = append(vals, err.Error())
					}
					vals = append(vals, fmt.Sprintf("%q", str))
				}
				csvContent += strings.Join(vals, *delimiter) + "\n"
			}
		}
		path := filepath.Join(*outputPath, sheet.Name+".txt")
		err := ioutil.WriteFile(path, []byte(csvContent), 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if len(os.Args) < 3 {
		flag.PrintDefaults()
		return
	}
	flag.Parse()
	if err := generateCSVFromXLSXFile(*xlsxPath); err != nil {
		fmt.Println(err)
	}
}
