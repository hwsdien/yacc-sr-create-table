package main

import (
	"errors"
	"strings"
	"text/scanner"
	"unicode"
)

type CreateStmtTable struct {
	Database           string             `json:"database,omitempty"`
	TableName          string             `json:"tableName,omitempty"`
	ColumnList         []ColumnField      `json:"columnList,omitempty"`
	PrimaryKeyList     []string           `json:"primaryKeyList,omitempty"`
	PartitionRangeData PartitionRangeInfo `json:"partitionRangeInfo"`
	DistributedData    DistributedInfo    `json:"distributedInfo"`
	PropertiesData     []KeyValueInfo     `json:"propertiesInfo"`
}

type KeyValueInfo struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type DistributedInfo struct {
	HashField   string `json:"hashField,omitempty"`
	BucketCount int    `json:"bucketCount,omitempty"`
}

type PartitionRangeInfo struct {
	RangeValue    string          `json:"rangeValue,omitempty"`
	PartitionList []PartitionInfo `json:"partitionList,omitempty"`
}

type PartitionInfo struct {
	Name   string   `json:"name,omitempty"`
	Values []string `json:"values,omitempty"`
}

type ColumnField struct {
	ColName string `json:"columnName,omitempty"`
	ColType string `json:"columnType,omitempty"`
	IsNull  bool   `json:"isNull,omitempty"`
}

type tableLex struct {
	scanner.Scanner
	tableInfo CreateStmtTable
	tokenDict map[string]int
	err       error
}

func (t *tableLex) Lex(lval *srSymType) int {
	s, lit := t.Scan(), t.TokenText()
	litLower := strings.ToLower(lit)

	switch s {
	case scanner.EOF:
		return scanner.EOF
	case scanner.String:
		lval.ident = strings.Trim(lit, "\"")
		return identifier
	case scanner.Char:
		lval.ident = strings.Trim(lit, "'")
		return identifier
	case scanner.Int:
		lval.num = lit
		return number
	case scanner.Ident:
		if r, ok := t.tokenDict[litLower]; ok {
			lval.word = lit
			return r
		} else {
			lval.ident = lit
			return identifier
		}
	default:
		switch lit {
		case "(":
			return int('(')
		case ")":
			return int(')')
		case ",":
			return int(',')
		case ".":
			return int('.')
		case "[":
			return int('[')
		case "-":
			return int('-')
		case "=":
			return int('=')
		case ";":
			return int(';')
		}
	}

	return 0
}

func (t *tableLex) Error(s string) {
	t.err = errors.New(s)
}

func Parse(stmt string) (CreateStmtTable, error) {
	s := new(tableLex)
	s.tokenDict = map[string]int{
		"create":      create,
		"table":       table,
		"if":          ifWord,
		"not":         not,
		"exists":      exists,
		"null":        null,
		"default":     defaultWord,
		"primary":     primaryWord,
		"key":         keyWord,
		"partition":   partitionWord,
		"by":          byWord,
		"range":       rangeWord,
		"values":      valuesWord,
		"distributed": distributedWord,
		"hash":        hashWord,
		"buckets":     bucketsWord,
		"properties":  propertiesWord,

		"date":    dateType,
		"tinyint": tinyintType,
		"int":     intType,
		"bigint":  bigintType,
		"string":  stringType,
	}
	s.tableInfo = CreateStmtTable{}
	s.Init(strings.NewReader(stmt))
	s.IsIdentRune = func(ch rune, i int) bool {
		return ch == '-' && i > 0 || ch == '_' || unicode.IsLetter(ch) || unicode.IsDigit(ch) && i > 0
	}
	// s.Mode = scanner.GoTokens ^ scanner.ScanChars
	srParse(s)

	// fmt.Printf("%+v\n", s.tableInfo)
	return s.tableInfo, s.err

}
