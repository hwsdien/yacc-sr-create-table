package main

import (
	"fmt"
	"regexp"
)

func isSnakeCase(str string) bool {
	re := regexp.MustCompile(`^[a-z0-9_]+$`)
	if !re.MatchString(str) {
		return false
	}
	// 判断下划线的位置是否符合蛇形命名规则
	for i := 1; i < len(str)-1; i++ {
		if str[i] == '_' && (str[i-1] == '_' || str[i+1] == '_') {
			return false
		}
	}
	return true
}

func inFrontOfColumnList(strs []string, targets []ColumnField) bool {
	positions := make([]int, len(strs))
	for i, target := range targets {
		found := false
		for j, str := range strs {
			if str == target.ColName {
				positions[i] = j
				found = true
				break
			}
		}
		if !found {
			break
		}
	}
	for i := 1; i < len(positions); i++ {
		if positions[i] <= positions[i-1] {
			return false
		}
	}
	return true
}

func containsString(list []string, target string) bool {
	for _, str := range list {
		if str == target {
			return true
		}
	}
	return false
}

func Review(tableInfo CreateStmtTable) {
	/**
	假设的规范如下：
	1, 字段的命名是蛇形命名规则;
	2, 建表时必须使用 DISTRIBUTED BY HASH 子句指定分桶键，否则建表失败；
	3, 主键通过 PRIMARY KEY 定义，且主键列必须在其他列之前;
	4, 必须指定分桶列;
	5, 分区列和分桶列必须在主键中;
	*/

	found := false
	fmt.Printf("发现问题并建议如下：\n")
	for _, column := range tableInfo.ColumnList {
		if !isSnakeCase(column.ColName) {
			found = true
			fmt.Printf("\"%s\"的命名建议使用蛇形命名规则。\n", column.ColName)
		}
	}

	if tableInfo.DistributedData.BucketCount == 0 {
		found = true
		fmt.Printf("必须使用 DISTRIBUTED BY HASH 子句指定分桶键。\n")
	}

	if len(tableInfo.PrimaryKeyList) > 0 {
		if !inFrontOfColumnList(tableInfo.PrimaryKeyList, tableInfo.ColumnList) {
			found = true
			fmt.Printf("主键列必须在其他列之前, 当前的主键列%+v不符合规则。\n", tableInfo.PrimaryKeyList)
		}
	} else {
		found = true
		fmt.Printf("貌似缺少主键列, 请指定主键列。\n")
	}

	if len(tableInfo.PartitionRangeData.RangeValue) > 0 {
		if !containsString(tableInfo.PrimaryKeyList, tableInfo.PartitionRangeData.RangeValue) {
			found = true
			fmt.Printf("分区列必须在主键中，当前的分区列是：%s。\n", tableInfo.PartitionRangeData.RangeValue)
		}
	}

	if len(tableInfo.DistributedData.HashField) > 0 {
		if !containsString(tableInfo.PrimaryKeyList, tableInfo.DistributedData.HashField) {
			found = true
			fmt.Printf("分桶列必须在主键中，当前的分桶列是：%s。\n", tableInfo.DistributedData.HashField)
		}
	} else {
		found = true
		fmt.Printf("缺少分桶列，请指定分桶列。\n")
	}

	if !found {
		fmt.Printf("无\n")
	}

}
