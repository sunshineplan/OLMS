$(document).on('click', '.default', function () {
    sort($(this), false);
});

$(document).on('click', '.asc', function () {
    sort($(this), false);
});

$(document).on('click', '.desc', function () {
    sort($(this));
});

function sort(element, asc = true) {
    loading();
    $('.sortable').removeClass('asc');
    $('.sortable').removeClass('desc');
    $('.sortable').addClass('default');
    if (asc) element.addClass('asc'); else element.addClass('desc');
    var mode = $('.pagination').data('mode');
    var type = $('.pagination').data('type');
    var data = JSON.parse($('.pagination').data('data'));
    data.sort = element.data('name')
    if (asc) data.order = 'asc'; else data.order = 'desc';
    var promise
    if (type == 'empl') promise = loadEmpls(mode, 1, data);
    else if (type == 'stat') promise = loadStats(mode, 1, data);
    else promise = loadRecords(mode, 1, data);
    promise.then(() => loading(false));
}
