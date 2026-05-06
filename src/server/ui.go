package server

import "net/http"

func serveMockUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(mockUIHTML))
}

const mockUIHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Gigamock Control</title>
  <style>
    :root {
      color-scheme: light;
      --bg: #f7f8fb;
      --panel: #ffffff;
      --text: #172033;
      --muted: #667085;
      --line: #d9deea;
      --accent: #1463ff;
      --accent-soft: #e8f0ff;
      --danger: #b42318;
      --ok: #027a48;
      --shadow: 0 10px 28px rgba(23, 32, 51, .08);
    }

    * { box-sizing: border-box; }
    body {
      margin: 0;
      background: var(--bg);
      color: var(--text);
      font: 14px/1.45 system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
    }
    header {
      position: sticky;
      top: 0;
      z-index: 2;
      border-bottom: 1px solid var(--line);
      background: rgba(247, 248, 251, .94);
      backdrop-filter: blur(10px);
    }
    .bar {
      max-width: 1180px;
      margin: 0 auto;
      padding: 18px 20px;
      display: grid;
      grid-template-columns: 1fr auto;
      gap: 16px;
      align-items: center;
    }
    h1 {
      margin: 0;
      font-size: 22px;
      line-height: 1.2;
      font-weight: 720;
      letter-spacing: 0;
    }
    .toolbar {
      display: grid;
      grid-template-columns: minmax(220px, 1fr) 160px 220px auto;
      gap: 10px;
      align-items: center;
    }
    input, select, button {
      height: 36px;
      border: 1px solid var(--line);
      border-radius: 6px;
      background: #fff;
      color: var(--text);
      font: inherit;
    }
    input {
      width: 100%;
      padding: 0 12px;
    }
    .toolbar select {
      width: 100%;
      padding: 0 10px;
    }
    button {
      padding: 0 12px;
      cursor: pointer;
      font-weight: 650;
    }
    button.primary {
      border-color: var(--accent);
      background: var(--accent);
      color: #fff;
    }
    main {
      max-width: 1180px;
      margin: 0 auto;
      padding: 22px 20px 40px;
    }
    .summary {
      color: var(--muted);
      margin-bottom: 14px;
    }
    .grid {
      display: grid;
      gap: 12px;
    }
    .endpoint {
      background: var(--panel);
      border: 1px solid var(--line);
      border-radius: 8px;
      box-shadow: var(--shadow);
      padding: 16px;
      display: grid;
      grid-template-columns: minmax(0, 1fr) auto;
      gap: 16px;
      align-items: start;
    }
    .meta {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
      margin-bottom: 8px;
      align-items: center;
    }
    .badge {
      display: inline-flex;
      align-items: center;
      height: 24px;
      padding: 0 8px;
      border-radius: 999px;
      background: var(--accent-soft);
      color: #0b4dcc;
      font-size: 12px;
      font-weight: 700;
      text-transform: uppercase;
    }
    .service {
      background: #fff7ed;
      color: #9a3412;
    }
    .method {
      background: #ecfdf3;
      color: var(--ok);
    }
    .path {
      margin: 0 0 4px;
      font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
      font-size: 15px;
      overflow-wrap: anywhere;
    }
    .desc {
      margin: 0;
      color: var(--muted);
    }
    .file {
      margin-top: 10px;
      color: var(--muted);
      font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
      font-size: 12px;
      overflow-wrap: anywhere;
    }
    .provenance {
      margin-top: 12px;
      display: grid;
      gap: 4px;
      color: var(--muted);
      font-size: 12px;
    }
    .provenance span {
      color: var(--text);
      font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    }
    .control {
      min-width: 260px;
      display: grid;
      gap: 8px;
    }
    .control label {
      color: var(--muted);
      font-size: 12px;
      font-weight: 700;
      text-transform: uppercase;
    }
    .control select {
      width: 100%;
      padding: 0 10px;
    }
    .status {
      min-height: 20px;
      color: var(--muted);
      font-size: 13px;
    }
    .status.error { color: var(--danger); }
    .status.ok { color: var(--ok); }
    .empty {
      border: 1px dashed var(--line);
      border-radius: 8px;
      padding: 28px;
      color: var(--muted);
      text-align: center;
      background: #fff;
    }
    @media (max-width: 760px) {
      .bar, .endpoint { grid-template-columns: 1fr; }
      .toolbar { grid-template-columns: 1fr; align-items: stretch; }
      input { width: 100%; }
      .control { min-width: 0; }
    }
  </style>
