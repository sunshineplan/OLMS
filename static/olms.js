function getParams(mode, query) {
    var data = {}, param;
    $('select').serializeArray().forEach(i => { if (i.value != '') data[i.name] = i.value });
    if ($.isEmptyObject(data)) param = 'mode=' + mode + '&query=' + query;
    else param = $.param(data) + '&mode=' + mode + '&query=' + query;
    return param;
};

function getDepts(element) {
    $.post('/get', 'mode=admin&query=depts', json =>
        $.each(json.rows, (i, item) => $(element).append($('<option>').text(item.Name).val(item.ID))));
};

function getEmpls(element, deptID, all = true) {
    var param
    if (deptID === undefined) param = 'mode=admin&query=empls'
    else param = 'mode=admin&query=empls&dept=' + deptID
    if (all) $(element).empty().append($('<option>').text('All').val(''));
    else $(element).empty().append($('<option>').text(' -- select an employee -- ').prop('disabled', true).val(''));
    $.post('/get', param, json =>
        $.each(json.rows, (i, item) => $(element).append($('<option>').text(item.Realname).val(item.ID))));
    $(element).val('');
};

function getYears(mode, userID, deptID) {
    var param
    if (userID !== undefined) param = 'mode=admin&query=years&empl=' + userID
    else if (deptID !== undefined) param = 'mode=admin&query=years&dept=' + deptID
    else if (mode === undefined) param = 'query=years'
    else param = 'mode=admin&query=years'
    $('#year').empty().append($('<option>').text('All').val(''));
    $.post('/get', param, json =>
        $.each(json.rows, (i, item) => $('#year').append($('<option>').text(item).val(item))));
    $('#year').val('');
};

