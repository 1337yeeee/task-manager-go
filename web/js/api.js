function buildApiUrl(url) {
    if (url.startsWith("/")) {
        return window.APP_CONFIG.API_URL + url;
    }

    return window.APP_CONFIG.API_URL + "/" + url;
}

function apiRequest(method, url, data, options) {
    const requestOptions = options || {};
    const deferred = $.Deferred();
    const headers = Object.assign({}, requestOptions.headers);
    const accessToken = getAccessToken();

    if (accessToken) {
        headers.Authorization = "Bearer " + accessToken;
    }

    $.ajax({
        url: buildApiUrl(url),
        method: method,
        contentType: "application/json",
        headers: headers,
        data: data ? JSON.stringify(data) : null
    }).done(function(response, textStatus, jqXHR) {
        deferred.resolve(response, textStatus, jqXHR);
    }).fail(function(jqXHR) {
        const shouldRetry = jqXHR.status === 401 && requestOptions.retryOnUnauthorized !== false;

        if (!shouldRetry) {
            deferred.reject.apply(deferred, arguments);
            return;
        }

        refreshAccessToken()
            .done(function() {
                apiRequest(method, url, data, Object.assign({}, requestOptions, {
                    retryOnUnauthorized: false
                }))
                    .done(function(response, textStatus, retryJqXHR) {
                        deferred.resolve(response, textStatus, retryJqXHR);
                    })
                    .fail(function() {
                        deferred.reject.apply(deferred, arguments);
                    });
            })
            .fail(function() {
                deferred.reject.apply(deferred, arguments);
            });
    });

    return deferred.promise();
}
