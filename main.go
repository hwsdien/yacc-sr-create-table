package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func main() {
	sqlCreateTable := `
create table orders (
    dt date NOT NULL,
    order_id bigint NOT NULL,
    user_id int NOT NULL,
    merchant_id int NOT NULL,
    good_id int NOT NULL,
    good_name string NOT NULL,
    price int NOT NULL,
    cnt int NOT NULL,
    revenue int NOT NULL,
    state tinyint NOT NULL
) PRIMARY KEY (dt, order_id)
PARTITION BY RANGE(dt) (
    PARTITION p20210820 VALUES [('2021-08-20'), ('2021-08-21')),
    PARTITION p20210821 VALUES [('2021-08-21'), ('2021-08-22')),
    PARTITION p20210929 VALUES [('2021-09-29'), ('2021-09-30')),
    PARTITION p20210930 VALUES [('2021-09-30'), ('2021-10-01'))
) DISTRIBUTED BY HASH(order_id) BUCKETS 4
PROPERTIES (
    "replication_num" = "3",
    "enable_persistent_index" = "true"
);
`

	fmt.Println("Begin Parse")
	fmt.Println("========================================================================")

	sqlContent := strings.Replace(sqlCreateTable, "\n", " ", -1)
	result, err := Parse(sqlContent)
	fmt.Println("原始的SQL建表语句：")
	fmt.Println(sqlCreateTable)
	if err != nil {
		fmt.Printf("原始的SQL建表语句有问题，错误：%s\n", err.Error())
		fmt.Println("========================================================================")
		fmt.Println("End Parse")
		return
	}

	fmt.Println("经过词法和语法分析后生成的JSON：")
	jsonInfo, _ := json.Marshal(result)
	fmt.Println(string(jsonInfo))

	fmt.Println("========================================================================")
	fmt.Println("End Parse")
	fmt.Println("========================================================================")
	fmt.Printf("Review 建表语句\n\n")
	Review(result)
	fmt.Printf("\nReview 结束\n")

}
