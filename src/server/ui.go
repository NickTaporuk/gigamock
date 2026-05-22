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
      --bg: #f6f7f9;
      --panel: #ffffff;
      --text: #172033;
      --muted: #5d6b7c;
      --subtle: #6b7788;
      --line: #d7dce2;
      --input-line: #c7d0db;
      --accent: #1e6bff;
      --accent-soft: #e8f0ff;
      --danger: #b42318;
      --ok: #15803d;
      --ok-soft: #ecfdf3;
      --warn: #c2410c;
      --warn-soft: #fff7ed;
      --broker: #f59e0b;
      --shadow: none;
    }

    * { box-sizing: border-box; }
    body {
      margin: 0;
      background: var(--bg);
      color: var(--text);
      font: 14px/1.45 system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
    }
    header {
      max-width: 1216px;
      margin: 28px auto 0;
      padding: 0 32px;
      border: 1px solid var(--line);
      border-bottom: 0;
      border-radius: 8px 8px 0 0;
      background: var(--panel);
    }
    .bar {
      margin: 0 auto;
      padding: 28px 0 22px;
      display: grid;
      gap: 20px;
    }
    h1 {
      margin: 0;
      font-size: 34px;
      line-height: 1.2;
      font-weight: 740;
      letter-spacing: 0;
    }
    .toolbar {
      display: grid;
      grid-template-columns: minmax(300px, 432px) minmax(160px, 176px) minmax(220px, 256px) 112px 84px;
      gap: 10px;
      align-items: center;
      min-width: 0;
    }
    input, select, button {
      height: 44px;
      border: 1px solid var(--input-line);
      border-radius: 6px;
      background: #fff;
      color: var(--text);
      font: inherit;
      min-width: 0;
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
      width: 100%;
      white-space: nowrap;
    }
    button:disabled {
      cursor: progress;
      opacity: .7;
    }
    button.primary {
      border-color: var(--accent);
      background: var(--accent);
      color: #fff;
    }
    .live-status {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      color: var(--text);
      font-size: 15px;
      font-weight: 700;
      white-space: nowrap;
      width: 84px;
    }
    .live-dot {
      width: 12px;
      height: 12px;
      border-radius: 999px;
      background: var(--ok);
      box-shadow: none;
    }
    .live-status.paused .live-dot {
      background: var(--muted);
      box-shadow: none;
    }
    .live-status.error {
      color: var(--danger);
    }
    .live-status.error .live-dot {
      background: var(--danger);
      box-shadow: 0 0 0 4px rgba(180, 35, 24, .1);
    }
    main {
      max-width: 1216px;
      margin: 0 auto;
      padding: 0 32px 28px;
      border: 1px solid var(--line);
      border-top: 0;
      border-radius: 0 0 8px 8px;
      background: var(--panel);
    }
    .summary {
      color: var(--text);
      margin-bottom: 24px;
      font-size: 18px;
      font-weight: 720;
    }
    .metrics {
      display: flex;
      overflow-x: auto;
      gap: 20px;
      margin-bottom: 38px;
      padding-bottom: 2px;
    }
    .metric-card {
      position: relative;
      flex: 0 0 248px;
      background: #f8fafc;
      border: 1px solid var(--line);
      border-radius: 8px;
      box-shadow: var(--shadow);
      padding: 22px;
      min-width: 0;
      min-height: 96px;
    }
    .metric-card::after {
      content: "";
      position: absolute;
      left: 22px;
      bottom: 16px;
      width: 92px;
      height: 6px;
      border-radius: 3px;
      background: var(--ok);
    }
    .metric-card:nth-child(2)::after {
      width: 128px;
      background: var(--accent);
    }
    .metric-card:nth-child(3)::after {
      width: 78px;
      background: #94a3b8;
    }
    .metric-card:nth-child(n+4)::after {
      width: 146px;
      background: var(--broker);
    }
    .metric-card.has-data::after {
      display: none;
    }
    .metric-card h2 {
      margin: 0 0 10px;
      font-size: 17px;
      line-height: 1.2;
      font-weight: 760;
      letter-spacing: 0;
    }
    .metric-list {
      display: grid;
      gap: 8px;
    }
    .metric-row {
      display: grid;
      gap: 8px;
      align-items: start;
      border-top: 1px solid var(--line);
      padding-top: 8px;
    }
    .metric-row:first-child {
      border-top: 0;
      padding-top: 0;
    }
    .metric-key {
      font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
      font-size: 12px;
      line-height: 1.35;
      overflow-wrap: break-word;
      word-break: break-word;
    }
    .metric-values {
      display: flex;
      flex-wrap: wrap;
      justify-content: flex-start;
      gap: 6px;
    }
    .metric-pill {
      display: inline-flex;
      align-items: center;
      height: 22px;
      padding: 0 7px;
      border-radius: 999px;
      background: #f2f4f7;
      color: #344054;
      font-size: 12px;
      font-weight: 680;
      max-width: 100%;
      min-width: 0;
      white-space: normal;
    }
    .metric-empty {
      color: var(--muted);
      font-size: 14px;
    }
    .grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
      gap: 36px;
    }
    .endpoint {
      background: var(--panel);
      border: 1px solid var(--line);
      border-radius: 8px;
      box-shadow: var(--shadow);
      padding: 20px;
      display: grid;
      grid-template-rows: minmax(150px, auto) auto;
      gap: 18px;
      min-height: 300px;
    }
    .meta {
      display: flex;
      flex-wrap: wrap;
      gap: 10px;
      margin-bottom: 20px;
      align-items: center;
    }
    .badge {
      display: inline-flex;
      align-items: center;
      height: 24px;
      padding: 0 8px;
      border-radius: 4px;
      background: var(--accent-soft);
      color: var(--accent);
      font-size: 12px;
      font-weight: 700;
      text-transform: uppercase;
    }
    .service {
      background: var(--warn-soft);
      color: #9a3412;
    }
    .method {
      background: var(--ok-soft);
      color: var(--ok);
    }
    .path {
      margin: 0 0 14px;
      font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
      font-size: 22px;
      font-weight: 720;
      line-height: 1.25;
      overflow-wrap: anywhere;
    }
    .desc {
      margin: 0;
      color: var(--muted);
      font-size: 14px;
      display: -webkit-box;
      overflow: hidden;
      -webkit-box-orient: vertical;
      -webkit-line-clamp: 2;
    }
    .file {
      margin-top: 10px;
      color: var(--muted);
      font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
      font-size: 12px;
      overflow-wrap: anywhere;
    }
    .provenance {
      margin-top: 18px;
      display: grid;
      gap: 4px;
      color: var(--subtle);
      font-size: 13px;
    }
    .provenance div {
      min-width: 0;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
    .provenance span {
      color: var(--subtle);
      font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    }
    .control {
      display: grid;
      grid-template-columns: minmax(0, 1fr) 72px;
      gap: 8px 12px;
      align-items: end;
    }
    .control label {
      grid-column: 1 / -1;
      color: var(--text);
      font-size: 13px;
      font-weight: 700;
    }
    .control select {
      width: 100%;
      height: 38px;
      padding: 0 12px;
    }
    .control button {
      height: 38px;
      border-color: var(--text);
      background: var(--text);
      color: #fff;
    }
    .status {
      grid-column: 1 / -1;
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
      header {
        margin: 12px;
        padding: 0 18px;
      }
      main {
        margin: 0 12px 12px;
        padding: 0 18px 20px;
      }
      h1 { font-size: 28px; }
      .metrics { grid-template-columns: 1fr; }
      .toolbar { grid-template-columns: 1fr; align-items: stretch; }
      input { width: 100%; }
      .grid { grid-template-columns: 1fr; gap: 16px; }
      .control { grid-template-columns: minmax(0, 1fr) 72px; }
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
        <span id="liveStatus" class="live-status"><span class="live-dot"></span><span>Live</span></span>
      </div>
    </div>
  </header>
  <main>
    <div id="summary" class="summary">Loading scenarios...</div>
    <section id="metrics" class="metrics" aria-live="polite"></section>
    <section id="endpoints" class="grid" aria-live="polite"></section>
  </main>
  <script>
    const endpointsEl = document.querySelector("#endpoints");
    const metricsEl = document.querySelector("#metrics");
    const summaryEl = document.querySelector("#summary");
    const searchEl = document.querySelector("#search");
    const typeFilterEl = document.querySelector("#typeFilter");
    const serviceFilterEl = document.querySelector("#serviceFilter");
    const refreshEl = document.querySelector("#refresh");
    const liveStatusEl = document.querySelector("#liveStatus");
    let endpoints = [];
    let metrics = {};
    let metricsPolling = false;
    let loading = false;
    let loadRequestId = 0;
    let metricsSocket = null;
    let metricsSocketReconnectTimer = null;
    let metricsFallbackTimer = null;
    const metricsIntervalMs = 2000;

    async function load() {
      if (loading) return;
      loading = true;
      const requestId = ++loadRequestId;
      const hadContent = endpoints.length > 0;
      refreshEl.disabled = true;
      refreshEl.textContent = "Refreshing";
      if (!hadContent) {
        summaryEl.textContent = "Loading scenarios...";
      } else {
        summaryEl.textContent = "Refreshing scenarios...";
      }
      try {
        const [response, metricsData] = await Promise.all([
          fetch("/internal/v1/scenarios"),
          loadMetrics(),
        ]);
        if (!response.ok) throw new Error(await response.text());
        const data = await response.json();
        if (requestId !== loadRequestId) return;
        endpoints = data.endpoints || [];
        metrics = metricsData;
        populateFilters();
        render();
      } catch (error) {
        if (requestId !== loadRequestId) return;
        summaryEl.textContent = "Failed to refresh scenarios: " + error.message;
        if (!hadContent) {
          endpointsEl.innerHTML = "";
          metricsEl.innerHTML = "";
        }
      } finally {
        if (requestId === loadRequestId) {
          loading = false;
          refreshEl.disabled = false;
          refreshEl.textContent = "Refresh";
        }
      }
    }

    async function refreshMetrics() {
      if (metricsSocket && metricsSocket.readyState === WebSocket.OPEN) return;
      if (metricsPolling || loading || document.hidden) return;
      metricsPolling = true;
      setLiveStatus("updating");
      try {
        metrics = await loadMetrics();
        renderMetrics();
        setLiveStatus("live");
      } catch (error) {
        setLiveStatus("error");
      } finally {
        metricsPolling = false;
      }
    }

    function connectMetricsSocket() {
      if (!("WebSocket" in window)) {
        startMetricsFallback();
        return;
      }
      if (metricsSocket && (metricsSocket.readyState === WebSocket.OPEN || metricsSocket.readyState === WebSocket.CONNECTING)) {
        return;
      }

      const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
      metricsSocket = new WebSocket(protocol + "//" + window.location.host + "/internal/v1/mock-ui/ws");

      metricsSocket.addEventListener("open", () => {
        stopMetricsFallback();
        setLiveStatus("live");
      });

      metricsSocket.addEventListener("message", (event) => {
        if (loading || document.hidden) return;
        try {
          metrics = JSON.parse(event.data || "{}");
          renderMetrics();
          setLiveStatus("live");
        } catch (_) {
          setLiveStatus("error");
        }
      });

      metricsSocket.addEventListener("close", () => {
        metricsSocket = null;
        if (!document.hidden) {
          startMetricsFallback();
          scheduleMetricsSocketReconnect();
        }
      });

      metricsSocket.addEventListener("error", () => {
        setLiveStatus("error");
      });
    }

    function scheduleMetricsSocketReconnect() {
      if (metricsSocketReconnectTimer) return;
      metricsSocketReconnectTimer = window.setTimeout(() => {
        metricsSocketReconnectTimer = null;
        connectMetricsSocket();
      }, 3000);
    }

    function closeMetricsSocket() {
      if (metricsSocketReconnectTimer) {
        window.clearTimeout(metricsSocketReconnectTimer);
        metricsSocketReconnectTimer = null;
      }
      if (metricsSocket) {
        metricsSocket.close();
        metricsSocket = null;
      }
    }

    function startMetricsFallback() {
      if (metricsFallbackTimer) return;
      metricsFallbackTimer = window.setInterval(refreshMetrics, metricsIntervalMs);
    }

    function stopMetricsFallback() {
      if (!metricsFallbackTimer) return;
      window.clearInterval(metricsFallbackTimer);
      metricsFallbackTimer = null;
    }

    async function loadMetrics() {
      const result = {};
      const configs = [
        ["grpc", "/internal/v1/grpc/metrics"],
        ["graphql", "/internal/v1/graphql/metrics"],
        ["kafka", "/internal/v1/kafka/metrics"],
        ["nats", "/internal/v1/nats/metrics"],
        ["rabbitmq", "/internal/v1/rabbitmq/metrics"],
        ["mqtt", "/internal/v1/mqtt/metrics"],
        ["websocket", "/internal/v1/websocket/metrics"],
        ["s3", "/internal/v1/s3/metrics"],
        ["sqs", "/internal/v1/sqs/metrics"],
        ["sns", "/internal/v1/sns/metrics"],
        ["pubsub", "/internal/v1/pubsub/metrics"],
        ["servicebus", "/internal/v1/servicebus/metrics"],
        ["soap", "/internal/v1/soap/metrics"],
      ];
      await Promise.all(configs.map(async ([type, url]) => {
        try {
          const response = await fetch(url);
          result[type] = response.ok ? await response.json() : {};
        } catch (_) {
          result[type] = {};
        }
      }));
      return result;
    }

    function setLiveStatus(state) {
      liveStatusEl.className = "live-status";
      if (state === "paused") {
        liveStatusEl.classList.add("paused");
        liveStatusEl.lastElementChild.textContent = "Paused";
        return;
      }
      if (state === "error") {
        liveStatusEl.classList.add("error");
        liveStatusEl.lastElementChild.textContent = "Live error";
        return;
      }
      if (state === "updating") {
        liveStatusEl.lastElementChild.textContent = "Updating";
        return;
      }
      liveStatusEl.lastElementChild.textContent = "Live";
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
      renderMetrics();
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

    function renderMetrics() {
      metricsEl.innerHTML = [
        metricCard("gRPC", metrics.grpc || {}),
        metricCard("GraphQL", metrics.graphql || {}),
        metricCard("Kafka", metrics.kafka || {}),
        metricCard("NATS", metrics.nats || {}),
        metricCard("RabbitMQ", metrics.rabbitmq || {}),
        metricCard("MQTT", metrics.mqtt || {}),
        metricCard("WebSocket", metrics.websocket || {}),
        metricCard("S3", metrics.s3 || {}),
        metricCard("SQS", metrics.sqs || {}),
        metricCard("SNS", metrics.sns || {}),
        metricCard("Pub/Sub", metrics.pubsub || {}),
        metricCard("Azure Service Bus", metrics.servicebus || {}),
        metricCard("SOAP", metrics.soap || {}),
      ].join("");
    }

    function metricCard(title, data) {
      const keys = Object.keys(data || {}).sort((a, b) => a.localeCompare(b));
      const rows = keys.map((key) => metricRow(key, data[key])).join("");
      return '<article class="metric-card' + (rows ? " has-data" : "") + '">' +
        '<h2>' + escapeHtml(title) + '</h2>' +
        (rows ? '<div class="metric-list">' + rows + '</div>' : '<div class="metric-empty">No runtime calls yet</div>') +
      '</article>';
    }

    function metricRow(key, values) {
      const pills = Object.keys(values || {}).sort((a, b) => a.localeCompare(b)).map((name) => {
        return '<span class="metric-pill">' + escapeHtml(name) + ': ' + escapeHtml(values[name]) + '</span>';
      }).join("");
      return '<div class="metric-row">' +
        '<div class="metric-key">' + escapeHtml(key) + '</div>' +
        '<div class="metric-values">' + pills + '</div>' +
      '</div>';
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
    document.addEventListener("visibilitychange", () => {
      if (document.hidden) {
        setLiveStatus("paused");
        closeMetricsSocket();
        stopMetricsFallback();
      } else {
        connectMetricsSocket();
        refreshMetrics();
      }
    });
    load();
    connectMetricsSocket();
  </script>
</body>
</html>`
