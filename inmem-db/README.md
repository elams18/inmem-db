# inmem-db
## In-Memory Database using GoLang

Simple version of an in memory key value store using the RESP (REdis Serialization Protocol).

Following commands are to be implemented:

1. GET - DONE
2. DEL - DONE
3. EXPIRE - DONE
4. KEYS - DONE 
5. SET - DONE
6. TTL - DONE
7. ZADD - DONE
8. ZRANGE - DONE
 
The implementation should follow the redis command standard. For example, for SET, it
is at: https://redis.io/commands/set


## Overall Functionality
* The server listens for incoming connections and handles commands from clients.
* The client connects to the server and sends commands for operations like getting, setting, deleting keys, setting expiration, and retrieving keys.
* The server processes the commands and sends back appropriate responses.

#### Concerns -
* Added -1 marking as End of response which needs to be taken care of.
* Multiple client trying to modify the data - **Used Mutex Lock to take care of this**
* For Complex data structures, this may not work as expected.

#### v0 -> Implemented basic SET, GET, EXPIRE, KEYS, TTL, ZADD, ZRANGE
