# ��������

## 1. ys_agent
> ys_agent��������ÿ̨��Ҫ�����ص������ϣ������ռ��������ݡ������¼���

```sh
export OPEN_OPERATION_ROOT=/home/admin/open_operation_root
mkdir $OPEN_OPERATION_ROOT/ys_agent -p
# open-operationĿ¼��Ҫ��ȫ
cd open-operation/modules/ys_agent
go install
# configure.json�����ļ����ݿ��Զ���
cp ./configure.json $OPEN_OPERATION_ROOT/ys_agent/
cp $GOBIN/ys_agent $OPEN_OPERATION_ROOT/ys_agent/
# rootȨ������
nohup $OPEN_OPERATION_ROOT/ys_agent/ys_agent &
```


## 2. monitor
> monitor��������ݵĹ����������ԵĹ���agent������ά�����¼��·��ȡ���������ģ�廯�����Զ��屨�������쳣��Ϣ�뼶�ϱ���

```sh
export OPEN_OPERATION_ROOT=/home/admin/open_operation_root
mkdir $OPEN_OPERATION_ROOT/monitor -p
# open-operationĿ¼��Ҫ��ȫ
cd open-operation/modules/monitor
go install
# configure.json�����ļ����ݿ��Զ���
cp ./configure.json $OPEN_OPERATION_ROOT/monitor/
cp $GOBIN/monitor $OPEN_OPERATION_ROOT/monitor/
# ����
nohup $OPEN_OPERATION_ROOT/monitor/monitor &
```

## 3. control
> control��ǰ�˶Խӷ������������û����������ת����

```sh
export OPEN_OPERATION_ROOT=/home/admin/open_operation_root
mkdir $OPEN_OPERATION_ROOT/control -p
# open-operationĿ¼��Ҫ��ȫ
cd open-operation/modules/control
go install
# configure.json�����ļ����ݿ��Զ���
cp ./configure.json $OPEN_OPERATION_ROOT/control/
cp $GOBIN/control $OPEN_OPERATION_ROOT/control/
# ����
nohup $OPEN_OPERATION_ROOT/control/control &
```


