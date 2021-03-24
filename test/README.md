
### Scenarios

We use docker containers to build scenarios and environments during development
to help building and testing new features - getting lost in too many logs and
manual terminal work is time consuming.

We probably need to run the tests within a docker container that is attached to
the same network as the cluster containers - so we can use the same service
names etc and call the RPC server for state.

There is a SendState RPC method in the node's RPC server, this will return us
a JSON of the current state of a node. This can be used in automating scenarios
by for instance:


#### Leader election scenario
- Starting the cluster
- Killing the leader
- Checking the state of both remaining containers until a new leader is elected
- pass the state if a new leader is elected
- fail if there is no leader elected in a certain time limit


When we add a new feature (next we need to add log replication) we can make
sure to run the various scenarios before committing our changes, this will help
be more confident that nothing has broken in the previous functionality.
