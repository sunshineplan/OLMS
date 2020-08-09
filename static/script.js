BootstrapButtons = Swal.mixin({
    customClass: {
        confirmButton: 'swal btn btn-primary'
    },
    buttonsStyling: false
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

function goback() {
    var last = document.cookie.split('LastVisit=')[1];
    show(last);
};

$(() => {
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
});
