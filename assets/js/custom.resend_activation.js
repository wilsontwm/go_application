$(document).ready(function(){
    
    $("#resend-activation-form").validate({
        rules: {
            email: {
                required: true,
                email: true
            }
        },
        messages: {
            email: {
                required: "Email is a mandatory field.",
                email: "Invalid email address."
            }
        },
        submitHandler: function(form) {
            form.submit();         
            toggleLoading();
        }
    });
});