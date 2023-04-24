package matrixUtil

import (
	"fmt"
	"github.com/3th1nk/easygo/util/convertor"
	"github.com/3th1nk/easygo/util/mapUtil"
)

// 对二位数组进行行列转置
func Transpose(a [][]interface{}) [][]interface{} {
	m := len(a)
	if m == 0 {
		return nil
	}
	n := len(a[0])

	b := make([][]interface{}, n)
	for i := range b {
		tmp := make([]interface{}, m)
		for j := range tmp {
			tmp[j] = a[j][i]
		}
		b[i] = tmp
	}
	return b
}

// 对二位数组进行行列转置
func TransposeInt(a [][]int) [][]int {
	m := len(a)
	if m == 0 {
		return nil
	}
	n := len(a[0])

	b := make([][]int, n)
	for i := range b {
		tmp := make([]int, m)
		for j := range tmp {
			tmp[j] = a[j][i]
		}
		b[i] = tmp
	}
	return b
}

// 对二位数组进行行列转置
func TransposeString(a [][]string) [][]string {
	m := len(a)
	if m == 0 {
		return nil
	}
	n := len(a[0])

	b := make([][]string, n)
	for i := range b {
		tmp := make([]string, m)
		for j := range tmp {
			tmp[j] = a[j][i]
		}
		b[i] = tmp
	}
	return b
}

// 对二位数组进行行列转置
func TransposeFloat(a [][]float64) [][]float64 {
	m := len(a)
	if m == 0 {
		return nil
	}
	n := len(a[0])

	b := make([][]float64, n)
	for i := range b {
		tmp := make([]float64, m)
		for j := range tmp {
			tmp[j] = a[j][i]
		}
		b[i] = tmp
	}
	return b
}

func TransposeStringTable(column string, rows []mapUtil.StringMap) ([]mapUtil.StringMap, error) {
	rowCount := len(rows)
	if rowCount == 0 {
		return rows, nil
	}

	columnMap := make(map[string]bool, len(rows[0]))
	for _, row := range rows {
		for _, column := range row.Keys() {
			columnMap[column] = true
		}
	}
	if !columnMap[column] {
		return nil, fmt.Errorf("列 '%v' 不存在", column)
	}
	delete(columnMap, column)
	columnList, columnCount := mapUtil.StringKeys(columnMap), len(columnMap)

	result := make([]mapUtil.StringMap, columnCount)
	for i, str := range columnList {
		result[i] = make(mapUtil.StringMap, rowCount)
		result[i][column] = str
	}
	for _, row := range rows {
		keyVal, _ := convertor.ToString(row[column])
		for idx, str := range columnList {
			if val, ok := row[str]; ok {
				result[idx][keyVal] = val
			}
		}
	}

	return result, nil
}

func TransposeStringObjectTable(column string, rows []mapUtil.StringObjectMap) ([]mapUtil.StringObjectMap, error) {
	rowCount := len(rows)
	if rowCount == 0 {
		return rows, nil
	}

	columnMap := make(map[string]bool, len(rows[0]))
	for _, row := range rows {
		for _, column := range row.Keys() {
			columnMap[column] = true
		}
	}
	if !columnMap[column] {
		return nil, fmt.Errorf("列 '%v' 不存在", column)
	}
	delete(columnMap, column)
	columnList, columnCount := mapUtil.StringKeys(columnMap), len(columnMap)

	result := make([]mapUtil.StringObjectMap, columnCount)
	for i, str := range columnList {
		result[i] = make(mapUtil.StringObjectMap, rowCount)
		result[i][column] = str
	}
	for _, row := range rows {
		keyVal, _ := convertor.ToString(row[column])
		for idx, str := range columnList {
			if val, ok := row[str]; ok {
				result[idx][keyVal] = val
			}
		}
	}

	return result, nil
}
