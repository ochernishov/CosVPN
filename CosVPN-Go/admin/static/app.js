/* ===== CosVPN Dashboard — app.js ===== */

// ---- Tab Switching ----
document.querySelectorAll('.tab').forEach(function(tab) {
    tab.addEventListener('click', function() {
        switchTab(tab.dataset.tab);
    });
});

function switchTab(name) {
    document.querySelectorAll('.tab').forEach(function(t) {
        t.classList.remove('active');
    });
    document.querySelectorAll('.tab-content').forEach(function(c) {
        c.classList.remove('active');
    });
    var tabBtn = document.querySelector('[data-tab="' + name + '"]');
    var tabContent = document.getElementById(name);
    if (tabBtn) tabBtn.classList.add('active');
    if (tabContent) tabContent.classList.add('active');
    loadTabData(name);
}

// ---- API Helper ----
async function api(method, url, body) {
    var opts = {
        method: method,
        headers: { 'Content-Type': 'application/json' }
    };
    if (body) {
        opts.body = JSON.stringify(body);
    }
    try {
        var res = await fetch(url, opts);
        if (res.status === 401) {
            location.href = '/login.html';
            return null;
        }
        return res;
    } catch (err) {
        console.error('API error:', err);
        return null;
    }
}

// ---- Utility Functions ----
function formatBytes(bytes) {
    if (!bytes || bytes === 0) return '0 B';
    var units = ['B', 'KB', 'MB', 'GB', 'TB'];
    var i = Math.floor(Math.log(bytes) / Math.log(1024));
    if (i >= units.length) i = units.length - 1;
    return (bytes / Math.pow(1024, i)).toFixed(1) + ' ' + units[i];
}

function formatTime(isoStr) {
    if (!isoStr || isoStr === 'never') return 'Never';
    try {
        var d = new Date(isoStr);
        if (isNaN(d.getTime())) return isoStr;
        var now = new Date();
        var diffMs = now - d;
        var diffSec = Math.floor(diffMs / 1000);

        if (diffSec < 60) return diffSec + 's ago';
        if (diffSec < 3600) return Math.floor(diffSec / 60) + 'm ago';
        if (diffSec < 86400) return Math.floor(diffSec / 3600) + 'h ago';
        return Math.floor(diffSec / 86400) + 'd ago';
    } catch (e) {
        return isoStr;
    }
}

// Sanitize strings for safe DOM insertion — prevents XSS
function escapeHtml(str) {
    if (!str) return '';
    var div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
}

function getProgressClass(val) {
    if (val >= 90) return 'danger';
    if (val >= 70) return 'warning';
    return '';
}

// Safe DOM helper: creates an element, sets textContent (not HTML), returns it
function createTextEl(tag, text, className) {
    var el = document.createElement(tag);
    el.textContent = text || '';
    if (className) el.className = className;
    return el;
}

// ---- Dashboard ----
async function loadDashboard() {
    var res = await api('GET', '/api/status');
    if (!res) return;

    var data = await res.json();

    // Server card
    document.getElementById('stat-uptime').textContent = data.uptime || '--';
    document.getElementById('stat-ip').textContent = data.publicIP || '--';
    document.getElementById('stat-cpu').textContent = data.cpu + '%';
    document.getElementById('stat-ram').textContent = data.ram + '%';
    document.getElementById('stat-disk').textContent = data.disk + '%';

    // VPN card — built with safe DOM methods
    var vpnStatusEl = document.getElementById('stat-vpn-status');
    vpnStatusEl.textContent = '';
    var badge = document.createElement('span');
    var dot = document.createElement('span');
    if (data.vpnUp) {
        badge.className = 'status-badge online';
        dot.className = 'status-dot online';
        badge.appendChild(dot);
        badge.appendChild(document.createTextNode('Running'));
    } else {
        badge.className = 'status-badge offline';
        dot.className = 'status-dot offline';
        badge.appendChild(dot);
        badge.appendChild(document.createTextNode('Stopped'));
    }
    vpnStatusEl.appendChild(badge);

    document.getElementById('stat-port').textContent = data.port || '--';

    // Clients card
    document.getElementById('stat-total-clients').textContent = data.totalClients;
    document.getElementById('stat-online-clients').textContent = data.onlineClients;

    // Obfuscation card
    var obfsText = data.obfsMode || 'none';
    document.getElementById('stat-obfs-mode').textContent = obfsText.charAt(0).toUpperCase() + obfsText.slice(1);

    // Progress bars
    updateProgressBar('cpu', data.cpu);
    updateProgressBar('ram', data.ram);
    updateProgressBar('disk', data.disk);
}

