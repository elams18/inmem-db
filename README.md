# inmem-db
## In-Memory Database using GoLang

Simple version of a in memory key value store using the RESP (REdis Serialization Protocol).

Following commands are to be implemented:

1. GET - DONE
2. DEL - DONE
3. EXPIRE - DONE
4. KEYS - DONE ( has a bug)
5. SET - DONE
6. TTL - DONE
7. ZADD - WIP ( has some parameter issues )
8. ZRANGE - WIP ( has a bug )
 
The implementation should follow the redis command standard. For example, for SET, it
is at: https://redis.io/commands/set


## Overall Functionality
* The server listens for incoming connections and handles commands from clients.
* The client connects to the server and sends commands for operations like getting, setting, deleting keys, setting expiration, and retrieving keys.
* The server processes the commands and sends back appropriate responses.

#### Concenns -
* 

#### v0 -> Implemented basic SET, GET, EXPIRE, KEYS, TTL, ZADD, ZRANGE
