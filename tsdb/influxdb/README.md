# Influxdb v1.x 客户端
* 本包内部不处理db、rp、measurement、tag key、field key特殊命名问题，由调用方保证符合influxdb规范，建议参考下面的命名规范
* 对于tag value、field value中特殊字符，基于经验做了一定的转义处理，详见 EscapeXXX 系列函数，如果使用 RawQuery|RawWrite 则需自行处理转义逻辑

## db、rp、measurement、tag key、field key命名规范
- 大小写敏感：统一使用小写字母
- 字符限制：
  - 以字母开头，可包含字母、数字、下划线、连接符
  - 避免使用其他特殊字符
  - 避免使用保留字
- 长度限制：没有严格的长度限制，但建议保持名称简短且具描述性
- 保持一致性：在整个数据库中保持一致的命名约定，便于管理和查询

## tag value、field value格式
- 由于行协议和influxQL语法差异，在读写数据时对值格式的要求不同，所以读写时转义逻辑也不同
- 写入数据：
  - tag value不能用双引号包裹，field value如果是字符值必须用双引号包裹，所以转义逻辑应该区分对待
- 查询数据：
  - tag value和field value都是以单引号包裹,time字段和数值型field value除外
- tag字段通常用于索引和过滤，因此它们的值应该具有较低的基数（即重复率高），基数太大的话会导致内存占用过高
- field字段用于存储实际的测量值，具有较高的基数

## 数据类型
* https://docs.influxdata.com/influxdb/v1/write_protocols/line_protocol_reference/#data-types
* 虽然 InfluxDB 支持多种数据类型，但是在实际使用中，建议将所有数值类型统一转换为浮点数（float）进行存储，以便于后续的计算和分析。
* 字符串类型无法使用聚合函数进行计算，应尽量避免使用。
* 字符串类型最大长度为 64KB。

## 常见问题
* https://docs.influxdata.com/influxdb/v1/troubleshooting/frequently-asked-questions