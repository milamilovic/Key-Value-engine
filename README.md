# Key-Value Engine for Data Storage
#### Limitation of User Requests
#### Caching Content
#### Implementation using the Golang programming language

## Components in the project:
#### Correct Write Path
![image](https://github.com/milamilovic/Key-Value-engine/assets/104532211/9196f44b-0324-4155-9903-8ca2cedf6f5e)

#### Correct Read Path
![image](https://github.com/milamilovic/Key-Value-engine/assets/104532211/ed5c2cce-f271-4924-9c10-1ea92ce6d92c)


## Operations supported by the system:
PUT - adding data (accepts a string key and a bit array value, and returns a boolean value)
GET - seeking information (accepts a string key and returns a bit array value)
DELETE - record deletion (accepts a string key and returns a boolean value)
LIST - seeking records by prefix (accepts a string prefix and returns a list of values whose keys start with the specified prefix)
RANGE SCAN - seeking records in a range (accepts minimum and maximum string key values, and returns a list of values whose keys are within the range)

## Data structures used:
1. WAL (Write ahead log)
2. Memtable
3. SSTable (contains Data, index and summary)
4. Cache
5. LSM tree
6. Tocken Bucket
7. Bloom filter
8. Merkle tree
9. B tree
10. Skip list
11. HyperLogLog
12. Count-min sketch
13. SimHash

## Authors
Sonja Baljicki, Dunja Matejić, Mila Milović
