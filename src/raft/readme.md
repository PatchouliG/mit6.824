# readme

## how to run test

```shell
cd src/raft
go test  # run all test
go test -run 2A # run test 2A
```

注意多测几次，fail可能会随机出现

## lab 2A

### goal

实现选举

实现上，除了选举外，由于心跳机制需要实现append entry，
只是只完成了log的replication，不负责committee，apply之类的逻辑

## 总体流程

每个角色（或者说状态，包括leader， follower，candidate）运作在一个死循环中
接受外部请求的chan（vote 和appendentry）和超时chan，进行处理，相应，角色变更

注：每个角色的相关公告函数以角色名开头

## lab 2B

### goal

Implement the leader and follower code to append new log entries, so that the go test -run 2B tests pass.

## lab 2c 
pass


### todo
go race 有数据竞争，但是应该没问题（pass test）,后面再解决

## lab 3c

实现幂等，每个操作包括clerkid和opid
### kv server 
每个kv server包含一个raft实例和一个kvdb

新增一个notify，
在clerk发起请求后，在notification中进行注册一个结果的chan
kvdb进行操作后，通过chan 将操作结果，操作id发送到notification，由notification将结果输出给clerk（通过注册的chan）

#### 写请求
1. 判断当前raft是否是leader
   否，返回错误，client去尝试其他
2. 通过raft start 发起写入命令 (注意加锁,包含raft的数据，但是可以有多个写on the fly)
3. raft 写入成功，通过applyMsg chan将命令下发给kvdb ,注意applymsg chan是blocked的，不处理会阻塞后面的raft请求处理
4. kvdb的常驻go routine，读取到apply msg，进行相应op，完成通过

注意，client失败会重试，
即client可能无法知道某个请求是否成功，raft log可能会重复，需要在kv层做幂等，
client给每个请求加一个递增id，kv端抛弃重复的id



####  读请求
同写请求
为了避免读到stale数据，读也以raft log的形式下发
  

### client 

#### 读写
每个请求有唯一id

### 选择leader
随机挑一个当做是leader，发送请求，失败则选择其他，如果成功，保持当前的选择
请求失败，持续重试
需要考虑请求超时，重试

### dev log

kb db
请求id去重
后面以此作为snapshot

kvserver

client



