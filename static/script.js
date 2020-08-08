function loading(show = true) {
    if (show) {
        $('.loading').css('display', 'flex');
        $('.content').css('opacity', 0.5);
    } else {
        $('.loading').hide();
        $('.content').css('opacity', 1);
    };
};

$(() => {
    $(document).on('change', '#dept', () => {
        getEmpls('#empl', $('#dept').val());
        getYears(deptID = $('#dept').val());
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
