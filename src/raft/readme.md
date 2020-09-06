# dev log

## how to run test

```shell
cd src/raft
go test  # run all test
go test -run 2A # run test 2A
```

## lab 2A
lab 2a git tag （todo)

实现上，除了选举外，由于心跳机制需要实现append entry，
只是只完成了log的replication，不负责committee，apply之类的逻辑

## 总体流程

每个角色（或者说状态，包括leader， follower，candidate）运作在一个死循环中
接受外部请求的chan（vote 和appendentry）和超时chan，进行处理，相应，角色变更

注：每个角色的相关公告函数以角色名开头

## lab 2B

todo 增加选举限制


