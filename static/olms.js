BootstrapButtons = Swal.mixin({
    customClass: {
        confirmButton: 'swal btn btn-primary'
    },
    buttonsStyling: false
});

function getParam(mode, query) {
    var data = {}, param;
    $('select').serializeArray().forEach(i => { if (i.value != '') data[i.name] = i.value });
    if ($.isEmptyObject(data)) param = 'mode=' + mode + '&query=' + query;
    else param = $.param(data) + '&mode=' + mode + '&query=' + query;
    return param;
}

function load(mode, query) {
    var param = getParam(mode, query);
    $.post('/get', param, json => {
        $('tbody').empty();
        $.each(json.rows, (i, item) => {
            var $tr = $('<tr></tr>');
            $.each(item, (k, v) => {
                $tr.append('<td>' + v + '</td>');
            });
            $tr.appendTo('tbody');
        });
    });
    getDepts('#dept')
};

function loadDepts(mode) {
    var param = getParam(mode, 'depts');
    $.post('/get', param, json => {
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

function loadEmpls(mode) {
    var param = getParam(mode, 'empls');
    $.post('/get', param, json => {
        $('tbody').empty();
        $.each(json.rows, (i, item) => {
            var $tr = $('<tr></tr>');
            $tr.append('<td>' + item.Username + '</td>');
            $tr.append('<td>' + item.Realname + '</td>');
            $tr.append('<td>' + item.DeptName + '</td>');
            if (item.Role == 0) $tr.append('<td>General Employee</td>');
            else $tr.append('<td>Administrator</td>');
            $tr.append('<td>' + item.Permission + '</td>');
            $tr.append("<td><a class='btn btn-outline-primary btn-sm' onclick='empl(" + item.ID + ")'>Edit</a></td>");
            $tr.appendTo('tbody');
        });
    });
    getDepts('#dept')
};

function loadRecords(mode) {
    var param = getParam(mode, 'records');
    $.post('/get', param, json => {
        $('tbody').empty();
        $.each(json.rows, (i, item) => {
            var $tr = $('<tr></tr>');
            $tr.append('<td>' + item.DeptName + '</td>');
            $tr.append('<td>' + item.Name + '</td>');
            $tr.append('<td>' + item.Date.split('T')[0] + '</td>');
            if (item.Type == true) $tr.append('<td>Overtime</td>');
            else $tr.append('<td>Leave</td>');
            $tr.append('<td>' + item.Duration + '</td>');
            $tr.append('<td>' + item.Describe + '</td>');
            $tr.append('<td>' + item.Created.split('T')[0] + '</td>');
            if (item.Status == 0) $tr.append("<td><a class='text-muted'>Unverified</a></td>");
            else if (item.Status == 1) $tr.append("<td><a class='text-success'>Verified</a></td>");
            else if (item.Status == 2) $tr.append("<td><a class='text-danger'>Rejected</a></td>");
            if (mode == 'admin') $tr.append("<td><a class='btn btn-outline-primary btn-sm' onclick='verify(" + item.ID + ")'>Verify</a></td>");
            else $tr.append("<td><a class='btn btn-outline-primary btn-sm' onclick='record(" + item.ID + ")'>Edit</a></td>");
            $tr.appendTo('tbody');
        });
    });
    getDepts('#dept')
};

function getDepts(element) {
    $.post('/get', 'mode=admin&query=depts', json => {
        $.each(json.rows, (i, item) => {
            $(element).append($('<option>').text(item.Name).prop('value', item.ID));
        });
    });
}

function getEmpls(element) {
    $.post('/get', 'mode=admin&query=empls', json => {
        $.each(json.rows, (i, item) => {
            $(element).append($('<option>').text(item.Realname).prop('value', item.ID));
        });
    });
}

function download(mode, query) {
    $.post('/export', getParam(mode, query));
};

function show(query) {
    var url;
    if (query == 'record') url = '/record';
    else if (query == 'stat') url = '/stat';
    else return false;
    loading();
    $.get(url).done(html => {
        loading(false);
        $('.content').html(html);
        if (query == 'record') {
            document.title = 'Employee Records - OLMS';
            $('.title').text('Employee Records');
        } else {
            document.title = 'Employee Statistics - OLMS';
            $('.title').text('Employee Statistics');
        };
    }).done(load('', query)).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login'; });
};

function showAdmin(query) {
    var url;
    if (query == 'depts') url = '/dept';
    else if (query == 'empls') url = '/empl';
    else if (query == 'records') url = '/record/dept';
    else if (query == 'stats') url = '/stat/dept';
    else return false;
    loading();
    $.get(url, html => {
        loading(false);
        $('.content').html(html);
        if (query == 'records') {
            document.title = 'Department Records - OLMS';
            $('.title').text('Department Records');
        } else if (query == 'stats') {
            document.title = 'Department Statistics - OLMS';
            $('.title').text('Department Statistics');
        }
        else if (query == 'depts') document.title = 'Departments List - OLMS';
        else if (query == 'empls') document.title = 'Employees List - OLMS';
    }).done(() => {
        if (query == 'empls') loadEmpls('admin');
        else if (query == 'depts') loadDepts('admin');
        else if (query == 'records') loadRecords('admin');
        else load('admin', query);
    }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login'; });
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
                $.each(json.rows[0], (k, v) => {
                    $('#' + k).val(v);
                })
                $('#dept').focus();
            });
        };
    }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login'; });
};

