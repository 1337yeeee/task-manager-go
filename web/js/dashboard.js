/* =========================
   RENDER
========================= */

function renderProjects() {
    const list = $("#projectList");
    list.empty();

    apiRequest("GET", "api/projects")
        .done(function(projects) {
            state.projects = projects;

            projects.forEach(p => {
                const active = p.ID === state.currentProjectId ? "active" : "";

                list.append(`
                    <li class="project ${active}" data-id="${p.ID}">
                        ${p.Name}
                    </li>
                `);
            });
        })
        .fail(function(err) {
            console.error("Failed to load projects", err);
        });
}

function renderTasks(projectId) {
    if (!projectId) return;

    // очищаем все колонки
    $(".column .task-list").empty();

    apiRequest("GET", `api/projects/${projectId}/tasks`)
        .done(function(tasks) {

            tasks.forEach(t => {
                const el = $(`
                    <div class="task-card" data-id="${t.ID}">
                        <div class="task-header">
                            <span class="task-title">${t.Name}</span>
                            <button class="menu-btn">⋮</button>
                        </div>

                        <div class="task-content">${t.Content}</div>

                        <div class="task-footer">
                            <span>Exec: ${t.ExecutiveID}</span>
                            <span>Audit: ${t.AuditorID}</span>
                        </div>

                        <div class="task-edit hidden">
                            <input type="text" value="${t.Name}" />
                            <textarea>${t.Content}</textarea>
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
                    "in_progress": "in-progress",
                    "done": "done"
                };

                const columnId = statusMap[t.Status];
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

$("#projectList").on("click", ".project", function () {
    state.currentProjectId = $(this).data("id");
    renderProjects();
    renderTasks();
});

/* =========================
   CREATE PROJECT
========================= */

$(".sidebar .primary").click(function () {
    const name = prompt("Project name");
    if (!name) return;

    state.projects.push({
        id: Date.now(),
        name
    });

    renderProjects();
});

/* =========================
   CREATE TASK
========================= */

$("#createTask").click(function () {
    const name = prompt("Task name");
    const content = prompt("Content");

    state.tasks.push({
        id: Date.now(),
        name,
        content,
        status: "created",
        projectId: state.currentProjectId,
        executive: "User1",
        auditor: "User2"
    });

    renderTasks();
});

/* =========================
   DRAG & DROP
========================= */

$(".task-list").sortable({
    connectWith: ".task-list",
    stop: function (event, ui) {
        const taskId = ui.item.data("id");
        const newStatus = ui.item.parent().attr("id");

        const task = state.tasks.find(t => t.id === taskId);
        task.status = newStatus;

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
    alert("Logout (stub)");
});

/* =========================
   INIT
========================= */

renderProjects();
renderTasks(state.currentProjectId);
