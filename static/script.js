BootstrapButtons = Swal.mixin({
    customClass: {
        confirmButton: 'swal btn btn-primary'
    },
    confirmButtonText: $.i18n('OK'),
    buttonsStyling: false
});

$(document).on('submit', '#login', () => {
    if ($('#username').val() != 'root') localStorage.setItem('username', $('#username').val());
});

$(document).on('change', '#lang', () => {
    document.cookie = `lang=${$('#lang').val()}; Path=/; max-age=31536000`;
    BootstrapButtons.fire($.i18n('Success'), $.i18n('LanguageChanged'), 'success')
        .then(() => window.location = '/');
});

$(document).on('input', '#email', () => $('#subscribe').prop('checked', false));

$(document).on('change', '#subscribe', () => {
    var data;
    if ($('#subscribe').is(':checked')) {
        var email = $('#email').val();
        if (validateEmail(email)) data = { subscribe: 1, email: email };
        else BootstrapButtons.fire($.i18n('Error'), $.i18n('EmailNotValid'), 'error').then(() => {
            $('#email').val('');
            $('#subscribe').prop('checked', false);
        });
    } else data = { subscribe: 0 };
    if (data === undefined) return false;
    $.post('/subscribe', data).done(json => {
        if (json.status == 1)
            BootstrapButtons.fire($.i18n('Success'), $.i18n('SubscribeChanged'), 'success');
        else BootstrapButtons.fire($.i18n('Error'), $.i18n('EmailNotValid'), 'error').then(() => {
            $('#email').val('');
            $('#subscribe').prop('checked', false);
        });
    });
});

$(document).on('click', 'li>a.nav-link', function () {
    $('li>a.nav-link').removeClass('selected');
    $(this).addClass('selected');
    if ($(window).width() <= 1200) $('.sidebar').toggle('slide');
});

$(document).on('change', '#dept', () => {
    if ($('#empl').length) getEmpls('#empl', $('#dept').val());
    if ($('#year').length) getYears(deptID = $('#dept').val());
});

$(document).on('change', '#Dept', () => getEmpls('#Empl', $('#Dept').val(), false));

$(document).on('change', '#empl', () => getYears(userID = $('#empl').val()));

$(document).on('change', '#Role', function () {
    if ($('#Role').val() == 0) {
        $('#permission-selector option:selected').prop('selected', false);
        $('#permission-selector').prop('hidden', true);
    } else {
        $('#permission-selector').prop('hidden', false);
    };
});

$(document).on('change', '#period', () => {
    if ($('#period').val() == 'year') {
        $('#month-selector').prop('hidden', true);
        $('#year').val('');
        $('#month').val('');
    } else {
        $('#month-selector').prop('hidden', false);
        $('#year').val('');
    };
});

$(document).on('change', '#year', () => {
    if ($('#year').val() == '') $('#month').prop('disabled', true).val('');
    else $('#month').prop('disabled', false);
});

$(document).on('click', '#filter', () => {
    $('.sortable').removeClass('asc');
    $('.sortable').removeClass('desc');
    $('.sortable').addClass('default');
});

$(document).on('click', '.toggle', () => $('.sidebar').toggle('slide'));

$(document).on('click', '.content', () => {
    if ($('.sidebar').is(':visible') && $(window).width() <= 1200)
        $('.sidebar').toggle('slide');
});

$(document).on('keyup', event => {
    if (event.key == 'Enter')
        if (!$('.swal2-container').length)
            $('#submit').click()
});

function loading(show = true) {
    if (show) {
        $('.loading').css('display', 'flex');
        $('.content').css('opacity', 0.5);
    } else {
        $('.loading').hide();
        $('.content').css('opacity', 1);
    };
};

function validate() {
    var result = true;
    $('input').each(function () {
        if ($(this)[0].checkValidity() === false) {
            $('.form').addClass('was-validated');
            result = false;
        };
    });
    return result;
};

function validateEmail(email) {
    //https://www.w3.org/TR/2016/REC-html51-20161101/sec-forms.html#email-state-typeemail
    const re = /^[a-zA-Z0-9.!#$%&'*+\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/;
    return re.test(email);
}

function goback(mode) {
    var last = document.cookie.split('Last=')[1];
    if (last == '/') window.location = '/';
    else if (last == 'dept') showDepts();
    else if (last == 'empl') showEmpls(mode);
    else showRecords(mode);
};

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
