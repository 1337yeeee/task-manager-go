window.APP_CONFIG = window.APP_CONFIG || {
    API_URL: ""
};

const ACCESS_TOKEN_KEY = "access_token";
const USER_ROLE_KEY = "user_role";

let refreshRequest = null;

function getAccessToken() {
    return localStorage.getItem(ACCESS_TOKEN_KEY);
}

function getUserRole() {
    const storedRole = localStorage.getItem(USER_ROLE_KEY);
    if (storedRole) {
        return storedRole;
    }

    const tokenRole = getRoleFromAccessToken();
    if (tokenRole) {
        localStorage.setItem(USER_ROLE_KEY, tokenRole);
    }

    return tokenRole;
}

function getRoleFromAccessToken() {
    const token = getAccessToken();
    if (!token) {
        return null;
    }

    const parts = token.split(".");
    if (parts.length < 2) {
        return null;
    }

    try {
        const payload = JSON.parse(atob(parts[1].replace(/-/g, "+").replace(/_/g, "/")));
        return payload.Role || payload.role || null;
    } catch (err) {
        return null;
    }
}

function setTokens(accessToken, role) {
    if (accessToken) {
        localStorage.setItem(ACCESS_TOKEN_KEY, accessToken);
    }

    if (role) {
        localStorage.setItem(USER_ROLE_KEY, role);
    }
}

function clearTokens() {
    localStorage.removeItem(ACCESS_TOKEN_KEY);
    localStorage.removeItem(USER_ROLE_KEY);
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
        setTokens(response.access_token, response.role);
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
        setTokens(response.access_token, response.role);
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
