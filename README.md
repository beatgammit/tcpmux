tcpmux
======

[![Build Status](https://drone.io/github.com/beatgammit/tcpmux/status.png)](https://drone.io/github.com/beatgammit/tcpmux/latest)

`tcpmux` is a TCP multiplexer and enables running multiple services on the same port by sniffing the first few bytes of incoming connections and forwarding them onto other services.

Once common use-case is running both an HTTP server and an SSH server on the same port.

Benchmarks
----------

Here are the benchmark results on my machine:

	$ go test -bench .
	PASS
	BenchmarkNoMuxParallel     20000            70020 ns/op
	BenchmarkMuxParallel        5000            207293 ns/op
	BenchmarkNoMuxSequential   20000            76017 ns/op
	BenchmarkMuxSequential     10000            201911 ns/op
	ok      tcpmux  7.518s

Example
-------

    go get github.com/beatgammit/tcpmux/examples/http_ssh
	$GOPATH/bin/http_ssh -addr :8080
	curl localhost:8080
	ssh -p 8080 localhost

License
=======

`tcpmux` is licensed under the BSD 3-clause license. See LICENSE.BSD3 for details.
