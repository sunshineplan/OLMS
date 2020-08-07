function init_selector(element, mode = 'init', text = 'All') {
    if (mode == 'unselect') {
        $(element).val('');
    } else if (mode == 'clear') {
        $(element).empty().append($('<option>').text(text).prop('value', '').prop('disabled', true));
    } else if (mode == 'init') {
        $(element).empty().append($('<option>').text(text).prop('value', ''));
    }
}

function loading(show = true) {
    if (show) {
        $('.loading').css('display', 'flex');
        $('.content').css('opacity', 0.5);
    } else {
        $('.loading').hide();
        $('.content').css('opacity', 1);
    }
};
