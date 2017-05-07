Starting bench for local server hitting static content being sent over the wire with dynamic headers

`go run $GOPATH/src/github.com/rakyll/hey/hey.go -n 10000 -c 100 -h2 https://0.0.0.0/api/categories`

23 requests done.
267 requests done.
1420 requests done.
4322 requests done.
All requests done.

Summary:
  Total:        3.3602 secs
  Slowest:      3.1118 secs
  Fastest:      0.0001 secs
  Average:      0.0330 secs
  Requests/sec: 2976.0451
  Total data:   8460300 bytes
  Size/request: 846 bytes

Status code distribution:
  [200] 10000 responses

Response time histogram:
  0.000 [1]     |
  0.311 [9899]  |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  0.622 [0]     |
  0.934 [0]     |
  1.245 [1]     |
  1.556 [3]     |
  1.867 [13]    |
  2.178 [28]    |
  2.489 [16]    |
  2.801 [15]    |
  3.112 [24]    |

Latency distribution:
  10% in 0.0009 secs
  25% in 0.0022 secs
  50% in 0.0050 secs
  75% in 0.0115 secs
  90% in 0.0259 secs
  95% in 0.0399 secs
  99% in 1.1787 secs

