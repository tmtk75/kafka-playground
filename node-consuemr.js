const kafka = require('kafka-node');

const topic = "my-topic-3";
const partition = 0;
const offset = 2421400;


const client = new kafka.KafkaClient({
  kafkaHost: "127.0.0.1:9192,127.0.0.1:9292,127.0.0.1:9392",
  sessionTimeout: 10000,
  requestTimeout: 30000*3,
});
const consumer = new kafka.Consumer(client, [{topic, partition}], {autoCommit: true, fromOffset: true});
consumer.setOffset(topic, partition, offset)


consumer.on('message', (msg) => { console.log(msg); });

new kafka.Offset(client).fetch([
    { topic, partition, time: Date.now(), maxNum: 1 },
  ], (err, data) => console.log(data));
