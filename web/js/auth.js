window.APP_CONFIG = window.APP_CONFIG || {
    API_URL: ""
};

const ACCESS_TOKEN_KEY = "access_token";

let refreshRequest = null;

function getAccessToken() {
    return localStorage.getItem(ACCESS_TOKEN_KEY);
}

function setTokens(accessToken) {
    if (accessToken) {
        localStorage.setItem(ACCESS_TOKEN_KEY, accessToken);
    }
}

function clearTokens() {
    localStorage.removeItem(ACCESS_TOKEN_KEY);
}

function redirectToLogin() {
    clearTokens();
    window.location.replace("/login");
}

function isAuthPage() {
    return window.location.pathname === "/login" || window.location.pathname.endsWith("/login.html");
}

function requireAuth() {
    if (!getAccessToken() && !isAuthPage()) {
        redirectToLogin();
        return false;
    }

    return true;
}

function redirectAuthenticatedUser() {
    if (isAuthPage() && getAccessToken()) {
        window.location.replace("/");
    }
}

function login(email, password) {
    return $.ajax({
        url: window.APP_CONFIG.API_URL + "/api/auth/login",
        method: "POST",
        contentType: "application/json",
        data: JSON.stringify({
            email: email,
            password: password
        })
    }).done(function(response) {
        setTokens(response.access_token);
    });
}

function refreshAccessToken() {
    if (refreshRequest) {
        return refreshRequest;
    }

    refreshRequest = $.ajax({
        url: window.APP_CONFIG.API_URL + "/api/refresh",
        method: "POST",
        contentType: "application/json",
    }).done(function(response) {
        setTokens(response.access_token);
    }).fail(function() {
        redirectToLogin();
    }).always(function() {
        refreshRequest = null;
    });

    return refreshRequest;
}

function logout() {
    clearTokens();
    window.location.replace("/login");
}

if (isAuthPage()) {
    redirectAuthenticatedUser();
} else {
    requireAuth();
}
