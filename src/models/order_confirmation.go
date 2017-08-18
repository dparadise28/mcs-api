package models

import (
	"errors"
	"strconv"
)

func (o *Order) UserOrderConfirmationEmail(storeConfirmation bool) Email {
	email, gratitudes, readyIn := "", "", ""
	fee := 0
	if storeConfirmation {
		gratitudes := "You've recieved an order!"
		email := o.Store.Email
		if o.OrderType == DELIVERY {
			fee := int64(o.Store.Delivery.Fee)
			readyIn := "Your delivery is expected to be ready in " + strconv.Itoa(int(o.Store.Delivery.MinTime)) +
				"-" + strconv.Itoa(int(o.Store.Delivery.MaxTime)) + " minutes"
		} else {
			fee := int64(0)
			readyIn := "Your delivery is expected to be ready in " + strconv.Itoa(int(o.Store.Pickup.MinTime)) +
				"-" + strconv.Itoa(int(o.Store.Pickup.MaxTime)) + " minutes"
		}
	} else {
		gratitudes := "Thank You For Your Order!"
		email := o.User.Email
		if o.OrderType == DELIVERY {
			fee := int64(o.Store.Delivery.Fee)
			readyIn := "Your delivery will be ready in " + strconv.Itoa(int(o.Store.Delivery.MinTime)) +
				"-" + strconv.Itoa(int(o.Store.Delivery.MaxTime)) + " minutes"
		} else {
			fee := int64(0)
			readyIn := "Your delivery will be ready in " + strconv.Itoa(int(o.Store.Pickup.MinTime)) +
				"-" + strconv.Itoa(int(o.Store.Pickup.MaxTime)) + " minutes"
		}
	}
	return Email{email, `
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
    .fixed{
      width: 50%;
      font-size: 0.875em;
      font-weight: bold;
      padding: 5px;
    }
    .flex-item{
      flex-grow: 1;
      font-size: 0.875em;
      padding: 5px;
    }
		.column-left{float:left; width: 9%; padding-left: 10px;}
		.column-center{display: inline-block; width: 64%; }
		.column-right{float: right; width: 24%;}
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
		<div style="margin:0 auto;max-width:600px;background:white;">
			<div style="background:#F5F7F9; color:#404040;">
				<h1 style="color:#484848; font-family:'Avenir Next', Avenir, sans-serif; margin:1 auto;" align="center" border="0">` + gratitudes + `</h1>
				<h2>` + o.Store.Name + `</h2>
				<h4>` + o.Store.Address.Line1 + `</h4>
				<br>
				<h4>` + readyIn + `</h4>
			</div>` + o.Cart.GetCartProductsOrderMarkup() + `
			<br><br><br>
		</div>
    <div style="margin:0 auto;max-width:600px;background:white;">
      <div class="container">
        <div class="container">
          <div class="fixed">
          	OrderID: <br><br> Payment: <br><br> Deliver To:
            <br><br><br><br>
          </div>
          <div class="flex-item">` + o.ID + `<br><br>` + o.PaymentMethod + `
						<br><br>` + o.Address.AptSuite + ", " + o.Address.Line1 + `
						<br><br><br>
          </div>
        </div>
        <div class="container">
          <div class="fixed">
            Tip: <br><br> Tax: <br><br> Fee: <br><br> Subtotal:
            <br><br><br>
            <font size="+2">Total:</font>
            <br><br>
          </div>
          <div class="flex-item">` + FormatPriceCents(int64(o.Tip)) + `
						<br><br>` + FormatPriceCents(int64(o.Cart.TotalsTax)) + `
						<br><br>` + FormatPriceCents(fee) + `
						<br><br>` + FormatPriceCents(int64(o.Cart.Totals.Subtotal)) + `<br>
						<br><br>
            <font size="+2">` + FormatPriceCents(int64(o.Cart.Totals.Total)) + `
						</font>
            <br><br>
          </div>
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
</body>`, gratitudes}
}
