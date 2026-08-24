# orbit-notebook__010 Docker 交付说明

## 项目概览
- Orbit Notebook 是一个本地运行的天文观测资料整理 CLI。观测员可以创建观测条目、追加描述，审核员可以处理审核，发布员可以把通过审核的条目发布成可检索资料。数据保存在 `ORBIT_NOTEBOOK_HOME` 指定的目录，默认是当前目录下的 `.orbit-notebook`。
- Go module: `example.com/orbit-notebook`

## 标准命令

```bash
go build ./...
go test ./...
```

## 实际启动入口

```bash
go run ./cmd/orbit-notebook
```

## Docker 构建

```bash
./build_benzhi_docker.sh orbit-notebook__010-benzhi linux/amd64
docker run --rm -it orbit-notebook__010-benzhi bash
```

## 环境

- 基础镜像: `golang:1.26.2`
- 依赖在镜像构建阶段预下载，容器内可直接执行 Go 构建和测试命令。
- 代码目录: `/app`
