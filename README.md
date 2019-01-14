# TimeShard
TimeShard is an open-source operational transform time series database. It was primarily designed, to store *OT-Operations* from web edited documents (by sending only the modified bits).

## Currently implemented
- [x] Insert Operation with Retain to move the cursor.
- [x] Delete Operation with Retain
- [x] Data snapshots
- [x] Squash (See how a document looked in a point in time)
- [ ] Error recovery
- [ ] Text Formatting
- [ ] Disk persistence with Snappy

*Inserting data*
```go
batch := NewBatch()

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