function updateProgressBar(name, value) {
    var fill = document.getElementById('progress-' + name);
    var valEl = document.getElementById('progress-' + name + '-val');
    if (!fill || !valEl) return;

    fill.style.width = value + '%';
    valEl.textContent = value + '%';

    fill.className = 'progress-fill';
    var extraClass = getProgressClass(value);
    if (extraClass) fill.classList.add(extraClass);
}

// ---- Clients ----
async function loadClients() {
    var res = await api('GET', '/api/clients');
    if (!res) return;

    var clients = await res.json();
    var tbody = document.getElementById('clients-tbody');
    var empty = document.getElementById('clients-empty');

    // Clear existing rows
    tbody.textContent = '';

    if (!clients || clients.length === 0) {
        empty.style.display = 'block';
        return;
    }

    empty.style.display = 'none';

    clients.forEach(function(c) {
        var tr = document.createElement('tr');

        // Name cell
        var tdName = document.createElement('td');
        var strong = document.createElement('strong');
        strong.textContent = c.name;
        tdName.appendChild(strong);
        tr.appendChild(tdName);

        // IP cell
        var tdIP = createTextEl('td', c.ip, 'text-mono');
        tr.appendChild(tdIP);

        // Status cell
        var tdStatus = document.createElement('td');
        var statusBadge = document.createElement('span');
        var statusDot = document.createElement('span');
        var statusClass = c.online ? 'online' : 'offline';
        var statusLabel = c.online ? 'Online' : 'Offline';
        statusBadge.className = 'status-badge ' + statusClass;
        statusDot.className = 'status-dot ' + statusClass;
        statusBadge.appendChild(statusDot);
        statusBadge.appendChild(document.createTextNode(statusLabel));
        tdStatus.appendChild(statusBadge);
        tr.appendChild(tdStatus);

        // Last Handshake cell
        var tdHandshake = createTextEl('td', formatTime(c.lastHandshake));
        tr.appendChild(tdHandshake);

        // Upload cell
        var tdUp = createTextEl('td', formatBytes(c.transferUp), 'text-right text-mono');
        tr.appendChild(tdUp);

        // Download cell
        var tdDown = createTextEl('td', formatBytes(c.transferDown), 'text-right text-mono');
        tr.appendChild(tdDown);

        // Actions cell
        var tdActions = document.createElement('td');
        var actionsDiv = document.createElement('div');
        actionsDiv.className = 'actions-cell';

        var btnQR = document.createElement('button');
        btnQR.className = 'btn-icon';
        btnQR.textContent = 'QR';
        btnQR.title = 'QR Code';
        btnQR.addEventListener('click', (function(name) {
            return function() { showQR(name); };
        })(c.name));

        var btnDL = document.createElement('button');
        btnDL.className = 'btn-icon';
        btnDL.textContent = 'DL';
        btnDL.title = 'Download Config';
        btnDL.addEventListener('click', (function(name) {
            return function() { downloadConf(name); };
        })(c.name));

        var btnDel = document.createElement('button');
        btnDel.className = 'btn-danger';
        btnDel.textContent = 'Delete';
        btnDel.addEventListener('click', (function(name) {
            return function() { deleteClient(name); };
        })(c.name));

        actionsDiv.appendChild(btnQR);
        actionsDiv.appendChild(btnDL);
        actionsDiv.appendChild(btnDel);
        tdActions.appendChild(actionsDiv);
        tr.appendChild(tdActions);

        tbody.appendChild(tr);
    });
}

