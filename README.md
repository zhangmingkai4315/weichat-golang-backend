# weichat-golang-backend
使用golang作为后端开发微信公众号

1.安装数据库 用于保存worker获取的数据

```
$ sudo apt-get install mysql-server mysql-client

mysql>create database SecurityNews;
mysql>GRANT ALL PRIVILEGES ON SecurityNews.* TO weixin@"%" IDENTIFIED BY '123456' WITH GRANT OPTION;
mysql>FLUSH PRIVILEGES;


$ mysql -uyour_database_user -p < install.sql

Enter password: 
Field   Type    Null    Key     Default Extra
...
```

2.执行测试

```
cd project
go test ./... -v -bench=. 

```

