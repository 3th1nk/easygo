package influxdb

// TODO 封装CQ操作
//	基本语法：
//		CREATE CONTINUOUS QUERY <cq_name> ON <database_name>
//		BEGIN
//  		<cq_query>
//		END
// 	- cq_name: 连续查询的名称
// 	- database_name: 连续查询的数据库
// 	- cq_query: 连续查询的SQL语句，语法参考：https://docs.influxdata.com/influxdb/v1/query_language/continuous_queries/#description-of-basic-syntax
