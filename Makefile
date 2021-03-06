.PHONY: run
run: ssh-config
	ansible-playbook playbook.yml

ssh-config: .vagrant/machines/node*/virtualbox/id
	vagrant ssh-config | tee ssh-config

topic_name := my-topic-3
create-topic:
	ansible -m shell --become -a "cd /opt/kafka/kafka_2.11-0.11.0.1; \
	  bin/kafka-topics.sh --create \
	    --partitions 1 --replication-factor 3 \
	    --topic $(topic_name) \
	    --zookeeper localhost:2181" \
	  node1


## Local
kafka_home := kafka_2.11-0.11.0.2
run-zk:
	cd $(kafka_home); ./bin/zookeeper-server-start.sh config/zookeeper.properties
dump-zk:
	echo dump | nc localhost 2181

run-kafka1:
	cd $(kafka_home); \
		LOG_DIR=`pwd`/../logs-1 ./bin/kafka-server-start.sh `pwd`/../kafka-1.properties
run-kafka2:
	cd $(kafka_home); \
		LOG_DIR=`pwd`/../logs-2 ./bin/kafka-server-start.sh `pwd`/../kafka-2.properties
run-kafka3:
	cd $(kafka_home); \
		LOG_DIR=`pwd`/../logs-3 ./bin/kafka-server-start.sh `pwd`/../kafka-3.properties
run-consumer:
	cd $(kafka_home); \
		./bin/kafka-console-consumer.sh ./config/consumer.properties \
		--bootstrap-server localhost:9192 \
		--topic $(topic_name)

run-producer:
	cd $(kafka_home); \
		./bin/kafka-console-producer.sh \
			--broker-list localhost:9192,localhost:9292,localhost:9392 \
			--topic my-topic-3

mk-topic:
	cd $(kafka_home); ./bin/kafka-topics.sh --zookeeper localhost:2181 \
		--create \
		--replication-factor 3 --partitions 1 \
		--topic $(topic_name)

list-topic:
	cd $(kafka_home); ./bin/kafka-topics.sh --zookeeper localhost:2181 \
		--describe

desc-topic:
	cd $(kafka_home); ./bin/kafka-topics.sh --zookeeper localhost:2181 \
		--describe --topic $(topic_name)

gen-config: kafka-x.properties.erb
	./gen-config.rb
