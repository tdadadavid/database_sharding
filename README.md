# Database Sharding
## Sharding is a process of breaking down table(s) into different database server instance called shards.
<p>It's more like *Horizontal Partitioning* </p>

***
### Things you should know to help you understand sharding.
 * Hashing Algorithm
 * HashMap Data structure
 * Consistent Hashing (The algorithm that helps to achieve this.)

***
### Key Points
* Sharding is basically taking the principles of Hashmaps over a network.
* It helps in the Shared-Nothing Architecture. Check: [DiceDB](https://github.com/DiceDB/dice)

### My Advice
* When optimizing your database let sharding be your last resort, it really brings complexity.

### Advantages
* It improves query performance due to small index of each table.(Scalability)
* It helps with fault-tolerance when a shard fails it does not affect other shards (Remember: `Sharded-Nothing-Architecture`)

### Disadvantages
* It is very complicated to do
* Transaction cannot be achieved 
* Updating table schema is really tideous since we have to effect that change on every shard.
* The client application has to handle the complex logic of knowing which shard to establish connection with, just has the hashmap(hash func) needs to know which node to work with.
* Joins are very impossible since shards are across different server instances
* Rollbacks are also impossible.