# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 1、项目概览
调度中心，是一款支持分布式定时任务调度，支持几千个任务同时调度。有着高性能、高可用（集群）、低延迟（0ms延迟）、低内存（45mb以内）等特点。

FSchedule运行的是调度逻辑，到达时间节点时，会通知你的应用程序执行任务。

它有分布式、高可用、解耦任务执行与调度逻辑、弹性伸缩、跨语言、一致性、数据分片执行、集群模式、广播模式、分布式日志等特性。

Server端依赖Redis、Mysql，通常使用Docker运行（支持集群HA），部署完后我们一般不需要再去维护。

客户端，就是我们写的应用程序，通过官方SDK（farseer-go/fschedule组件）集成到你的应用，比如我需要凌晨2点执行数据库清理操作，在接入后就能按计划执行任务了。
