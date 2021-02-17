# NaiveDistributedMemCache

## Keywords
* LRU
* Consistent Hash
* Cache breakdown

## Test
1. `go build -o cache`
2. `./cache -port=8081` `./cache -port=8082` `./cache -port=8083`
3. `curl http://localhost:8081/score/Sam & curl http://localhost:8081/score/Sam & curl http://localhost:8081/score/Sam & `

The cache on 8081 will try to load value from 8083. Three concurrent client requests will merge into one request from 8081 to 8083.