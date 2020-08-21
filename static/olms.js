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
    if (all) $(element).empty().append($('<option>').text('All').val(''));
    else $(element).empty().append($('<option>').text(' -- select an employee -- ').prop('disabled', true).val(''));
    return postJSON('/get', cleanObj({ mode: 'admin', query: 'empls', dept: deptID }), json => {
        $.each(json.rows, (i, item) => $(element).append($('<option>').text(item.Realname).val(item.ID)));
        $(element).val('');
    });
};

function getYears(mode, userID, deptID) {
    var data;
    if (mode === undefined) data = { query: 'years' }
    else data = cleanObj({ mode: 'admin', query: 'years', empl: userID, dept: deptID });
    $('#year').empty().append($('<option>').text('All').val(''));
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
    }).catch(jqXHR => { if (jqXHR.status == 404) BootstrapButtons.fire('Info', 'No result.', 'info'); });
};

function loadDepts(mode) {
    return postJSON('/get', getData(mode, 'depts'), json => {
        $('tbody').empty();
        $.each(json.rows, (i, item) => {
            var $tr = $('<tr></tr>');
            $tr.append('<td>' + item.ID + '</td>');
            $tr.append('<td>' + item.Name + '</td>');
            $tr.append("<td><a class='btn btn-outline-primary btn-sm' onclick='dept(" + item.ID + ")'>Edit</a></td>");
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
        $.each(json.rows, (i, item) => {
            var $tr = $('<tr></tr>');
            $tr.append('<td>' + item.Username + '</td>');
            $tr.append('<td>' + item.Realname + '</td>');
            $tr.append('<td>' + item.DeptName + '</td>');
            if (mode == 'super') {
                if (item.Role == 0) $tr.append('<td>General Employee</td>');
                else $tr.append('<td>Administrator</td>');
                $tr.append('<td>' + item.Permission + '</td>');
                $tr.append("<td><a class='btn btn-outline-primary btn-sm' onclick='empl(" + item.ID + ")'>Edit</a></td>");
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
        $.each(json.rows, (i, item) => {
            var $tr = $('<tr></tr>');
            if (mode != '') {
                $tr.append('<td>' + item.DeptName + '</td>');
                $tr.append('<td>' + item.Name + '</td>');
            };
            $tr.append('<td>' + item.Date.replace(':00Z', '').replace(/-/g, '/').replace('T', ' ') + '</td>');
            if (item.Type == true) $tr.append('<td>Overtime</td>');
            else $tr.append('<td>Leave</td>');
            if (item.Duration == 0 || Math.abs(item.Duration) == 1)
                $tr.append('<td>' + item.Duration + ' Hour</td>');
            else $tr.append('<td>' + item.Duration + ' Hours</td>');
            $tr.append("<td class='describe'>" + item.Describe + '</td>');
            $tr.append('<td>' + item.Created.split('T')[0] + '</td>');
            if (item.Status == 0) $tr.append("<td><a class='text-muted'>Unverified</a></td>");
            else if (item.Status == 1) $tr.append("<td><a class='text-success'>Verified</a></td>");
            else if (item.Status == 2) $tr.append("<td><a class='text-danger'>Rejected</a></td>");
            if (mode == 'admin')
                if (item.Status == 0)
                    $tr.append("<td><a class='btn btn-outline-primary btn-sm' onclick='verify(" + item.ID + ")'>Verify</a></td>");
                else $tr.append("<td><a class='btn btn-outline-primary btn-sm disabled'>Verify</a></td>");
            else if (mode == '' && item.Status != 0)
                $tr.append("<td><a class='btn btn-outline-primary btn-sm disabled'>Edit</a></td>");
            else
                $tr.append("<td><a class='btn btn-outline-primary btn-sm' onclick='record(\"" + mode + "\"," + item.ID + ")'>Edit</a></td>");
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
        $.each(json.rows, (i, item) => {
            var $tr = $('<tr></tr>');
            $tr.append('<td>' + item.Period + '</td>');
            if (mode != '') {
                $tr.append('<td>' + item.DeptName + '</td>');
                $tr.append('<td>' + item.Name + '</td>');
            };
            if (item.Overtime == 0 || Math.abs(item.Overtime) == 1)
                $tr.append('<td>' + item.Overtime + ' Hour</td>');
            else $tr.append('<td>' + item.Overtime + ' Hours</td>');
            if (item.Leave == 0 || Math.abs(item.Leave) == 1)
                $tr.append('<td>' + item.Leave + ' Hour</td>');
            else $tr.append('<td>' + item.Leave + ' Hours</td>');
            if (item.Summary == 0 || Math.abs(item.Summary) == 1)
                $tr.append('<td>' + item.Summary + ' Hour</td>');
            else $tr.append('<td>' + item.Summary + ' Hours</td>');
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
        document.title = 'Departments List - OLMS';
    }).done(() => loadDepts('admin').then(() => loading(false)))
        .fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login' });
};

function showEmpls(mode) {
    document.cookie = "Last=empl; Path=/";
    var url = '/empl';
    loading();
    $.get(url, html => {
        $('.content').html(html);
        document.title = 'Employees List - OLMS';
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
        if (mode == 'admin') {
            document.title = 'Department Records - OLMS';
            $('.title').text('Department Records');
        } else if (mode == 'super') {
            document.title = 'All Records - OLMS';
            $('.title').text('All Records');
        } else {
            document.title = 'Employee Records - OLMS';
            $('.title').text('Employee Records');
        };
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
        if (mode == '') {
            document.title = 'Employee Statistics - OLMS';
            $('.title').text('Employee Statistics');
        } else {
            document.title = 'Department Statistics - OLMS';
            $('.title').text('Department Statistics');
        };
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
        if (id == 0) {
            document.title = 'Add Department - OLMS';
            $('.title').text('Add Department');
            loading(false);
        } else {
            document.title = 'Edit Department - OLMS';
            $('.title').text('Edit Department');
            postJSON('/get', { mode: 'admin', query: 'depts', id: id }, json =>
                $.each(json.dept, (k, v) => $('#' + k).val(v)))
                .then(() => loading(false));
        };
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
        if (id == 0) {
            document.title = 'Add Employee - OLMS';
            $('.title').text('Add Employee');
        } else {
            document.title = 'Edit Employee - OLMS';
            $('.title').text('Edit Employee');
        };
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
        if (id == 0) {
            document.title = 'Add Record - OLMS';
            $('.title').text('Add Record');
        } else {
            document.title = 'Edit Record - OLMS';
            $('.title').text('Edit Record');
        };
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
        document.title = 'Verify Record - OLMS';
        postJSON('/get', { mode: 'admin', id: id }, json => {
            $.each(json.record, (k, v) => $('#' + k).val(v));
            $('#Date').val(json.record.Date.replace(':00Z', '').replace('T', ' '));
            if (json.record.Type) $('#Type').val('Overtime');
            else $('#Type').val('Leave');
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
        document.title = 'Setting - OLMS';
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
                BootstrapButtons.fire('Error', json.message, 'error').then(() => $('#dept').val(''));
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
                BootstrapButtons.fire('Error', json.message, 'error').then(() => {
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
                BootstrapButtons.fire('Error', json.message, 'error')
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
        title: 'Are you sure?',
        text: 'This ' + type + ' will be deleted permanently.',
        icon: 'warning',
        confirmButtonText: 'Delete',
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
                    BootstrapButtons.fire('Error', json.message, 'error');
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
                BootstrapButtons.fire('Success', 'Your password has changed. Please Re-login!', 'success')
                    .then(() => window.location = '/auth/login');
            else BootstrapButtons.fire('Error', json.message, 'error').then(() => {
                if (json.error == 1) $('#password').val('');
                else if (json.error == 2) {
                    $('#password1').val('');
                    $('#password2').val('');
                };
            });
        }).catch(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login' });
    };
};
