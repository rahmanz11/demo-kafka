package alpha.kafka.demo.api;

import alpha.kafka.demo.payload.*;
import alpha.kafka.demo.producer.DemoKafkaProducer;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;

@RequestMapping("api")
@RestController
public class DemoKafkaApi {

    private DemoKafkaProducer producer;

    public DemoKafkaApi(DemoKafkaProducer producer) {
        this.producer = producer;
    }

    @PostMapping("/order")
    @ResponseStatus(HttpStatus.CREATED)
    public int postOrderMessage(@RequestBody Order order) {
        return producer.publish(order);
    }

    @PostMapping("/matched-order")
    @ResponseStatus(HttpStatus.CREATED)
    public int postMatchedOrder(@RequestBody MatchedOrder matchedOrder) {
        return producer.publish(matchedOrder);
    }

    @PostMapping("/paid-order")
    @ResponseStatus(HttpStatus.CREATED)
    public int postMatchedOrder(@RequestBody PaidOrder paidOrder) {
        return producer.publish(paidOrder);
    }

    @PostMapping("/delivered-order")
    @ResponseStatus(HttpStatus.CREATED)
    public int postMatchedOrder(@RequestBody DeliveredOrder deliveredOrder) {
        return producer.publish(deliveredOrder);
    }

    @PostMapping("/receipt")
    @ResponseStatus(HttpStatus.CREATED)
    public int postMatchedOrder(@RequestBody Receipt receipt) {
        return producer.publish(receipt);
    }
}
