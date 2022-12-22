# Cockroach Trigrams
Let's use CockroachDB's trigram indexes to search USDA's FDCD!

A naive approach to using trigram indexes and the similarity operator (`%`) on
the FDCD dataset has proved to be slower and less useful than would be ideal.

This repository explores different techniques of indexing and querying data to
find the ideal way to make the FDCD dataset discoverable.

Each technique will be evaluated by its speed and the relevance of the search
results it returns.


## Running locally

### Start up CockroachDB
```bash
❯❯❯ cockroach version
Build Tag:        v22.2.0
Build Time:       2022/12/05 16:56:56
Distribution:     CCL
Platform:         darwin arm64 (aarch64-apple-darwin21.2)
Go Version:       go1.19.1
C Compiler:       Clang 10.0.0
Build Commit ID:  77667a1b0101cd323090011f50cf910aaa933654
Build Type:       release

❯❯❯ cockroach start-single-node --insecure
```

### Index some data
```bash
# This will automatically download the FDCD dataset to ./data and load it into
# CRDB.
make crdb-trgrm && ./crdb-trgrm load
```

### Explore the querying techniques
The `query` command will query the dataset with each configured querier.
```
make crdb-trgrm && ./crdb-trgrm query "melk"
```