function showAddClient() {
    var modal = document.getElementById('modal-content');
    modal.textContent = '';

    // Close button
    var closeBtn = document.createElement('button');
    closeBtn.className = 'modal-close';
    closeBtn.textContent = '\u00D7';
    closeBtn.addEventListener('click', hideModal);
    modal.appendChild(closeBtn);

    // Title
    var title = createTextEl('h3', 'Add New Client', 'modal-title');
    modal.appendChild(title);

    // Form group
    var formGroup = document.createElement('div');
    formGroup.className = 'form-group';

    var label = document.createElement('label');
    label.textContent = 'Client Name';
    label.setAttribute('for', 'new-client-name');
    formGroup.appendChild(label);

    var input = document.createElement('input');
    input.type = 'text';
    input.id = 'new-client-name';
    input.className = 'form-control';
    input.placeholder = 'e.g. my-phone';
    input.addEventListener('keydown', function(e) {
        if (e.key === 'Enter') addClient();
    });
    formGroup.appendChild(input);

    var hint = createTextEl('span', 'Alphanumeric and dashes only', 'form-hint');
    formGroup.appendChild(hint);

    modal.appendChild(formGroup);

    // Actions
    var actions = document.createElement('div');
    actions.className = 'modal-actions';

    var cancelBtn = document.createElement('button');
    cancelBtn.className = 'btn-ghost';
    cancelBtn.textContent = 'Cancel';
    cancelBtn.addEventListener('click', hideModal);
    actions.appendChild(cancelBtn);

    var createBtn = document.createElement('button');
    createBtn.className = 'btn-primary';
    createBtn.textContent = 'Create Client';
    createBtn.addEventListener('click', function() { addClient(); });
    actions.appendChild(createBtn);

    modal.appendChild(actions);

    document.getElementById('modal-overlay').classList.remove('hidden');

    setTimeout(function() { input.focus(); }, 100);
}

async function addClient() {
    var input = document.getElementById('new-client-name');
    if (!input) return;

    var name = input.value.trim();
    if (!name) {
        input.style.borderColor = 'var(--danger)';
        return;
    }

    var res = await api('POST', '/api/clients', { name: name });
    if (!res) return;

    if (!res.ok) {
        var err = await res.json().catch(function() { return {}; });
        alert(err.error || 'Failed to create client');
        return;
    }

    var data = await res.json();

    // Show result modal with config and QR — built with safe DOM methods
    var modal = document.getElementById('modal-content');
    modal.textContent = '';

    var closeBtn = document.createElement('button');
    closeBtn.className = 'modal-close';
    closeBtn.textContent = '\u00D7';
    closeBtn.addEventListener('click', hideModal);
    modal.appendChild(closeBtn);

    var title = createTextEl('h3', 'Client Created: ' + data.name, 'modal-title');
    modal.appendChild(title);

    if (data.ip) {
        var ipInfo = document.createElement('p');
        ipInfo.style.cssText = 'color:var(--text-secondary);margin-bottom:1rem;';
        ipInfo.textContent = 'IP: ';
        var ipSpan = createTextEl('span', data.ip, 'text-mono');
        ipInfo.appendChild(ipSpan);
        modal.appendChild(ipInfo);
    }

    if (data.qr) {
        var qrDiv = document.createElement('div');
        qrDiv.className = 'qr-container';
        var qrImg = document.createElement('img');
        qrImg.src = 'data:image/png;base64,' + data.qr;
        qrImg.alt = 'QR Code';
        qrDiv.appendChild(qrImg);
        modal.appendChild(qrDiv);
    }

    if (data.config) {
        var configBlock = document.createElement('div');
        configBlock.className = 'config-block';
        configBlock.textContent = data.config;
        modal.appendChild(configBlock);
    }

    var actions = document.createElement('div');
    actions.className = 'modal-actions';

    var dlBtn = document.createElement('button');
    dlBtn.className = 'btn-ghost';
    dlBtn.textContent = 'Download .conf';
    dlBtn.addEventListener('click', function() { downloadConf(name); });
    actions.appendChild(dlBtn);

    var doneBtn = document.createElement('button');
    doneBtn.className = 'btn-primary';
    doneBtn.textContent = 'Done';
    doneBtn.addEventListener('click', hideModal);
    actions.appendChild(doneBtn);

    modal.appendChild(actions);

    loadClients();
}

