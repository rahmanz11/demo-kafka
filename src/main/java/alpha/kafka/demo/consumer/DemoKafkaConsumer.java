package alpha.kafka.demo.consumer;

import alpha.kafka.demo.payload.*;
import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import org.apache.kafka.clients.consumer.*;
import org.apache.kafka.clients.producer.ProducerConfig;
import org.apache.kafka.common.TopicPartition;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import java.time.Duration;
import java.util.*;

@Service
public class DemoKafkaConsumer {

    private static final Logger logger = LoggerFactory.getLogger(DemoKafkaConsumer.class);
    Properties properties = new Properties();
    // Create Kafka producer
    KafkaConsumer<String, String> consumer;
    Gson gson;

    public DemoKafkaConsumer() {
        properties.put("bootstrap.servers"          , "localhost:9092");
        properties.put("key.deserializer"           , "org.apache.kafka.common.serialization.StringDeserializer");
        properties.put("value.deserializer"         , "org.apache.kafka.common.serialization.StringDeserializer");
        properties.put("max.partition.fetch.bytes"  , "2097152");
        properties.put("auto.offset.reset"          , "earliest");
        properties.put("enable.auto.commit"         , "false");

        gson = new GsonBuilder().setPrettyPrinting().create();
    }

    public List<Order> getOrder() {
        List<Order> orders = new ArrayList<>();

        properties.put("group.id", "order-group");
        consumer = new KafkaConsumer<>(properties);

        try {

            consumer.subscribe(Arrays.asList("order"));
            for (int i = 0; i < 5; i++) {
                ConsumerRecords<String, String> records = consumer.poll(Duration.ofMillis(100));
                for (ConsumerRecord<String, String> record : records) {
                    Order order = gson.fromJson(record.value(), Order.class);
                    orders.add(order);
                }
            }
        } finally {
            consumer.close();
        }

        logger.debug("All Orders ~~~");
        logger.debug(gson.toJson(orders));

        return orders;
    }

    public List<MatchedOrder> getMatchedOrder() {
        List<MatchedOrder> matchedOrders = new ArrayList<>();

        properties.put("group.id", "matched-order-group");
        consumer = new KafkaConsumer<>(properties);

        try {
            consumer.subscribe(Arrays.asList("matched-order"));
            for (int i = 0; i < 5; i++) {
                ConsumerRecords<String, String> records = consumer.poll(Duration.ofMillis(100));
                for (ConsumerRecord<String, String> record : records) {
                    MatchedOrder matchedOrder = gson.fromJson(record.value(), MatchedOrder.class);
                    matchedOrders.add(matchedOrder);
                }
            }
        } finally {
            consumer.close();
        }

        logger.debug("All Matched Orders ~~~");
        logger.debug(gson.toJson(matchedOrders));

        return matchedOrders;
    }

    public List<PaidOrder> getPaidOrder() {
        List<PaidOrder> paidOrders = new ArrayList<>();

        properties.put("group.id", "paid-order-group");
        consumer = new KafkaConsumer<>(properties);

        try {
            consumer.subscribe(Arrays.asList("paid-order"));
            for (int i = 0; i < 5; i++) {
                ConsumerRecords<String, String> records = consumer.poll(Duration.ofMillis(100));
                for (ConsumerRecord<String, String> record : records) {
                    PaidOrder paidOrder = gson.fromJson(record.value(), PaidOrder.class);
                    paidOrders.add(paidOrder);
                }
            }
        } finally {
            consumer.close();
        }

        logger.debug("All Paid Orders ~~~");
        logger.debug(gson.toJson(paidOrders));

        return paidOrders;
    }

    public List<DeliveredOrder> getDeliveredOrder() {
        List<DeliveredOrder> deliveredOrders = new ArrayList<>();

        properties.put("group.id", "delivered-order-group");
        consumer = new KafkaConsumer<>(properties);

        try {
            consumer.subscribe(Arrays.asList("delivered-order"));
            for (int i = 0; i < 5; i++) {
                ConsumerRecords<String, String> records = consumer.poll(Duration.ofMillis(100));
                for (ConsumerRecord<String, String> record : records) {
                    DeliveredOrder deliveredOrder = gson.fromJson(record.value(), DeliveredOrder.class);
                    deliveredOrders.add(deliveredOrder);
                }
            }
        } finally {
            consumer.close();
        }

        logger.debug("All Delivered Orders ~~~");
        logger.debug(gson.toJson(deliveredOrders));

        return deliveredOrders;
    }

    public List<Receipt> getReceipt() {
        List<Receipt> receipts = new ArrayList<>();

        properties.put("group.id", "receipt-group");
        consumer = new KafkaConsumer<>(properties);

        try {
            consumer.subscribe(Arrays.asList("receipt"));
            for (int i = 0; i < 5; i++) {
                ConsumerRecords<String, String> records = consumer.poll(Duration.ofMillis(100));
                for (ConsumerRecord<String, String> record : records) {
                    Receipt receipt = gson.fromJson(record.value(), Receipt.class);
                    receipts.add(receipt);
                }
            }
        } finally {
            consumer.close();
        }

        logger.debug("All Receipts ~~~");
        logger.debug(gson.toJson(receipts));

        return receipts;
    }
}
