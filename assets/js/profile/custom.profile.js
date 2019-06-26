$(document).ready(function(){
    $.validator.addMethod("password",function(value,element){
        return this.optional(element) || /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{8,16}$/i.test(value);
    },"Passwords are 8-16 characters with uppercase letters, lowercase letters and at least one number.");
    
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

    $("#edit-password-form").validate({
        rules: {
            password: {
                required: true,
                password: true
            },
            retype_password: {
                equalTo: "#password_input"
            }
        },
        messages: {
            password: {
                required: "Password is a mandatory field.",
                password: "Passwords are 8-16 characters with uppercase letters, lowercase letters and at least one number."
            },
            retype_password: {
                equalTo: "Retype password does not match password."
            }
        },
        submitHandler: function(form) {
            form.submit();            
            toggleLoading();
        }
    });
});