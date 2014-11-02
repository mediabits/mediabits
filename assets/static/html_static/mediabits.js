function buildentry(dir, dname, name) {
	if (name === undefined) {
		name = dname;
	}

	var li = $('<li />');
	var a = $('<a />');
	a.attr('class', 'dir-entry');
	a.attr('href', '#');
	a.attr('data-dir', dir + '/' + dname);
	a.text(name);
	li.html(a);
	return li;
}

function buildfileentry(dir, name) {
	var li = $('<li />');
	var a = $('<a />');
	a.attr('class', 'file-entry');
	a.attr('href', '#');
	a.attr('data-file', dir + '/' + name);
	a.attr('data-dismiss', 'modal');
	a.text(name);
	li.html(a);
	return li;
}

function listdir(dir) {
	if (dir === undefined) {
		dir = '';
	}
	$.get('/listfiles', { dir: dir })
		.done(function(data) {
			data = JSON.parse(data);

			var fp = $('#filepicker');
			if (!data.Exists) {
				if (dir === '') {
					listdir('/');
				} else {
					listdir();
				}
			} else {
				var dlist = $('<ul />');

				if (dir !== '/') {
					dlist.append(buildentry(data.Directory, '..', '(Up one level)'));
				}

				for (i = 0; i < data.Directories.length; i++) {
					dlist.append(buildentry(data.Directory, data.Directories[i]));
				}

				var flist = $('<ul />');
				for (i = 0; i < data.Files.length; i++) {
					flist.append(buildfileentry(data.Directory, data.Files[i]));
				}
				
				var nav = $('<span />');
				var dirparts = [];
				var sd = data.Directory.split('/');
				for (i = 0; i < sd.length; i++) {
					nav.append(' / ');
					dirparts.push(sd[i]);
					var path = dirparts.join('/');

					if (sd[i] === '') {
						if (i > 0) {
							break;
						}

						sd[i] = 'root (/ or C:)';
						path = '/';
					}

					var a = $('<a />');
					a.attr('href', '#');
					a.attr('class', 'dir-entry');
					a.attr('data-dir', path);
					a.text(sd[i]);
					nav.append(a);
				}

				var t = $('<dir />');
				t.append(nav);
				t.append('<br />');
				t.append('Directories:')
				t.append(dlist);
				t.append('Files:');
				t.append(flist);

				fp.html(t);
			}
		})
		.fail(function() {
			alert('failed to list files');
		});
}

$(function() {
	$('.loading').hide();

	listdir();

	$('#filepicker').on('click', '.dir-entry', function() {
		listdir($(this).attr('data-dir'));
	});

	$('#filepicker').on('click', '.file-entry', function() {
		$('#mediafile').val($(this).attr('data-file'));
	});

	$('#movieForm').submit(function(event) {
		$('.loading').show();
		$.post('/movie', { title: $("#movieTitle").val(), year: $("#movieYear").val(), file: $('#mediafile').val() })
			.done(function(data) {
				$('.loading').hide();
				data = JSON.parse(data);
				
				$('#movieOutTitle').val(data.Title);
				$('#movieOutSource').val(data.Info.Mediainfo.GeneralSection.Source);
				$('#movieOutVideoCodec').val(data.Info.Mediainfo.VideoStream.Format);
				$('#movieOutAudioCodec').val(data.Info.Mediainfo.AudioStream.Format);
				$('#movieOutContainer').val(data.Info.Mediainfo.GeneralSection.Container);
				$('#movieOutResolution').val(data.Info.Mediainfo.VideoStream.Resolution);
				$('#movieOutYear').val(data.Info.IMDB.Year);
				$('#movieOutDescription').val(data.Description);
				$('#movieOutMediainfo').val(data.Info.Mediainfo.Raw);
				$('#movieOutScreenshots').val(data.Info.Screenshots.join('\n'));
				$('#movieOutImage').val(data.Info.Image);
			})
			.fail(function(data) {
				$('.loading').hide();
				alert('Error: ' + data.responseText);
			});
		event.preventDefault();
	});

	$('#tvForm').submit(function(event) {
		$('.loading').show();
		$.post('/tv', { show: $("#tvShow").val(), year: $("#tvYear").val(), season: $('#tvSeason').val(), episode: $('#tvEpisode').val(), file: $('#mediafile').val() })
			.done(function(data) {
				$('.loading').hide();
				data = JSON.parse(data);
				
				$('#tvOutTitle').val(data.Title);
				$('#tvOutDescription').val(data.Description);
				$('#tvOutImage').val(data.Info.Image);
			})
			.fail(function(data) {
				$('.loading').hide();
				alert('Error: ' + data.responseText);
			});
		event.preventDefault();
	});
});