package models

func ConirmationEmail(uid string, c_code string, u_email string, c_link string) string {
	// uid = user id
	// u_email = user email
	// c_code = confirmation code
	// c_link = cofirmation link generated from the method above
	return `
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
                                      <td style="width: 60%;"><a href="about:blank" target="_blank"><img alt="auth0" title="" height="auto" src="` + DOMAIN + `/img/fulllogo.00df53e.png" style="border:none;border-radius:;display:block;outline:none;text-decoration:none;width:100%;height:auto;" width="60%"></a></td>
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
        			<table cellpadding="0" cellspacing="0" style="font-size:0px;width:100%;background:white;" align="center" border="0">
        				<tbody>
        					<tr>
        						<td style="text-align:center;vertical-align:top;font-size:0px;padding:40px 25px 0px;">
        							<div aria-labelledby="mj-column-per-100" class="mj-column-per-100" style="vertical-align:top;display:inline-block;font-size:13px;text-align:left;width:100%;">
        								<table cellpadding="0" cellspacing="0" width="100%" border="0">
        									<tbody>
        										<tr>
        											<td style="word-break:break-word;font-size:0px;padding:0px 0px 25px;" align="left">
        												<div style="cursor:auto;color:#0f1f38;font-family:'Avenir Next', Avenir, sans-serif;font-size:18px;font-weight:500;line-height:30px;">
        													Your account information
        												</div>
        											</td>
        										</tr>
        									</tbody>
        								</table>
        							</div>
        							<div aria-labelledby="mj-column-per-30" class="mj-column-per-30" style="vertical-align:top;display:inline-block;font-size:13px;text-align:left;width:100%;">
        								<table cellpadding="0" cellspacing="0" width="100%" border="0">
        									<tbody>
        										<tr>
        											<td style="word-break:break-word;font-size:0px;padding:0px 0px 10px;" align="left">
        												<div style="cursor:auto;color:#0f1f38;font-family:'Avenir Next', Avenir, sans-serif;font-size:16px;line-height:30px;">
        													<strong style="font-weight: 500; white-space: nowrap;">Email</strong>
        												</div>
        											</td>
        										</tr>
        									</tbody>
        								</table>
        							</div>
        							<div aria-labelledby="mj-column-per-70" class="mj-column-per-70" style="vertical-align:top;display:inline-block;font-size:13px;text-align:left;width:100%;">
        								<table cellpadding="0" cellspacing="0" width="100%" border="0">
        									<tbody>
        										<tr>
        											<td style="word-break:break-word;font-size:0px;padding:0px 0px 10px;" align="left">
        												<div style="cursor:auto;color:#0f1f38;font-family:'Avenir Next', Avenir, sans-serif;font-size:16px;line-height:30px;">
        													` + u_email + `
        												</div>
        											</td>
        										</tr>
        									</tbody>
        								</table>
        							</div>
									<div aria-labelledby="mj-column-per-30" class="mj-column-per-30" style="vertical-align:top;display:inline-block;font-size:13px;text-align:left;width:100%;">
        								<table cellpadding="0" cellspacing="0" width="100%" border="0">
        									<tbody>
        										<tr>
        											<td style="word-break:break-word;font-size:0px;padding:0px 0px 10px;" align="left">
        												<div style="cursor:auto;color:#0f1f38;font-family:'Avenir Next', Avenir, sans-serif;font-size:16px;line-height:30px;">
        													<strong style="font-weight: 500; white-space: nowrap;">User ID</strong>
        												</div>
        											</td>
        										</tr>
        									</tbody>
        								</table>
        							</div>
        							<div aria-labelledby="mj-column-per-70" class="mj-column-per-70" style="vertical-align:top;display:inline-block;font-size:13px;text-align:left;width:100%;">
        								<table cellpadding="0" cellspacing="0" width="100%" border="0">
        									<tbody>
        										<tr>
        											<td style="word-break:break-word;font-size:0px;padding:0px 0px 25px;" align="left">
        												<div style="cursor:auto;color:#0f1f38;font-family:'Avenir Next', Avenir, sans-serif;font-size:16px;line-height:30px;"> ` + uid + `
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
					<div align="center" style="margin:0 auto;max-width:600px;background:white;"><br>
						Please click the link below and login to verify your account.
					</div>
        			<table cellpadding="0" cellspacing="0" style="font-size:0px;width:100%;background:white;" align="center" border="0">
        				<tbody>
        					<tr>
        						<td style="text-align:center;vertical-align:top;font-size:0px;padding:20px 0px;">
        							<div aria-labelledby="mj-column-per-100" class="mj-column-per-100" style="vertical-align:top;display:inline-block;font-size:13px;text-align:left;width:100%;">
        								<table cellpadding="0" cellspacing="0" width="100%" border="0">
        									<tbody>
        										<tr>
        											<td style="word-break:break-word;font-size:0px;padding:10px 25px;" align="center">
        												<table cellpadding="0" cellspacing="0" align="center" border="0">
        													<tbody>
        														<tr>
        															<td style="border-radius:3px;color:white;cursor:auto;" align="center" valign="middle" bgcolor="#000bb2c">
                                        <a href="` + c_link + `" style="display:inline-block;text-decoration:none;background:#00bb2c;border-radius:3px;color:white;font-family:'Avenir Next', Avenir, sans-serif;font-size:14px;font-weight:500;line-height:35px;padding:10px 25px;margin:0px;" target="_blank">
        																  VERIFY YOUR ACCOUNT
        																</a>
        															</td>
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
        		<div></div>
        	</div>
        </body>
    `
}
