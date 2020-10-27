package main

import (
	"bufio"
	"database/sql"
	"os"
)

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
func GetDB(addr string) *sql.DB {
	db, err := sql.Open("mysql", addr)
	if err != nil {
		return nil
	}
	return db
}

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
func reverse(data []byte, i int) {
	j := len(data) - 1
	for i < j {
		data[i], data[j] = data[j], data[i]
		i++
		j--
	}
}

func NextSeqs(data []byte) []byte {
	if len(data) <= 1 {
		return data
	}
	i := len(data) - 2
	for ; i >= 0 && data[i] >= data[i+1]; i-- {
	}
	if i >= 0 {
		j := len(data) - 1
		for ; j >= 0 && data[j] <= data[i]; j-- {
		}
		data[i], data[j] = data[j], data[i]
	}
	reverse(data, i+1)
	return data
}
