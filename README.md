# charg-go-microservice
Expense Calculator/Organizer Microservice written in Go-lang


## How-to

cd rest-api/
go run main.go


*Open a second terminal*

To get the current Expense List :

```sh

curl localhost:9090/ | jq 

```
jq: optional JSON Beautifier


To update an Expense with a specific ID:

```sh

curl localhost:9090/update/2

```


To add an Expense:

```sh

curl localhost:9090/add -d '{"vendor": vendorName, "value": amountOfMoney}'

```

### First Goal 

A Basic REST API that Gets, Puts, Posts entries for an Expense list
Maybe have some sort of data processing to get per month avg expenses + top expenses

### Second Goal

Use a DB

### Third Goal 

Make it into a microservice (Docker/Kubernetes imlpementation)

### Fourth Goal

Refactor with GRPC