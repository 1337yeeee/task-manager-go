const API_URL = "http://localhost:8080";

function getToken() {
    return localStorage.getItem("accessToken");
}

function apiRequest(method, url, data) {
    return $.ajax({
        url: API_URL + url,
        method: method,
        contentType: "application/json",
        headers: {
            Authorization: "Bearer " + getToken()
        },
        data: data ? JSON.stringify(data) : null
    });
}
