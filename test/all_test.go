package test

import (
	"fmt"
	"testing"

	"github.com/vedadiyan/sqlparser/pkg/sqlparser"
)

func TestParse(t *testing.T) {
	queries := []string{
		"SELECT * FROM table_name;",
		"SELECT column1, column2 FROM table_name;",
		"SELECT * FROM table_name WHERE column_name = 'value';",
		"SELECT COUNT(*) FROM table_name;",
		"SELECT column1, COUNT(*) FROM table_name GROUP BY column1;",
		"SELECT column1, AVG(column2) FROM table_name GROUP BY column1;",
		"SELECT column1, MAX(column2) FROM table_name GROUP BY column1;",
		"SELECT column1, MIN(column2) FROM table_name GROUP BY column1;",
		"SELECT * FROM table_name ORDER BY column1 DESC;",
		"SELECT column1, column2 FROM table_name WHERE column1 LIKE 'a%';",
		"SELECT column1, column2 FROM table_name WHERE column1 IN (1, 2, 3);",
		"SELECT column1, column2 FROM table_name WHERE column1 BETWEEN 1 AND 10;",
		"SELECT column1, column2 FROM table_name WHERE column1 IS NULL;",
		"SELECT column1, column2 FROM table_name WHERE column1 IS NOT NULL;",
		"SELECT column1, column2 FROM table_name WHERE column1 = (SELECT column1 FROM other_table WHERE column2 = 'value');",
		"WITH cte_name AS (SELECT column1, column2 FROM table_name WHERE column1 = 'value') SELECT * FROM cte_name;",
		"SELECT column1, column2, CASE WHEN column1 = 'value' THEN 'yes' ELSE 'no' END AS new_column FROM table_name;",
		"SELECT column1, column2 FROM table_name WHERE BINARY column1 = 'Value';",
		"SELECT column1, column2 FROM table_name WHERE column1 + column2 > 100;",
		"SELECT column1, column2 FROM table_name WHERE NOT (column1 = 'value' AND column2 = 100);",
		"SELECT * FROM table1 JOIN table2 ON table1.column = table2.column;",
		"SELECT * FROM table1 LEFT JOIN table2 ON table1.column = table2.column;",
		"SELECT * FROM table1 RIGHT JOIN table2 ON table1.column = table2.column;",
		"SELECT * FROM table1 JOIN table2 ON table1.column = table2.column;",
		"SELECT column1, column2 FROM table1 UNION SELECT column1, column2 FROM table2;",
		"SELECT * FROM table_name LIMIT 10 OFFSET 5;",
		"SELECT column1, COUNT(*) FROM table_name GROUP BY column1 HAVING COUNT(*) > 10;",
		"SELECT column1, column2 FROM table1 INNER JOIN table2 ON table1.column = table2.column WHERE table2.column = 'value';",
		"SELECT column1, column2 FROM table1 LEFT JOIN table2 ON table1.column = table2.column WHERE table2.column IS NULL;",
		"SELECT column1, column2 FROM table1 RIGHT JOIN table2 ON table1.column = table2.column WHERE table1.column IS NULL;",
	}
	for _, query := range queries {
		_, err := sqlparser.Parse(query)
		if err != nil {
			fmt.Println(query, err.Error())
			t.FailNow()
		}
	}
}
