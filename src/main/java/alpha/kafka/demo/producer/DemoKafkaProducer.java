package alpha.kafka.demo.producer;

import alpha.kafka.demo.payload.*;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.gson.Gson;
import org.apache.kafka.clients.producer.KafkaProducer;
import org.apache.kafka.clients.producer.ProducerRecord;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;

import java.util.Properties;
import java.util.UUID;

@Service
public class DemoKafkaProducer {

    private static final Logger logger = LoggerFactory.getLogger(DemoKafkaProducer.class);
    Properties properties = new Properties();
    // Create Kafka producer
    KafkaProducer<String, String> producer;
    Gson gson;

    public DemoKafkaProducer() {
        properties.put("bootstrap.servers", "localhost:9092");
        properties.put("acks"             , "0");
        properties.put("retries"          , "1");
        properties.put("batch.size"       , "20971520");
        properties.put("linger.ms"        , "100");
        properties.put("max.request.size" , "2097152");
        properties.put("compression.type" , "gzip");
        properties.put("key.serializer"   , "org.apache.kafka.common.serialization.StringSerializer");
        properties.put("value.serializer" , "org.apache.kafka.common.serialization.StringSerializer");

        gson = new Gson();
    }

    public int publish(Order order) {
        print("Order", order);
        UUID id = UUID.randomUUID();
        try {
            producer = new KafkaProducer<>(properties);
            producer.send(new ProducerRecord<>("order", id.toString(), gson.toJson(order)));
        } finally {
            producer.close();
        }

        return HttpStatus.CREATED.value();
    }

    public int publish(MatchedOrder matchedOrder) {
        print("MatchedOrder", matchedOrder);
        UUID id = UUID.randomUUID();
        try {
            producer = new KafkaProducer<>(properties);
            producer.send(new ProducerRecord<>("matched-order", id.toString(), gson.toJson(matchedOrder)));
        } finally {
            producer.close();
        }
        return HttpStatus.CREATED.value();
    }

    public int publish(PaidOrder paidOrder) {
        print("PaidOrder", paidOrder);
        UUID id = UUID.randomUUID();
        try {
            producer = new KafkaProducer<>(properties);
            producer.send(new ProducerRecord<>("paid-order", id.toString(), gson.toJson(paidOrder)));
        } finally {
            producer.close();
        }
        return HttpStatus.CREATED.value();
    }

    public int publish(DeliveredOrder deliveredOrder) {
        print("DeliveredOrder", deliveredOrder);
        UUID id = UUID.randomUUID();
        try {
            producer = new KafkaProducer<>(properties);
            producer.send(new ProducerRecord<>("delivered-order", id.toString(), gson.toJson(deliveredOrder)));
        } finally {
            producer.close();
        }
        return HttpStatus.CREATED.value();
    }

    public int publish(Receipt receipt) {
        print("Receipt", receipt);
        UUID id = UUID.randomUUID();
        try {
            producer = new KafkaProducer<>(properties);
            producer.send(new ProducerRecord<>("receipt", id.toString(), gson.toJson(receipt)));
        } finally {
            producer.close();
        }
        return HttpStatus.CREATED.value();
    }

    private void print(String prefix, Object object) {
        ObjectMapper mapper = new ObjectMapper();
        String json = null;
        try {
            json = mapper.writerWithDefaultPrettyPrinter().writeValueAsString(object);
        } catch (JsonProcessingException e) {
            e.printStackTrace();
        }
        logger.debug("{} ~~ {}", prefix, json);
    }
}
