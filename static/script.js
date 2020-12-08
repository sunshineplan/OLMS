$(document).on('click', '#filter', () => {
    $('.sortable').removeClass('asc');
    $('.sortable').removeClass('desc');
    $('.sortable').addClass('default');
});

function cleanObj(obj) {
    $.each(obj, (k, v) => {
        if (v === undefined || v === '')
            delete obj[k]
    });
    return obj;
};

function postJSON(url, data, success) {
    if ($('.g-recaptcha-response').length)
        return grecaptcha.execute(sitekey, { action: url.replace('/', '') })
            .then(token => {
                data['g-recaptcha-response'] = token;
                return $.ajax({
                    url: url,
                    type: 'POST',
                    contentType: 'application/json',
                    data: JSON.stringify(data),
                    success: success,
                });
            });
    return $.ajax({
        url: url,
        type: 'POST',
        contentType: 'application/json',
        data: JSON.stringify(data),
        success: success
    });
};

function showEmpls(mode) {
    document.cookie = "Last=empl; Path=/";
    var url = '/empl';
    loading();
    $.get(url, html => {
        $('.content').html(html);
        document.title = $('.title').text() + ' - ' + $.i18n('OLMS');
    }).done(() => {
        loadEmpls(mode).then(() => loading(false));
        $('.sortable').addClass('default');
        getDepts('#dept');
    }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login' });
};