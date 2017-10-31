package models

func (o *Order) CompletedEmail() Email {
	msgIntro, punctuation := "", ""
	if o.OrderStatus == REJECTED || o.OrderStatus == CANCELED {
		msgIntro, punctuation = "We regret to inform you that your order has been", "."
	} else {
		msgIntro, punctuation = "Your order has been", "!"
	}

	body := `
<html>
<head>
    <title>Order Update</title>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
</head>
<body style="margin:0 auto;max-width:600px; background: #E3E5E7;">
    <div style="margin:0 auto;max-width:600px;background-color:#E3E5E7;">
    		<br>
        <div style="margin:0 auto;max-width:600px;background:white;" align="center">
            <div style="margin:0 auto;max-width:600px;background:#edf0f5;">
                <img alt="MyCorner"
                         title="MyCorner"
                         height="auto"
                         src="https://mycorner.store/img/fulllogo.00df53e.png"
                         style="display:block;width:100%;height:auto;"
                         width="60%" \>
            </div>
            <div style="width: 95%; cursor:auto;color:#0f1f38;font-family:'Avenir Next', Avenir, sans-serif;font-size:16px;line-height:30px;">
                <h3>OrderID: ` + o.ID.Hex() + `</h3>
                ` + msgIntro + ` <i>` + o.OrderStatus + `</i>` + punctuation + `
            </div>
            <br><br>
        </div>
        <div style="margin:0 auto;max-width:600px;background:white;">
            <div style="text-align:center;vertical-align:top;font-size:0px;padding:0px 30px;">
                <p style="font-size:1px;margin:0 auto;border-top:1px solid #E3E5E7;width:100%;"></p>
            </div>
        </div>
        <div style="margin:0 auto;max-width:600px;background:white;">
            <div style="padding:0px 25px 15px; cursor:auto;color:#0f1f38;font-family:'Avenir Next', Avenir, sans-serif;font-size:16px;line-height:30px;">
                If you are having any issues with your account, please don't hesitate to contact us by replying to this mail. Please set the subject of email
                inquiries about <i>` + o.OrderStatus + `</i> orders to
                <h4>[` + o.OrderStatus + ` Order Inquiry]: ` + o.ID.Hex() + `</h4>
                This will help us service your questions as quickly as possible.
                <br>Thanks!
            </div>
        </div>
        <div style="margin:0 auto;max-width:600px;background:#F5F7F9;" align="center">
	        <div style="cursor:auto;color:#0f1f38;font-family:'Avenir Next', Avenir, sans-serif;font-size:13px;line-height:20px; padding:20px 3px;">
	            You are receiving this email because you have an account on MyCorner.
	            If you are not sure why you're receiving this, please contact us.
	        </div>
        </div>
    </div>
</body>
</html>
	`
	return Email{
		o.DestinationEmail,
		body,
		"Your order has been updated",
	}
}
