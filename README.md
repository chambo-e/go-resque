# go-resque

A simple go package to enqueue jobs to a resque queue

## Install

```shell
$ go get "github.com/chambo-e/go-resque"
```

### Usage
```go

import (
    "time"
    "log"

    "github.com/chambo-e/go-resque"
)

package main

func logErr(err error) {
    if err != nil {
        log.Println(err)
    }
}

func main() {
    cfg := resque.Configuration{
        RedisURI: "redis://locahost:6379/0",
        Namespace: "test",
    }
    // or
    cfg := resque.Configuration{
        Redis: resque.RedisOptions{
            Addr: "localhost:6379",
            DB: 0,
            Password: "hello",
            DialTimeout: 100 * time.Millisecond,
        },
        Namespace: "test",
    }

    // client is goroutine safe
    cli, err := resque.New(cfg)
    if err != nil {
        log.Fatalln(err)
    }

    // Single Enqueue
    err := cli.Enqueue("test_queue", resque.NewJob("Klass", 1, 2, 3, "args"))
    logErr(err)

    // Batch Enqueue
    batch, err := cli.NewBatch("test_queue")
    if err != nil {
        log.Fatalln(err)
    }

    for x := 0; x < 100; x++ {
        batch.Enqueue(resque.NewJob("Klass", x, "iter", "so easy"))
    }

    // Execute enqueue every jobs in a MULTI to avoid uneeded redis roundtrip
    // Can greatly improve perfs for large enqueue
    err := batch.Execute()
    logErr(err)

    // Batch is cleared, you can reuse it

    for x := 0; x < 100; x++ {
        batch.Enqueue(resque.NewJob("Klass", x, "iter", "so easy"))
    }

    err := batch.Execute()
    logErr(err)
}

```
