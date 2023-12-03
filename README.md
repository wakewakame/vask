# vask

タスクの見積と実績をバージョン管理するアプリ

# ビルド & 起動

```
git clone https://github.com/wakewakame/vask.git
cd vask
go build
./vask
```

# 使い方

WIP

# API

- `curl -X GET    http://localhost:8080/api/project`
- `curl -X POST   http://localhost:8080/api/project -d '{"name": "test"}'`
- `curl -X GET    http://localhost:8080/api/project/1`
- `curl -X PUT    http://localhost:8080/api/project/1 -d '{"name": "test2"}'`
- `curl -X DELETE http://localhost:8080/api/project/1`
