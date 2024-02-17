# inmem-db
## In-Memory Database using GoLang

Simple version of a in memory key value store using the RESP (REdis Serialization Protocol).

Following commands are to be implemented:

1. GET
2. DEL
3. EXPIRE
4. KEYS
5. SET
6. TTL
7. ZADD
8. ZRANGE
 
The implementation should follow the redis command standard. For example, for SET, it
is at: https://redis.io/commands/set


## Overall Functionality
* The server listens for incoming connections and handles commands from clients.
* The client connects to the server and sends commands for operations like getting, setting, deleting keys, setting expiration, and retrieving keys.
* The server processes the commands and sends back appropriate responses.

### To implement RESP Protocol:
    1. Used bufio package to create a reader and writer instance to parse the string starting with "$" and add the values and end the statement with "\r\n"

### Create a database:
    1. Create a simple database with interface and expiry to handle TTL and EXPIRE and sortedSet for ZADD


#### v0 -> Implemented basic SET, GET, EXPIRE, KEYS