function exportCSV(mode, query) {
    if (mode == 'super') mode = 'admin'
    $.post('/export', getParams(mode, query), (data, status, jqXHR) => {
        var blob = new Blob([data]);
        var link = document.createElement('a');
        link.href = window.URL.createObjectURL(blob);
        link.download = jqXHR.getResponseHeader('Content-Disposition').split('filename=')[1].replace(/"/g, '');
        link.click();
    });
};

function loadDepts(mode) {
    $.post('/get', getParams(mode, 'depts'), json => {
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

function loadEmpls(mode, page = 1) {
    if (mode == 0) mode = 'super';
    else mode = 'admin'
    $.post('/get', getParams(mode, 'empls') + '&page=' + page, json => {
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
    });
};

function loadRecords(mode, page = 1) {
    var param
    if (mode == 'super') param = getParams('admin', 'records');
    else param = getParams(mode, 'records')
    $.post('/get', param + '&page=' + page, json => {
        $('tbody').empty();
        $("#total").text(json.total);
        $.each(json.rows, (i, item) => {
            var $tr = $('<tr></tr>');
            if (mode != '') {
                $tr.append('<td>' + item.DeptName + '</td>');
                $tr.append('<td>' + item.Name + '</td>');
            };
            $tr.append('<td>' + item.Date.split('T')[0] + '</td>');
            if (item.Type == true) $tr.append('<td>Overtime</td>');
            else $tr.append('<td>Leave</td>');
            $tr.append('<td>' + item.Duration + '</td>');
            $tr.append('<td>' + item.Describe + '</td>');
            $tr.append('<td>' + item.Created.split('T')[0] + '</td>');
            if (item.Status == 0) $tr.append("<td><a class='text-muted'>Unverified</a></td>");
            else if (item.Status == 1) $tr.append("<td><a class='text-success'>Verified</a></td>");
            else if (item.Status == 2) $tr.append("<td><a class='text-danger'>Rejected</a></td>");
            if (mode == 'admin')
                if (item.Status == 0)
                    $tr.append("<td><a class='btn btn-outline-primary btn-sm' onclick='verify(" + item.ID + ")'>Verify</a></td>");
                else
                    $tr.append("<td><a class='btn btn-outline-primary btn-sm disabled'>Verify</a></td>");
            else if (mode == '' && item.Status != 0)
                $tr.append("<td><a class='btn btn-outline-primary btn-sm disabled'>Edit</a></td>");
            else
                $tr.append("<td><a class='btn btn-outline-primary btn-sm' onclick='record(\"" + mode + "\"," + item.ID + ")'>Edit</a></td>");
            $tr.appendTo('tbody');
        });
    });
};

function loadStats(mode, page = 1) {
    $.post('/get', getParams(mode, 'stats') + '&page=' + page, json => {
        $('tbody').empty();
        $.each(json.rows, (i, item) => {
            var $tr = $('<tr></tr>');
            $tr.append('<td>' + item.Period + '</td>');
            if (mode != '') {
                $tr.append('<td>' + item.DeptName + '</td>');
                $tr.append('<td>' + item.Name + '</td>');
            };
            $tr.append('<td>' + item.Overtime + '</td>');
            $tr.append('<td>' + item.Leave + '</td>');
            $tr.append('<td>' + item.Summary + '</td>');
            $tr.appendTo('tbody');
        });
    });
};

function showDepts() {
    document.cookie = "Last=dept; Path=/";
    var url = '/dept';
    loading();
    $.get(url, html => {
        loading(false);
        $('.content').html(html);
        document.title = 'Departments List - OLMS';
    }).done(() => loadDepts('admin')).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login' });
};

function showEmpls(mode) {
    document.cookie = "Last=empl; Path=/";
    var url = '/empl';
    loading();
    $.get(url, html => {
        loading(false);
        $('.content').html(html);
        document.title = 'Employees List - OLMS';
    }).done(() => {
        loadEmpls(mode);
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
        loading(false);
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
        loadRecords(mode);
        if (mode == '') getYears();
        else {
            getDepts('#dept');
            getEmpls('#empl');
            getYears('admin');
        };
    }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login' });
};

function showStats(mode) {
    var url;
    if (mode == '') url = '/stat';
    else url = '/stat/admin';
    loading();
    $.get(url, html => {
        loading(false);
        $('.content').html(html);
        if (mode == '') {
            document.title = 'Employee Statistics - OLMS';
            $('.title').text('Employee Statistics');
        } else {
            document.title = 'Department Statistics - OLMS';
            $('.title').text('Department Statistics');
        };
    }).done(() => {
        loadStats(mode);
        if (mode == '') getYears();
        else {
            getDepts('#dept');
            getEmpls('#empl');
            getYears('admin');
        };
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
            $.post('get', 'mode=admin&query=depts&id=' + id, json => {
                loading(false);
                $.each(json.dept, (k, v) => $('#' + k).val(v));
                $('#dept').focus();
            });
        };
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
            $.post('get', 'mode=super&query=empls&id=' + id, json => {
                $.each(json.empl, (k, v) => $('#' + k).val(v));
                $('#Dept').val(json.empl.DeptID);
                if (json.empl.Role) {
                    $('#Role').val(1);
                    $('#permission-selector').prop('hidden', false);
                    $('#Permission').val(json.empl.Permission.split(',')); // need improve run order
                } else $('#Role').val(0);
            });
        };
        loading(false);
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
        if (id != 0) $.post('get', 'id=' + id + '&mode=' + mode, json => {
            $.each(json.record, (k, v) => $('#' + k).val(v));
            $('#Dept').val(json.record.DeptID);
            getEmpls('#Empl', json.record.DeptID, false);
            $('#Empl').val(json.record.UserID);
            $('#Date').val(json.record.Date.split('T')[0]);
            if (json.record.Type) $('#Type').val('1');
            else $('#Type').val('0');
        }, 'json');
        else getEmpls('#Empl', undefined, false);
        loading(false);
    }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login' });
};

function verify(id) {
    loading();
    $.get('/record/verify/' + id, html => $('.content').html(html)).done(() => {
        loading(false);
        document.title = 'Verify Record - OLMS';
        $.post('get', '&mode=admin&id=' + id, json => {
            $.each(json.record, (k, v) => $('#' + k).val(v));
            $('#Date').val(json.record.Date.split('T')[0]);
            if (json.record.Type) $('#Type').val('Overtime');
            else $('#Type').val('Leave');
        }, 'json');
    }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login' });
};

function setting() {
    document.cookie = "Last=/; Path=/";
    loading();
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
    var url;
    if (id == 0) url = '/record/add';
    else if (mode == '') url = '/record/edit/' + id;
    else url = '/record/super/edit/' + id;
    if (valid())
        $.post(url, $('input, select, textarea').serialize(), json => {
            $('.form').removeClass('was-validated');
            if (json.status == 0)
                BootstrapButtons.fire('Error', json.message, 'error').then(() => $('#duration').val(''));
            else showRecords(mode);
        }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login' });
};

function doVerify(id, status) {
    $.post('/record/verify/' + id, 'status=' + status, () =>
        showRecords('admin')).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login' });
};

function doDelete(mode, id) {
    var url;
    if (mode == 'dept') url = '/dept/delete/' + id;
    else if (mode == 'empl') url = '/empl/delete/' + id;
    else if (mode == 'record') url = '/record/delete/' + id;
    else return false;
    Swal.fire({
        title: 'Are you sure?',
        text: 'This ' + mode + ' will be deleted permanently.',
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
                if (json.status == 1)
                    if (mode == 'dept') showDepts();
                    else if (mode == 'empl') showEmpls;
                    else window.location = '/';
            }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login' });
    });
};

function doSetting() {
    if (valid())
        $.post('/auth/setting', $('input').serialize(), json => {
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
        }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login' });
};
