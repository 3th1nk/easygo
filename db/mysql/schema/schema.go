package schema

import (
	"github.com/3th1nk/easygo/util/convertor"
	"github.com/3th1nk/easygo/util/mapUtil"
	"github.com/3th1nk/easygo/util/strUtil"
	"github.com/modern-go/reflect2"
	"github.com/mohae/deepcopy"
	"regexp"
	"sort"
	"strings"
	"time"
)

type DbQuery interface {
	QueryString(sqlOrArgs ...interface{}) ([]mapUtil.StringMap, error)
}

type TableSchema struct {
	Database     string         `json:"database"`
	Table        string         `json:"table"`
	Type         string         `json:"type"`
	StructName   string         `json:"structName,omitempty"`
	TableComment string         `json:"tableComment,omitempty"`
	Fields       []*FieldSchema `json:"fields,omitempty"`
	Indexes      []*IndexSchema `json:"indexes,omitempty"`
}

type FieldSchema struct {
	ColumnName     string `json:"columnName" description:"列名"`
	ColumnType     string `json:"columnType" description:"列的完整类型（DataType + DataLength）"`
	ColumnNullable bool   `json:"columnNullable,omitempty" description:""`
	ColumnDefault  string `json:"columnDefault,omitempty" description:"列默认值"`
	ColumnComment  string `json:"columnComment,omitempty" description:"列注释"`
	DataType       string `json:"dataType,omitempty" description:"列数据类型"`
	FieldName      string `json:"fieldName,omitempty" description:"对应结构体的字段名"`
	FieldType      string `json:"fieldType,omitempty" description:"对应结构体的字段类型"`
}

type IndexSchema struct {
	Name    string   `json:"name"`
	Primary bool     `json:"primary,omitempty"`
	Unique  bool     `json:"unique,omitempty"`
	Columns []string `json:"columns,omitempty"`
}

type StructFormatter struct {
	NameFormatter      func(database, table string) string
	FieldNameFormatter func(database, table string, column *FieldSchema) string
	FieldTypeFormatter func(database, table string, column *FieldSchema) string
}

var (
	customTypePattern = regexp.MustCompile(`{type\s*[:=]\s*(\*?[\w\d-_]+)}`)
)

func (this *TableSchema) GetFieldSchema(column string) *FieldSchema {
	for _, schema := range this.Fields {
		if strings.EqualFold(column, schema.ColumnName) {
			return schema
		}
	}
	return nil
}

func (this *TableSchema) GetUniqueSchema(columns []string) *IndexSchema {
	columnsCount := len(columns)
	columnsCopy := deepcopy.Copy(columns).([]string)
	sort.Strings(columnsCopy)
	columnsStr := strings.Join(columnsCopy, ",")
	for _, schema := range this.Indexes {
		if schema.Unique && columnsCount == len(schema.Columns) {
			columnsCopy = deepcopy.Copy(schema.Columns).([]string)
			sort.Strings(columnsCopy)
			if strings.EqualFold(columnsStr, strings.Join(columnsCopy, ",")) {
				return schema
			}
		}
	}
	return nil
}

func (this *TableSchema) GetIndexSchema(columns []string, checkOrder bool) *IndexSchema {
	columnsCount := len(columns)
	if checkOrder {
		for _, schema := range this.Indexes {
			if columnsCount == len(schema.Columns) && strings.EqualFold(strings.Join(columns, ","), strings.Join(schema.Columns, ",")) {
				return schema
			}
		}
	} else {
		columnsCopy := deepcopy.Copy(columns).([]string)
		sort.Strings(columnsCopy)
		columnsStr := strings.Join(columnsCopy, ",")
		for _, schema := range this.Indexes {
			if columnsCount == len(schema.Columns) {
				columnsCopy = deepcopy.Copy(schema.Columns).([]string)
				sort.Strings(columnsCopy)
				if strings.EqualFold(columnsStr, strings.Join(columnsCopy, ",")) {
					return schema
				}
			}
		}
	}
	return nil
}

func (this *FieldSchema) ColumnDefaultVal() interface{} {
	return this.ParseValue(this.ColumnDefault)
}

func (this *FieldSchema) ParseValue(val string) interface{} {
	switch this.FieldType {
	case "bool":
		return convertor.ToBoolNoError(val)
	case "int":
		return convertor.ToIntNoError(val)
	case "int64":
		return convertor.ToInt64NoError(val)
	case "time.Time":
		v, _ := time.ParseInLocation("2006-01-02 15:04:05", val, time.Local)
		return v
	default:
		return val
	}
}

