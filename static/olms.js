function getData(mode, type) {
    var obj = {};
    $('select, input').serializeArray().forEach(i => { if (i.value != '') obj[i.name] = i.value });
    return cleanObj($.extend(obj, { mode: mode, query: type }));
};

function getDepts(element) {
    return postJSON('/get', { mode: 'admin', query: 'depts' }, json =>
        $.each(json.rows, (i, item) => $(element).append($('<option>').text(item.Name).val(item.ID))));
};

function getEmpls(element, deptID, all = true) {
    if (all) $(element).empty().append($('<option>').text($.i18n('All')).val(''));
    else $(element).empty().append($('<option>').text(` -- ${$.i18n('SelectEmployee')} -- `).prop('disabled', true).val(''));
    return postJSON('/get', cleanObj({ mode: 'admin', query: 'empls', dept: deptID }), json => {
        $.each(json.rows, (i, item) => $(element).append($('<option>').text(item.Realname).val(item.ID)));
        $(element).val('');
    });
};

function getYears(mode, userID, deptID) {
    var data;
    if (mode === undefined) data = { query: 'years' }
    else data = cleanObj({ mode: 'admin', query: 'years', empl: userID, dept: deptID });
    $('#year').empty().append($('<option>').text($.i18n('All')).val(''));
    return postJSON('/get', data, json => {
        $.each(json.rows, (i, item) => $('#year').append($('<option>').text(item).val(item)));
        $('#year').val('');
    });
};

function getInfo() {
    return postJSON('/get', { mode: 'admin', query: 'info' }, json => {
        $.each(json.depts, (i, item) => $("#dept").append($('<option>').text(item.Name).val(item.ID)));
        $.each(json.empls, (i, item) => $("#empl").append($('<option>').text(item.Realname).val(item.ID)));
        $.each(json.years, (i, item) => $('#year').append($('<option>').text(item).val(item)));
    });
};

