# GoFuse

A simple circuit breaker for Go.

## Usage

**Development is still very active and everything below might change over
night.**

To use GoFuse in your application, after importing it where needed of course,
you will need to provide a couple of things:

```go
func Action(in []byte, out chan []byte) {
	// A query function that describes how we pass a query to the service
	// we are trying to contact and pushes the answer back into a channel.
}

type LogWriter struct {
	// An IO.Writer that is used for logging.
}

func (l LogWriter) Write(p []byte) (int, err) {
	// To implement IO.Writer, we need this write function.
}
```

After providing the things above, you can go ahead and create a new fuse:

```go
// NewFuse takes:
// - the action function to use
// - the logwriter
// - the size of the queue, which controls how many concurrent attempts are
//   allowed
// - the timeout to use, a time.Duration
// - after how many failed attempts the fuse blows
// - the time interval to use when trying to reestablish contact, again a
//   time.Duration
// - how many successful contacts we want to consider the connection stable
//   again
f := fuse.NewFuse(Action, writer, 1, time.Second, 3, 2 * time.Second, 5)
```

Then you can use the fuse like this:

```go
// prepare a channel to return the data on
retchan := make(chan []byte)

// start the query
go f.Query([]byte("This is some data"), retchan)

// see what happens
select {
case data := <- retchan:
	// Yay, we can evaluate data.
case <-f.timeout:
	// Host is down or unresponsive, let's call this a failure and move on.
}
```

## How it works

The concept is quite simple and not really original. We proxy a query through
the fuse, and if there is no or a really slow response, the fuse will
eventually blow, and then returning the failure status immediately until the
host we try to contact is reasonable stable again.

## Dependencies

- A not too old version of Go

