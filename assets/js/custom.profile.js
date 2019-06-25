$(document).ready(function(){
    $("#edit-profile-form").validate({
        rules: {
            name: {
                required: true
            }
        },
        messages: {
            name: {
                required: "Name is a mandatory field."
            }
        },
        submitHandler: function(form) {
            form.submit();            
            toggleLoading();
        }
    });
});