# README

## Testing locally
This section is based on https://kafka.apache.org/quickstart.
You will need six terminals.

### Install kafka locally
```
$ ansible-playbook playbook-local.yml
...
$ ls kafka_2.11-0.11.0.1
...
```


### Kafka cluster & create a topic
Let's you have four terminals to launch kafka cluster.
One ZK process and three kafka processes.

```
[1]$ make run-zk
...

[2]$ make run-kafka1
...

[3]$ make run-kafka2
...

[4]$ make run-kafka3
...
```

```
[5]$ make mk-topic
...
```

### Kafka consumer
```
[5]$ make run-consumer
...
```

### Kafka producer
```
[6]$ make run-producer
...
>
```

You see a prompt `>`, then type something.
It should appear on the `[5]` terminal of consumer.
