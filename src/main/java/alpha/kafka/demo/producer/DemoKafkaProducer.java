package alpha.kafka.demo.producer;

import alpha.kafka.demo.payload.*;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;

@Service
public class DemoKafkaProducer {

    private static final Logger logger = LoggerFactory.getLogger(DemoKafkaProducer.class);

    public int publish(Order order) {
        print("Order", order);
        //TODO PUBLISH TO KAFKA
        return HttpStatus.CREATED.value();
    }

    public int publish(MatchedOrder matchedOrder) {
        print("MatchedOrder", matchedOrder);
        //TODO PUBLISH TO KAFKA
        return HttpStatus.CREATED.value();
    }

    public int publish(PaidOrder paidOrder) {
        print("PaidOrder", paidOrder);
        //TODO PUBLISH TO KAFKA
        return HttpStatus.CREATED.value();
    }

    public int publish(DeliveredOrder deliveredOrder) {
        print("DeliveredOrder", deliveredOrder);
        //TODO PUBLISH TO KAFKA
        return HttpStatus.CREATED.value();
    }

    public int publish(Receipt receipt) {
        print("Receipt", receipt);
        //TODO PUBLISH TO KAFKA
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
