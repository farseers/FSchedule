# CLAUDE.md

## 1、项目概览
调度中心，是一款支持分布式定时任务调度，支持几千个任务同时调度。有着高性能、高可用（集群）、低延迟（0ms延迟）、低内存（45mb以内）等特点。

## 2、运行原理
多个客户端实例请求同一个服务端不同实例时，服务端只会有一个实例拿到redis锁并一直与该客户端保持通讯，直到任意一方断开连接。

其他客户端同样的任务组的请求会进入等待锁。（仅保持连接但不再发送消息）

./domain/monitorTaskGroupService.go文件的MonitorTaskGroupPush负责接收新的客户端请求，并决定是否拿到锁后开始负责这次请求。

## 任务组
任务组定义了下一次的执行时间。

每个任务组taskGroup.DomainObject对象，由monitorTaskGroupService.go负责跟踪。

拿到锁的实例负责同步这个对象到Redis

## 客户端注册
首次注册进来通过wss调用./application/ws/Connect -> ./domain/MonitorTaskGroupPush

注册后会与服务端保持长连接模式，通过wss协议来收发消息，不需要每次都重新访问ws.Connect。