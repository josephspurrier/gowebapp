$(function() {
    $(document).foundation();
		
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
		var message = '<div data-alert id="flash-message" class="alert-box-fixed'
		+ rnum + ' alert-box-fixed alert-box '+v.cssclass+'">'
		+ v.message + '<a href="#" class="close">&times;</a></div>';
        $('#flash-container').prepend(message);
        hideFlash(rnum);
    });
}