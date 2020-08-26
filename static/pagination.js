function pagination(total, current = 1) {
    $('.pagination').empty();
    if (total > 0) {
        var page = Math.ceil(total / 10);
        var structure = Array.from(new Set([1, 2, page, page - 1].concat(new Array(5).fill().map((d, i) =>
            i + current - 2)).sort((a, b) => { return a - b })))
        var $prev = $("<li class='page-item'></li>");
        if (current > 1) $prev.append(`<a class='page-link' data-page='${current - 1}'>${$.i18n('Previous')}</a>`);
        else $prev.addClass('disabled').append(`<a class='page-link'>${$.i18n('Previous')}</a>`);
        $('.pagination').append($prev);
        var flag = 0;
        $.each(structure, i => {
            var $li = $("<li class='page-item'></li>");
            if (i >= 1 && i <= page) {
                if (i - flag != 1) $li.append("<span class='page-link'>...</span>");
                else if (i == current) $li.addClass('active').append(`<span class='page-link'>${i}</span>`);
                else $li.append(`<span class='page-link' data-page='${i}'>${i}</span>`);
                $('.pagination').append($li);
            };
            flag = i;
        });
        var $next = $("<li class='page-item'></li>");
        if (current < page) $next.append(`<a class='page-link' data-page='${current + 1}'>${$.i18n('Next')}</a>`);
        else $next.addClass('disabled').append(`<a class='page-link'>${$.i18n('Next')}</a>`);
        $('.pagination').append($next);
        if (total == 1 && current == 1) return `${$.i18n('Showing')} 1 ${$.i18n('row')}`
        else if (page == current)
        return `${$.i18n('Showing')} ${(current - 1) * 10 + 1}-${total} ${$.i18n('of')} ${total} ${$.i18n('rows')}`
        else if (page >= current && current >= 1)
            return `${$.i18n('Showing')} ${(current - 1) * 10 + 1}-${current * 10} ${$.i18n('of')} ${total} ${$.i18n('rows')}`
    };
};

$(document).on('click', '.page-link', function () {
    var page = $(this).data('page');
    if (page !== undefined) {
        loading();
        $('.page-item').removeClass('active');
        var mode = $('.pagination').data('mode');
        var type = $('.pagination').data('type');
        var data = JSON.parse($('.pagination').data('data'));
        var promise
        if (type == 'empl') promise = loadEmpls(mode, page, data);
        else if (type == 'stat') promise = loadStats(mode, page, data);
        else promise = loadRecords(mode, page, data);
        promise.then(() => loading(false));
    };
});
