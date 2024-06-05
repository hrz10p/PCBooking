package models

const EmailTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Booking Confirmation</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f4f4f4;
        }
        .email-container {
            max-width: 600px;
            margin: 20px auto;
            background-color: #ffffff;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }
        .email-header {
            text-align: center;
            border-bottom: 1px solid #dddddd;
            padding-bottom: 10px;
            margin-bottom: 20px;
        }
        .email-header h1 {
            margin: 0;
            color: #333333;
        }
        .email-body {
            color: #555555;
        }
        .email-body h2 {
            color: #333333;
        }
        .email-body p {
            line-height: 1.6;
        }
        .email-footer {
            text-align: center;
            border-top: 1px solid #dddddd;
            padding-top: 10px;
            margin-top: 20px;
            color: #999999;
        }
    </style>
</head>
<body>
    <div class="email-container">
        <div class="email-header">
            <h1>Booking Confirmation</h1>
        </div>
        <div class="email-body">
            <h2>Hi, {{.RecipientEmail}}</h2>
            <p>Thank you for your booking. Here are your booking details:</p>
            <p><strong>Booking Date:</strong> {{.BookingDate}}</p>
            <p><strong>Computer Number:</strong> {{.ComputerNumber}}</p>
        </div>
        <div class="email-footer">
            <p>Thank you for choosing our service.</p>
        </div>
    </div>
</body>
</html>
`
