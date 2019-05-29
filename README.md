# TimeShard
TimeShard is an open-source operational transform time series database. It was primarily designed, to store *OT-Operations* from web edited documents (by sending only the modified bits).

## Currently implemented
- [x] Insert Operation with Retain to move the cursor.
- [x] Delete Operation with Retain
- [x] Data snapshots
- [x] Squash (See how a document looked in a point in time)
- [x] JSON Export
- [ ] Error recovery
- [ ] Text Formatting
- [ ] Disk persistence with Snappy

*Inserting data*
```go
batch := timeshard.NewBatch()

text := []byte("this is a long text")
batch.Insert(0, text)
```

*Deleting data*
```go
batch.Delete(0, 10)
```

*Iterators*
```go
snap := batch.Snapshot()

iter := snap.Iterator(true)
for iter.HasNext() {
	//iter.Value()
}
```

# FAQ
*Q:* _Is it safe to iterate while adding new data?_

*A:* You can only iterate over snapshots (a copy of a document, at a certain period of time). Snapshots cannot be edited, but you can create Batches from them.



*Q:* _Is it production ready?_

*A:* No. Only if you treat bugs as features.
