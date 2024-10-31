$(document).ready(function() {
    $(".btn-ballast").click(function(){        
        $.ajax({
            url : '/set-ballast',
            type : 'POST',
            data : JSON.stringify({
                'server' : Number($(this).attr("server")),
                'results': $(this).attr("results")
            }),
            async: true,
            dataType:'json',
            beforeSend:  function(){$("#loadingoverlay").fadeIn();},
            complete:  function(){$("#loadingoverlay").fadeOut();},
            success : function() {              
                alert('Ballast saved correctly');
            },
            error : function()
            {
                alert("Error saving ballast");
            }
        });
    });
});

$(document).ready(function() {
    $(".btn-clear-ballast").click(function(){        
        $.ajax({
            url : '/clear-ballast',
            type : 'POST',
            data : JSON.stringify({
                'server' : Number($(this).attr("server"))
            }),
            async: true,
            dataType:'json',
            beforeSend:  function(){$("#loadingoverlay").fadeIn();},
            complete:  function(){$("#loadingoverlay").fadeOut();},
            success : function() {              
                alert('Ballast clear correctly');
            },
            error : function()
            {
                alert("Error clearing ballast");
            }
        });
    });
});

$(document).ready(function() {    
    $(".btn-set-general-ballast").click(function(){        
        $.ajax({
            url : '/set-general-ballast',
            type : 'POST',
            data : JSON.stringify({
                'server' : Number($(this).attr("server"))
            }),
            async: true,
            dataType:'json',
            beforeSend:  function(){$("#loadingoverlay").fadeIn();},
            complete:  function(){$("#loadingoverlay").fadeOut();},
            success : function() {              
                alert('Ballast clear correctly');
            },
            error : function()
            {
                alert("Error setting general ballast");
            }
        });
    }); 
});