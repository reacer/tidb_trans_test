package main

import (
	"bufio"
	"database/sql"
	"os"
)

// ReadSQLFile 从文件中获取sql语句
func ReadSQLFile(name string) []string {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(f)
	sqls := make([]string, 0)
	for {
		if sql, _, err := reader.ReadLine(); err == nil {
			sqls = append(sqls, string(sql))
		} else {
			break
		}
	}
	return sqls
}

// GetDB 获取DB连接
func GetDB(addr string) *sql.DB {
	db, err := sql.Open("mysql", addr)
	if err != nil {
		return nil
	}
	return db
}

/* GenerateFirstSeq
*  构造第一个序列 m=3,n=3,返回 "111222"
 */
func GenerateFirstSeq(m, n int) []byte {
	ans := make([]byte, m+n)
	for i := 0; i < m; i++ {
		ans[i] = '1'
	}
	for i := m; i < m+n; i++ {
		ans[i] = '2'
	}
	return ans
}

/* reverse: 字符串 索引i开始到最后逆序
*  reverse([]byte("abce",1)) 字符串为 aecb
 */

func reverse(data []byte, i int) {
	j := len(data) - 1
	for i < j {
		data[i], data[j] = data[j], data[i]
		i++
		j--
	}
}

/* NextSeqs 返回下一个字典序列比data的大的字符串,如果没有，则返回最小的
* NextSeqs([]byte("111222")) => "112122"
* NextSeqs([]byte("222111")) => "111222"
 */
func NextSeqs(data []byte) []byte {
	if len(data) <= 1 {
		return data
	}
	i := len(data) - 2
	for ; i >= 0 && data[i] >= data[i+1]; i-- { // 从后面向前搜索第一个逆序的字符
	}
	if i >= 0 {
		j := len(data) - 1
		for ; j >= 0 && data[j] <= data[i]; j-- { // 从后面向前搜索第一个大于 data[i]的字符
		}
		data[i], data[j] = data[j], data[i]
	}
	reverse(data, i+1)
	return data
}
