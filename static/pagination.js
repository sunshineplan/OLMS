function pagination(total, current = 1) {
    var page = Math.ceil(total / 10);
    var prev;
    var structure = Array.from(new Set([1, 2, page, page - 1].concat(new Array(5).fill().map((d, i) =>
        i + current - 2)).sort((a, b) => { return a - b })))
    var $prev = $("<li class='page-item'></li>");
    if (current > 1) {
        if (current > 2) prev = current - 1; else 1;
        $prev.append("<a class='page-link' data-page='" + prev + "'>Previous</a>");
    } else $prev.addClass('disabled').append("<a class='page-link'>Previous</a>");
    $('.pagination').append($prev);
    var flag = 0;
    $.each(structure, i => {
        var $li = $("<li class='page-item'></li>");
        if (i >= 1 && i <= page) {
            if (i - flag != 1) $li.append("<span class='page-link'>...</span>");
            else if (i == current)
                $li.addClass('active').append("<span class='page-link' data-page='" + i + "'>" + i + "</span>");
            else $li.append("<span class='page-link' data-page='" + i + "'>" + i + "</span>");
            $('.pagination').append($li);
        };
        flag = i;
    });
    var $next = $("<li class='page-item'></li>");
    if (current < page) {
        $next.append("<a class='page-link' data-page='" + current + 1 + "'>Next</a>");
    } else $next.addClass('disabled').append("<a class='page-link'>Next</a>");
    $('.pagination').append($next);
};