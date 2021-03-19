### Randomization of leader election initialization 

We need to randomize the time between the cutoff timeout from hearing the
last heartbeat to the node initializing a leader election (if a node
initializes a leader election it automatically votes for itself).

Q: When being requested for votes, does the node just vote for the first
  voterequest it receives? 

  A requestVote call from a candidate to each follower also contains
  information about the candidate's log. A follower won't vote for a candidate
  whose log is less up to date than it's own. This prevents an out of date node
  becoming leader and overwriting previously committed entries in follower's
  logs.
  The voting rule: Compare terms of last entries, break ties with log length.

Q: What happens if two machines timeout and send heartbeats at the same time?
  This means that they both vote for themselves (have 1 vote), and if there are
  4/5 servers alive for example then they both might get one more vote each,
  meaning that both have 2 votes and neither get elected.

A: What we can do in this situation is just wait for another timeout (as each
  term's timeout is randomized per node), and eventually there will be a round
  in which a single node times out first.


What about partial failure of a leader?

I've seen that a leader can go down temporarily, a new election happens and
a new leader is elected, but the old leader comes back (resolved network partition)
and still thinks it is the leader.


### Replicated state machine 
Log entries:
* Command: the command itself (depending on what the client requested)
* Term Number: the term in which the command was received.


We probably need markers for the next index to write to because if a leader and
follower disagree on their log, the index can be reset to where they do agree
and the follower's log overwritten up til it matches the leader

Log markers ->
* Next index: the next index in a follower's log that the leader will append to.
* Match index: the index up to which the leader knows that it and the follower
agree. This can be piggybacked as the reply to the AppendEntries RPC from the
leader, so the master is keeping track of the match index of each follower.
* Commit index: This is also part of the AppendEntries RPC, informing each
  follower of the log index up to which the master has committed the entries to
  it's state.

Normal case ->
A leader only marks a log entry as committed when a majority of nodes in the
cluster have written that entry into their log (so when the match index in all
followers is up to that point). When a leader commits an entry in it's own log,
the state will be updated and client informed. 

Q: What happens if the leader dies before it has replicated the log append
to a majority of followers? If there is a new leader election, then the
chances that the follower elected leader is the one with the append is lower
than 50%.
A: if the entry has not been replicated to a majority of followers, then it
will not be committed by the leader itself meaning it won't propogate to the
state before the leader dies. This means that there is no loss of consistency,
although the uncommitted entries will be lost (?)


Q: what happens if a client reads from a replica after a new log entry has been
committed into the master's state but not that replica's state. Would this read
count as incorrect?

### Consensus Module 

This is what receives commands from clients

### Safety 

Only a server with an up-to-date log can become leader. 
Q: How is this ensured?
A: By making sure the votes are only given to a candidate if it's log is
contains all the voter's committed log entries.


ADDITIONAL
New leaders only commit entries from prior terms until it's sent out it's own
appendEntry.
