# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 1、项目概览
调度中心，是一款支持分布式定时任务调度，支持几千个任务同时调度。有着高性能、高可用（集群）、低延迟（0ms延迟）、低内存（45mb以内）等特点。

FSchedule运行的是调度逻辑，到达时间节点时，会通知你的应用程序执行任务。

## 2、运行原理
这是一个高可用的解决方案。即服务端（FSchedule）会运行多个实例，客户端也会多个实例。

多个客户端实例请求同一个任务组到服务端不同实例时，服务端只会有一个实例拿到分布式锁并一直与该客户端保持通讯，直到任意一方断开连接。

而其他客户端的请求会进入等待锁。（仅保持连接但不再发送消息）

./domain/monitorTaskGroupService.go文件的MonitorTaskGroupPush就是负责接收新的客户端请求，并决定是否拿到锁后开始负责这次请求。

## 任务组
任务组定义了下一次的执行时间。

每个任务组，即taskGroup.DomainObject对象，由monitorTaskGroupService.go负责跟踪。

拿到锁的实例负责同步这个对象到Redis

## 客户端注册
客户端首次注册进来时，会通过wss调用./application/ws/Connect -> ./domain/MonitorTaskGroupPush 函数

注册后会与服务端保持长连接模式，直接通过wss协议来发送或接收消息，而不需要每次都重新访问ws.Connect函数。