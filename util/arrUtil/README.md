# 切片操作
- Convertor
  - 从 []T1 转换为 []T2，比如将 []string 转化为 []int。
  
- Find
  - 从切片中查找符合条件的子集，返回一个新的切片。
  - 派生出来的有 First、FirstInt、FirstString 等，返回符合条件的第一个元素。
  
- IndexOf
  - 从切片中查找符合条件的第一个元素的索引
  - 派生出来的有 IndexOfInt、IndexOfString、IndexOfInt64 等。
  
- Remove/RemoveAt
  - 从切片中删除符合条件的元素，并返回新的切片。
  
- Repeat
  - 重复一个对象N次，形成一个切片数组。循环过程中可以通过自定义函数控制生成的对象。
  - 派生出来的有 RepeatInt、RepeatString 等。