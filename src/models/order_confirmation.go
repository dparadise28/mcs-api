package models

import (
	//"errors"
	"encoding/json"
	"log"
	"strconv"
)

func (o *Order) UserOrderConfirmationEmail(storeConfirmation bool) Email {
	if o.Cart.Totals.Subtotal <= uint32(0) {
		o.Cart.UpdateCartTotals()
	}
	email, gratitudes, readyIn, subject := "", "", "", ""
	fee := int64(0)
	if storeConfirmation {
		gratitudes = "You've recieved <br> an order!"
		subject = "You've recieved an order!"
		email = o.Store.Email
		if o.OrderType == DELIVERY {
			fee = int64(o.Store.Delivery.Fee)
			readyIn = "Your delivery is expected to arrive in " + strconv.Itoa(int(o.Store.Delivery.MinTime)) +
				"-" + strconv.Itoa(int(o.Store.Delivery.MaxTime)) + " minutes"
		} else {
			fee = int64(0)
			o.Cart.ApplyFee = false
			readyIn = "Your pickup is expected to be ready in " + strconv.Itoa(int(o.Store.Pickup.MinTime)) +
				"-" + strconv.Itoa(int(o.Store.Pickup.MaxTime)) + " minutes"
		}
	} else {
		gratitudes = "Thank You For <br> Your Order!"
		subject = "Thank You For Your Order!"
		email = o.User.Email
		if o.OrderType == DELIVERY {
			fee = int64(o.Store.Delivery.Fee)
			readyIn = "Your delivery will arrive " + strconv.Itoa(int(o.Store.Delivery.MinTime)) +
				"-" + strconv.Itoa(int(o.Store.Delivery.MaxTime)) + " minutes"
		} else {
			fee = int64(0)
			o.Cart.ApplyFee = false
			readyIn = "Your pickup will be ready in " + strconv.Itoa(int(o.Store.Pickup.MinTime)) +
				"-" + strconv.Itoa(int(o.Store.Pickup.MaxTime)) + " minutes"
		}
	}
	b, _ := json.Marshal(o.Cart)
	log.Println(string(b))
	subtotal := FormatPriceCents(int64(o.Cart.Totals.Subtotal))
	log.Println("asldkjfhasd", o.Cart.Totals.Subtotal, FormatPriceCents(int64(o.Cart.Totals.Subtotal)))
	body := `
<html>
<head>
	<title></title>
	<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
	<style type="text/css">
		#outlook a { padding: 0; }
		.ReadMsgBody { width: 100%; }
		.ExternalClass { width: 100%; }
		.ExternalClass * { line-height:100%; }
		body { margin: 0; padding: 0; -webkit-text-size-adjust: 100%; -ms-text-size-adjust: 100%; }
		table, td { border-collapse:collapse; mso-table-lspace: 0pt; mso-table-rspace: 0pt; }
		img { border: 0; height: auto; line-height: 100%; outline: none; text-decoration: none; -ms-interpolation-mode: bicubic; }
		p { display: block; margin: 13px 0; }
	</style>
	<style type="text/css">
		@media only screen and (max-width:480px) {
		@-ms-viewport { width:320px; }
		@viewport { width:320px; }
		}
	</style>
	<style type="text/css">
		@media only screen and (min-width:480px) {
  		.mj-column-per-100, * [aria-labelledby="mj-column-per-100"] { width:100%!important; }
  		.mj-column-per-80, * [aria-labelledby="mj-column-per-80"] { width:80%!important; }
  		.mj-column-per-30, * [aria-labelledby="mj-column-per-30"] { width:30%!important; }
  		.mj-column-per-70, * [aria-labelledby="mj-column-per-70"] { width:70%!important; }
		}
		.container{
      display:flex;
      background: #F5F7F9;
    }
		.innercontainer{
      display:flex;
      background: #F5F7F9;
			display: table;
  		margin: 0 auto;
			width:100%!important;
    }
    .fixed{
      width: 50%;
      font-weight: bold;
      padding: 5px;
    }
    .flex-item{
      flex-grow: 1;
      padding: 5px;
    }
		.column-left{float:left; width: 9%; padding-left: 5px;}
		.column-center{display: inline-block; width: 72%; }
		.column-right{float:right; width: 15%; padding-left: 5px;}
	</style>
</head>
<body style="background: #E3E5E7;">
	<div style="background-color:#E3E5E7;">
    <br>
		<div style="margin:0 auto;max-width:600px;background:#0f1f38;">
			<table cellpadding="0" cellspacing="0" style="font-size:0px;width:100%;background:#edf0f5;" align="center" border="0">
				<tbody>
					<tr>
						<td style="text-align:center;vertical-align:top;font-size:0px;padding:20px 0px;">
							<div aria-labelledby="mj-column-per-80" class="mj-column-per-80" style="vertical-align:top;display:inline-block;font-size:13px;text-align:left;width:100%;">
								<table cellpadding="0" cellspacing="0" style="vertical-align:top;" width="100%" border="0">
									<tbody>
										<tr>
											<td style="word-break:break-word;font-size:0px;padding:10px 25px;padding-top:30px;" align="center">
												<table cellpadding="0" cellspacing="0" style="border-collapse:collapse;border-spacing:0px;" align="center" border="0">
													<tbody>
														<tr>
                              <td style="width: 60%;"><a href="about:blank" target="_blank"><img alt="auth0" title="" height="auto" src="http://mycorner.store:8003/img/fulllogo.00df53e.png" style="border:none;border-radius:;display:block;outline:none;text-decoration:none;width:100%;height:auto;" width="60%"></a></td>
														</tr>
													</tbody>
												</table>
											</td>
										</tr>
									</tbody>
								</table>
							</div>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
		<div style="margin:0 auto;max-width:600px;background:white;" style="word-break:break-word;font-size:0px;padding:10px 25px;padding-top:30px;" align="center">
			<div style="background:#F5F7F9; color:#404040;">
				<h1 style="color:#484848; font-family:'Avenir Next', Avenir, sans-serif; margin:1 auto;" align="center" border="0">` + gratitudes + `</h1>
				<h2 align="left">` + o.Store.Name + `</h2>
				<h4 align="left">` + o.Store.Address.Line1 + `</h4>
			</div>

			<br><br>` + o.Cart.GetCartProductsOrderMarkup() + `
			<br><br><br>
		</div>
    <div style="margin:0 auto;max-width:600px;background:white;" style="word-break:break-word;font-size:0px;padding:10px 25px;padding-top:30px;" align="center">
      <div class="container">
        <div class="fixed">OrderID:
          <br><br>Payment:
          <br><br>Order Type:
          <br><br>Expected By:
          <br><br>Phone:
          <br><br>Apt/Suite:
          <br><br>Deliver To:
          <br><br><br>
        </div>
        <div class="fixed">` + o.ID.Hex() + `
					<br><br>` + o.OrderType + `
					<br><br>` + o.PaymentMethod + `
					<br><br>` + readyIn + `
          			<br><br>` + o.Address.Phone + `
          			<br><br>` + o.Address.AptSuite + `
          			<br><br>` + o.Address.Line1 + `
					<br><br>
			<br>
        </div>
      </div>
			<div class="container" align="left">
        <div class="fixed">Tip:
          <br><br>Tax:
          <br><br>Fee:<br><br>
        </div>
        <div class="fixed">` + FormatPriceCents(int64(o.Tip)) + `
					<br><br>` + FormatPriceCents(int64(o.Cart.Totals.Tax)) + `
					<br><br>` + FormatPriceCents(int64(fee)) + `
					<br><br>
        </div>
        <div class="fixed">Subtotal:
          <br><br><br>Total: <br><br>
        </div>
        <div class="fixed">` + subtotal + `
          <br><br><br>` + FormatPriceCents(int64(o.Cart.Totals.Total)) + `<br><br>
        </div>
      </div>
    </div>
		<div style="margin:0 auto;max-width:600px;background:white;">
			<table cellpadding="0" cellspacing="0" style="font-size:0px;width:100%;background:white;" align="center" border="0">
				<tbody>
					<tr>
						<td style="text-align:center;vertical-align:top;font-size:0px;padding:0px 30px;">
							<p style="font-size:1px;margin:0 auto;border-top:1px solid #E3E5E7;width:100%;"></p>
						</td>
					</tr>
				</tbody>
			</table>
		</div>

		<div style="margin:0 auto;max-width:600px;background:white;">
			<table cellpadding="0" cellspacing="0" style="font-size:0px;width:100%;background:white;" align="center" border="0">
				<tbody>
					<tr>
						<td style="text-align:center;vertical-align:top;font-size:0px;padding:20px 0px;">
							<div aria-labelledby="mj-column-per-100" class="mj-column-per-100" style="vertical-align:top;display:inline-block;font-size:13px;text-align:left;width:100%;">
								<table cellpadding="0" cellspacing="0" style="vertical-align:top;" width="100%" border="0">
									<tbody>
										<tr>
											<td style="word-break:break-word;font-size:0px;padding:0px 25px 15px;" align="left">
												<div style="cursor:auto;color:#0f1f38;font-family:'Avenir Next', Avenir, sans-serif;font-size:16px;line-height:30px;">
													If you are having any issues with your account, please don't hesitate to contact us by replying to this mail.
													<br>Thanks!
												</div>
											</td>
										</tr>
									</tbody>
								</table>
							</div>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
		<div style="margin:0 auto;max-width:600px;background:#F5F7F9;">
			<table cellpadding="0" cellspacing="0" style="font-size:0px;width:100%;background:#F5F7F9;" align="center" border="0">
				<tbody>
					<tr>
						<td style="text-align:center;vertical-align:top;font-size:0px;padding:20px 0px;">
							<div aria-labelledby="mj-column-per-100" class="mj-column-per-100" style="vertical-align:top;display:inline-block;font-size:13px;text-align:left;width:100%;">
								<table cellpadding="0" cellspacing="0" style="vertical-align:top;" width="100%" border="0">
									<tbody>
										<tr>
											<td style="word-break:break-word;font-size:0px;padding:0px 20px;" align="center">
												<div style="cursor:auto;color:#0f1f38;font-family:'Avenir Next', Avenir, sans-serif;font-size:13px;line-height:20px;">
													You are receiving this email because you have an account on MyCorner.
													If you are not sure why you're receiving this, please contact us.
												</div>
											</td>
										</tr>
									</tbody>
								</table>
							</div>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</body>
</html>`
	return Email{email, body, subject}
}
