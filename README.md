# haproxy-go
haproxy-go is a free, very fast and reliable reverse-proxy offering high availability, load balancing, and proxying for TCP and UDP applications.

configs.json
[
  {
    "inport": 8080,
    "outport": 8080,
    "server": "1.1.1.1"
  },
  {
    "inport": 8081,
    "outport": 8081,
    "server": "1.1.1.1"
  },
  {
    "inport": 8082,
    "outport": 8082,
    "server": "1.1.1.1"
  }
]
