$(function(){
    var fileInput = $('#myfiles');
    var maxSize = fileInput.data('max-size');
    $('.upload-form').submit(function(e){
        if(fileInput.get(0).files.length){
            var fileSize = fileInput.get(0).files[0].size; // in bytes
            if(fileSize>maxSize){
                alert('Размер файлов слишком большой');
                return false;
            }
        }
    });
});

$(document).ready(function(){
    $(".delete-file button").click(function(){
        var img_id = $(this).parent().attr('id');
        console.log(img_id)
            $.ajax({
                type: "POST",
                url: "delete-image",
                data: {filepath: img_id},                        
                success: function() {
                        console.log("success");
                        document.getElementById(`${img_id}`).remove()
                }
            }); 
        });
    });
