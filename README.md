# Orbit Notebook

Orbit Notebook 是一个本地运行的天文观测资料整理 CLI。观测员可以创建观测条目、追加描述，审核员可以处理审核，发布员可以把通过审核的条目发布成可检索资料。数据保存在 `ORBIT_NOTEBOOK_HOME` 指定的目录，默认是当前目录下的 `.orbit-notebook`。

## 目录结构

- `cmd/orbit-notebook`：命令行入口。
- `internal/model`：观测条目、仪器、审核和摘要模型。
- `internal/storage`：本地 JSON 存储与原子写入。
- `internal/validation`：输入和状态校验。
- `internal/workflow`：跨模块业务流程编排。
- `internal/command`：参数解析和输出。

## 运行

```bash
go run ./cmd/orbit-notebook capture --target M42 --instrument scope-a --observed-at 2026-08-24T10:00:00Z --description "清晰可见"
go run ./cmd/orbit-notebook annotate --id obs-000001 --description "核心区域有细微结构" --tags nebula,visual
go run ./cmd/orbit-notebook review --id obs-000001 --reviewer alice --decision approve --note "信息完整"
go run ./cmd/orbit-notebook release --id obs-000001 --publisher bob
go run ./cmd/orbit-notebook list --state released
```

`ORBIT_NOTEBOOK_HOME` 可设置数据目录；未设置时使用 `.orbit-notebook`。所有命令输出 JSON，错误输出到标准错误。
