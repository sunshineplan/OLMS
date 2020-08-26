BootstrapButtons = Swal.mixin({
    customClass: {
        confirmButton: 'swal btn btn-primary'
    },
    buttonsStyling: false
});

$(document).on('submit', '#login', () => {
    if ($('#username').val() != 'root') localStorage.setItem('username', $('#username').val());
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

function valid() {
    var result = true;
    $('input').each(function () {
        if ($(this)[0].checkValidity() === false) {
            $('.form').addClass('was-validated');
            result = false;
        };
    });
    return result;
};

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
