# 字符串操作
- case: 首字母大小写转换、在 CamelCase、SnakeCase、PascalCase、KebabCase 几种命名格式之间做转换。
- join: 扩展 strings.Join，允许将 []int、[]interface 类型拼接为字符串。
- rand: 随机字符串生成
- split: 扩展 strings.Split，允许将字符串拆分为 []int，且允许设置是否过滤空字符串等更多参数。
- stringBuilder: 当需要频繁拼接字符串时，可用此类代替字符串相加操作。备注：如果只是少量字符串的拼接，建议仍使用字符串相加，因少量拼接的情况下两者性能并无差别。
- 其他常用函数的封装