# go-redis-streams

Demo using messaging with Redis Streams in Golang

**Using messaging in practice with Redis streams and Golang**

# Redis

**What is?**

Redis is an OpenSource non-relational database whose structure includes key-value storage.
Redis has strategies for storing data in memory and on disc, guaranteeing fast response and data persistence. The main use cases for Redis include caching, session management, PUB/SUB.

# Redis Streams for Messaging

![Design of flow](/assets/redis-streams-flow.png)

**Positive Points**

- Supports Topics and Queues
- Persistence on disc (via RDB files)
- High availability (with Clustering)
- High Throughput
- Allows reprocessing
- Consumer Groups
- Minimal latency
- No need for zookeper
- Takes up far fewer resources than Kafka/RabbitMQ

**Negative Points**

- No guaranteed delivery order (yet)
- Processed messages with error not returned for redistribution

# Links

https://redis.io/topics/streams-intro

https://redislabs.com/blog/use-redis-streams-apps/

# Quickstart

Start dockerized execution

```
make startdocker
```

Stop dockerized execution

```
make stopdocker
```

Download dependencies

```
go mod tidy
```

Build project

```
make build
```

Start consumer

```
make consumer
```

Start producer

```
make producer
```