function empl(id = 0) {
    var url;
    if (id == 0) url = '/empl/add';
    else url = '/empl/edit/' + id;
    loading();
    $.get(url, html => $('.content').html(html)).done(() => {
        if (id == 0) {
            document.title = 'Add Employee - OLMS';
            $('.title').text('Add Employee');
        } else {
            document.title = 'Edit Employee - OLMS';
            $('.title').text('Edit Employee');

        };
    }).done(() => {
        getDepts('#Dept, #Permission');
        if (id != 0) {
            $.post('get', 'mode=admin&query=empls&id=' + id, json => {
                $.each(json.rows[0], (k, v) => {
                    $('#' + k).val(v);
                })
                $('#Username').focus();
            });
        }
        loading(false);
    }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login'; });
};

function record(mode, id = 0) {
    var url;
    if (id == 0) {
        if (mode == 'depts') url = '/record/dept/add';
        else url = '/record/add';
    } else {
        if (mode == 'depts') url = '/record/dept/edit/' + id;
        else url = '/record/edit/' + id;
    }
    loading();
    $.get(url, html => $('.content').html(html)).done(() => {
        if (id == 0) {
            document.title = 'Add Record - OLMS';
            $('.title').text('Add Record');
        } else {
            document.title = 'Edit Record - OLMS';
            $('.title').text('Edit Record');
            $.post('get', 'id=' + id + '&mode=' + mode, json => {
                $.each(json.record, (k, v) => {
                    $('#' + k).val(v);
                });
                $('#Date').val(json.record.Date.split('T')[0]);
                if (json.record.Type) $('#Type').val('Overtime');
                else $('#Type').val('Leave');
            }, 'json');
        };
    }).done(() => {
        if (mode == 'depts') {
            getDepts('#dept');
        }
        loading(false);
    }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login'; });
};

function verify(id) {
    loading();
    $.get('/record/verify/' + id, html => $('.content').html(html)).done(() => {
        loading(false);
        document.title = 'Verify Record - OLMS';
        $.post('get', '&mode=admin&id=' + id, json => {
            $.each(json.record, (k, v) => {
                $('#' + k).val(v);
            });
            $('#Date').val(json.record.Date.split('T')[0]);
            if (json.record.Type) $('#Type').val('Overtime');
            else $('#Type').val('Leave');
        }, 'json');
    }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login'; });
};

function setting() {
    loading();
    $.get('/auth/setting').done(html => {
        loading(false);
        $('.content').html(html);
        document.title = 'Setting - OLMS';
        $('#password').focus();
    }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login'; });
};

function doDept(id) {
    var url;
    if (id == 0) url = '/dept/add';
    else url = '/dept/edit/' + id;
    if (valid())
        $.post(url, $('input').serialize(), json => {
            $('.form').removeClass('was-validated');
            if (json.status == 0)
                BootstrapButtons.fire('Error', json.message, 'error').then(() => {
                    $('#dept').val('');
                });
            else showAdmin('depts');
        }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login'; });
};

function doEmpl(id) {
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
            else showAdmin('empls');
        }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login'; });
};

function doRecord(id) {
    var url;
    if (id == 0) url = '/record/add';
    else url = '/record/edit/' + id;
    if (valid())
        $.post(url, $('input, select, textarea').serialize(), json => {
            $('.form').removeClass('was-validated');
            if (json.status == 0)
                BootstrapButtons.fire('Error', json.message, 'error').then(() => {
                    $('#duration').val('');
                });
            else show('record');
        }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login'; });
};

function doVerify(id, status) {
    $.post('/verify/' + id, 'status=' + status, () => show('record')).fail(xhr => { if (xhr.status == 401) window.location = '/auth/login'; });
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
                if (json.status == 1) show(mode);
            }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login'; });
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
                if (json.error == 1)
                    $('#password').val('');
                else if (json.error == 2) {
                    $('#password1').val('');
                    $('#password2').val('');
                };
            });
        }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login'; });
};

function valid() {
    var result = true
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

