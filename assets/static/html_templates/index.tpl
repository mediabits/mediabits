<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Mediabits Web UI</title>

    <!-- Bootstrap -->
    <link href="/static/bootstrap.min.css" rel="stylesheet">
  </head>
  <body>
	<div class="row">
		<div class="col-md-2"></div>
		<div class="col-md-8">
			<h1>Mediabits Web UI</h1>
		</div>
		<div class="col-md-2"></div>
	</div>

	<div class="row">
		<div class="col-md-2"></div>
		<div class="col-md-8">
			<!-- File Picker -->
			<button type="button" class="btn btn-primary btn-lg" data-toggle="modal" data-target="#myModal">
				Choose File
			</button>

			<input type="text" disabled="disabled" id="mediafile" />

			<!-- Modal -->
			<div class="modal fade" id="myModal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
			  <div class="modal-dialog">
			    <div class="modal-content">
			      <div class="modal-header">
			        <button type="button" class="close" data-dismiss="modal"><span aria-hidden="true">&times;</span><span class="sr-only">Close</span></button>
			        <h4 class="modal-title" id="myModalLabel">File Picker</h4>
			      </div>
			      <div class="modal-body">
			        <div id="filepicker">
			        </div>
			      </div>
			      <div class="modal-footer">
			        <button type="button" class="btn btn-default" data-dismiss="modal">Close</button>
			      </div>
			    </div>
			  </div>
			</div>
		</div>
		<div class="col-md-2"></div>
	</div>

	<div class="row">
		<div class="col-md-2"></div>
		<div class="col-md-8">
			<!-- Nav tabs -->
			<ul class="nav nav-tabs" role="tablist">
			  <li role="presentation" class="active"><a href="#movie" role="tab" data-toggle="tab">Movie</a></li>
			  <li role="presentation"><a href="#tv" role="tab" data-toggle="tab">TV</a></li>
			</ul>

			<!-- Tab panes -->
			<div class="tab-content">
			  <div role="tabpanel" class="tab-pane active" id="movie">
			  	<form role="form" id="movieForm">
			  		<div class="form-group">
						<label for="movieTitle">Title</label>
						<input type="text" name="title" id="movieTitle">
						<label for="movieYear">Year</label>
						<input type="number" id="movieYear">
						<p class="help-block">The year may be left blank.</p>
						<p class="help-block">DO NOT CLICK THE GENERATE BUTTON MORE THAN ONCE. Generating the description may take up to 30 seconds depending on your computer and internet speed.</p>
					</div>
					<button type="submit" class="btn btn-default">Generate</button>
					<span class="loading">Fetching data...</span>
			  	</form>
			  	<div>
			  		<h3>Copy the following to the upload form:</h3>
			  		<strong>Source:</strong> <input id="movieOutSource" readonly /><br />
			  		<strong>Video Codec:</strong> <input id="movieOutVideoCodec" readonly /><br />
			  		<strong>Audio Codec:</strong> <input id="movieOutAudioCodec" readonly /><br />
			  		<strong>Container:</strong> <input id="movieOutContainer" readonly /><br />
			  		<strong>Resolution:</strong> <input id="movieOutResolution" readonly /><br />
			  		<strong>Year:</strong> <input id="movieOutYear" readonly /><br />
			  		<strong>Description:</strong><br />
			  		<textarea id="movieOutDescription" rows="10" cols="80" readonly></textarea><br />
			  		<strong>Mediainfo:</strong><br />
			  		<textarea id="movieOutMediainfo" rows="10" cols="80" readonly></textarea><br />
			  		<strong>Screenshots:</strong><br />
			  		<textarea id="movieOutScreenshots" rows="3" cols="60" readonly></textarea><br />
			  		<strong>Image:</strong> <input id="movieOutImage" size="60" readonly /><br />
			  	</div>
			  </div>
			  <div role="tabpanel" class="tab-pane" id="tv">
			  	<form role="form" id="tvForm">
			  		<div class="form-group">
						<label for="tvShow">Show</label>
						<input type="text" name="title" id="tvShow">
						<label for="tvYear">Year</label>
						<input type="number" id="tvYear">
						<label for="tvSeason">Season</label>
						<input type="number" id="tvSeason">
						<label for="tvEpisode">Episode</label>
						<input type="number" id="tvEpisode">
						<p class="help-block">Leave the episode field blank when doing a season. The year may be left blank (not zero) if there is only one TV show by that name.</p>
						<p class="help-block">DO NOT CLICK THE GENERATE BUTTON MORE THAN ONCE. Generating the description may take up to 30 seconds depending on your computer and internet speed.</p>
					</div>
					<button type="submit" class="btn btn-default">Generate</button>
					<span class="loading">Fetching data...</span>
			  	</form>
			  	<div>
			  		<h3>Copy the following to the upload form:</h3>
			  		<strong>Title:</strong> <input id="tvOutTitle" readonly /><br />
			  		<strong>Description:</strong><br />
			  		<textarea id="tvOutDescription" rows="10" cols="80" readonly></textarea><br />
			  		<strong>Image:</strong> <input id="tvOutImage" size="60" readonly /><br />
			  	</div>
			  </div>
			</div>
		</div>
		<div class="col-md-2"></div>
	</div>

    <!-- jQuery (necessary for Bootstrap's JavaScript plugins) -->
    <script src="/static/jquery.min.js"></script>
    <!-- Include all compiled plugins (below), or include individual files as needed -->
    <script src="/static/bootstrap.min.js"></script>

    <script src="/static/mediabits.js"></script>
  </body>
</html>