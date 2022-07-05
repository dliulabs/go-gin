# RabbitMQ Tutorials

[RabbitMQ Tutorials](https://www.rabbitmq.com/tutorials/tutorial-two-go.html)

# Test Run

```
export RABBITMQ_URI=amqp://user:password@localhost:5672/
export RABBITMQ_QUEUE=rss_urls

cd producer
go run main.go

curl -X POST http://localhost:5000/parse -d '{"url":"https://www.reddit.com/r/recipes/.rss"}'
```