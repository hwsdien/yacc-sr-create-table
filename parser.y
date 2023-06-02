%{
package main

import (
	"fmt"
	"strconv"
)

var tmpValues []string
%}

%union {
	stmt string
	col	string
	agg string
	word string
	ident string
	num string
	values []string
	boolean bool
}

%type <stmt>
	stmt
	externalStmt
	ifNotExistsStmt
	tableNameStmt

	columnListStmt
	columnFieldListStmt
	columnStmt

	colTypeStmt
	aggTypeStmt
	defaultStmt

	primaryKeyStmt
	primaryKeyColumnListStmt

	partitionStmt
	partitionRangeStmt
	partitonRangeListStmt

	strStmt

	distributedStmt

	propertiesStmt
	keyValueListStmt
	keyValueStmt

%type <values>
	partitionValueStmt

%type <boolean>
	notNullStmt

%type <stmt> main

%token <ident> identifier

%token <num> number

%token <word> '[' '(' ')' ',' '.' '-' ';'

%token <word>
	create
	external
	table
	ifWord
	not
	exists
	null
	defaultWord

	primaryWord
	keyWord

	partitionWord
	byWord
	rangeWord
	valuesWord

	distributedWord
	hashWord
	bucketsWord

	propertiesWord

%token <col>
	tinyintType
	smallintType
	intType
	bigintType
	largeIntType
	floatType
	doubleType
	decimalType
	dateType
	dateTimeType
	charType
	varcharType
	stringType
	hllType
	bitmapType

%token <agg>
	sum
	max
	min
	replace
	hllUnion
	bitmapUnion
	replaceIfNotNull

%start main

%%

main	:
			{}
		|	stmt
			{ $$ = $1 }
		|	stmt ';'
			{ $$ = $1 }
		;

stmt	:	create externalStmt table ifNotExistsStmt tableNameStmt columnListStmt
primaryKeyStmt partitionStmt distributedStmt propertiesStmt
			{}
		;
externalStmt	:
					{}
				|	external
					{ $$ = $1 }
				;

ifNotExistsStmt	:
					{}
				|	ifWord not exists
					{}
				;

tableNameStmt	:	identifier
					{
						srlex.(*tableLex).tableInfo.TableName = $1
					}
				|	identifier '.' identifier
					{
						srlex.(*tableLex).tableInfo.Database = $1
						srlex.(*tableLex).tableInfo.TableName = $3
					}
				;

columnListStmt	:
					{}
				|	'(' columnFieldListStmt ')'
					{ $$ = $2 }
				;
columnFieldListStmt	:
						{}
					| columnStmt
						{ $$ = $1 }
					|	columnFieldListStmt ',' columnStmt
						{ $$ = $1 }
					;
columnStmt	:	identifier colTypeStmt aggTypeStmt notNullStmt defaultStmt
			{
				srlex.(*tableLex).tableInfo.ColumnList = append(
					srlex.(*tableLex).tableInfo.ColumnList,
					ColumnField{ColName: $1, ColType: $<word>2, IsNull: $4},
				)
			}
		;
colTypeStmt	:
				{}
			|	tinyintType
				{ $$ = $1 }
			|	intType
				{ $$ = $1 }
			|	bigintType
				{ $$ = $1 }
			|	dateType
				{ $$ = $1 }
			|	stringType
				{ $$ = $1 }
			;
aggTypeStmt	:
				{}
			|	sum
				{}
			|	max
				{}
			;
notNullStmt	:
				{$$ = true}
			|	null
				{$$ = true}
			|	not null
				{$$ = false}
			;
defaultStmt	:
				{}
			|	defaultWord '"' identifier '"'
				{}
			;
primaryKeyStmt	:
					{}
				|	primaryWord keyWord '(' primaryKeyColumnListStmt ')'
					{}
				;
primaryKeyColumnListStmt	:
								{}
							|	identifier
								{
									srlex.(*tableLex).tableInfo.PrimaryKeyList = append(
										srlex.(*tableLex).tableInfo.PrimaryKeyList,
										$1,
									)

									$$ = $1
								}
							|	primaryKeyColumnListStmt ',' identifier
								{
									srlex.(*tableLex).tableInfo.PrimaryKeyList = append(
										srlex.(*tableLex).tableInfo.PrimaryKeyList,
										$3,
									)
									$$ = $3
								}
							;
partitionStmt	:
					{}
				|	partitionWord byWord rangeWord '(' identifier ')' '(' partitionRangeStmt ')'
					{
						srlex.(*tableLex).tableInfo.PartitionRangeData.RangeValue = $5
						$$ = $5
					}
				;
partitionRangeStmt	:
						{}
					|	partitonRangeListStmt
						{ fmt.Println($1) }
					|	partitionRangeStmt ',' partitonRangeListStmt
						{ fmt.Println($3) }
					;
partitonRangeListStmt	:
							{}
						|	partitionWord identifier valuesWord '[' partitionValueStmt ')'
							{
								srlex.(*tableLex).tableInfo.PartitionRangeData.PartitionList = append(
									srlex.(*tableLex).tableInfo.PartitionRangeData.PartitionList,
									PartitionInfo{
										Name: $2,
										Values: $5,
									},
								)

								tmpValues = make([]string, 0)
							}
						;
partitionValueStmt	:
						{}
					|	strStmt
						{
							tmpValues = append(tmpValues, $1)
							$$ = tmpValues
						}

					|	partitionValueStmt ',' strStmt
						{
							tmpValues = append(tmpValues, $3)
							$$ = tmpValues
						}
					;
strStmt	:
			{}
		|	'(' identifier ')'
			{ $$ = $2 }
		;
distributedStmt	:
					{}
				|	distributedWord byWord hashWord '(' identifier ')' bucketsWord number
					{
						srlex.(*tableLex).tableInfo.DistributedData.HashField = $5
						srlex.(*tableLex).tableInfo.DistributedData.BucketCount, _ = strconv.Atoi($8)
					}
				;
propertiesStmt	:
					{}
				|	propertiesWord '(' keyValueListStmt ')'
					{ $$ = $3 }
				;

keyValueListStmt	:
						{}
					|	keyValueStmt
						{ $$ = $1 }
					|	keyValueListStmt ',' keyValueStmt
						{ $$ = $3 }
					;
keyValueStmt	:
					{}
				|	identifier '=' identifier
					{
						srlex.(*tableLex).tableInfo.PropertiesData = append(
							srlex.(*tableLex).tableInfo.PropertiesData,
							KeyValueInfo{
								Key: $1,
								Value: $3,
							},
						)
					 $$ = $1
					}
				;
%%


