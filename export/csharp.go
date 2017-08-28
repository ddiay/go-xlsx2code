package export

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type CSharpExporter struct {
}

func (c *CSharpExporter) makeFieldStr(fi *FieldInfo) string {
	str := fmt.Sprintf("\t\t//%s\n", fi.Desc)
	switch fi.Type {
	case "number":
		str += fmt.Sprintf("\t\tpublic int %s;\n", fi.Name)
	case "float":
		str += fmt.Sprintf("\t\tpublic double %s;\n", fi.Name)
	case "string":
		str += fmt.Sprintf("\t\tpublic string %s;\n", fi.Name)
	case "list":
		str += fmt.Sprintf("\t\tpublic List<%s> %s = new List<%s>();\n", fi.Value, fi.Name, fi.Value)
	case "map":
		str += fmt.Sprintf("\t\tpublic Dictionary<%s, %s> %s = new Dictionary<%s, %s>();\n", fi.Key, fi.Value, fi.Name, fi.Key, fi.Value)
	}
	return str
}

func (c *CSharpExporter) makeFieldsStr(t *Table) string {
	str := ""
	for _, fi := range t.FieldInfos {
		str += c.makeFieldStr(&fi)
	}
	str += "\n"

	for _, fi := range t.FieldInfos {
		if fi.Index {
			switch fi.Type {
			case "number":
				str += fmt.Sprintf("\t\tpublic static Dictionary<int, %s> %sMap = new Dictionary<int, %s>();\n", t.Name, fi.Name, t.Name)
			case "string":
				str += fmt.Sprintf("\t\tpublic static Dictionary<string, %s> %sMap = new Dictionary<string, %s>();\n", t.Name, fi.Name, t.Name)
			}
		}
	}
	str += fmt.Sprintf("\t\tpublic static List<%s> Rows = new List<%s>();", t.Name, t.Name)
	return str + "\n"
}

func (c *CSharpExporter) makeLoadMethodStr(t *Table) string {
	str := "\t\tpublic static void Load(IDatatable dt)\n"
	str += "\t\t{\n"
	str += "\t\t\tstring tbname = dt.GetName();\n"
	str += "\t\t\tint numRows = dt.GetRowCount();\n"
	str += "\t\t\tfor (int i = 0; i < numRows; ++i)\n"
	str += "\t\t\t{\n"
	str += fmt.Sprintf("\t\t\t\t%s o = new %s();\n", t.Name, t.Name)

	for i, fi := range t.FieldInfos {
		switch fi.Type {
		case "number":
			str += fmt.Sprintf("\t\t\t\to.%s = dt.GetInt(%d)\n", fi.Name, i)
		case "string":
			str += fmt.Sprintf("\t\t\t\to.%s = dt.GetString(%d)\n", fi.Name, i)
		case "float":
			str += fmt.Sprintf("\t\t\t\to.%s = dt.GetDouble(%d)\n", fi.Name, i)
		}
	}

	str += "\t\t\t\tRows.Add(o);\n"

	for _, fi := range t.FieldInfos {
		if fi.Index {
			str += fmt.Sprintf("\t\t\t\t%sMap.Add(o.%s, o);\n", fi.Name, fi.Name)
		}
	}

	str += "\t\t\t}\n"
	str += "\t\t}\n"

	return str + "\n"
}

/*
	public static void Map(IDataTable dt)
	{
		string tbname = dt.GetName();
		int numRows = dt.GetRowCount();
		for (int i = 0; i < numRows; ++i)
		{
			Unit unit = Rows[i];

			int CClassId = dt.GetInt(i, 2);
			unit.CClass = CharacterClass.IdMap[CClassId];

			Dictionary<string, int> SkinMapFields = dt.GetStringIntDict();
			foreach (KeyValuePair<string, int> kv in SkinMapFields)
			{
				Skin skin = Skin.IdMap[kv.Value];
				unit.SkinMap.Add(kv.Key, skin);
			}
		}
	}
*/
func (c *CSharpExporter) getCSharpTypeStr(string typestr) []string {
	switch typestr {
		case "number":
			return []string { "int", "Int" }
		case "string":
			return []string { "string", "String" }
	}
	return []string { typestr, typestr }
}

func (c *CSharpExporter) makeMapMethodStr(t *Table) string {
	str := "\t\tpublic static void Map(IDataTable dt)\n"
	str += "\t\t{\n"

	str += "\t\t\tstring tbname = dt.GetName();\n"
	str += "\t\t\tint numRows = dt.GetRowCount();\n"
	str += "\t\t\tfor (int i = 0; i < numRows; ++i)\n"
	str += "\t\t\t{\n"

	str += fmt.Sprintf("\t\t\t\t%s o = Rows[i];\n", t.Name)
	
	for i, fi := range t.FieldInfos {
		switch fi.Type {
		case "list":
			str += fmt.Sprintf("\t\t\t\to.%s = dt.GetInt(%d)\n", fi.Name, i)
		case "map":
			kStrs := c.getCSharpTypeStr(fi.Key)
			vStrs := c.getCSharpTypeStr(fi.Value)
			str += fmt.Sprintf("\t\t\t\tDictionary<%s, %s> __%sMap = dt.Get%s%sDict();\n", kStrs[0], vStrs[0], fi.Value, kStrs[1], vStrs[1])
			str += fmt.Sprintf("\t\t\t\tforeach (KeyValuePair<%s, %s> kv in __%sMap)\n", kStrs[0], vStrs[0], fi.Value)
			str += "\t\t\t{\n"
			str += fmt.Sprintf("\t\t\t%s __temp = Skin.IdMap[kv.Value];\n", fi.)
			str += "\t\t\t\to.SkinMap.Add(kv.Key, __temp);
			}
		default:
			
		}
	}
	

	str += "\t\t\t}\n"
	str += "\t\t}\n"
		{
			int CClassId = dt.GetInt(i, 2);
			unit.CClass = CharacterClass.IdMap[CClassId];

			Dictionary<string, int> SkinMapFields = dt.GetStringIntDict();
			foreach (KeyValuePair<string, int> kv in SkinMapFields)
			{
				Skin skin = Skin.IdMap[kv.Value];
				unit.SkinMap.Add(kv.Key, skin);
			}
		}
	}

	return str + "\n"
}

func (c *CSharpExporter) makeMethodsStr(t *Table) string {
	str := c.makeLoadMethodStr(t)
	str += c.makeMapMethodStr(t)
	return str + "\n"
}

func (c *CSharpExporter) makeClassStr(fieldsStr, methodsStr string, t *Table) string {
	str := "using System;\n"
	str += "using System.Collections.Generic;\n"
	str += "\n"
	str += "namespace Datatable\n"
	str += "{\n"
	str += "\tpublic Class " + t.Name + "\n"
	str += "\t{\n"
	str += fieldsStr + "\n"
	str += methodsStr + "\n"
	str += "\t}\n"
	str += "}\n"
	return str
}

func (c *CSharpExporter) Save(path string, table *Table) error {
	fieldsStr := c.makeFieldsStr(table)
	methodsStr := c.makeMethodsStr(table)
	str := c.makeClassStr(fieldsStr, methodsStr, table)

	fullpath := filepath.Join(path, table.Name+".cs")
	ioutil.WriteFile(fullpath, []byte(str), 0666)
	return nil
}