async function deleteClient(name) {
    if (!confirm('Delete client "' + name + '"? This action cannot be undone.')) return;

    var res = await api('DELETE', '/api/clients/' + encodeURIComponent(name));
    if (!res) return;

    if (!res.ok) {
        var err = await res.json().catch(function() { return {}; });
        alert(err.error || 'Failed to delete client');
        return;
    }

    loadClients();
}

function showQR(name) {
    var modal = document.getElementById('modal-content');
    modal.textContent = '';

    var closeBtn = document.createElement('button');
    closeBtn.className = 'modal-close';
    closeBtn.textContent = '\u00D7';
    closeBtn.addEventListener('click', hideModal);
    modal.appendChild(closeBtn);

    var title = createTextEl('h3', 'QR Code: ' + name, 'modal-title');
    modal.appendChild(title);

    var qrDiv = document.createElement('div');
    qrDiv.className = 'qr-container';
    var qrImg = document.createElement('img');
    qrImg.src = '/api/clients/' + encodeURIComponent(name) + '/qr';
    qrImg.alt = 'QR Code';
    qrImg.onerror = function() { this.alt = 'Failed to load QR code'; };
    qrDiv.appendChild(qrImg);
    modal.appendChild(qrDiv);

    var actions = document.createElement('div');
    actions.className = 'modal-actions';

    var dlBtn = document.createElement('button');
    dlBtn.className = 'btn-ghost';
    dlBtn.textContent = 'Download .conf';
    dlBtn.addEventListener('click', function() { downloadConf(name); });
    actions.appendChild(dlBtn);

    var closeModalBtn = document.createElement('button');
    closeModalBtn.className = 'btn-primary';
    closeModalBtn.textContent = 'Close';
    closeModalBtn.addEventListener('click', hideModal);
    actions.appendChild(closeModalBtn);

    modal.appendChild(actions);
    document.getElementById('modal-overlay').classList.remove('hidden');
}

function downloadConf(name) {
    window.open('/api/clients/' + encodeURIComponent(name) + '/conf', '_blank');
}

// ---- Settings ----
async function loadSettings() {
    var res = await api('GET', '/api/settings');
    if (!res) return;

    var data = await res.json();

    var obfsSelect = document.getElementById('setting-obfs');
    if (obfsSelect) {
        obfsSelect.value = data.obfuscationMode || 'none';
    }

    var dnsInput = document.getElementById('setting-dns');
    if (dnsInput) {
        dnsInput.value = (data.dns || []).join(', ');
    }

    var mtuInput = document.getElementById('setting-mtu');
    if (mtuInput) {
        mtuInput.value = data.mtu || '';
    }

    var portInput = document.getElementById('setting-port');
    if (portInput) {
        portInput.value = data.listenPort || '';
    }

    var subnetInput = document.getElementById('setting-subnet');
    if (subnetInput) {
        subnetInput.value = data.subnet || '';
    }

    // Clear save status
    var statusEl = document.getElementById('save-status');
    if (statusEl) statusEl.textContent = '';
}

function handleSaveSettings(e) {
    e.preventDefault();
    saveSettings();
    return false;
}

async function saveSettings() {
    var btn = document.getElementById('btn-save-settings');
    var statusEl = document.getElementById('save-status');

    var obfsMode = document.getElementById('setting-obfs').value;
    var dnsRaw = document.getElementById('setting-dns').value;
    var mtuRaw = document.getElementById('setting-mtu').value;

    var dns = dnsRaw.split(',').map(function(s) { return s.trim(); }).filter(function(s) { return s !== ''; });
    var mtu = parseInt(mtuRaw, 10) || 0;

    if (btn) {
        btn.disabled = true;
        btn.textContent = 'Applying...';
    }
    if (statusEl) statusEl.textContent = '';

    var res = await api('PUT', '/api/settings', {
        obfuscationMode: obfsMode,
        dns: dns,
        mtu: mtu
    });

    if (btn) {
        btn.disabled = false;
        btn.textContent = 'Apply Settings';
    }

    if (!res) {
        if (statusEl) {
            statusEl.style.color = 'var(--danger)';
            statusEl.textContent = 'Connection error';
        }
        return;
    }

    if (res.ok) {
        if (statusEl) {
            statusEl.style.color = 'var(--success)';
            statusEl.textContent = 'Settings saved';
        }
        setTimeout(function() {
            if (statusEl) statusEl.textContent = '';
        }, 3000);
    } else {
        var err = await res.json().catch(function() { return {}; });
        if (statusEl) {
            statusEl.style.color = 'var(--danger)';
            statusEl.textContent = err.error || 'Failed to save';
        }
    }
}

