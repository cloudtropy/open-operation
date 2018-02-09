# 单机部署

## 1. ys_agent
> ys_agent，部署在每台需要管理监控的主机上，负责收集主机数据、处理事件。

```sh
export OPEN_OPERATION_ROOT=/home/admin/open_operation_root
mkdir $OPEN_OPERATION_ROOT/ys_agent -p
# open-operation目录需要补全
cd open-operation/modules/ys_agent
go install
# configure.json配置文件内容可自定义
cp ./configure.json $OPEN_OPERATION_ROOT/ys_agent/
cp $GOBIN/ys_agent $OPEN_OPERATION_ROOT/ys_agent/
# root权限启动
nohup $OPEN_OPERATION_ROOT/ys_agent/ys_agent &
```


## 2. monitor
> monitor，监控数据的管理、报警策略的管理、agent心跳的维护，事件下发等。监控项组合模板化管理，自定义报警规则，异常信息秒级上报。

```sh
export OPEN_OPERATION_ROOT=/home/admin/open_operation_root
mkdir $OPEN_OPERATION_ROOT/monitor -p
# open-operation目录需要补全
cd open-operation/modules/monitor
go install
# configure.json配置文件内容可自定义
cp ./configure.json $OPEN_OPERATION_ROOT/monitor/
cp $GOBIN/monitor $OPEN_OPERATION_ROOT/monitor/
# 启动
nohup $OPEN_OPERATION_ROOT/monitor/monitor &
```

## 3. control
> control，前端对接服务器，负责用户管理，请求的转发。

```sh
export OPEN_OPERATION_ROOT=/home/admin/open_operation_root
mkdir $OPEN_OPERATION_ROOT/control -p
# open-operation目录需要补全
cd open-operation/modules/control
go install
# configure.json配置文件内容可自定义
cp ./configure.json $OPEN_OPERATION_ROOT/control/
cp $GOBIN/control $OPEN_OPERATION_ROOT/control/
# 启动
nohup $OPEN_OPERATION_ROOT/control/control &
```


