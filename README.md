## 使用说明

### 安装tidb
```shell
./install.sh
```
### 编译
```
go build
```
### client使用
```shell
./trans_client --user root --port 4000 --host 127.0.0.1 --db zenos --sql1=data/1.sql --sql2=data/2.sql --init_sql=data/init.sql
```
