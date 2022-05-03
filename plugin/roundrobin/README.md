# RoundRobin

## Name
*roundrobin* - plugin provides several round-robin strategies that mix A and AAAA 
records in a DNS response.

## Description
The roundrobin plugin implements [round-robin](https://en.wikipedia.org/wiki/Round-robin_scheduling)
strategies on returned A and AAAA records. 
```
loadbalance [STRATEGY]
```
`[STRATEGY]` defines how A and AAAA will be shuffled. If there are records other than A or AAAA in the
answer, they are moved to the end of the answer section in the exact order as they were in the answer 
section before the roundrobin applied. 
- [stateful](#stateful)
- [stateless](#stateless)
- [random](#random)

## RoundRobin
### stateful
```
. {
    log
    roundrobin stateful
}
```
Stateful will probably be the strategy you expect from a consistent round-robin. Each hit rotates the A or AAAA records 
by one position. The RoundRobin plugin remembers the positions of the last query for ten minutes and manages changes 
to the answers: clears non-existing records and adds new ones, shuffling (rotation by one position) and garbage collection.
_NOTE: key to the request state is the pair `EDNS0_SUBNET` and request Question, [see more](https://en.wikipedia.org/wiki/EDNS_Client_Subnet)._

#### Usage


### stateless
```
. {
    log
    roundrobin stateless
}
```
Stateless is useful where you require extremely high scalability and performance, or you cannot use stateful. The state 
is stored on the CoreDNS client and like in HTTP, you send the records back in the DNS query having the section Extra 
with the OPT field `EDNS0_LOCAL`. The `stateless` plugin takes care of shuffling (rotation by one position), 
clears non-existing records and adds new ones. As in HTTP, the client must store the response in its memory for the next request.

#### Usage
```go
// stateless-client.go

```

### random
```
. {
    log
    roundrobin random
}
```
Random is the simplest strategy available. It is simple, fast, and does not work with any state. On the other hand, 
it does not offer consistent results. Responses may be repeated, and you have no guarantee that the returned IP 
addresses will be provided equally. 

#### Usage
