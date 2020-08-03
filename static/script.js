function show_empl(element, dept_id, value = '') {
    $.post('/manage/empl/get', { dept: dept_id }, function (json) {
        $.each(json, function (id, realname) {
            $(element).append($('<option>').text(realname).prop('value', id));
        });
        $(element).val(value)
    }, 'json').fail(function (jqXHR) {
        if (jqXHR.status == 200) {
            alert('Session timeout. Please re-login!');
            $(location).attr('href', '/auth/login');
        } else {
            alert('Getting employee list failed.');
        };
    });
}

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
