$(document).ready(function(){
    $.validator.addMethod("password",function(value,element){
        return this.optional(element) || /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{8,16}$/i.test(value);
    },"Passwords are 8-16 characters with uppercase letters, lowercase letters and at least one number.");
    
    $("#login-form").validate({
        rules: {
            email: {
                required: true,
                email: true
            },
            password: {
                required: true,
                password: true
            }
        },
        messages: {
            email: {
                required: "Email is a mandatory field.",
                email: "Invalid email address."
            },
            password: {
                required: "Password is a mandatory field.",
                password: "Passwords are 8-16 characters with uppercase letters, lowercase letters and at least one number."
            }
        },
        submitHandler: function(form) {
            form.submit();            
            toggleLoading();
        }
    });
});