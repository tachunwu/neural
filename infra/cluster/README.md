# Cluster

## 如何啟動 Cluster
```
goreman start
```

## 啟動的 Cluster 資訊
這個 Cluster 為三個 node 組成，達到 HA Raft 最低要求的設計。
|  Node | Client Post     | Cluster Port  | Store Dir    |
|  ---- | ----            |----           |----          |
| n-0   | localhost:20000 |localhost:10000| /data/node-0 |
| n-1   | localhost:20001 |localhost:10001| /data/node-1 |
| n-2   | localhost:20002 |localhost:10002| /data/node-2 |


## 檢查 Cluster 狀態
```
nats server report jetstream -s localhost:10000 --user admin --password secret
```

## 預期 Cluster 狀態
```
╭───────────────────────────────────────────────────────────────────────────────────────────────╮
│                                       JetStream Summary                                       │
├────────┬─────────┬─────────┬───────────┬──────────┬───────┬────────┬──────┬─────────┬─────────┤
│ Server │ Cluster │ Streams │ Consumers │ Messages │ Bytes │ Memory │ File │ API Req │ API Err │
├────────┼─────────┼─────────┼───────────┼──────────┼───────┼────────┼──────┼─────────┼─────────┤
│ n-1    │ cluster │ 0       │ 0         │ 0        │ 0 B   │ 0 B    │ 0 B  │ 0       │ 0       │
│ n-0    │ cluster │ 0       │ 0         │ 0        │ 0 B   │ 0 B    │ 0 B  │ 0       │ 0       │
│ n-2*   │ cluster │ 0       │ 0         │ 0        │ 0 B   │ 0 B    │ 0 B  │ 0       │ 0       │
├────────┼─────────┼─────────┼───────────┼──────────┼───────┼────────┼──────┼─────────┼─────────┤
│        │         │ 0       │ 0         │ 0        │ 0 B   │ 0 B    │ 0 B  │ 0       │ 0       │
╰────────┴─────────┴─────────┴───────────┴──────────┴───────┴────────┴──────┴─────────┴─────────╯

╭─────────────────────────────────────────────────╮
│           RAFT Meta Group Information           │
├──────┬────────┬─────────┬────────┬────────┬─────┤
│ Name │ Leader │ Current │ Online │ Active │ Lag │
├──────┼────────┼─────────┼────────┼────────┼─────┤
│ n-0  │        │ true    │ true   │ 0.31s  │ 0   │
│ n-1  │        │ true    │ true   │ 0.32s  │ 0   │
│ n-2  │ yes    │ true    │ true   │ 0.00s  │ 0   │
╰──────┴────────┴─────────┴────────┴────────┴─────╯
```