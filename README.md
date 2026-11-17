# A12 铁路信号设备事件湖

面向铁路信号设备的事件数据工程基线，处理设备事件乱序、序号缺口、重复上报和迟到数据，重建信号机、轨道电路、道岔和联锁状态，并输出告警去重、告警风暴、拓扑关联和工单匹配结果。

## 运行

```bash
go test ./...
go run ./cmd/railway-lake demo --data-dir .demo-data
go run ./cmd/railway-lake ingest --events testdata/events.jsonl --links testdata/topology_links.jsonl --work-orders testdata/work_orders.jsonl --data-dir data
go run ./cmd/railway-lake compute --data-dir data
```

## 职责拆分

- `model`：设备事件、拓扑边、工单、快照和关联结果。
- `sequence`：事件去重、设备序号缺口与重复序号检测。
- `statemachine`：按设备序号重建状态，识别非法状态转换。
- `alarm`：告警去重、抖动和告警风暴策略。
- `correlation`：依据拓扑边做跨设备告警关联。
- `workorder`：按设备与时间窗口匹配工单。
- `analysis`：无存储的纯规则编排。
- `store`：原始事件与派生快照/关联结果分区保存。
- `service`：接入与计算用例编排。
- `cmd`：参数解析和输出。

## 可继续扩展的题目方向

1. 加入联锁进路状态机和进路锁闭冲突检测。
2. 实现按设备类型分区的状态快照表与增量重放。
3. 增加跨站拓扑版本、维修封锁和告警抑制窗口。
