$(function() {
    //$(document).foundation();
		
	// Hide any messages after a few seconds
    hideFlash();
});

function hideFlash(rnum)
{    
    if (!rnum) rnum = '0';
    
    _.delay(function() {
        $('.alert-box-fixed' + rnum).fadeOut(300, function() {
            $(this).css({"visibility":"hidden",display:'block'}).slideUp();
            
            var that = this;
            
            _.delay(function() { that.remove(); }, 400);
        });
    }, 4000);
}

function showFlash(obj)
{
    $('#flash-container').html();
    $(obj).each(function(i, v) {
        var rnum = _.random(0, 100000);
		var message = '<div id="flash-message" class="alert-box-fixed'
		+ rnum + ' alert-box-fixed alert alert-dismissible '+v.cssclass+'">'
		+ '<button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>'
		+ v.message + '</div>';
        $('#flash-container').prepend(message);
        hideFlash(rnum);
    });
}

function flashError(message) {
	var flash = [{Class: "alert-danger", Message: message}];
	showFlash(flash);
}

function flashSuccess(message) {
	var flash = [{Class: "alert-success", Message: message}];	
	showFlash(flash);
}

function flashNotice(message) {
	var flash = [{Class: "alert-info", Message: message}];
	showFlash(flash);
}

function flashWarning(message) {
	var flash = [{Class: "alert-warning", Message: message}];
	showFlash(flash);
}