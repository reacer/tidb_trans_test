package main

import (
	"database/sql"
	"flag"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	user     = flag.String("user", "root", "db user")
	port     = flag.Int("port", 4000, "db port")
	passwd   = flag.String("passwd", "", "db password")
	host     = flag.String("host", "127.0.0.1", "db host")
	dbName   = flag.String("db", "zenos", "db name")
	sqlfile1 = flag.String("sql1", "data/1.sql", "the first sql file")
	sqlfile2 = flag.String("sql2", "data/2.sql", "the second sql file")
	initsql  = flag.String("init_sql", "data/init.sql", "init database")
)

// ExecSqls 执行sqls
func ExecSqls(db *sql.DB, sqls chan string) chan error {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	ans := make(chan error)
	go func() {
		for sql := range sqls {
			fmt.Printf("exec sql=%s\n", sql)
			_, err := tx.Exec(sql)
			if err != nil {
				tx.Rollback()
			}
			ans <- err
		}
		err := tx.Commit()
		if err != nil {
			panic(err)
		}
	}()
	return ans

}

// InitDBData 初始化DB数据 todo
func InitDBData() error {
	dbSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", *user, *passwd, *host, *port, *dbName)
	db := GetDB(dbSource)
	if db == nil {
		panic("GetDBSource failed")
	}
	defer db.Close()
	sqls := ReadSQLFile(*initsql)

	for _, sql := range sqls {
		_, err := db.Exec(sql)
		if err != nil {

			panic(err)
		}
	}
	return nil
}

// SaveDBData 保存数据 todo
func SaveDBData(db *sql.DB) {}
func main() {

	InitDBData()
	sqls1 := ReadSQLFile(*sqlfile1)
	sqls2 := ReadSQLFile(*sqlfile2)
	firstSeq := GenerateFirstSeq(len(sqls1), len(sqls2))
	seq := make([]byte, len(firstSeq))
	dbSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", *user, *passwd, *host, *port, *dbName)
	db1 := GetDB(dbSource)
	db2 := GetDB(dbSource)
	defer db1.Close()
	defer db2.Close()
	copy(seq, firstSeq)
	for {
		fmt.Printf("current seques=%v\n", seq)
		chSqls1 := make(chan string)
		chSqls2 := make(chan string)
		ch1 := ExecSqls(db1, chSqls1)
		ch2 := ExecSqls(db2, chSqls2)
		i := 0
		j := 0
		for _, item := range seq {
			//todo 可能会死锁

			sql := ""
			var ch chan error
			if item == '1' {
				sql = sqls1[i]
				chSqls1 <- sql
				ch = ch1
				i++
			} else {
				sql = sqls2[j]
				chSqls2 <- sql
				ch = ch2
				j++
			}

			select {
			case <-time.After(1 * time.Second):
				panic("timeout")
			case err := <-ch:
				if err != nil {
					panic(err)
				}
			}
		}

		//todo save current result
		SaveDBData(db1)
		//todo reset db data
		InitDBData()
		seq = NextSeqs(seq)
		if string(seq) == string(firstSeq) {
			break
		}
	}
}
