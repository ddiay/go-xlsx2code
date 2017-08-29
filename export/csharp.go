package export

import (
	"fmt"
	"io/ioutil"
	"os"
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
	case "bool":
		str += fmt.Sprintf("\t\tpublic bool %s;\n", fi.Name)
	case "list":
		str += fmt.Sprintf("\t\tpublic List<%s> %s = new List<%s>();\n", fi.Value, fi.Name, fi.Value)
	case "map":
		str += fmt.Sprintf("\t\tpublic Dictionary<%s, %s> %s = new Dictionary<%s, %s>();\n", fi.Key, fi.Value, fi.Name, fi.Key, fi.Value)
	default:
		str += fmt.Sprintf("\t\tpublic %s %s;\n", fi.Type, fi.Name)
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
	str := "\t\tpublic static void PreLoad(IDatatable dt)\n"
	str += "\t\t{\n"
	str += "\t\t\tstring tbname = dt.GetName();\n"
	str += "\t\t\tint numRows = dt.GetRowCount();\n"
	str += "\t\t\tfor (int i = 0; i < numRows; ++i)\n"
	str += "\t\t\t{\n"
	str += fmt.Sprintf("\t\t\t\t%s o = new %s();\n", t.Name, t.Name)

	for i, fi := range t.FieldInfos {
		switch fi.Type {
		case "number":
			str += fmt.Sprintf("\t\t\t\to.%s = dt.GetInt(i, %d);\n", fi.Name, i)
		case "string":
			str += fmt.Sprintf("\t\t\t\to.%s = dt.GetString(i, %d);\n", fi.Name, i)
		case "float":
			str += fmt.Sprintf("\t\t\t\to.%s = dt.GetDouble(i, %d);\n", fi.Name, i)
		case "bool":
			str += fmt.Sprintf("\t\t\t\to.%s = dt.GetBool(i, %d);\n", fi.Name, i)
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
func (c *CSharpExporter) getCSharpTypeStr(typestr string) []string {
	switch typestr {
	case "number":
		return []string{"int", "Int"}
	case "string":
		return []string{"string", "String"}
	}
	return []string{typestr, typestr}
}

func (c *CSharpExporter) makePostLoadMethodStr(t *Table) string {
	str := "\t\tpublic static void PostLoad(IDatatable dt)\n"
	str += "\t\t{\n"

	str += "\t\t\tstring tbname = dt.GetName();\n"
	str += "\t\t\tint numRows = dt.GetRowCount();\n"
	str += "\t\t\tfor (int i = 0; i < numRows; ++i)\n"
	str += "\t\t\t{\n"

	str += fmt.Sprintf("\t\t\t\t%s o = Rows[i];\n", t.Name)

	for i, fi := range t.FieldInfos {
		if len(fi.RefClassFieldName) == 0 {
			continue
		}

		switch fi.Type {
		case "list":
			typeStr := FindIndexType(fi.Value + "." + fi.RefClassFieldName)
			keyStrs := c.getCSharpTypeStr(typeStr)
			str += fmt.Sprintf("\t\t\t\tvar __%sMap = dt.Get%sList(i, %d);\n", fi.Value, keyStrs[1], i)
			str += fmt.Sprintf("\t\t\t\tforeach (var kv in __%sMap)\n", fi.Value)
			str += "\t\t\t\t{\n"
			str += fmt.Sprintf("\t\t\t\t\to.%s.Add(kv.Key, %s.%sMap[kv.Value]);\n", fi.Name, fi.Value, fi.RefClassFieldName)
			str += "\t\t\t\t}\n"

		case "map":
			keyStrs := c.getCSharpTypeStr(fi.Key)
			typeStr := FindIndexType(fi.Value + "." + fi.RefClassFieldName)
			valStrs := c.getCSharpTypeStr(typeStr)
			str += fmt.Sprintf("\t\t\t\tvar __%sMap = dt.Get%s%sDict(i, %d);\n", fi.Value, keyStrs[1], valStrs[1], i)
			str += fmt.Sprintf("\t\t\t\tforeach (var kv in __%sMap)\n", fi.Value)
			str += "\t\t\t\t{\n"
			str += fmt.Sprintf("\t\t\t\t\to.%s.Add(kv.Key, %s.%sMap[kv.Value]);\n", fi.Name, fi.Value, fi.RefClassFieldName)
			str += "\t\t\t\t}\n"

		default:
			typeStr := FindIndexType(fi.Type + "." + fi.RefClassFieldName)
			if len(typeStr) > 0 {
				csTypeStrs := c.getCSharpTypeStr(typeStr)
				str += fmt.Sprintf("\t\t\t\tvar __%s = dt.Get%s(i, %d);\n", fi.Name, csTypeStrs[1], i)
				str += fmt.Sprintf("\t\t\t\to.%s = %s.%sMap[__%s];\n", fi.Name, fi.Type, fi.RefClassFieldName, fi.Name)
			}
		}
	}

	str += "\t\t\t}\n"
	str += "\t\t}\n"
	// 	{
	// 		int CClassId = dt.GetInt(i, 2);
	// 		unit.CClass = CharacterClass.IdMap[CClassId];

	// 		Dictionary<string, int> SkinMapFields = dt.GetStringIntDict();
	// 		foreach (KeyValuePair<string, int> kv in SkinMapFields)
	// 		{
	// 			Skin skin = Skin.IdMap[kv.Value];
	// 			unit.SkinMap.Add(kv.Key, skin);
	// 		}
	// 	}
	// }

	return str + "\n"
}

func (c *CSharpExporter) makeMethodsStr(t *Table) string {
	str := c.makeLoadMethodStr(t)
	str += c.makePostLoadMethodStr(t)
	return str + "\n"
}

func (c *CSharpExporter) makeClassStr(fieldsStr, methodsStr string, t *Table) string {
	str := "using System;\n"
	str += "using System.Collections.Generic;\n"
	str += "\n"
	str += "namespace Datatable\n"
	str += "{\n"
	str += "\tpublic class " + t.Name + "\n"
	str += "\t{\n"
	str += fieldsStr + "\n"
	str += methodsStr + "\n"
	str += "\t}\n"
	str += "}\n"
	return str
}

func (c *CSharpExporter) makeTableLoaderStr(preLoadStr string, postLoadStr string) string {
	str := "using System;\n"
	str += "using System.Collections.Generic;\n"
	str += "\n"
	str += "namespace Datatable\n"
	str += "{\n"
	str += "\tpublic class TableLoader\n"
	str += "\t{\n"

	str += "\t}\n"
	str += "}\n"
	return str
}

func (c *CSharpExporter) Save(path string, tables []Table) error {
	// preLoadStr := ""
	// postLoadStr := ""
	str := ""
	fullpath := ""
	cspath := filepath.Join(path, "cs")
	os.MkdirAll(cspath, 0777)
	for _, t := range tables {
		fieldsStr := c.makeFieldsStr(&t)
		methodsStr := c.makeMethodsStr(&t)
		str = c.makeClassStr(fieldsStr, methodsStr, &t)

		fullpath = filepath.Join(cspath, t.Name+".cs")
		ioutil.WriteFile(fullpath, []byte(str), 0666)

		// preLoadStr += c.makePreLoadStr(t.Name)
		// postLoadStr += c.makePostLoadStr(t.Name)
	}

	// str = c.makeTableLoaderStr(preLoadStr, postLoadStr)
	// path = filepath.Join(path, "TableLoader.cs")
	// ioutil.WriteFile(path, []byte(str), 0666)

	return nil
}