// ---- Logs ----
async function loadLogs() {
    var res = await api('GET', '/api/logs?limit=100');
    if (!res) return;

    var logs = await res.json();
    var tbody = document.getElementById('logs-tbody');
    var empty = document.getElementById('logs-empty');

    // Clear existing rows
    tbody.textContent = '';

    if (!logs || logs.length === 0) {
        empty.style.display = 'block';
        return;
    }

    empty.style.display = 'none';

    logs.forEach(function(entry) {
        var tr = document.createElement('tr');

        // Time cell
        var tdTime = document.createElement('td');
        tdTime.className = 'text-mono';
        tdTime.style.fontSize = '0.8rem';
        if (entry.time) {
            try {
                var d = new Date(entry.time);
                tdTime.textContent = d.toLocaleString();
            } catch (e) {
                tdTime.textContent = entry.time;
            }
        }
        tr.appendChild(tdTime);

        // Type cell
        var tdType = document.createElement('td');
        var typeSpan = document.createElement('span');
        var typeClass = 'default';
        if (entry.type === 'client_add') typeClass = 'client_add';
        else if (entry.type === 'client_remove') typeClass = 'client_remove';
        else if (entry.type === 'settings') typeClass = 'settings';
        typeSpan.className = 'log-type ' + typeClass;
        typeSpan.textContent = entry.type || '--';
        tdType.appendChild(typeSpan);
        tr.appendChild(tdType);

        // Client cell
        var tdClient = createTextEl('td', entry.client || '--');
        tr.appendChild(tdClient);

        // Details cell
        var tdDetails = createTextEl('td', entry.details || '--');
        tr.appendChild(tdDetails);

        tbody.appendChild(tr);
    });
}

// ---- Modals ----
function showModal(contentHtml) {
    // This function is kept for backward compat but modals are built via DOM methods now
    document.getElementById('modal-overlay').classList.remove('hidden');
}

function hideModal() {
    document.getElementById('modal-overlay').classList.add('hidden');
    document.getElementById('modal-content').textContent = '';
}

function handleOverlayClick(e) {
    if (e.target === document.getElementById('modal-overlay')) {
        hideModal();
    }
}

// Close modal on Escape
document.addEventListener('keydown', function(e) {
    if (e.key === 'Escape') {
        var overlay = document.getElementById('modal-overlay');
        if (overlay && !overlay.classList.contains('hidden')) {
            hideModal();
        }
    }
});

// ---- Auth ----
function logout() {
    document.cookie = 'token=; Max-Age=0; path=/';
    location.href = '/login.html';
}

// ---- Tab Data Loading ----
function loadTabData(tab) {
    if (tab === 'dashboard') loadDashboard();
    else if (tab === 'clients') loadClients();
    else if (tab === 'settings') loadSettings();
    else if (tab === 'logs') loadLogs();
}

// ---- Init ----
loadDashboard();

// Auto-refresh: dashboard every 15s, logs every 10s if active
setInterval(function() {
    var dashboardTab = document.querySelector('[data-tab="dashboard"]');
    if (dashboardTab && dashboardTab.classList.contains('active')) {
        loadDashboard();
    }
}, 15000);

setInterval(function() {
    var logsTab = document.querySelector('[data-tab="logs"]');
    if (logsTab && logsTab.classList.contains('active')) {
        loadLogs();
    }
}, 10000);
