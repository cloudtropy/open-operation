# 环境准备

## 1. golang
* 安装golang，版本 >= 1.8

## 2. mysql
* 安装mysql，版本 >= 5.0
* 初始化mysql数据库   

```sh
mysql > create database operation;
mysql > use operation;
mysql > source open-operation/init/mysql/init.sql
```

## 3. redis
* 安装redis
* 启动redis

## 4. rrdtool
* 安装rrdtool