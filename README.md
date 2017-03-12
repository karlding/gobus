# GoBus

GoBus is a super-simple **Go**lang implementation of a Publish/Subscribe Message **Bus** with Topics (using async callbacks).

## Interface

This implementation supports the following operations:

* `Publish(topic, args...)`: publish to a topic
* `Subscribe(topic, callback)`: subscribe to a topic, and run an async callback
* `Unsubscribe(topic, callback)`: unsubscribe from a topic
* `Wait`: blocks until all async callbacks are run

## Example Usage

```go
import github.com/karlding/gobus

gobus := gobus.New()

addition := func(a int, b int) {
  fmt.Printf("%d + %d = %d\n", a, b, a + b)
}

multiplication := func(a int, b int) {
  fmt.Printf("%d * %d = %d\n", a, b, a * b)
}

// create new subscribers on a topic
gobus.Subscribe("math", addition)
gobus.Subscribe("math", multiplication)

// publish on the bus
// both addition and multiplication will run
gobus.Publish("math", 1, 2)

// wait for any async callbacks to finish
gobus.Wait()

// unsubscribe from the topic
gobus.Unsubscribe("math", addition)

// only multiplication will run
gobus.Publish("math", 1, 2)

gobus.Wait()
```
