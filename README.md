# Project budget-app

One Paragraph of project description goes here

##  Kill service on port

```bash
sudo lsof -i :9090
kill -9 <PID>
```

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

run all make commands with clean tests
```bash
make all build
```

build the application
```bash
make build
```

run the application
```bash
make run
```

Create DB container
```bash
make docker-run
```

Shutdown DB container
```bash
make docker-down
```

live reload the application
```bash
make watch
```

run the test suite
```bash
make test
```

clean up binary from the last build
```bash
make clean
```

# Deployment

```bash
ssh root@192.168.0.4

cd /app/budget-app

systemctl stop budget-app

git pull

make build

systemctl restart budget-app
```

## TODO

- [ ] Look into loading html templates globally
- [ ] More UI styling
- [x] Sort Out Routing on FE
- [ ] Look at FE Hot Reloading
- [ ] Add tests
- [ ] Auto download and import transactions from bank
- [ ] Add month/date selector for transactions
- [ ] Add edit functionality in the CRUD ui's
- [ ] Add delete functionality in the CRUD ui's
- [ ] Add a login page
- [ ] Scan receipts and add to transaction
- [ ] Use AI to categorize transactions
- [ ] Add a debt tracker (outside of the banking transactions)
- [ ] Add a savings tracker
- [ ] Add a wallet scanner for crypto holdings
- [x] Add a crypto holdings history tracker (Adda table to store portfolio value over time - perhaps just save the total evrythime the amounts are updated)
- [ ] Auto update crypto holdings once a day
- [ ] Add email report functionality (weekly, monthly)
- [ ] Add alert if budget item is not yet paid at a certain date