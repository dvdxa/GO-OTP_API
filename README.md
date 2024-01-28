# OTP API
The OTP API is responsible for generating, sending, and validating one-time passwords (OTPs)

## OTP Generation
An external service initiates a request to generate an OTP, which is then sent to the SMS-Gateway for further delivery to the client. The application utilizes a specific algorithm and a predefined expiration time to generate the OTP. Subsequently, the generated OTP is transmitted to the SMS-Gateway for delivery.


## OTP Status
In certain scenarios, an external service may require information about the status of an OTP. To obtain this status, the service sends a request, and the OTP's status is determined by the "state" field in the response body. Refer to the table below for the possible "state" values and their descriptions

## OTP Validation
When an external service sends an OTP value, the API first checks the sender's account for any bans. Following this, the service proceeds to validate the OTP by checking its expiration and overall validity. The corresponding response is then sent back to the external service.

## OTP Statuses and Description

| State value | Description | 
|----------|----------|
| 201 | OTP successfully created|
| 206 | OTP successfully sended to SMS-Gateway/client |
|400 | Failed to send OTP to SMS-Gateway/client |
|200 | OTP is valid and processed |
|304 | OTP is invalid |

