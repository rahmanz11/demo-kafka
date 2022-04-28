package alpha.kafka.demo.payload;

import java.io.Serializable;
import java.time.LocalDateTime;

public class AssetReceipt implements Serializable {
    private String deliveryStatus;
    private String from;
    private String to;
    private String assetAmt;
    private LocalDateTime assetDeliveredTime;

    public AssetReceipt() {}

    public String getDeliveryStatus() {
        return deliveryStatus;
    }

    public void setDeliveryStatus(String deliveryStatus) {
        this.deliveryStatus = deliveryStatus;
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

    public String getAssetAmt() {
        return assetAmt;
    }

    public void setAssetAmt(String assetAmt) {
        this.assetAmt = assetAmt;
    }

    public LocalDateTime getAssetDeliveredTime() {
        return assetDeliveredTime;
    }

    public void setAssetDeliveredTime(LocalDateTime assetDeliveredTime) {
        this.assetDeliveredTime = assetDeliveredTime;
    }
}
