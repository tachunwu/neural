# Neural

用 Application 去找可以用的 Setup(不要去研究單獨的參數設定，因為沒有使用情境很難定義)

# Application
目標/基礎設施設置/應用程式邏輯code
<!-- --------------------------------------------------------------------------------------------------------- -->
## Synchronous Request
<!-- --------------------------------------------------------------------------------------------------------- -->
## Vectored I/O
<!-- --------------------------------------------------------------------------------------------------------- -->
## Distirbuted Concurrency Control
### Objective
主要作為分散式的 Concurrency Control，在底層的 Record 都會對應到上層的 Stream Subject，我們可以限制在 Subject 中訊息的數量，只要超過 Queueing Size 就會拒絕服務。
<!-- --------------------------------------------------------------------------------------------------------- -->
## Workflow (Multi-Stage Transaction)
<!-- --------------------------------------------------------------------------------------------------------- -->






## Deterministic Transaction (Heterogeneous Transaction)
### Objective
應用 Calvin Protocol，先將 Transaction 在一個 Stream 排序，然後 Consumer Service 然後一一照順序執行。假設有兩個服務 Service A 和 Service B，Request 可以森成三種分類 A/B/A&B，可以設計成以下的 Subject。
* Calvin: 將 Distributed Transaction 排序。
* <ServiceA>: 只服務這個 Local Service 的 Transaction
* <ServiceB>: 只服務這個 Local Service 的 Transaction

### Implementation
首先切換到目錄 ```infra/jetstream``` 啟動 single node Disk-based JetStream。
```
goreman start
```

實驗的 code 在目錄 ```application/calvin``` 用以下指令啟動實驗。
```
goreman start
```
首先 Server 會建立兩個 Consumer(ConsumerA, ConsumerB) 分別模擬兩個服務，會用 Pull 的方式消化 Calvin 這個 Subject 的 Transaction。Client 會先配一個 UUID 放入 message Header，等下游處理完畢之後會將 commit 發送回去 Coordinator，等收到所有參與服務的 commit message 就意味著完成。

### Result
這個例子的 Transaction Id 是 ```xrf0f0IWJBhQTYrZAIYFrP```，client 會發起訂閱 ```Calvin.{transaction_id}``` 也就是 ```Calvin.xrf0f0IWJBhQTYrZAIYFrP``` 所以我們可以看到有兩個 Coordinator commit 的訊息。

```
17:02:00 calvin | 2022/12/29 17:02:00 Service: CalvinConsumerA local commit: xrf0f0IWJBhQTYrZAIYFrP
17:02:00 calvin | 2022/12/29 17:02:00 Service: CalvinConsumerB local commit: xrf0f0IWJBhQTYrZAIYFrP
17:02:00 calvin | 2022/12/29 17:02:00 Coordinator commit: CalvinConsumerA
17:02:00 calvin | 2022/12/29 17:02:00 Coordinator commit: CalvinConsumerB
```





<!-- --------------------------------------------------------------------------------------------------------- -->
## HA Writer
<!-- --------------------------------------------------------------------------------------------------------- -->
## Application Level Sharding (Notion-based Method)
<!-- --------------------------------------------------------------------------------------------------------- -->
## Worker Pool
<!-- --------------------------------------------------------------------------------------------------------- -->


# Infrastructure Operation
## Core NATS
## JetStream

## Embedded NATS

## Cluster

## Super Cluster

## Leaf


# Stream Management
Stream 的命名採用以下格式，實體本身表示操作的物件，動詞則是操作的方式，由於本身是事件只會紀錄發生事情，語意上則是過去式。

Stream 本身可以視為一個 Table 的概念，我們紀錄了它所有的 event。

* Workflow-based: <Entity>.<verb+ed>
* Table-based: <Entity>.<Primary-Key>.<verb+ed>


# Clean Architecture
雖然我們使用 NATS 作為基礎設施，但是還是希望能夠 Decouple。

有 Consumer struct 本身提供
* Start() error: 開始訂閱，把 onMessage 當作變數傳入 helper，helper 本身負責 Ack 相關的控制邏輯。
* onMessage(ctx, msgs) bool: 處理訊息的實作。


# Reference