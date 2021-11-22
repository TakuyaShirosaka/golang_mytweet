# 開発

## Go Version
go version go1.17.3

## 依存ライブラリの導入
```
go mod tidy
```

## 環境変数の設定
```
./config 配下のyml、development.ymlはコミットしていないので各自作成する
```

## hot reloadの導入
https://github.com/cosmtrek/air

GoPath/bin配下にexeを設置

``` 
versionの確認
air -v
```

```
該当プロジェクトのルートパスで起動
air
```