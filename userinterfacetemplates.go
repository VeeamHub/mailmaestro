package main

const restoreTemplate = `{{$sid := .Sessionid}}<!DOCTYPE html>
<html lang="en">
	<head>
		<!-- Required meta tags -->
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

		<!-- Bootstrap CSS -->
		<link rel="stylesheet" href="/static/css/mailmaestro.css">
		<link rel="stylesheet" href="/static/css/tether.min.css">
		<link rel="stylesheet" href="/static/css/bootstrap.min.css">
	</head>
	<body>
		<nav class="navbar fixed-top navbar-toggleable-md navbar-light bg-faded" style="background-color: #4e4e4e">
			<div class="navbar-brand mr-auto">
				<img src="/static/img/logo.png"/>
			</div>
			<div class="collapse navbar-collapse" id="navbarSupportedContent">
				<ul class="navbar-nav mr-auto mt-2 mt-md-0">
				</ul>

				<form class="form-inline my-2 my-lg-0" action="/{{.Sessionid}}/stop" method="post" >
					<button class="btn btn-danger my-2 my-sm-0" type="submit">Stop Explorer</button>
				</form>
			</div>
		</nav>
		<div class="block" id="block" style="">
		</div>
		<div class="container-fluid fixed-top" style="margin-top:0px;margin-left:160px;margin-right:160px;height:80px">
			<div id="dangermsg" class="alert alert-danger mt-2" role="alert" style="display:none">
			</div>
			<div id="successmsg" class="alert alert-success mt-2" role="alert" style="display:none">
			</div>
		</div>
		<div class="container-fluid" >
			<div class="row" style="margin-top:85px;">
				<div class="col-md-2">
					<div class="nav flex-column nav-pills" id="folders" >
						
					</div>
				</div>
				<div class="col-md-10">
					<table class="table table-striped" id="table-mail" style="display:none;" >
						<thead>
						<tr>
							<th>Restore</th>
							<th>From</th>
							<th>To</th>
							<th>Subject</th>
						</tr>
						</thead>
						<tbody id="table-mail-body">
						</tbody>
					</table>
					<table class="table table-striped" id="table-cal" style="display:none;" >
						<thead>
						<tr>
							<th>Restore</th>
							<th>Organizer</th>
							<th>Location</th>
							<th>StartTime</th>
						</tr>
						</thead>
						<tbody id="table-cal-body">
						</tbody>
					</table>
					<table class="table table-striped" id="table-contact" style="display:none;" >
						<thead>
							<tr>
								<th>Restore</th>
								<th>Contact Display</th>
								<th>Contact Email</th>
							</tr>
						</thead>
						<tbody id="table-contact-body">
						</tbody>
					</table>
					<table class="table table-striped" id="table-task" style="display:none;" >
						<thead>
							<tr>
								<th>Restore</th>
								<th>Owner</th>
							</tr>
						</thead>
						<tbody id="table-task-body">
						</tbody>
					</table>
				</div>
			</div>
		</div>
		
		<script src="/static/js/jquery-3.2.1.js"></script>
		<script src="/static/js/popper.min.js"></script>
		<script src="/static/js/tether.min.js"></script>
		<script src="/static/js/bootstrap.min.js"></script>
		<script>
		function selectFolder(id) {
			clearTimeout(window.loadAutoClose)
			clearTimeout(window.delayedBlack)

			var b = $("#block")
			var mails = $("#table-mail")
			var cals = $("#table-cal")
			var contacts = $("#table-contact")
			var tasks = $("#table-task")

			mails.hide()
			cals.hide()
			contacts.hide()
			tasks.hide()

			b.show()

			console.log("Selecting folder: "+id)

			var newSelectedInbox = "nav-"+id
			if (window.selectedInbox) {
				//id limitation of jquery
				$(document.getElementById(window.selectedInbox)).removeClass("active")
			}
			window.selectedInbox = newSelectedInbox

			$(document.getElementById(newSelectedInbox)).addClass("active")

			var req = {
				dataType: "json", 
				url:"/{{$sid}}/folderItems/"+id,
				success:function( data ) {
					var s = $("#successmsg")
					var e = $("#dangermsg")
					var b = $("#block")

					//delay so that it doesn't flash if the get action is short
					window.delayedBlack = setTimeout(function(){ 
						b.hide()
					},100)
					
					if (data.ErrorString == "") {
						console.log(data)

						if(data.Mails.length > 0) {
							var tbody = document.getElementById("table-mail-body")
							while (tbody.firstChild) {
								tbody.removeChild(tbody.firstChild);
							}

							for (var i=0; i < data.Mails.length ;i++ ) {
								var m = data.Mails[i]

								var tr = document.createElement("tr")


								var td = document.createElement("td")
								var div = document.createElement("div")
								div.appendChild(document.createTextNode("Restore"))
								div.setAttribute("onClick","execRestore('"+m.Id+"')")
								div.setAttribute("class","btn btn-success")
								td.appendChild(div)
								tr.appendChild(td)

								var td = document.createElement("td")
								td.appendChild(document.createTextNode(m.From))
								tr.appendChild(td)

								var td = document.createElement("td")
								td.appendChild(document.createTextNode(m.To))
								tr.appendChild(td)

								var td = document.createElement("td")
								td.appendChild(document.createTextNode(m.Subject))
								tr.appendChild(td)


								tbody.appendChild(tr)
								
							}
							mails.show()
						}


						if(data.Appointments.length > 0) {
							var tbody = document.getElementById("table-cal-body")
							while (tbody.firstChild) {
								tbody.removeChild(tbody.firstChild);
							}

							for (var i=0; i < data.Appointments.length ;i++ ) {
								var m = data.Appointments[i]

								var tr = document.createElement("tr")


								var td = document.createElement("td")
								var div = document.createElement("div")
								div.appendChild(document.createTextNode("Restore"))
								div.setAttribute("onClick","execRestore('"+m.Id+"')")
								div.setAttribute("class","btn btn-success")
								td.appendChild(div)
								tr.appendChild(td)

								var td = document.createElement("td")
								td.appendChild(document.createTextNode(m.Organizer))
								tr.appendChild(td)

							
								var td = document.createElement("td")
								td.appendChild(document.createTextNode(m.Location))
								tr.appendChild(td)

								var td = document.createElement("td")
								td.appendChild(document.createTextNode(m.StartTime))
								tr.appendChild(td)


								tbody.appendChild(tr)
								
							}
							cals.show()
						}

						if(data.Contacts.length > 0) {
							var tbody = document.getElementById("table-contact-body")
							while (tbody.firstChild) {
								tbody.removeChild(tbody.firstChild);
							}

							for (var i=0; i < data.Contacts.length ;i++ ) {
								var m = data.Contacts[i]

								var tr = document.createElement("tr")


								var td = document.createElement("td")
								var div = document.createElement("div")
								div.appendChild(document.createTextNode("Restore"))
								div.setAttribute("onClick","execRestore('"+m.Id+"')")
								div.setAttribute("class","btn btn-success")
								td.appendChild(div)
								tr.appendChild(td)

								var td = document.createElement("td")
								td.appendChild(document.createTextNode(m.Display))
								tr.appendChild(td)

								var td = document.createElement("td")
								td.appendChild(document.createTextNode(m.Email))
								tr.appendChild(td)

								tbody.appendChild(tr)
								
							}
							contacts.show()
						}

						if(data.Tasks.length > 0) {
							var tbody = document.getElementById("table-task-body")
							while (tbody.firstChild) {
								tbody.removeChild(tbody.firstChild);
							}

							for (var i=0; i < data.Tasks.length ;i++ ) {
								var m = data.Tasks[i]

								var tr = document.createElement("tr")


								var td = document.createElement("td")
								var div = document.createElement("div")
								div.appendChild(document.createTextNode("Restore"))
								div.setAttribute("onClick","execRestore('"+m.Id+"')")
								div.setAttribute("class","btn btn-success")
								td.appendChild(div)
								tr.appendChild(td)

								var td = document.createElement("td")
								td.appendChild(document.createTextNode(m.Owner))
								tr.appendChild(td)

								tbody.appendChild(tr)
								
							}
							tasks.show()
						}

					} else {
						e.html("Error Loading "+data.ErrorString)
						e.show()
						window.loadAutoClose = setTimeout(function(){ 
							e.hide()
						}, 2000);
					}
				 }
			}
			$.ajax(req)
		}
		function execRestore(id) {
			clearTimeout(window.restoringAutoClose)
			clearTimeout(window.delayedBlack)

			var s = $("#successmsg")
			var e = $("#dangermsg")
			var b = $("#block")

			s.hide()
			e.hide()
			
			
			b.show()


			var req = {
				dataType: "json", 
				url:"/{{$sid}}/restorejson/"+id,
				success:function( data ) {
					var s = $("#successmsg")
					var e = $("#dangermsg")
					var b = $("#block")

					//delay so that it doesn't flash if the get action is short
					window.delayedBlack = setTimeout(function(){ 
						b.hide()
					},100)
					
					if (data.result == "success") {
						s.html("Succesfull restore (Restored = "+data.createdItemsCount+", Merged = "+data.mergedItemsCount+", Skipped = "+data.skippedItemsCount+", Failure  = "+data.failedItemsCount+")")
						s.show()
					} else {
						console.log(data)
						e.html("Restore status : "+data.result+" Error : "+data.message)
						e.show()
					}
					window.restoringAutoClose = setTimeout(function(){ 
						e.hide()
						s.hide()
					}, 2000);
	
				 }
			};
			$.ajax(req)

		}
		$( document ).ready(function() {
			//Select inbox at load, otherwise the first folder
			window.MailFolder = {}
			window.ReverseMailFolder = {}

			{{range .Folders}}window.MailFolder["{{.Id}}"] = "{{.Name}}"
			window.ReverseMailFolder["{{.Name}}"] = "{{.Id}}"
			{{end}}




			var folders = document.getElementById("folders")
			for (key in window.MailFolder) {
				var a = document.createElement("a")
				a.appendChild(document.createTextNode(window.MailFolder[key]))
				a.setAttribute("href","#")
				a.setAttribute("onClick","selectFolder('"+key+"')")
				a.setAttribute("class","nav-link")
				a.setAttribute("id","nav-"+key)
				folders.appendChild(a)
			}
			


			if ("Inbox" in window.ReverseMailFolder) {
				selectFolder(window.ReverseMailFolder["Inbox"])
			} else {
				selectFolder(Object.keys(window.MailFolder)[0]);
			}
		});
		</script>
	</body>
</html>	`

