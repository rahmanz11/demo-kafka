package alpha.kafka.demo.payload;

import alpha.kafka.demo.config.LocalDateTimeSerializer;
import com.fasterxml.jackson.annotation.JsonFormat;
import com.fasterxml.jackson.databind.annotation.JsonSerialize;

import java.io.Serializable;
import java.time.LocalDateTime;

public class FundsReceipt implements Serializable {
    private String status;
    private String from;
    private String to;
    private Double fundAmt;
    @JsonFormat(shape = JsonFormat.Shape.STRING, pattern = "yyyy-MM-dd'T'HH:mm:ss.SSSXXX")
    @JsonSerialize(using = LocalDateTimeSerializer.class)
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
