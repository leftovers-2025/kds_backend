# kds_backend

## 環境構築

- 言語: go 1.24.4

### sql-migrate

```sh
go install github.com/rubenv/sql-migrate/...@latest
```

## 実行方法

```sh
go run ./cmd/api/
```

## 認証URL

https://accounts.google.com/o/oauth2/v2/auth?
scope=openid%20https%3A//www.googleapis.com/auth/userinfo.profile%20https%3A//www.googleapis.com/auth/userinfo.email&
access_type=offline&
include_granted_scopes=true&
response_type=code&
redirect_uri=http%3A//localhost:8630/oauth/google/redirect&
client_id=532447272997-9dfmmst462j9okkg893nmidhhi6v94mn.apps.googleusercontent.com
