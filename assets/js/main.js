var loadingOverlay; 
$(document).ready(function(){
    loadingOverlay = document.querySelector('.loading');
    // Display the flash message
    window.Flash.create('.flash-message');

    $(".sidebar").slimScroll({
        height: 'auto',
        color: '#CCCCCC'
    });
});

// Show/hide the loading screen
function toggleLoading(){    
    document.activeElement.blur();
    if (loadingOverlay.classList.contains('hidden')){
        loadingOverlay.classList.remove('hidden');
    } else {
        loadingOverlay.classList.add('hidden');
    }
}
