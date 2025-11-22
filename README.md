# Go-Redis-Clone

This project is a **multi-threaded**, in-memory data store built with **Go**, designed to mimic the functionality of Redis. While it shares many similarities with Redis, the key differentiator is its multi-threaded architecture, leveraging Go's concurrency model.

## Features implemented
- **In-Memory Data Storage**: Fast key-value store.
- **Multi-Threading**: Handles multiple client connections concurrently using Go routines.
- **Persistence**: Implements AOF (Append Only File) persistence to ensure data durability across restarts.
- **RESP Protocol**: Speaks the Redis Serialization Protocol, making it compatible with standard Redis clients (like `redis-cli`).

## Supported Commands
The following Redis commands are currently supported:

*   **Basic**: `PING`, `QUIT`, `COMMAND`
*   **String Operations**: `SET`, `GET`, `SETNX`, `MSET`, `MGET`, `INCR`, `DECR`
*   **Key Management**: `DEL`, `KEYS`, `RENAME`
*   **Database**: `SELECT`, `FLUSHDB`, `FLUSHALL`

## Future Roadmap
I am actively working on expanding the capabilities of this project. Here are the things I'm most interested in implementing next:

*   **Key Expiry**: Implementing TTL (Time To Live) for keys.
*   **Redis Streams**: Redis's append-only log
*   **More Data Structures**: Lists, Sets, Hashes, and Sorted Sets.
*   **Vector Database**: A stretch goal to explore vector similarity search and embeddings.
