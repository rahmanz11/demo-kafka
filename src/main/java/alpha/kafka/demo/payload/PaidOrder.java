package alpha.kafka.demo.payload;

import java.io.Serializable;

public class PaidOrder implements Serializable {
    private String orderId;
    private String orderType;
    private Integer amt;
    private String from;
    private String to;
    private String pmtMethod;
    private String sellOrderId;
    private FundsReceipt fundsReceipt;

    public PaidOrder() {}

    public String getOrderId() {
        return orderId;
    }

    public void setOrderId(String orderId) {
        this.orderId = orderId;
    }

    public String getOrderType() {
        return orderType;
    }

    public void setOrderType(String orderType) {
        this.orderType = orderType;
    }

    public Integer getAmt() {
        return amt;
    }

    public void setAmt(Integer amt) {
        this.amt = amt;
    }

    public String getFrom() {
        return from;
    }

    public void setFrom(String from) {
        this.from = from;
    }

    public String getTo() {
        return to;
    }

    public void setTo(String to) {
        this.to = to;
    }

    public String getPmtMethod() {
        return pmtMethod;
    }

    public void setPmtMethod(String pmtMethod) {
        this.pmtMethod = pmtMethod;
    }

    public String getSellOrderId() {
        return sellOrderId;
    }

    public void setSellOrderId(String sellOrderId) {
        this.sellOrderId = sellOrderId;
    }

    public FundsReceipt getFundsReceipt() {
        return fundsReceipt;
    }

    public void setFundsReceipt(FundsReceipt fundsReceipt) {
        this.fundsReceipt = fundsReceipt;
    }
}
