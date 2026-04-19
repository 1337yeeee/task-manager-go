const state = {
    projects: [],
    tasks: [],
    currentProjectId: null
};

function getProjectsFromResponse(response) {
    if (response && Array.isArray(response.data)) {
        return response.data;
    }

    if (Array.isArray(response)) {
        return response;
    }

    return [];
}

function getProjectFromResponse(response) {
    if (response && response.data) {
        return response.data;
    }

    return response || null;
}

function getTasksFromResponse(response) {
    if (response && Array.isArray(response.data)) {
        return response.data;
    }

    if (Array.isArray(response)) {
        return response;
    }

    return [];
}

function getProjectIdFromUrl() {
    const normalizedPath = window.location.pathname.replace(/^\/+|\/+$/g, "");
    if (!normalizedPath) {
        return null;
    }

    const pathParts = normalizedPath.split("/");
    return pathParts.length === 1 ? decodeURIComponent(pathParts[0]) : null;
}

function updateProjectUrl(projectId, shouldReplace) {
    const nextPath = "/" + encodeURIComponent(projectId);
    if (window.location.pathname === nextPath) {
        return;
    }

    if (shouldReplace) {
        window.history.replaceState({}, "", nextPath);
        return;
    }

    window.history.pushState({}, "", nextPath);
}

function updateTaskStatus(task) {
    apiRequest("PATCH", `api/tasks/${task.id}/status`, {"status": task.status})
        .done(function() {
            console.log("status updated successfully")
        })
        .fail(function(err) {
            console.error("Failed to load projects", err);
        });
}

function renderProjectList() {
    const list = $(".project-list");
    list.empty();

    state.projects.forEach(p => {
        const active = p.id === state.currentProjectId ? "active" : "";
        const href = "/" + encodeURIComponent(p.id);

        list.append(`
            <li class="project ${active}" data-id="${p.id}">
                <a class="project-link" href="${href}">${p.name}</a>
            </li>
        `);
    });
}

function renderProjectHeader() {
    const activeProject = state.projects.find(p => p.id === state.currentProjectId);

    if (!activeProject) {
        $(".main-header h1").text("No project selected");
        $(".desc-header p").text("");
        return;
    }

    $(".main-header h1").text(activeProject.name || "Untitled project");
    $(".desc-header p").text(activeProject.description || "");
}

function setActiveProject(projectId, options) {
    const settings = options || {};

    state.currentProjectId = projectId;
    renderProjectList();
    renderProjectHeader();
    renderTasks(projectId);

    if (settings.syncUrl) {
        updateProjectUrl(projectId, settings.replaceUrl === true);
    }
}

/* =========================
   RENDER
========================= */

function renderProjects() {
    apiRequest("GET", "api/projects")
        .done(function(response) {
            state.projects = getProjectsFromResponse(response);

            if (!state.projects.length) {
                state.currentProjectId = null;
                renderProjectList();
                renderProjectHeader();
                renderTasks(null);
                return;
            }

            const projectIdFromUrl = getProjectIdFromUrl();
            const projectFromUrl = state.projects.find(p => p.id === projectIdFromUrl);
            if (projectFromUrl) {
                setActiveProject(projectFromUrl.id, { syncUrl: false });
                return;
            }

            setActiveProject(state.projects[0].id, { syncUrl: true, replaceUrl: true });
        })
        .fail(function(err) {
            console.error("Failed to load projects", err);
        });
}

function renderTasks(projectId) {
    // очищаем все колонки
    $(".column .task-list").empty();
    if (!projectId) return;

    apiRequest("GET", `api/projects/${projectId}/tasks`)
        .done(function(response) {
            const tasks = getTasksFromResponse(response);
            state.tasks = tasks

            tasks.forEach(t => {
                const el = $(`
                    <div class="task-card" data-id="${t.id}">
                        <div class="task-header">
                            <span class="task-title">${t.name}</span>
                            <button class="menu-btn">⋮</button>
                        </div>

                        <div class="task-content">${t.content}</div>

                        <div class="task-footer">
                            <span>Exec: ${t.executive_id}</span>
                            <span>Audit: ${t.auditor_id}</span>
                        </div>

                        <div class="task-edit hidden">
                            <input type="text" value="${t.name}" />
                            <textarea>${t.content}</textarea>
                            <button class="btn save">Save</button>
                        </div>

                        <div class="context-menu hidden">
                            <label>Executive</label>
                            <select>
                                <option>User1</option>
                                <option>User2</option>
                            </select>

                            <label>Auditor</label>
                            <select>
                                <option>User1</option>
                                <option>User2</option>
                            </select>

                            <button class="btn save">Save</button>
                        </div>
                    </div>
                `);

                // важно: привести статус к id колонки
                // например: created / in_progress / done
                const statusMap = {
                    "created": "created",
                    "in_progress": "in_progress",
                    "done": "done"
                };

                const columnId = statusMap[t.status];
                if (!columnId) return;

                $("#" + columnId + " .task-list").append(el);
            });
        })
        .fail(function(err) {
            console.error("Failed to load tasks", err);
        });
}

/* =========================
   PROJECTS
========================= */

$(".project-list").on("click", ".project-link", function (event) {
    event.preventDefault();
    const projectId = String($(this).closest(".project").data("id"));
    setActiveProject(projectId, { syncUrl: true, replaceUrl: false });
});

/* =========================
   CREATE PROJECT
========================= */

