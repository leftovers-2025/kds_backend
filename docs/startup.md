# 実行手順書

## 前提条件

- goがインストールされていること
- dockerの環境があること

## 実行手順

### 1. envをコピー

```sh
cp example.env .env
```

### 2. dockerコンテナ起動

```sh
docker-compose up -d
```

### 3. airをインストール

```sh
go install github.com/air-verse/air@latest
```

### 4. 実行

```sh
air .
```
