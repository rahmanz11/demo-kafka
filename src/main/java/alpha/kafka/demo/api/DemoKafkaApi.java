package alpha.kafka.demo.api;

import alpha.kafka.demo.consumer.DemoKafkaConsumer;
import alpha.kafka.demo.payload.*;
import alpha.kafka.demo.producer.DemoKafkaProducer;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RequestMapping("api")
@RestController
public class DemoKafkaApi {

    private DemoKafkaProducer producer;
    private DemoKafkaConsumer consumer;

    public DemoKafkaApi(DemoKafkaProducer producer,
                        DemoKafkaConsumer consumer) {
        this.producer = producer;
        this.consumer = consumer;
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

    @GetMapping("/order")
    public List<Order> getOrder() {
        return consumer.getOrder();
    }

    @GetMapping("/matched-order")
    public List<MatchedOrder> getMatchedOrder() {
        return consumer.getMatchedOrder();
    }

    @GetMapping("/paid-order")
    public List<PaidOrder> getPaidOrder() {
        return consumer.getPaidOrder();
    }

    @GetMapping("/delivered-order")
    public List<DeliveredOrder> getDeliveredOrder() {
        return consumer.getDeliveredOrder();
    }

    @GetMapping("/receipt")
    public List<Receipt> getReceipt() {
        return consumer.getReceipt();
    }
}
