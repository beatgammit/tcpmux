tcpmux
======

[![Build Status](https://drone.io/github.com/beatgammit/tcpmux/status.png)](https://drone.io/github.com/beatgammit/tcpmux/latest)

`tcpmux` is a TCP multiplexer and enables running multiple services on the same port by sniffing the first few bytes of incoming connections and forwarding them onto other services.

Once common use-case is running both an HTTP server and an SSH server on the same port.

Note
----

This approach is slower than a direct connection to the given service because all data from the socket needs to be piped to the other service. Please prefer a direct connection if at all possible.

LICENSE
=======

`tcpmux` is licensed under the BSD 3-clause license. See LICENSE.BSD3 for details.
