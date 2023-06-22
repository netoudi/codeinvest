# Kafka

## start kafka
```bash
docker-compose up -d
```

## stop kafka
```bash
docker-compose down
```

## verify logs
```bash
docker-compose logs -f kafka
```

## access container kafka
```bash
docker-compose exec kafka bash
```

## send message to kafka (process message in golang)
```bash
kafka-console-producer --bootstrap-server=host.docker.internal:9094 --topic=input
```

## read message from kafka (result message from golang)
```bash
kafka-console-consumer --bootstrap-server=host.docker.internal:9094 --topic=output
```

# input sell
```json
{
  "order_id": "1",
  "investor_id": "Mari",
  "asset_id": "asset1",
  "current_shares": 10,
  "shares": 5,
  "price": 3.8,
  "order_type": "SELL"
}

{"order_id":"1","investor_id":"Mari","asset_id":"asset1","current_shares":10,"shares":5,"price":3.8,"order_type":"SELL"}
```

# input buy
```json
{
  "order_id": "2",
  "investor_id": "Celia",
  "asset_id": "asset1",
  "current_shares": 0,
  "shares": 5,
  "price": 5.0,
  "order_type": "BUY"
}

{"order_id":"2","investor_id":"Celia","asset_id":"asset1","current_shares":0,"shares":5,"price":5.0,"order_type":"BUY"}
```
