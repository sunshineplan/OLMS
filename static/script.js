function cleanObj(obj) {
    $.each(obj, (k, v) => {
        if (v === undefined || v === '')
            delete obj[k]
    });
    return obj;
};

function postJSON(url, data, success) {
    if ($('.g-recaptcha-response').length)
        return grecaptcha.execute(sitekey, { action: url.replace('/', '') })
            .then(token => {
                data['g-recaptcha-response'] = token;
                return $.ajax({
                    url: url,
                    type: 'POST',
                    contentType: 'application/json',
                    data: JSON.stringify(data),
                    success: success,
                });
            });
    return $.ajax({
        url: url,
        type: 'POST',
        contentType: 'application/json',
        data: JSON.stringify(data),
        success: success
    });
};
