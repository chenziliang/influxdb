## Summary

When session expiration, it means the Master Node has lost connection with ETCD for some time(exceeds session TTL). The local metadata may not up-to-date with ETCD.

## Master Node
- Stop the control logic.
- Try to connect to the ETCD just as the data node join cluster for the first time.
- During the connection loss period, the master node should reject all requests directly(requests from client and other data nodes).  

## Connection Loss Period
### option1
Through some flags to tell the http server and grpc server to reject requests. 
This option can reduce restart time but will introduce complexity. For each request, it will check the related flags.

### option2
Shutdown the http server and grpc server. Restart them when data node connect to the ETCD again. 
This option may have a long recovery time.

## etcd client investigation
Please take a look at the pod part.



