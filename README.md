# demo-kafka

Kafka downloaded from the following link:
https://dlcdn.apache.org/kafka/3.1.0/kafka_2.13-3.1.0.tgz

kafka located at server in the following path: 
/root/kafka_2.13-3.1.0

zookeeper is running under the screen named "zookeeper"
zookeeper was started using the following command:
/root/kafka_2.13-3.1.0/bin/zookeeper-server-start.sh config/zookeeper.properties


kafka is running under the screen named "kafka"
kafka was started using the following command:
/root/kafka_2.13-3.1.0/bin/kafka-server-start.sh config/server.properties

rest api is running under the screen named "api"
rest api was started using the following command:
java -jar demo-kafka.jar --server.port=89
