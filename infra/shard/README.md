# Shard

## 如何啟動 Shard
```
goreman start
```

## 啟動的 Cluster 資訊
這個 Cluster 為三個 node 組成，彼此是獨立的 server，專門用來實驗 partition 相關的實驗。
|  Node | Client Post     | Store Dir    |
|  ---- | ----            |----          |
| n-0   | localhost:20000 | /data/node-0 |
| n-1   | localhost:20001 | /data/node-1 |
| n-2   | localhost:20002 | /data/node-2 |

