Poële
=====

Poële is a very simple queue processor for Redis
written in Go.


Usage
-----

```go
package main

import (
	"github.com/thesides/poele"
	"github.com/vmihailenco/redis"
)

func main() {
	client := redis.NewTCPClient(host, password, -1)
	defer client.Close()
	
	pan := poele.New(client, "channel", process)
	pan.Serve()
}

func process(message string) interface{} {
	return someLongComputation()
}
```

Now, just __LPUSH__ some messages to the `poele.channel.todo`
list in Redis and they will get processed by the `process()`
function above. While they are being processed, they are
stored in the `poele.channel.processing` list, and when each
is completed it is moved to the `poele.channel.done` list.

It is recommended to just pass in a database ID and retrieve
and update the records from the Go `process()` function to
keep the Redis queues as clean and slim as possible.