Reddit backend


To run app: 

- Run app ```go run cmd/reddit_clone/main.go```

To run tests:

- Run tests ```go test  ./... -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html```

To run DBs:

- Run DBs ```docker-compose up redis mysql mongodb -d```
- Down DBs ```docker-compose down```