</head>
<body>
  <header>
    <div class="bar">
      <h1>Gigamock Control</h1>
      <div class="toolbar">
        <input id="search" type="search" placeholder="Search path, service, file, description">
        <select id="typeFilter" aria-label="Filter by type"></select>
        <select id="serviceFilter" aria-label="Filter by service"></select>
        <button id="refresh" class="primary" type="button">Refresh</button>
      </div>
    </div>
  </header>
  <main>
    <div id="summary" class="summary">Loading scenarios...</div>
    <section id="endpoints" class="grid" aria-live="polite"></section>
  </main>
  <script>
    const endpointsEl = document.querySelector("#endpoints");
    const summaryEl = document.querySelector("#summary");
    const searchEl = document.querySelector("#search");
    const typeFilterEl = document.querySelector("#typeFilter");
    const serviceFilterEl = document.querySelector("#serviceFilter");
    const refreshEl = document.querySelector("#refresh");
    let endpoints = [];

    async function load() {
      summaryEl.textContent = "Loading scenarios...";
      endpointsEl.innerHTML = "";
      try {
        const response = await fetch("/internal/v1/scenarios");
        if (!response.ok) throw new Error(await response.text());
        const data = await response.json();
        endpoints = data.endpoints || [];
        populateFilters();
        render();
      } catch (error) {
        summaryEl.textContent = "Failed to load scenarios: " + error.message;
        endpointsEl.innerHTML = "";
      }
    }

    function render() {
      const query = searchEl.value.trim().toLowerCase();
      const selectedType = typeFilterEl.value;
      const selectedService = serviceFilterEl.value;
      const visible = endpoints.filter((endpoint) => {
        const haystack = [
          endpoint.path,
          endpoint.method,
          endpoint.type,
          endpoint.service,
          endpoint.fileName,
          endpoint.directory,
          endpoint.name,
          endpoint.description,
          endpoint.filePath,
        ].join(" ").toLowerCase();
        return haystack.includes(query) &&
          (!selectedType || endpoint.type === selectedType) &&
          (!selectedService || endpoint.service === selectedService);
      });

      summaryEl.textContent = visible.length + " of " + endpoints.length + " endpoints" + summarySuffix(selectedType, selectedService);
      if (visible.length === 0) {
        endpointsEl.innerHTML = '<div class="empty">No mock endpoints found</div>';
        return;
      }

      endpointsEl.innerHTML = visible.map(endpointTemplate).join("");
      for (const card of endpointsEl.querySelectorAll("[data-key]")) {
        const endpoint = endpoints.find((candidate) => candidate.key === card.dataset.key);
        const select = card.querySelector("select");
        const button = card.querySelector("button");
        button.addEventListener("click", () => updateScenario(endpoint, Number(select.value), card));
      }
    }

    function endpointTemplate(endpoint) {
      const options = endpoint.scenarios.map((scenario) => {
        const selected = scenario.index === endpoint.currentScenario ? " selected" : "";
        return '<option value="' + scenario.index + '"' + selected + '>' +
          escapeHtml(String(scenario.index) + " - " + scenario.name) +
          '</option>';
      }).join("");

      return '<article class="endpoint" data-key="' + escapeHtml(endpoint.key) + '">' +
        '<div>' +
          '<div class="meta">' +
            '<span class="badge">' + escapeHtml(endpoint.type) + '</span>' +
            '<span class="badge method">' + escapeHtml(endpoint.method) + '</span>' +
            '<span class="badge service">' + escapeHtml(endpoint.service || "unknown") + '</span>' +
          '</div>' +
          '<p class="path">' + escapeHtml(endpoint.path) + '</p>' +
          '<p class="desc">' + escapeHtml(endpoint.description || endpoint.name || "No description") + '</p>' +
          '<div class="provenance">' +
            '<div>Service: <span>' + escapeHtml(endpoint.service || "unknown") + '</span></div>' +
            '<div>File: <span>' + escapeHtml(endpoint.fileName || "unknown") + '</span></div>' +
            '<div>Directory: <span>' + escapeHtml(endpoint.directory || "unknown") + '</span></div>' +
          '</div>' +
        '</div>' +
        '<div class="control">' +
          '<label>Active scenario</label>' +
          '<select>' + options + '</select>' +
          '<button type="button">Apply</button>' +
          '<div class="status"></div>' +
        '</div>' +
      '</article>';
    }

    function populateFilters() {
      const currentType = typeFilterEl.value;
      const currentService = serviceFilterEl.value;
      const types = uniqueSorted(endpoints.map((endpoint) => endpoint.type));
      const services = uniqueSorted(endpoints.map((endpoint) => endpoint.service));
      typeFilterEl.innerHTML = '<option value="">All types</option>' +
        types.map((type) => '<option value="' + escapeHtml(type) + '">' + escapeHtml(type) + '</option>').join("");
      serviceFilterEl.innerHTML = '<option value="">All services/directories</option>' +
        services.map((service) => '<option value="' + escapeHtml(service) + '">' + escapeHtml(service) + '</option>').join("");
      typeFilterEl.value = types.includes(currentType) ? currentType : "";
      serviceFilterEl.value = services.includes(currentService) ? currentService : "";
    }

    function uniqueSorted(values) {
      return [...new Set(values.filter(Boolean))].sort((a, b) => a.localeCompare(b));
    }

    function summarySuffix(type, service) {
      const filters = [];
      if (type) filters.push("type " + type);
      if (service) filters.push("service " + service);
      return filters.length ? " filtered by " + filters.join(", ") : "";
    }

    async function updateScenario(endpoint, scenarioNumber, card) {
      const status = card.querySelector(".status");
      status.className = "status";
      status.textContent = "Saving...";
      try {
        const response = await fetch("/internal/v1/in-memory", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            path: endpoint.path,
            method: endpoint.method,
            scenarioNumber: scenarioNumber,
          }),
        });
        if (!response.ok) throw new Error(await response.text());
        endpoint.currentScenario = scenarioNumber;
        status.className = "status ok";
        status.textContent = "Scenario updated";
      } catch (error) {
        status.className = "status error";
        status.textContent = "Update failed: " + error.message;
      }
    }

    function escapeHtml(value) {
      return String(value || "")
        .replaceAll("&", "&amp;")
        .replaceAll("<", "&lt;")
        .replaceAll(">", "&gt;")
        .replaceAll('"', "&quot;")
        .replaceAll("'", "&#039;");
    }

    searchEl.addEventListener("input", render);
    typeFilterEl.addEventListener("change", render);
    serviceFilterEl.addEventListener("change", render);
    refreshEl.addEventListener("click", load);
    load();
  </script>
</body>
</html>`
