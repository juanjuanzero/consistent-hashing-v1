# Juan man's attempt at consistent hashing

## What is consistent-hashing?

Consistent hashing is a technique that is widely used in distributed systems. It's a mechanism used to evenly distribute items in the context of network faults. A network fault is simply when communicate between two nodes fall apart.

- Imagine you are in a place where you store your data in different machines, although those different machines must act like a single machine. How would you direct requests to the correct machine? Consistent hashing solves this problem by having a consistent way to identify where the data lives
- Image in a distributed system where one of the nodes fail, you now have to add nodes, and remove the old node, and migrate the data into the new node. How do you know what data should go into the new node? Consistent hashing solves this problem by easily identifying the scope of data responsibility a server has.

## How does consistent hashing solve this?

- CH solves this by creating a hash-ring and assigning specific hashed keys to each sub-server.
- the severs have an assigned hash value, and all you need to do is create a hash ring and go clockwise to identify the first server you can see in the hash space.
- A hash is just a large number to is representative of the data passed into the hashing function.

## What shall we implement?

We will implement a consistent-hashing algorithm for a key-value store, requestors will make a request with a key, that key will get mapped to the appropriate server which houses the key. For simplicity we'll implement a server instance containing an in-memory map.

## What are the components that we would need?

- a server that manages requests (some kind of master)
- a few nodes server that contain data

## What happens on a request?

A client will make a get request for a particular key the master node will make a request to the leaf node, some thing for sets.

## What happens when a server get's killed?

The data from replicas will need to get migrated to the replacement node.

## Building

- We'll use the cli for this first to act as a single client that we make to the master node, this master node will contain most of the context and information about the sub-servers.
