package alpha.kafka.demo.payload;

import java.io.Serializable;
import java.time.LocalDateTime;

public class Receipt implements Serializable {
    private String receiptId;
    private String orderType;
    private Integer amt;
    private String from;
    private String to;
    private String sellOrderId;
    private String fundsReceiptId;
    private String assetReceiptId;
    private LocalDateTime receiptCreated;

    public Receipt() {}

    public String getReceiptId() {
        return receiptId;
    }

    public void setReceiptId(String receiptId) {
        this.receiptId = receiptId;
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

    public String getSellOrderId() {
        return sellOrderId;
    }

    public void setSellOrderId(String sellOrderId) {
        this.sellOrderId = sellOrderId;
    }

    public String getFundsReceiptId() {
        return fundsReceiptId;
    }

    public void setFundsReceiptId(String fundsReceiptId) {
        this.fundsReceiptId = fundsReceiptId;
    }

    public String getAssetReceiptId() {
        return assetReceiptId;
    }

    public void setAssetReceiptId(String assetReceiptId) {
        this.assetReceiptId = assetReceiptId;
    }

    public LocalDateTime getReceiptCreated() {
        return receiptCreated;
    }

    public void setReceiptCreated(LocalDateTime receiptCreated) {
        this.receiptCreated = receiptCreated;
    }
}
