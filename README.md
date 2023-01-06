# Neural

用 Application 去找可以用的 Setup(不要去研究單獨的參數設定，因為沒有使用情境很難定義)

# Application
目標/基礎設施設置/應用程式邏輯code
<!-- --------------------------------------------------------------------------------------------------------- -->
## Synchronous Request
<!-- --------------------------------------------------------------------------------------------------------- -->
## Vectored I/O
### Objective
Vectored I/O 算是 Pub/Sub 的一種應用，我們會將 Request 用廣播的方式發送出去，然後下游所有 Worker 可以平行處理，等我們收集到需要的結果就可以直接回傳。應用範圍可以是有 Sharding 過後的 Search Service



![](https://github.com/tachunwu/neural/blob/main/doc/img/vectored%20i_o.png)
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
  
![](https://github.com/tachunwu/neural/blob/main/doc/img/calvin.png)  

上圖分別描述了 Request 和 Commit 的行為，原始論文中有提到如何 Scale Sequencer 的做法，我們可以將 Transaction Queue Sharding，下游服務則用 Round-Robin 的方式分別從各個 Shard Pull Transaction 來執行。 

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
### Objective
HA Writer 的目標是將 Command 轉換成 Event 然後寫到下游的系統中，廣泛用在 Middleware 之類的應用之上。
### Implementation
實作的方式會用到 JetStream 的 Queue，為了達到 HA 我們會分配 Worker Pool 來達到 HA，同時之間我們需要保有寫如的順序性，所以我們要把 MaxAckPending 設為 1，這樣前一個訊息沒有處理完成之前，下一個就不能處裡。注意，這裡指的是邏輯上的順序，可以是整個 Database 或是 Table 甚至可以是 Row-Level (藉由 Entity ID)。



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


# NATS SDK
## Stream
```
type StreamConfig struct {
	Name                 string          `json:"name"`
	Description          string          `json:"description,omitempty"`
	Subjects             []string        `json:"subjects,omitempty"`
	Retention            RetentionPolicy `json:"retention"`
	MaxConsumers         int             `json:"max_consumers"`
	MaxMsgs              int64           `json:"max_msgs"`
	MaxBytes             int64           `json:"max_bytes"`
	Discard              DiscardPolicy   `json:"discard"`
	DiscardNewPerSubject bool            `json:"discard_new_per_subject,omitempty"`
	MaxAge               time.Duration   `json:"max_age"`
	MaxMsgsPerSubject    int64           `json:"max_msgs_per_subject"`
	MaxMsgSize           int32           `json:"max_msg_size,omitempty"`
	Storage              StorageType     `json:"storage"`
	Replicas             int             `json:"num_replicas"`
	NoAck                bool            `json:"no_ack,omitempty"`
	Template             string          `json:"template_owner,omitempty"`
	Duplicates           time.Duration   `json:"duplicate_window,omitempty"`
	Placement            *Placement      `json:"placement,omitempty"`
	Mirror               *StreamSource   `json:"mirror,omitempty"`
	Sources              []*StreamSource `json:"sources,omitempty"`
	Sealed               bool            `json:"sealed,omitempty"`
	DenyDelete           bool            `json:"deny_delete,omitempty"`
	DenyPurge            bool            `json:"deny_purge,omitempty"`
	AllowRollup          bool            `json:"allow_rollup_hdrs,omitempty"`

	// Allow republish of the message after being sequenced and stored.
	RePublish *RePublish `json:"republish,omitempty"`

	// Allow higher performance, direct access to get individual messages. E.g. KeyValue
	AllowDirect bool `json:"allow_direct"`
	// Allow higher performance and unified direct access for mirrors as well.
	MirrorDirect bool `json:"mirror_direct"`
}
```
## Consumer
```
type ConsumerConfig struct {
	Durable         string          `json:"durable_name,omitempty"`
	Name            string          `json:"name,omitempty"`
	Description     string          `json:"description,omitempty"`
	DeliverPolicy   DeliverPolicy   `json:"deliver_policy"`
	OptStartSeq     uint64          `json:"opt_start_seq,omitempty"`
	OptStartTime    *time.Time      `json:"opt_start_time,omitempty"`
	AckPolicy       AckPolicy       `json:"ack_policy"`
	AckWait         time.Duration   `json:"ack_wait,omitempty"`
	MaxDeliver      int             `json:"max_deliver,omitempty"`
	BackOff         []time.Duration `json:"backoff,omitempty"`
	FilterSubject   string          `json:"filter_subject,omitempty"`
	ReplayPolicy    ReplayPolicy    `json:"replay_policy"`
	RateLimit       uint64          `json:"rate_limit_bps,omitempty"` // Bits per sec
	SampleFrequency string          `json:"sample_freq,omitempty"`
	MaxWaiting      int             `json:"max_waiting,omitempty"`
	MaxAckPending   int             `json:"max_ack_pending,omitempty"`
	FlowControl     bool            `json:"flow_control,omitempty"`
	Heartbeat       time.Duration   `json:"idle_heartbeat,omitempty"`
	HeadersOnly     bool            `json:"headers_only,omitempty"`

	// Pull based options.
	MaxRequestBatch    int           `json:"max_batch,omitempty"`
	MaxRequestExpires  time.Duration `json:"max_expires,omitempty"`
	MaxRequestMaxBytes int           `json:"max_bytes,omitempty"`

	// Push based consumers.
	DeliverSubject string `json:"deliver_subject,omitempty"`
	DeliverGroup   string `json:"deliver_group,omitempty"`

	// Inactivity threshold.
	InactiveThreshold time.Duration `json:"inactive_threshold,omitempty"`

	// Generally inherited by parent stream and other markers, now can be configured directly.
	Replicas int `json:"num_replicas"`
	// Force memory storage.
	MemoryStorage bool `json:"mem_storage,omitempty"`
}
```

# Clean Architecture
雖然我們使用 NATS 作為基礎設施，但是還是希望能夠 Decouple。

有 Consumer struct 本身提供
* Start() error: 開始訂閱，把 onMessage 當作變數傳入 helper，helper 本身負責 Ack 相關的控制邏輯。
* onMessage(ctx, msgs) bool: 處理訊息的實作。


# Reference
