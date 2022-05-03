# RoundRobin

## Name
*roundrobin* - plugin provides several round-robin strategies that mix A and AAAA 
records in a DNS response.

## Description
The roundrobin plugin implements [round-robin](https://en.wikipedia.org/wiki/Round-robin_scheduling)
strategies on returned A and AAAA records. 
```
roundrobin [STRATEGY]
```
`[STRATEGY]` defines how A and AAAA will be shuffled. If there are records other than A or AAAA in the
answer, they are moved to the end of the answer section in the exact order as they were in the answer 
section before the roundrobin applied. 
- [stateful](#stateful)
- [stateless](#stateless)
- [random](#random)

## RoundRobin
### stateful
Stateful will probably be the strategy you expect from a round-robin. Each hit rotates the A or AAAA records 
by one position. The RoundRobin plugin remembers the positions of the last query for ten minutes and manages changes 
to the answers: clears non-existing records, adds new ones, shuffling (rotation by one position) and garbage collection.
_NOTE: key to the request state is the pair `EDNS0_SUBNET` and request Question, [see more](https://en.wikipedia.org/wiki/EDNS_Client_Subnet)._

#### Usage
```
# Corefile 
.:5053 {
    log
    roundrobin stateful
}
```
After you run the dig command repeatedly, the records for the `myhosts.com` domain rotate. The roundrobin plugin can 
handle the number of records or individual records changing. In addition, garbage collector clears the records every ten 
minutes of inactivity.

```shell
# runnig against a local hosts plugin `hosts etchosts` 

dig @localhost -p 5053 myhost.com                                 dig @localhost -p 5053 myhost.com                        
myhost.com.             3600    IN      A       200.0.0.2         myhost.com.             3600    IN      A       200.0.0.3
myhost.com.             3600    IN      A       200.0.0.3         myhost.com.             3600    IN      A       200.0.0.4
myhost.com.             3600    IN      A       200.0.0.4         myhost.com.             3600    IN      A       200.0.0.1
myhost.com.             3600    IN      A       200.0.0.1         myhost.com.             3600    IN      A       200.0.0.2

dig @localhost -p 5053 myhost.com
myhost.com.             3600    IN      A       200.0.0.4
myhost.com.             3600    IN      A       200.0.0.1
myhost.com.             3600    IN      A       200.0.0.2
myhost.com.             3600    IN      A       200.0.0.3
```

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
clears non-existing records and adds new ones. As in HTTP, the client must store the response in its memory for the next 
request. Sending client data is the state in the form of text encoded as a byte-array where `_rr_state=` is a constant 
identifying the state followed by the json containing the state, see: `_rr_state={"ip":["10.0.0.1","10.2.2.1","10.1.1.2"]}` 

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
