BootstrapButtons = Swal.mixin({
    customClass: {
        confirmButton: 'swal btn btn-primary'
    },
    buttonsStyling: false
});

$(document).on('click', 'li>a.nav-link', function () {
    $('li>a.nav-link').removeClass('selected');
    $(this).addClass('selected');
    if ($(window).width() <= 900) $('.sidebar').toggle('slide');
});

$(document).on('change', '#dept', () => {
    if ($('#empl').length) getEmpls('#empl', $('#dept').val());
    if ($('#year').length) getYears(deptID = $('#dept').val());
});

$(document).on('change', '#empl', () => { getYears(userID = $('#empl').val()) });

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

$(document).on('click', '.toggle', () => $('.sidebar').toggle('slide'));

$(document).on('click', '.content', () => {
    if ($('.sidebar').is(':visible') && $(window).width() <= 900)
        $('.sidebar').toggle('slide');
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
