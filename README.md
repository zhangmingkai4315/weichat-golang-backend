# weichat-golang-backend
使用golang作为后端开发微信公众号

1.安装数据库 用于保存worker获取的数据

```
$ mysql -uyour_database_user -p < install.sql

Enter password: 
Field   Type    Null    Key     Default Extra
id      int(10) unsigned        NO      PRI     NULL    auto_increment
uuid    varchar(32)     NO      UNI     NULL    
title   varchar(256)    NO              NULL    
link    varchar(256)    NO              NULL    
post_date       timestamp       NO              CURRENT_TIMESTAMP       
score   varchar(64)     NO              NULL    
user_name       varchar(32)     NO              NULL    
user_profile    varchar(128)    NO              NULL    
md5     varchar(32)     NO      MUL     NULL    
...
```

2.执行测试

```
cd project
go test ./... -v -bench=. 

```

