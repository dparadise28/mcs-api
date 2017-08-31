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
		.column-left{float:left; width: 9%; padding-left: 10px;}
		.column-center{display: inline-block; width: 64%; }
		.column-right{float:right; width: 11%; padding-left: 8px;}
	</style>
</head>
<body style="margin:0 auto;max-width:600px; background: #E3E5E7;">
	<div style="margin:0 auto;max-width:600px;background-color:#E3E5E7;">
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
		<div style="margin:0 auto;max-width:600px;background:white;" align="center">
			<div style="width: 95%; cursor:auto;color:#0f1f38;font-family:'Avenir Next', Avenir, sans-serif;font-size:16px;line-height:30px;">
				<h2>OrderID: o.ID</h2>
				` + msgIntro + ` <i>` + o.OrderStatus + `</i>` + punctuation + ` As always, dont
				hesitate to reach out by replying to this email. Please set the subject of email
				inqueries about <i>` + o.OrderStatus + `</i> orders to
				<h4>Canceled Order Questions: {OrderID}</h4>
				where {OrderID} is the id found above. This will help us service your questions as quickly as possible.
			</div>
			<br><br>
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
	return Email{
		o.DestinationEmail,
		body,
		"Your order has been updated",
	}
}
