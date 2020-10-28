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
	tx, err := db.Begin() //事务开始
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

	InitDBData()                                         // 执行DB初始化语句
	sqls1 := ReadSQLFile(*sqlfile1)                      // 读取事务一的sql语句列表
	sqls2 := ReadSQLFile(*sqlfile2)                      // 读取事务二的sql语句列表
	firstSeq := GenerateFirstSeq(len(sqls1), len(sqls2)) // 生成第一个 两个事务语句执行的序列
	dbSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", *user, *passwd, *host, *port, *dbName)
	db1 := GetDB(dbSource) // 获取事务一的DB连接
	db2 := GetDB(dbSource) // 获取事务二的DB连接
	defer db1.Close()
	defer db2.Close()
	curSeq := make([]byte, len(firstSeq)) // 当前两个事务执行的sql序列
	copy(curSeq, firstSeq)
	for {
		fmt.Printf("current seques=%s\n", string(curSeq))
		chSqls1 := make(chan string)
		chSqls2 := make(chan string)
		ch1 := ExecSqls(db1, chSqls1) // 开始事务一
		ch2 := ExecSqls(db2, chSqls2) // 开始事务二
		i := 0
		j := 0
		for _, item := range curSeq {
			//todo 可能会死锁

			sql := ""
			var ch chan error
			if item == '1' { // 当前为事务一的SQL语句
				sql = sqls1[i]
				chSqls1 <- sql // 将事务一执行的SQL语句发送给事务一执行
				ch = ch1
				i++
			} else { // 当前为事务二的SQL语句
				sql = sqls2[j]
				chSqls2 <- sql // 将事务二执行的SQL语句发送给事务一执行
				ch = ch2
				j++
			}

			select {
			case <-time.After(1 * time.Second): // 事务一或二SQL语句执行超时
				panic("timeout")
			case err := <-ch:
				if err != nil { // 事务一或二 某个SQL语句执行失败
					panic(err)
				}
			}
		}

		//todo 保存当前结果
		SaveDBData(db1)
		//todo 重新初始化DB数据
		InitDBData()
		curSeq = NextSeqs(curSeq)               // 获取下一个事务SQL语句执行顺序
		if string(curSeq) == string(firstSeq) { // 所有SQL执行顺序全部执行
			break
		}
	}
}
