BootstrapButtons = Swal.mixin({
    customClass: {
        confirmButton: 'swal btn btn-primary'
    },
    buttonsStyling: false
});

function load(mode, query) {
    var data = {};
    $('select').serializeArray().forEach(i => data[i.name] = i.value);
    $.post('/get', $.param(data) + '&mode=' + mode + '&query=' + query, json => {
        $('tbody').empty();
        $.each(json.rows, (i, item) => {
            var $tr = $('<tr><tr>');
            $.each(item, (k, v) => {
                $tr.append('<td>' + v + '</td>');
            });
            $tr.appendTo('tbody');
        });
    });
};

function download(mode, query) {
    var data = {}
    $('select').serializeArray().forEach(i => data[i.name] = i.value);
    $.post('/export', $.param(data) + '&mode=' + mode + '&query=' + query);
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

function showDept(query) {
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
        };
    }).done(load('admin', query)).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login'; });
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
            $.post('get', 'mode=admin&query=dept&id=' + id, json => {
                loading(false);
                $.each(json, (k, v) => {
                    $("#" + k).val(v);
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
            loading(false);
        } else {
            document.title = 'Edit Employee - OLMS';
            $('.title').text('Edit Employee');
            $.post('get', 'mode=admin&query=empl&id=' + id, json => {
                loading(false);
                $.each(json, (k, v) => {
                    $("#" + k).val(v);
                })
                $('#username').focus();
            });
        };
    }).fail(jqXHR => { if (jqXHR.status == 401) window.location = '/auth/login'; });
};

function record(mode, id = 0) {
    var url;
    if (id == 0) {
        if (mode == 'dept') url = '/record/dept/add';
        else url = '/record/add';
    } else {
        if (mode == 'dept') url = '/record/dept/edit/' + id;
        else url = '/record/edit/' + id;
    }
    loading();
    $.get(url, html => $('.content').html(html)).done(() => {
        if (id == 0) {
            document.title = 'Add Record - OLMS';
            $('.title').text('Add Record');
            loading(false);
        } else {
            document.title = 'Edit Record - OLMS';
            $('.title').text('Edit Record');
            $.post('get', 'id=' + id + '&mode=' + mode, json => {
                loading(false);
                $.each(json, (k, v) => {
                    $("#" + k).val(v);
                })
            }, 'json');
        };
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
            else show('dept');
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
            else show('empl');
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

