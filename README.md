## 限流

### 限流场景
在系统中，由于高并发大流量的场景，如秒杀、抢购、热点、刷流、爬虫等导致系统资源不足，负载过高，导致服务出现问题。

### 限流算法
- 固定窗口限流(Fixed Window)

```
常用的为计数器算法(Counter Limiter)
设计限流对象，如将请求url+用户id/商品id等作为限流对象，对该对象的访问进行全局计数，当在统计窗口期内达到阈值则限流
对于计数的存储，单机的话使用全局计数器即可，分布式系统可以使用redis等存储

优点：原理和实现简单
缺点：限流不均匀，临界区的突发流量无法控制
```

- 滑动窗口限流(Sliding Window)

```
滑动窗口是固定窗口的一个改进算法，将一个大的时间窗口分为多个小窗口，但从根本上并没有解决固定窗口算法的临界突发流量问题

优点：滑动窗口限流比固定窗口限流更平滑
缺点：还是会存在临界区突刺流量问题
```

- 漏桶算法(Leaky Bucket)

![leaky-bucket](https://github.com/bluesaka/limiter/blob/master/file/leaky-bucket.png)

```
算法内部维护一个漏桶容器，当新请求到来时，尝试加水，若容器未满则加水并处理请求，若满了则拒绝请求

优点：控制流量速率，使流量均匀，没有了流量突刺问题
缺点：无法应对流量突发问题，大量突发请求需排队等待，或者被直接拒绝
```

- 令牌桶算法(Token Bucket)

![token-bucket](https://github.com/bluesaka/limiter/blob/master/file/token-bucket.jpg)

```
算法以恒定速率往桶里放入令牌，请求到来时从桶里取令牌进行处理，若没有令牌可取时则拒绝请求

优点：控制流量的平均速率，没有流量突刺问题，能应对部分流量突发情况
缺点：能应对部分流量突发情况，但某些场景的流量突发还是会被拒绝
```

### 流量控制

- 阿里的Sentinel
    > https://github.com/alibaba/sentinel-golang
    >
    > https://github.com/alibaba/Sentinel/wiki/%E7%B3%BB%E7%BB%9F%E8%87%AA%E9%80%82%E5%BA%94%E9%99%90%E6%B5%81
    
- 奈飞的Hystrix
  
    > https://github.com/afex/hystrix-go

