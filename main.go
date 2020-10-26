package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	user   = flag.String("user", "root", "db user")
	port   = flag.Int("port", 4000, "db port")
	passwd = flag.String("passwd", "", "db password")
	host   = flag.String("host", "127.0.0.1", "db host")
	dbName = flag.String("db", "zenos", "db name")
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

type TransactionExector struct {
	db *sql.DB
	tx *sql.Tx
}

func NewTransactionExector(db *sql.DB) *TransactionExector {
	return &TransactionExector{
		db: db,
	}
}

func (t *TransactionExector) Begin() error {
	var err error
	t.tx, err = t.db.Begin()
	return err
}
func (t *TransactionExector) Exec(sql string) error {
	_, err := t.tx.Exec(sql)
	return err
}
func (t *TransactionExector) Commit() error {
	return t.tx.Commit()
}
func (t *TransactionExector) Rollback() error {
	return t.tx.Rollback()
}
func GetDB(addr string) *sql.DB {
	db, err := sql.Open("mysql", addr)
	if err != nil {
		return nil
	}
	return db
}
func ExecSql(tx *TransactionExector, sql string) error {
	err := tx.Exec(sql)
	return err
}
func main() {

	sqls1 := ReadSQLFile("data/1.sql")
	sqls2 := ReadSQLFile("data/2.sql")
	firstSeq := GenerateFirstSeq(len(sqls1), len(sqls2))
	seq := make([]byte, len(firstSeq))
	dbSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", *user, *passwd, *host, *port, *dbName)
	db1 := GetDB(dbSource)
	db2 := GetDB(dbSource)
	defer db1.Close()
	defer db2.Close()
	copy(seq, firstSeq)
	for {
		fmt.Printf("current seq=%s\n", string(seq))
		tx1 := NewTransactionExector(db1)
		tx2 := NewTransactionExector(db2)
		tx1.Begin()
		tx2.Begin()
		i := 0
		j := 0
		for _, item := range seq {
			//todo 可能会死锁
			var tx *TransactionExector
			sql := ""
			if item == '1' {
				tx = tx1
				sql = sqls1[i]
				i++
			} else {
				tx = tx2
				sql = sqls2[j]
				j++
			}
			tx.Exec(sql)
		}
		tx1.Commit()
		tx2.Commit()

		seq = NextSeqs(seq)
		if string(seq) == string(firstSeq) {
			break
		}
	}
}
