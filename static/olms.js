BootstrapButtons = Swal.mixin({
    customClass: {
        confirmButton: 'swal btn btn-primary'
    },
    buttonsStyling: false
});

function load(mode, obj) {
    $.post('/get', $.param(obj) + '&query=' + mode, json => {
        $('tbody').empty();
        $.each(json.rows, (i, item) => {
            var $tr = $('<tr><tr>');
            $.each(item, (k, v) => {
                $tr.append('<td>' + v + '</td>');
            })
            $tr.appendTo('tbody');
        });
    });
};

function show(mode) {
    var url;
    if (mode == 'dept') url = '/dept';
    else if (mode == 'empl') url = '/empl';
    else if (mode == 'record') url = '/';
    else return false;
    window.location = url;
};

function dept(id = 0) {
    if (id == 0) {
        document.title = 'Add Department - OLMS';
    } else {
        document.title = 'Edit Department - OLMS';
        loading();
        $.post('get', 'mode=admin&query=dept&id=' + id, json => {
            loading(false);
            $.each(json, (k, v) => {
                $("#" + k).val(v);
            })
            $('#dept').focus();
        });
    };
};

function empl(id = 0) {
    if (id == 0) {
        document.title = 'Add Employee - OLMS';
    } else {
        document.title = 'Edit Employee - OLMS';
        loading();
        $.post('get', 'mode=admin&query=empl&id=' + id, json => {
            loading(false);
            $.each(json, (k, v) => {
                $("#" + k).val(v);
            })
            $('#username').focus();
        });
    };
};

function record(id = 0) {
    if (id == 0) {
        document.title = 'Add Record - OLMS';
    } else {
        document.title = 'Edit Record - OLMS';
        loading();
        $.post('get', 'id=' + id, json => {
            loading(false);
            $.each(json, (k, v) => {
                $("#" + k).val(v);
            })
        }, 'json');
    };
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
    if (id == 0) url = '/add';
    else url = '/edit/' + id;
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
    else if (mode == 'record') url = '/delete/' + id;
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

