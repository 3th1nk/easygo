# 类型转换
- 用于将 interface 对象转化为 bool|int|float|string 等类型。
- 在对象与字符串的转换过程中会判断对象上是否实现了 StringMarshal、fmt.Stringer 接口，如果都没有则使用 JSON 格式进行序列化和反序列化。
- 判断字符串值的类型 GetStrValueType
  - "123" 会判断出这是个 Int
  - "123.456" 会判断是个浮点数
  - "abc" 会判断是个字符串
  - "{"a":123}" 会判断是个 Map
  - "[1,2,3]"、"[{"a":123}]" 会判断是个数组
