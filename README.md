# go协程任务池

## 为什么要做这个项目

> go协程自身调度管理已经很完善了。本身go协程也很轻量级，遇到需要异步的任务直接`go func`就可以了。
>
> 这种粗暴的直接开启go协程的方法有一个前提，就是并发执行的任务之前没有先后依赖关系，因为它们是并发执行的，必须保证先后顺序是不影响结果的。
>
> 在同一个客户端的多次执行需要有序时，直接`go func`的方式就不行了，该项目主要解决的就是这种场景下的问题
>

## 概念介绍

### GoPool

> 协程任务池，接收任务并执行。
>
> 主要方法：
> 1. DoWorker 执行一个任务
> 2. StopAll 停止所有任务
> 3. Running 是否运行中（非停止状态即为运行中）

### GoWorker

> 协程池任务Worker
>
> 描述：在`GoPool`中，`Worker`可以有很多个，但给定相同的`id`，其执行的任务一定在同一个Worker上，且按调用顺序执行，但不同`id`不保证一定在不同`Worker`中
>
> 主要方法:
> 1. start 启动，不对外暴露。启动`worker`任务，执行器从任务队列获取任务并执行
> 2. stop 停止，不对外暴露。停止`worker`任务
> 3. Running 是否运行中。当执行`stop`方法后，变为非运行中，此时无法往`worker`里丢任务
> 4. WorkerInfo 工作节点信息
> 5. doWork 往`worker`里添加一个任务

### TaskQ

> 任务队列
>
> 描述： 作为任务节点的内置变量，主要接收任务和往`worker`发送任务
>
> 主要方法：
> 1. Add 往任务队列添加任务
> 2. Shutdown 停止接收任务
> 3. Get 获取任务
> 4. Size 任务队列长度
> 5. QLen 任务队列中任务数
> 6. Clear 清空任务队列

## 使用