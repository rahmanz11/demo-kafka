package alpha.kafka.demo.payload;

import java.io.Serializable;
import java.time.LocalDateTime;

public class FundsReceipt implements Serializable {
    private String status;
    private String from;
    private String to;
    private Double fundAmt;
    private LocalDateTime fundsTransferredTime;

    public FundsReceipt() {}

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
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

    public Double getFundAmt() {
        return fundAmt;
    }

    public void setFundAmt(Double fundAmt) {
        this.fundAmt = fundAmt;
    }

    public LocalDateTime getFundsTransferredTime() {
        return fundsTransferredTime;
    }

    public void setFundsTransferredTime(LocalDateTime fundsTransferredTime) {
        this.fundsTransferredTime = fundsTransferredTime;
    }
}
