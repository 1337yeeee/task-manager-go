const usersState = {
    users: []
};

function ensureAdminAccess() {
    if (getUserRole() !== "admin") {
        window.location.replace("/");
    }
}

function getUsersFromResponse(response) {
    if (response && Array.isArray(response.users)) {
        return response.users;
    }

    if (response && Array.isArray(response.data)) {
        return response.data;
    }

    if (Array.isArray(response)) {
        return response;
    }

    return [];
}

function normalizeUser(user) {
    return {
        id: user.id || user.ID,
        name: user.name || user.Name,
        email: user.email || user.Email,
        role: user.role || user.Role,
        is_active: user.is_active !== undefined ? user.is_active : user.IsActive
    };
}

function renderUsers() {
    const list = $("#usersList");
    list.empty();

    usersState.users.forEach(function(user) {
        const statusLabel = user.is_active ? "Active" : "Suspended";
        const statusClass = user.is_active ? "active" : "suspended";

        list.append(`
            <article class="user-card" data-id="${user.id}">
                <div class="user-card-header">
                    <h3>${user.email}</h3>
                    <span class="user-status ${statusClass}">${statusLabel}</span>
                </div>
                <p class="user-meta">ID: ${user.id}</p>
                <label>Role</label>
                <select class="user-role-select">
                    <option value="viewer" ${user.role === "viewer" ? "selected" : ""}>viewer</option>
                    <option value="editor" ${user.role === "editor" ? "selected" : ""}>editor</option>
                    <option value="admin" ${user.role === "admin" ? "selected" : ""}>admin</option>
                </select>
                <label class="user-active-toggle">
                    <input type="checkbox" class="user-active-checkbox" ${user.is_active ? "checked" : ""}>
                    <span>is_active</span>
                </label>
                <button class="btn primary user-save-btn">Save</button>
                <div class="user-card-error"></div>
            </article>
        `);
    });
}

function loadUsers() {
    return apiRequest("GET", "api/users")
        .done(function(response) {
            const users = getUsersFromResponse(response).map(normalizeUser);
            usersState.users = users;
            renderUsers();
        })
        .fail(function(jqXHR) {
            const message = jqXHR.responseJSON && jqXHR.responseJSON.error
                ? jqXHR.responseJSON.error
                : "Failed to load users";
            $("#usersError").text(message);
        });
}

function openAddUserModal() {
    $("#newUserEmail").val("");
    $("#newUserPassword").val("");
    $("#newUserRole").val("viewer");
    $("#addUserError").text("");
    $("#addUserModal").removeClass("hidden");
}

function closeAddUserModal() {
    $("#addUserModal").addClass("hidden");
}

function buildNameFromEmail(email) {
    const cleanEmail = email.trim();
    const atIndex = cleanEmail.indexOf("@");
    if (atIndex > 0) {
        return cleanEmail.slice(0, atIndex);
    }
    return cleanEmail;
}

function createUser() {
    const email = $("#newUserEmail").val().trim();
    const password = $("#newUserPassword").val();
    const role = $("#newUserRole").val();
    const modalError = $("#addUserError");

    modalError.text("");

    if (!email) {
        modalError.text("Email is required");
        return;
    }
    if (password.length < 8) {
        modalError.text("Password must be at least 8 characters long");
        return;
    }

    apiRequest("POST", "api/users", {
        name: buildNameFromEmail(email),
        email: email,
        password: password,
        role: role
    }).done(function() {
        closeAddUserModal();
        loadUsers();
    }).fail(function(jqXHR) {
        const message = jqXHR.responseJSON && jqXHR.responseJSON.error
            ? jqXHR.responseJSON.error
            : "Failed to create user";
        modalError.text(message);
    });
}

function updateUser(card) {
    const userId = card.data("id");
    const role = card.find(".user-role-select").val();
    const isActive = card.find(".user-active-checkbox").is(":checked");
    const errorBox = card.find(".user-card-error");

    errorBox.text("");

    apiRequest("PUT", "api/users/" + userId, {
        role: role,
        is_active: isActive
    }).done(function(response) {
        const updatedUser = normalizeUser(response.user || response.data || {});
        usersState.users = usersState.users.map(function(existingUser) {
            if (existingUser.id === userId) {
                return Object.assign({}, existingUser, updatedUser);
            }
            return existingUser;
        });
        renderUsers();
    }).fail(function(jqXHR) {
        const message = jqXHR.responseJSON && jqXHR.responseJSON.error
            ? jqXHR.responseJSON.error
            : "Failed to update user";
        errorBox.text(message);
    });
}

$("#addUserBtn").click(function() {
    openAddUserModal();
});

$("#saveUserBtn").click(function() {
    createUser();
});

$("#usersList").on("click", ".user-save-btn", function() {
    const card = $(this).closest(".user-card");
    updateUser(card);
});

$("#projectsBtn").click(function() {
    window.location.href = "/";
});

$("#profileBtn").click(function() {
    alert("Go to profile page (stub)");
});

$("#usersBtn").click(function() {
    window.location.href = "/users.html";
});

$("#logoutBtn").click(function() {
    logout();
});

$("#addUserModal").on("click", function(event) {
    if (event.target === this) {
        closeAddUserModal();
    }
});

$(document).on("keydown", function(event) {
    if (event.key === "Escape") {
        closeAddUserModal();
    }
});

ensureAdminAccess();
loadUsers();