function openCreateProjectModal() {
    $("#projectNameInput").val("");
    $("#projectDescInput").val("");
    $("#createProjectError").text("");
    $("#createProjectModal").removeClass("hidden");
    $("#projectNameInput").trigger("focus");
}

function closeCreateProjectModal() {
    $("#createProjectModal").addClass("hidden");
}

$("#newProjectBtn").click(function () {
    openCreateProjectModal();
});

$("#saveProjectBtn").click(function () {
    const name = $("#projectNameInput").val().trim();
    const desc = $("#projectDescInput").val().trim();

    if (!name) {
        $("#createProjectError").text("Project name is required");
        return;
    }

    apiRequest("POST", "api/projects", {
        name: name,
        description: desc
    }).done(function(response) {
        const createdProject = getProjectFromResponse(response);

        if (!createdProject || !createdProject.id) {
            renderProjects();
            closeCreateProjectModal();
            return;
        }

        state.projects.unshift(createdProject);
        closeCreateProjectModal();
        setActiveProject(createdProject.id, { syncUrl: true, replaceUrl: false });
    }).fail(function(err) {
        console.error("Failed to create project", err);
        $("#createProjectError").text(err.message)
        renderProjects();
    });
});

$("#projectNameInput, #projectDescInput").on("keydown", function (event) {
    if (event.key === "Enter" && !event.shiftKey) {
        event.preventDefault();
        $("#saveProjectBtn").trigger("click");
    }
});

$("#createProjectModal").on("click", function (event) {
    if (event.target === this) {
        closeCreateProjectModal();
    }
});

/* =========================
   CREATE TASK
========================= */

function openCreateTaskModal() {
    $("#createTaskError").text("")
    $("#taskNameInput").val("");
    $("#taskContentInput").val("");
    $("#createProjectError").text("");
    $("#createTaskModal").removeClass("hidden");
    $("#taskNameInput").trigger("focus");
}

function closeCreateTaskModal() {
    $("#createTaskModal").addClass("hidden");
}

$("#newTaskBtn").click(function () {
    openCreateTaskModal();
});

$("#saveTaskBtn").click(function () {
    let createTaskError = $("#createTaskError")
    createTaskError.text("")

    const name = $("#taskNameInput").val().trim();
    const content = $("#taskContentInput").val().trim();

    if (!name) {
        createTaskError.text("Task name is required");
        return;
    }

    let currentProjectId = state.currentProjectId

    apiRequest("POST", "api/projects/" + currentProjectId + "/tasks", {
        name: name,
        content: content
    }).done(function() {
        renderTasks(currentProjectId);
        closeCreateProjectModal();
    }).fail(function(err) {
        console.error("Failed to create project", err);
        createTaskError.text(err.message)
    });
});

$("#taskNameInput, #taskContentInput").on("keydown", function (event) {
    if (event.key === "Enter" && !event.shiftKey) {
        event.preventDefault();
        $("#saveTaskBtn").trigger("click");
    }
});

/* =========================
   DRAG & DROP
========================= */

$(".task-list").sortable({
    connectWith: ".task-list",
    stop: function (event, ui) {
        const taskId = ui.item.data("id");
        const newStatus = ui.item.parent().data("status-id");
        console.log(ui.item.parent())

        const task = state.tasks.find(t => t.id === taskId);
        task.status = newStatus;
        updateTaskStatus(task)

        console.log("Updated task status (stub)", task);
    }
});

/* =========================
   INLINE EDIT
========================= */

$(document).on("dblclick", ".task-title, .task-content", function () {
    const card = $(this).closest(".task-card");

    card.find(".task-edit").removeClass("hidden");
    card.find(".task-title, .task-content").hide();
});

$(document).on("click", ".task-edit .save", function () {
    const card = $(this).closest(".task-card");
    const id = card.data("id");

    const name = card.find("input").val();
    const content = card.find("textarea").val();

    const task = state.tasks.find(t => t.id === id);
    task.name = name;
    task.content = content;

    renderTasks();
});

/* =========================
   CONTEXT MENU
========================= */

$(document).on("click", ".menu-btn", function (e) {
    e.stopPropagation();
    $(".context-menu").addClass("hidden");

    $(this).closest(".task-card").find(".context-menu").toggleClass("hidden");
});

$(document).on("click", function () {
    $(".context-menu").addClass("hidden");
});

$(document).on("click", ".context-menu .save", function () {
    const card = $(this).closest(".task-card");
    const id = card.data("id");

    const selects = card.find("select");
    const exec = $(selects[0]).val();
    const audit = $(selects[1]).val();

    const task = state.tasks.find(t => t.id === id);
    task.executive = exec;
    task.auditor = audit;

    renderTasks();
});

/* =========================
   PROFILE / LOGOUT
========================= */

$("#profileBtn").click(function () {
    alert("Go to profile page (stub)");
});

$("#logoutBtn").click(function () {
    logout();
});

/* =========================
   CLOSE MODALS
========================= */

$(document).on("keydown", function (event) {
    if (event.key === "Escape") {
        closeCreateProjectModal();
        closeCreateTaskModal();
    }
});

$(window).on("popstate", function () {
    if (!state.projects.length) {
        return;
    }

    const projectIdFromUrl = getProjectIdFromUrl();
    const projectFromUrl = state.projects.find(p => p.id === projectIdFromUrl);

    if (projectFromUrl) {
        setActiveProject(projectFromUrl.id, { syncUrl: false });
        return;
    }

    setActiveProject(state.projects[0].id, { syncUrl: true, replaceUrl: true });
});

/* =========================
   INIT
========================= */

renderProjects();