function exportCSV(mode, type) {
    if (mode == 'super') mode = 'admin';
    postJSON('/export', getData(mode, type), (data, status, jqXHR) => {
        var blob = new Blob([new Uint8Array([0xEF, 0xBB, 0xBF]), data], { type: 'text/csv;charset=utf-8' });
        var link = document.createElement('a');
        link.href = window.URL.createObjectURL(blob);
        link.download = decodeURI(jqXHR.getResponseHeader('Content-Disposition').split('filename=')[1].replace(/"/g, ''));
        link.click();
    }).catch(jqXHR => { if (jqXHR.status == 404) BootstrapButtons.fire($.i18n('Info'), $.i18n('NoResult'), 'info'); });
};

function loadDepts(mode) {
    return postJSON('/get', getData(mode, 'depts'), json => {
        $('tbody').empty();
        $.each(json.rows, (index, i) => {
            var $tr = $('<tr></tr>');
            $tr.append('<td>' + i.ID + '</td>');
            $tr.append('<td>' + i.Name + '</td>');
            $tr.append(
                `<td><a class='btn btn-outline-primary btn-sm' onclick='dept("${i.ID}")'>${$.i18n('Edit')}</a></td>`);
            $tr.appendTo('tbody');
        });
    });
};

function loadEmpls(mode, page = 1, data) {
    if (mode == 0 || mode == 'super') mode = 'super'; else mode = 'admin';
    if (data === undefined) data = getData(mode, 'empls');
    return postJSON('/get', $.extend(data, { page: page }), json => {
        pagination(json.total, page);
        $('tbody').empty();
        $("#total").text(json.total);
        $.each(json.rows, (index, i) => {
            var $tr = $('<tr></tr>');
            $tr.append('<td>' + i.Username + '</td>');
            $tr.append('<td>' + i.Realname + '</td>');
            $tr.append('<td>' + i.DeptName + '</td>');
            if (mode == 'super') {
                if (i.Role == 0) $tr.append('<td>' + $.i18n('GeneralEmployee') + '</td>');
                else $tr.append('<td>' + $.i18n('Administrator') + '</td>');
                $tr.append('<td>' + i.Permission + '</td>');
                $tr.append(
                    `<td><a class='btn btn-outline-primary btn-sm' onclick='empl("${i.ID}")'>${$.i18n('Edit')}</a></td>`);
            };
            $tr.appendTo('tbody');
        });
    }).then(() => {
        $('.pagination').data('mode', mode);
        $('.pagination').data('type', 'empl');
        $('.pagination').data('data', JSON.stringify(data));
    });
};

function loadRecords(mode, page = 1, data) {
    if (data === undefined)
        if (mode == 'super') data = getData('admin', 'records');
        else data = getData(mode, 'records');
    return postJSON('/get', $.extend(data, { page: page }), json => {
        pagination(json.total, page);
        $('tbody').empty();
        $("#total").text(json.total);
        $.each(json.rows, (index, i) => {
            var $tr = $('<tr></tr>');
            if (mode != '') {
                $tr.append('<td>' + i.DeptName + '</td>');
                $tr.append('<td>' + i.Name + '</td>');
            };
            $tr.append('<td>' + i.Date.replace(':00Z', '').replace(/-/g, '/').replace('T', ' ') + '</td>');
            if (i.Type == true) $tr.append('<td>' + $.i18n('Overtime') + '</td>');
            else $tr.append('<td>' + $.i18n('Leave') + '</td>');
            if (i.Duration == 0 || Math.abs(i.Duration) == 1)
                $tr.append('<td>' + i.Duration + ' ' + $.i18n('Hour') + '</td>');
            else $tr.append('<td>' + i.Duration + ' ' + $.i18n('Hours') + '</td>');
            $tr.append("<td class='describe'>" + i.Describe + '</td>');
            $tr.append('<td>' + i.Created.split('T')[0] + '</td>');
            if (i.Status == 0) $tr.append(`<td><a class='text-muted'>${$.i18n('Unverified')}</a></td>`);
            else if (i.Status == 1) $tr.append(`<td><a class='text-success'>${$.i18n('Verified')}</a></td>`);
            else if (i.Status == 2) $tr.append(`<td><a class='text-danger'>${$.i18n('Rejected')}</a></td>`);
            if (mode == 'admin')
                if (i.Status == 0)
                    $tr.append(
                        `<td><a class='btn btn-outline-primary btn-sm' onclick='verify("${i.ID}")'>${$.i18n('Verify')}</a></td>`);
                else $tr.append(`<td><a class='btn btn-outline-primary btn-sm disabled'>${$.i18n('Verify')}</a></td>`);
            else if (mode == '' && i.Status != 0)
                $tr.append(`<td><a class='btn btn-outline-primary btn-sm disabled'>${$.i18n('Edit')}</a></td>`);
            else
                $tr.append(
                    `<td><a class='btn btn-outline-primary btn-sm' onclick='record("${mode}","${i.ID}")'>${$.i18n('Edit')}</a></td>`);
            $tr.appendTo('tbody');
        });
    }).then(() => {
        $('.pagination').data('mode', mode);
        $('.pagination').data('type', 'record');
        $('.pagination').data('data', JSON.stringify(data));
    });
};

function loadStats(mode, page = 1, data) {
    if (data === undefined) data = getData(mode, 'stats');
    return postJSON('/get', $.extend(data, { page: page }), json => {
        pagination(json.total, page);
        $('tbody').empty();
        $.each(json.rows, (index, i) => {
            var $tr = $('<tr></tr>');
            $tr.append('<td>' + i.Period + '</td>');
            if (mode != '') {
                $tr.append('<td>' + i.DeptName + '</td>');
                $tr.append('<td>' + i.Name + '</td>');
            };
            if (i.Overtime == 0 || Math.abs(i.Overtime) == 1)
                $tr.append('<td>' + i.Overtime + ' ' + $.i18n('Hour') + '</td>');
            else $tr.append('<td>' + i.Overtime + ' ' + $.i18n('Hours') + '</td>');
            if (i.Leave == 0 || Math.abs(i.Leave) == 1)
                $tr.append('<td>' + i.Leave + ' ' + $.i18n('Hour') + '</td>');
            else $tr.append('<td>' + i.Leave + ' ' + $.i18n('Hours') + '</td>');
            if (i.Summary == 0 || Math.abs(i.Summary) == 1)
                $tr.append('<td>' + i.Summary + ' ' + $.i18n('Hour') + '</td>');
            else $tr.append('<td>' + i.Summary + ' ' + $.i18n('Hours') + '</td>');
            $tr.appendTo('tbody');
        });
    }).then(() => {
        $('.pagination').data('mode', mode);
        $('.pagination').data('type', 'stat');
        $('.pagination').data('data', JSON.stringify(data));
    });
};

function showDepts() {
    document.cookie = "Last=dept; Path=/";
    var url = '/dept';
    loading();
    $.get(url, html => {
        $('.content').html(html);
        document.title = $('.title').text() + ' - ' + $.i18n('OLMS');
    }).done(() => loadDepts('admin').then(() => loading(false)))
        .fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login' });
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

function showRecords(mode) {
    document.cookie = "Last=record; Path=/";
    var url;
    if (mode == 'admin') url = '/record/admin';
    else if (mode == 'super') url = '/record/super';
    else url = '/record';
    loading();
    $.get(url, html => {
        $('.content').html(html);
        document.title = $('.title').text() + ' - ' + $.i18n('OLMS');
    }).done(() => {
        loadRecords(mode).then(() => loading(false));
        $('.sortable').addClass('default');
        if (mode == '') getYears();
        else getInfo();
    }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login' });
};

function showStats(mode) {
    var url;
    if (mode == '') url = '/stat';
    else url = '/stat/admin';
    loading();
    $.get(url, html => {
        $('.content').html(html);
        document.title = $('.title').text() + ' - ' + $.i18n('OLMS');
    }).done(() => {
        loadStats(mode).then(() => loading(false));
        $('.sortable').addClass('default');
        if (mode == '') getYears();
        else getInfo();
    }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login' });
};

function dept(id = 0) {
    var url;
    if (id == 0) url = '/dept/add';
    else url = '/dept/edit/' + id;
    loading();
    $.get(url, html => $('.content').html(html)).done(() => {
        document.title = $('.title').text() + ' - ' + $.i18n('OLMS');
        if (id == 0) loading(false);
        else postJSON('/get', { mode: 'admin', query: 'depts', id: id }, json =>
            $.each(json.dept, (k, v) => $('#' + k).val(v)))
            .then(() => loading(false));
        $('#Name').focus();
    }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login' });
};

function empl(id = 0) {
    var url;
    if (id == 0) url = '/empl/add';
    else url = '/empl/edit/' + id;
    loading();
    $.get(url, html => $('.content').html(html)).done(() => {
        getDepts('#Dept, #Permission');
        document.title = $('.title').text() + ' - ' + $.i18n('OLMS');
    }).done(() => {
        if (id != 0) {
            postJSON('/get', { mode: 'super', query: 'empls', id: id }, json => {
                $.each(json.empl, (k, v) => $('#' + k).val(v));
                $('#Dept').val(json.empl.DeptID);
                if (json.empl.Role) {
                    $('#Role').val(1);
                    $('#permission-selector').prop('hidden', false);
                    $('#Permission').val(json.empl.Permission.split(','));
                } else $('#Role').val(0);
            }).then(() => loading(false));
        } else loading(false);
        $('#Username').focus();
    }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login' });
};

function record(mode = '', id = 0) {
    var url;
    if (id == 0) {
        if (mode == '') url = '/record/add';
        else url = '/record/admin/add';
    } else {
        if (mode == '') url = '/record/edit/' + id;
        else url = '/record/super/edit/' + id;
    };
    loading();
    $.get(url, html => $('.content').html(html)).done(() => {
        if (mode != '') getDepts('#Dept');
        document.title = $('.title').text() + ' - ' + $.i18n('OLMS');
    }).done(() => {
        if (mode != '') mode = 'admin';
        if (id != 0)
            postJSON('/get', cleanObj({ id: id, mode: mode }), json => {
                $.each(json.record, (k, v) => $('#' + k).val(v));
                if (mode != '') {
                    $('#Dept').val(json.record.DeptID);
                    getEmpls('#Empl', json.record.DeptID, false)
                        .then(() => $('#Empl').val(json.record.UserID))
                        .then(() => loading(false));
                } else loading(false);
                $('#Date').val(json.record.Date.replace(':00Z', ''));
                if (json.record.Type) $('#Type').val('1');
                else $('#Type').val('0');
            });
        else if (mode != '') getEmpls('#Empl', undefined, false).then(() => loading(false));
        else loading(false);
    }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login' });
};

function verify(id) {
    loading();
    $.get('/record/verify/' + id, html => $('.content').html(html)).done(() => {
        document.title = $('.title').text() + ' - ' + $.i18n('OLMS');
        postJSON('/get', { mode: 'admin', id: id }, json => {
            $.each(json.record, (k, v) => $('#' + k).val(v));
            $('#Date').val(json.record.Date.replace(':00Z', '').replace('T', ' '));
            if (json.record.Type) $('#Type').val($.i18n('Overtime'));
            else $('#Type').val($.i18n('Leave'));
        }).then(() => loading(false));
    }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login' });
};

function setting() {
    document.cookie = "Last=/; Path=/";
    loading();
    $('li>a.nav-link').removeClass('selected');
    $.get('/auth/setting', html => {
        loading(false);
        $('.content').html(html);
        document.title = $('.title').text() + ' - ' + $.i18n('OLMS');
        $('#password').focus();
    }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login' });
};

function doDept(id) {
    var url;
    if (id == 0) url = '/dept/add';
    else url = '/dept/edit/' + id;
    if (valid())
        $.post(url, $('input').serialize(), json => {
            $('.form').removeClass('was-validated');
            if (json.status == 0)
                BootstrapButtons.fire($.i18n('Error'), json.message, 'error')
                    .then(() => $('#dept').val(''));
            else showDepts();
        }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login' });
};

function doEmpl(mode, id) {
    var url;
    if (id == 0) url = '/empl/add';
    else url = '/empl/edit/' + id;
    if (valid())
        $.post(url, $('input, select').serialize(), json => {
            $('.form').removeClass('was-validated');
            if (json.status == 0)
                BootstrapButtons.fire($.i18n('Error'), json.message, 'error').then(() => {
                    if (json.error == 1) $('#username').val('');
                });
            else showEmpls(mode);
        }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login' });
};

function doRecord(mode, id) {
    var url, jqXHR;
    if (id == 0) url = '/record/add';
    else url = '/record/edit/' + id;
    if (valid()) {
        if ($('.g-recaptcha-response').length)
            jqXHR = grecaptcha.execute(sitekey, { action: 'record' })
                .then(() => $.post(url, $('input, select, textarea').serialize()));
        else jqXHR = $.post(url, $('input, select, textarea').serialize());
        jqXHR.then(json => {
            $('.form').removeClass('was-validated');
            if (json.status == 0)
                BootstrapButtons.fire($.i18n('Error'), json.message, 'error')
                    .then(() => $('#duration').val(''));
            else showRecords(mode);
        }).catch(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login' });
    };
};

function doVerify(id, status) {
    var jqXHR;
    if ($('.g-recaptcha-response').length)
        jqXHR = grecaptcha.execute(sitekey, { action: 'verify' })
            .then(token => $.post('/record/verify/' + id, { status: status, 'g-recaptcha-response': token }));
    else jqXHR = $.post('/record/verify/' + id, { status: status });
    jqXHR.then(() => showRecords('admin'))
        .catch(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login' });
};

function doDelete(type, id) {
    var url;
    if (type == 'dept') url = '/dept/delete/' + id;
    else if (type == 'empl') url = '/empl/delete/' + id;
    else if (type == 'record') url = '/record/delete/' + id;
    Swal.fire({
        title: $.i18n('AreYouSure'),
        text: $.i18n('DeleteWarning').replace('%s', $.i18n(type)),
        icon: 'warning',
        confirmButtonText: $.i18n('Delete'),
        cancelButtonText: $.i18n('Cancel'),
        showCancelButton: true,
        focusCancel: true,
        customClass: {
            confirmButton: 'swal btn btn-danger',
            cancelButton: 'swal btn btn-primary'
        },
        buttonsStyling: false
    }).then(confirm => {
        if (confirm.isConfirmed)
            $.post(url, json => {
                if (json.status == 0)
                    BootstrapButtons.fire($.i18n('Error'), json.message, 'error');
                else window.location = '/';
            }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login' });
    });
};

function doSetting() {
    if (valid()) {
        var jqXHR;
        if ($('.g-recaptcha-response').length)
            jqXHR = grecaptcha.execute(sitekey, { action: 'setting' })
                .then(() => $.post('/auth/setting', $('input, textarea').serialize()));
        else jqXHR = $.post('/auth/setting', $('input').serialize());
        jqXHR.then(json => {
            $('.form').removeClass('was-validated');
            if (json.status == 1)
                BootstrapButtons.fire($.i18n('Success'), $.i18n('PasswordChanged'), 'success')
                    .then(() => window.location = '/auth/login');
            else BootstrapButtons.fire($.i18n('Error'), json.message, 'error').then(() => {
                if (json.error == 1) $('#password').val('');
                else if (json.error == 2) {
                    $('#password1').val('');
                    $('#password2').val('');
                };
            });
        }).catch(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login' });
    };
};
