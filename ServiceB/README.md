# FC-Cloud-Run

Build
```shell
docker-compose build
```
Run
```shell
docker-compose up
```

Run test
```shell
go test -race ./...
```

Application url cloudrun

```
https://temperature-zipcode-mquvk55p3a-uc.a.run.app/temperature?zipcode=06835100
```
change zipcode in path
```
GET /temperature?zipcode=06835100
```