const loginPageTemplate = `<!DOCTYPE html>
<html lang="en">
	<head>
	<!-- Required meta tags -->
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

	<!-- Bootstrap CSS -->
	<link rel="stylesheet" href="/static/css/mailmaestro.css">
	<link rel="stylesheet" href="/static/css/tether.min.css">
	<link rel="stylesheet" href="/static/css/bootstrap.min.css">
	</head>
	<body>
	<div class="container-fluid">
	<div class="row align-items-center vh100">
	  <div class="col-md-8 green-background vh100">
	  </div>
	  <div class="col-md-4">
		<form action="/auth" method="post" >
			<div class="mt-2">
				<h1>Mail Maestro</h1>
			</div>
			<div class="mt-2">
				<input name="username" id="username"/>
			</div>
			<div class="mt-2">
				<input name="password" id="password" type="password"/>
			</div>
			<div class="mt-2">
				<button type="submit" class="btn btn-success" >Login</button>
			</div>
		</form>
	  </div>
	</div>
	</div>


	<script src="/static/js/jquery-3.2.1.slim.min.js"></script>
	<script src="/static/js/popper.min.js"></script>
	<script src="/static/js/tether.min.js"></script>
	<script src="/static/js/bootstrap.min.js"></script>
	</body>
</html>	
		`