func GetTableSchema(db DbQuery, database, table string, formatter ...*StructFormatter) (*TableSchema, error) {
	tableRows, err := db.QueryString("SELECT * FROM information_schema.TABLES where TABLE_SCHEMA = ? and TABLE_NAME = ?", database, table)
	if !reflect2.IsNil(err) {
		return nil, err
	} else if len(tableRows) == 0 {
		return nil, nil
	}

	columnRows, err := db.QueryString("SELECT * FROM information_schema.COLUMNS where TABLE_SCHEMA = ? and TABLE_NAME = ?", database, table)
	if !reflect2.IsNil(err) {
		return nil, err
	}

	indexRows, err := db.QueryString("SELECT INDEX_NAME, NON_UNIQUE, GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) COLUMN_NAMES FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? GROUP BY INDEX_NAME", database, table)
	if !reflect2.IsNil(err) {
		return nil, err
	}

	var theFormatter *StructFormatter
	if len(formatter) != 0 {
		theFormatter = formatter[0]
	}

	tableSchema := &TableSchema{
		Database: database,
		Table:    table,
		Fields:   make([]*FieldSchema, 0),
	}
	if theFormatter != nil && theFormatter.FieldNameFormatter != nil {
		tableSchema.StructName = theFormatter.NameFormatter(database, table)
	}
	if tableSchema.StructName == "" {
		tableSchema.StructName = strUtil.UcFirst(strUtil.SnakeToCamel(table))
	}

	row := mapUtil.NewCaseInsensitiveStringMap(tableRows[0])
	tableSchema.Type = row.MustGet("TABLE_TYPE")
	tableSchema.TableComment = row.MustGet("TABLE_COMMENT")
	for _, item := range columnRows {
		row := mapUtil.NewStringMap(item)
		field := &FieldSchema{
			ColumnName:     row.MustGet("COLUMN_NAME"),
			ColumnType:     row.MustGet("COLUMN_TYPE"),
			ColumnNullable: strings.EqualFold("YES", row.MustGet("IS_NULLABLE")),
			ColumnDefault:  row.MustGet("COLUMN_DEFAULT"),
			ColumnComment:  row.MustGet("COLUMN_COMMENT"),
			DataType:       row.MustGet("DATA_TYPE"),
		}
		tableSchema.Fields = append(tableSchema.Fields, field)
		// FieldName
		if theFormatter != nil && theFormatter.FieldNameFormatter != nil {
			field.FieldName = theFormatter.FieldNameFormatter(database, table, field)
		}
		if field.FieldName == "" {
			field.FieldName = strUtil.UcFirst(strUtil.SnakeToCamel(field.ColumnName))
		}
		// FieldType
		if theFormatter != nil && theFormatter.FieldTypeFormatter != nil {
			field.FieldType = theFormatter.FieldTypeFormatter(database, table, field)
		}
		if field.FieldType == "" {
			if matches := customTypePattern.FindAllStringSubmatch(field.ColumnComment, -1); len(matches) != 0 {
				field.FieldType = matches[0][1]
				field.ColumnComment = strings.TrimSpace(strings.Replace(field.ColumnComment, matches[0][0], "", -1))
			}
			if field.FieldType == "" {
				switch field.DataType {
				case "int", "tinyint":
					dataLen := convertor.ToIntNoError(strings.Trim(field.ColumnType[len(field.DataType):], "()"))
					if dataLen == 1 {
						field.FieldType = "bool"
					} else {
						field.FieldType = "int"
					}
				case "bigint":
					field.FieldType = "int64"
				case "varchar", "char", "longtext", "text":
					field.FieldType = "string"
				case "date", "datetime", "timestamp":
					field.FieldType = "time.Time"
				case "float":
					field.FieldType = "float64"
				default:
					field.FieldType = field.DataType
				}
			}
		}
	}

	indexMap := make(map[string]*IndexSchema)
	for _, item := range indexRows {
		row := mapUtil.NewCaseInsensitiveStringMap(item)
		schema := &IndexSchema{
			Name:    row.MustGet("INDEX_NAME"),
			Unique:  !row.MustGetBool("NON_UNIQUE"),
			Columns: strUtil.Split(row.MustGet("COLUMN_NAMES"), ",", true),
		}
		schema.Primary = schema.Name == "PRIMARY"
		tableSchema.Indexes = append(tableSchema.Indexes, schema)
		indexMap[schema.Name] = schema
	}

	return tableSchema, nil
}
