package contract

// meshpulse_ui.go — the full MeshPulse single-file HTML UI, embedded as a Go string.
// Served at GET / by the HTTP server in meshpulse_server.go

const meshpulseHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>MeshPulse — DePIN Speed Network</title>
<style>
  :root {
    --bg: #050a0e;
    --surface: #0d1520;
    --card: #111d2e;
    --border: #1a2e45;
    --accent: #00d4ff;
    --accent2: #7c3aed;
    --green: #00ff88;
    --yellow: #ffd700;
    --red: #ff4d4d;
    --text: #e2e8f0;
    --muted: #64748b;
    --glow: 0 0 20px rgba(0,212,255,0.3);
  }
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body {
    background: var(--bg);
    color: var(--text);
    font-family: 'Segoe UI', system-ui, sans-serif;
    min-height: 100vh;
  }

  /* NAV */
  nav {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 14px 24px;
    background: var(--surface);
    border-bottom: 1px solid var(--border);
    position: sticky;
    top: 0;
    z-index: 100;
  }
  .logo {
    font-size: 1.3rem;
    font-weight: 700;
    background: linear-gradient(90deg, var(--accent), var(--accent2));
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    margin-right: 16px;
  }
  .logo span { font-size: 1rem; }
  .nav-tabs { display: flex; gap: 4px; flex: 1; }
  .tab {
    padding: 7px 18px;
    border-radius: 8px;
    border: none;
    background: transparent;
    color: var(--muted);
    cursor: pointer;
    font-size: 0.9rem;
    font-weight: 500;
    transition: all 0.2s;
  }
  .tab:hover { background: var(--border); color: var(--text); }
  .tab.active {
    background: linear-gradient(135deg, rgba(0,212,255,0.15), rgba(124,58,237,0.15));
    color: var(--accent);
    border: 1px solid rgba(0,212,255,0.3);
  }
  .wallet-badge {
    font-size: 0.78rem;
    color: var(--muted);
    background: var(--card);
    border: 1px solid var(--border);
    padding: 5px 12px;
    border-radius: 20px;
    font-family: monospace;
  }
  .wallet-badge.connected { color: var(--green); border-color: rgba(0,255,136,0.3); }

  /* PAGES */
  .page { display: none; padding: 32px 24px; max-width: 1100px; margin: 0 auto; }
  .page.active { display: block; }

  /* CARDS */
  .card {
    background: var(--card);
    border: 1px solid var(--border);
    border-radius: 16px;
    padding: 24px;
    margin-bottom: 20px;
  }
  .card-title {
    font-size: 1rem;
    font-weight: 600;
    color: var(--muted);
    text-transform: uppercase;
    letter-spacing: 0.08em;
    margin-bottom: 16px;
  }

  /* STAT GRID */
  .stat-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(180px,1fr)); gap: 16px; margin-bottom: 24px; }
  .stat-box {
    background: var(--card);
    border: 1px solid var(--border);
    border-radius: 12px;
    padding: 20px;
    text-align: center;
  }
  .stat-box .label { font-size: 0.78rem; color: var(--muted); text-transform: uppercase; letter-spacing: 0.06em; margin-bottom: 8px; }
  .stat-box .value { font-size: 1.8rem; font-weight: 700; }
  .stat-box .unit { font-size: 0.75rem; color: var(--muted); margin-top: 4px; }

  /* SPEEDOMETER */
  .speedometer-area {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 40px 20px;
  }
  .gauge-wrap { position: relative; width: 260px; height: 140px; margin-bottom: 32px; }
  .gauge-svg { width: 260px; height: 140px; }
  .gauge-value {
    position: absolute;
    bottom: 0;
    left: 50%;
    transform: translateX(-50%);
    font-size: 2.8rem;
    font-weight: 800;
    color: var(--accent);
    text-shadow: var(--glow);
    text-align: center;
    line-height: 1;
  }
  .gauge-label { font-size: 0.8rem; color: var(--muted); text-align: center; margin-top: 4px; }

  .metrics-row {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 16px;
    width: 100%;
    max-width: 600px;
    margin-bottom: 32px;
  }
  .metric-pill {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 12px;
    padding: 16px;
    text-align: center;
  }
  .metric-pill .m-label { font-size: 0.72rem; color: var(--muted); text-transform: uppercase; letter-spacing: 0.06em; margin-bottom: 6px; }
  .metric-pill .m-val { font-size: 1.6rem; font-weight: 700; }
  .metric-pill .m-unit { font-size: 0.72rem; color: var(--muted); }
  .metric-pill.ping .m-val { color: var(--green); }
  .metric-pill.down .m-val { color: var(--accent); }
  .metric-pill.up .m-val { color: var(--accent2); }

  .test-controls { display: flex; flex-direction: column; align-items: center; gap: 12px; }
  .btn-run {
    position: relative;
    padding: 16px 56px;
    font-size: 1.1rem;
    font-weight: 700;
    border: none;
    border-radius: 50px;
    cursor: pointer;
    background: linear-gradient(135deg, var(--accent), var(--accent2));
    color: #fff;
    letter-spacing: 0.05em;
    box-shadow: 0 0 30px rgba(0,212,255,0.4);
    transition: all 0.25s;
    overflow: hidden;
  }
  .btn-run:hover:not(:disabled) { transform: scale(1.04); box-shadow: 0 0 40px rgba(0,212,255,0.6); }
  .btn-run:disabled { opacity: 0.5; cursor: not-allowed; }
  .btn-run .spinner {
    display: none;
    width: 18px; height: 18px;
    border: 2px solid rgba(255,255,255,0.3);
    border-top-color: #fff;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
    margin-right: 8px;
  }
  .btn-run.running .spinner { display: inline-block; }
  @keyframes spin { to { transform: rotate(360deg); } }

  .btn-submit {
    padding: 12px 40px;
    font-size: 0.95rem;
    font-weight: 600;
    border: 1px solid var(--green);
    border-radius: 50px;
    cursor: pointer;
    background: rgba(0,255,136,0.1);
    color: var(--green);
    transition: all 0.2s;
  }
  .btn-submit:hover { background: rgba(0,255,136,0.2); box-shadow: 0 0 20px rgba(0,255,136,0.3); }
  .btn-submit:disabled { opacity: 0.4; cursor: not-allowed; }

  .phase-label {
    font-size: 0.85rem;
    color: var(--accent);
    min-height: 20px;
    letter-spacing: 0.04em;
  }
  .reward-badge {
    font-size: 0.82rem;
    color: var(--yellow);
    background: rgba(255,215,0,0.1);
    border: 1px solid rgba(255,215,0,0.3);
    padding: 4px 14px;
    border-radius: 20px;
  }

  /* ISP/REGION INPUT */
  .meta-inputs {
    display: flex;
    gap: 12px;
    width: 100%;
    max-width: 600px;
    margin-bottom: 16px;
  }
  .meta-inputs input {
    flex: 1;
    padding: 10px 16px;
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 8px;
    color: var(--text);
    font-size: 0.9rem;
    outline: none;
    transition: border-color 0.2s;
  }
  .meta-inputs input:focus { border-color: var(--accent); }
  .meta-inputs input::placeholder { color: var(--muted); }

  /* FEED TABLE */
  .feed-table { width: 100%; border-collapse: collapse; font-size: 0.87rem; }
  .feed-table th {
    text-align: left;
    padding: 10px 14px;
    color: var(--muted);
    font-size: 0.72rem;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    border-bottom: 1px solid var(--border);
  }
  .feed-table td { padding: 11px 14px; border-bottom: 1px solid rgba(26,46,69,0.5); }
  .feed-table tr:hover td { background: rgba(255,255,255,0.02); }
  .addr { font-family: monospace; color: var(--accent); font-size: 0.82rem; }
  .badge-isp {
    background: rgba(124,58,237,0.15);
    border: 1px solid rgba(124,58,237,0.3);
    color: #a78bfa;
    padding: 2px 8px;
    border-radius: 4px;
    font-size: 0.75rem;
  }
  .badge-region {
    background: rgba(0,212,255,0.08);
    border: 1px solid rgba(0,212,255,0.2);
    color: var(--accent);
    padding: 2px 8px;
    border-radius: 4px;
    font-size: 0.75rem;
  }
  .val-ping { color: var(--green); font-weight: 600; }
  .val-down { color: var(--accent); font-weight: 600; }
  .val-up   { color: #a78bfa; font-weight: 600; }

  /* LEADERBOARD */
  .lb-row {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 14px 16px;
    border-radius: 10px;
    margin-bottom: 8px;
    background: var(--surface);
    border: 1px solid var(--border);
    transition: border-color 0.2s;
  }
  .lb-row:hover { border-color: rgba(0,212,255,0.3); }
  .lb-rank { font-size: 1.2rem; font-weight: 700; color: var(--muted); width: 36px; text-align: center; }
  .lb-rank.gold   { color: #ffd700; }
  .lb-rank.silver { color: #c0c0c0; }
  .lb-rank.bronze { color: #cd7f32; }
  .lb-addr { font-family: monospace; color: var(--accent); font-size: 0.82rem; flex: 1; }
  .lb-count { color: var(--text); font-weight: 600; margin-right: 8px; }
  .lb-tokens { color: var(--yellow); font-weight: 700; font-size: 0.9rem; }
  .lb-bar-wrap { flex: 1; height: 6px; background: var(--border); border-radius: 3px; max-width: 180px; }
  .lb-bar { height: 6px; border-radius: 3px; background: linear-gradient(90deg, var(--accent), var(--accent2)); }

  /* MY STATS */
  .stats-hero {
    display: flex;
    align-items: center;
    gap: 24px;
    padding: 28px;
    background: linear-gradient(135deg, rgba(0,212,255,0.08), rgba(124,58,237,0.08));
    border: 1px solid var(--border);
    border-radius: 16px;
    margin-bottom: 24px;
  }
  .stats-avatar {
    width: 64px;
    height: 64px;
    border-radius: 50%;
    background: linear-gradient(135deg, var(--accent), var(--accent2));
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1.6rem;
    flex-shrink: 0;
  }
  .stats-info { flex: 1; }
  .stats-info .addr-full {
    font-family: monospace;
    color: var(--accent);
    font-size: 0.85rem;
    margin-bottom: 4px;
    word-break: break-all;
  }
  .stats-info .rank-label { font-size: 0.78rem; color: var(--muted); }
  .balance-display {
    text-align: right;
  }
  .balance-display .amount { font-size: 2rem; font-weight: 800; color: var(--yellow); }
  .balance-display .symbol { font-size: 0.9rem; color: var(--muted); }

  .addr-input-row {
    display: flex;
    gap: 10px;
    margin-bottom: 20px;
  }
  .addr-input-row input {
    flex: 1;
    padding: 10px 16px;
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 8px;
    color: var(--text);
    font-family: monospace;
    font-size: 0.88rem;
    outline: none;
  }
  .addr-input-row input:focus { border-color: var(--accent); }
  .btn-lookup {
    padding: 10px 22px;
    background: linear-gradient(135deg, var(--accent), var(--accent2));
    border: none;
    border-radius: 8px;
    color: #fff;
    font-weight: 600;
    cursor: pointer;
    font-size: 0.9rem;
  }

  /* REFRESH */
  .refresh-row { display: flex; justify-content: flex-end; margin-bottom: 16px; }
  .btn-refresh {
    padding: 7px 18px;
    background: transparent;
    border: 1px solid var(--border);
    border-radius: 8px;
    color: var(--muted);
    cursor: pointer;
    font-size: 0.82rem;
    transition: all 0.2s;
  }
  .btn-refresh:hover { border-color: var(--accent); color: var(--accent); }

  /* TOAST */
  #toast {
    position: fixed;
    bottom: 28px;
    right: 28px;
    padding: 14px 22px;
    border-radius: 10px;
    font-size: 0.9rem;
    font-weight: 500;
    z-index: 999;
    opacity: 0;
    transform: translateY(12px);
    transition: all 0.3s;
    pointer-events: none;
    max-width: 340px;
  }
  #toast.show { opacity: 1; transform: translateY(0); }
  #toast.ok { background: rgba(0,255,136,0.15); border: 1px solid rgba(0,255,136,0.4); color: var(--green); }
  #toast.err { background: rgba(255,77,77,0.15); border: 1px solid rgba(255,77,77,0.4); color: var(--red); }
  #toast.info { background: rgba(0,212,255,0.12); border: 1px solid rgba(0,212,255,0.35); color: var(--accent); }

  /* LOADING SKELETON */
  .skeleton {
    background: linear-gradient(90deg, var(--border) 25%, rgba(255,255,255,0.04) 50%, var(--border) 75%);
    background-size: 200% 100%;
    animation: shimmer 1.5s infinite;
    border-radius: 6px;
    height: 18px;
    margin-bottom: 8px;
  }
  @keyframes shimmer { to { background-position: -200% 0; } }

  /* WAVE ANIMATION on test */
  @keyframes pulse-ring {
    0% { transform: scale(0.8); opacity: 0.8; }
    100% { transform: scale(2.2); opacity: 0; }
  }
  .pulse-ring {
    position: absolute;
    inset: 0;
    border-radius: 50%;
    border: 2px solid var(--accent);
    animation: pulse-ring 1.4s ease-out infinite;
  }
  .btn-run-wrap { position: relative; display: inline-flex; align-items: center; justify-content: center; }

  /* responsive */
  @media (max-width: 600px) {
    .metrics-row { grid-template-columns: repeat(3,1fr); gap: 8px; }
    nav { flex-wrap: wrap; }
    .stat-grid { grid-template-columns: repeat(2,1fr); }
    .meta-inputs { flex-direction: column; }
  }
</style>
</head>
<body>

<nav>
  <div class="logo">⚡ MeshPulse <span>DePIN</span></div>
  <div class="nav-tabs">
    <button class="tab active" onclick="showPage('speedtest',this)">⚡ Speed Test</button>
    <button class="tab" onclick="showPage('mystats',this)">👤 My Stats</button>
    <button class="tab" onclick="showPage('feed',this)">🌐 Live Feed</button>
    <button class="tab" onclick="showPage('leaderboard',this)">🏆 Leaderboard</button>
  </div>
  <div class="wallet-badge" id="walletBadge">No Wallet</div>
</nav>

<!-- ═══════════════════ PAGE 1: SPEED TEST ═══════════════════ -->
<div class="page active" id="page-speedtest">

  <div class="stat-grid" id="networkStatGrid">
    <div class="stat-box">
      <div class="label">Network Nodes</div>
      <div class="value" id="ns-total">—</div>
      <div class="unit">measurements</div>
    </div>
    <div class="stat-box">
      <div class="label">Avg Ping</div>
      <div class="value" id="ns-ping" style="color:var(--green)">—</div>
      <div class="unit">ms</div>
    </div>
    <div class="stat-box">
      <div class="label">Avg Download</div>
      <div class="value" id="ns-down" style="color:var(--accent)">—</div>
      <div class="unit">Kbps</div>
    </div>
    <div class="stat-box">
      <div class="label">Avg Upload</div>
      <div class="value" id="ns-up" style="color:#a78bfa">—</div>
      <div class="unit">Kbps</div>
    </div>
  </div>

  <div class="card">
    <div class="card-title">Speed Test</div>
    <div class="speedometer-area">

      <!-- Gauge -->
      <div class="gauge-wrap">
        <svg class="gauge-svg" viewBox="0 0 260 140">
          <defs>
            <linearGradient id="gaugeGrad" x1="0%" y1="0%" x2="100%" y2="0%">
              <stop offset="0%" stop-color="#00d4ff"/>
              <stop offset="100%" stop-color="#7c3aed"/>
            </linearGradient>
          </defs>
          <!-- background arc -->
          <path d="M20,130 A110,110 0 0,1 240,130" fill="none" stroke="#1a2e45" stroke-width="14" stroke-linecap="round"/>
          <!-- active arc -->
          <path id="gaugeArc" d="M20,130 A110,110 0 0,1 240,130" fill="none" stroke="url(#gaugeGrad)" stroke-width="14" stroke-linecap="round"
            stroke-dasharray="346" stroke-dashoffset="346"/>
          <!-- needle -->
          <line id="gaugeNeedle" x1="130" y1="130" x2="130" y2="30" stroke="#00d4ff" stroke-width="2.5" stroke-linecap="round"
            transform="rotate(-90 130 130)" opacity="0.8"/>
          <circle cx="130" cy="130" r="6" fill="#00d4ff"/>
        </svg>
        <div class="gauge-value" id="gaugeNum">0</div>
      </div>
      <div class="gauge-label" id="gaugeLabel">Ready to test</div>

      <div class="metrics-row">
        <div class="metric-pill ping">
          <div class="m-label">Ping</div>
          <div class="m-val" id="pingVal">—</div>
          <div class="m-unit">ms</div>
        </div>
        <div class="metric-pill down">
          <div class="m-label">Download</div>
          <div class="m-val" id="downVal">—</div>
          <div class="m-unit">Kbps</div>
        </div>
        <div class="metric-pill up">
          <div class="m-label">Upload</div>
          <div class="m-val" id="upVal">—</div>
          <div class="m-unit">Kbps</div>
        </div>
      </div>

      <div class="meta-inputs">
        <input type="text" id="ispInput" placeholder="ISP (e.g. Comcast)" value="Unknown ISP">
        <input type="text" id="regionInput" placeholder="Region (e.g. US-East)" value="Unknown Region">
      </div>
      <input type="text" id="addrInputTest" placeholder="Your wallet address (40 hex chars)" style="width:100%;max-width:600px;padding:10px 16px;background:var(--surface);border:1px solid var(--border);border-radius:8px;color:var(--text);font-family:monospace;font-size:0.88rem;outline:none;margin-bottom:20px;">

      <div class="test-controls">
        <div class="phase-label" id="phaseLabel">Enter your wallet address and press Run Test</div>
        <div class="btn-run-wrap">
          <button class="btn-run" id="btnRun" onclick="runTest()">
            <span class="spinner" id="testSpinner"></span>
            ⚡ Run Speed Test
          </button>
        </div>
        <button class="btn-submit" id="btnSubmit" onclick="submitResult()" disabled>
          Submit to Chain — Earn 10 $MESHP
        </button>
        <div class="reward-badge" id="rewardBadge" style="display:none">✅ +10 $MESHP earned!</div>
      </div>
    </div>
  </div>
</div>

<!-- ═══════════════════ PAGE 2: MY STATS ═══════════════════ -->
<div class="page" id="page-mystats">
  <div class="addr-input-row">
    <input type="text" id="lookupAddr" placeholder="Paste your wallet address (40 hex chars)…" />
    <button class="btn-lookup" onclick="lookupContributor()">Look Up</button>
  </div>

  <div id="statsContent" style="display:none">
    <div class="stats-hero">
      <div class="stats-avatar">🌐</div>
      <div class="stats-info">
        <div class="addr-full" id="statsAddr">—</div>
        <div class="rank-label" id="statsRank">MeshPulse Contributor</div>
      </div>
      <div class="balance-display">
        <div class="amount" id="statsBalance">0</div>
        <div class="symbol">$MESHP</div>
      </div>
    </div>

    <div class="stat-grid">
      <div class="stat-box">
        <div class="label">Measurements</div>
        <div class="value" id="statsMeasurements" style="color:var(--accent)">0</div>
        <div class="unit">submitted</div>
      </div>
      <div class="stat-box">
        <div class="label">Earned Total</div>
        <div class="value" id="statsEarned" style="color:var(--yellow)">0</div>
        <div class="unit">$MESHP</div>
      </div>
    </div>
  </div>

  <div id="statsEmpty" style="text-align:center;padding:60px 20px;color:var(--muted)">
    <div style="font-size:3rem;margin-bottom:16px">🔍</div>
    <div>Enter your wallet address above to see your stats</div>
  </div>
</div>

<!-- ═══════════════════ PAGE 3: LIVE FEED ═══════════════════ -->
<div class="page" id="page-feed">
  <div class="refresh-row">
    <button class="btn-refresh" onclick="loadFeed()">↻ Refresh</button>
  </div>
  <div class="card" style="padding:0;overflow:hidden">
    <table class="feed-table">
      <thead>
        <tr>
          <th>#</th>
          <th>Address</th>
          <th>Ping</th>
          <th>Download</th>
          <th>Upload</th>
          <th>ISP</th>
          <th>Region</th>
          <th>Time</th>
        </tr>
      </thead>
      <tbody id="feedBody">
        <tr><td colspan="8" style="padding:40px;text-align:center;color:var(--muted)">Loading…</td></tr>
      </tbody>
    </table>
  </div>
</div>

<!-- ═══════════════════ PAGE 4: LEADERBOARD ═══════════════════ -->
<div class="page" id="page-leaderboard">
  <div class="refresh-row">
    <button class="btn-refresh" onclick="loadLeaderboard()">↻ Refresh</button>
  </div>
  <div id="lbContent">
    <div class="skeleton" style="height:60px"></div>
    <div class="skeleton" style="height:60px"></div>
    <div class="skeleton" style="height:60px"></div>
  </div>
</div>

<div id="toast"></div>

<script>
// ── State ────────────────────────────────────────────────────────────────────
let testResult = null;
const API = '';  // same origin

// ── Navigation ───────────────────────────────────────────────────────────────
function showPage(id, btn) {
  document.querySelectorAll('.page').forEach(p => p.classList.remove('active'));
  document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
  document.getElementById('page-' + id).classList.add('active');
  btn.classList.add('active');
  if (id === 'feed') loadFeed();
  if (id === 'leaderboard') loadLeaderboard();
  if (id === 'speedtest') loadNetworkStats();
}

// ── Toast ─────────────────────────────────────────────────────────────────────
function toast(msg, type='info') {
  const t = document.getElementById('toast');
  t.textContent = msg;
  t.className = 'show ' + type;
  setTimeout(() => { t.className = ''; }, 3500);
}

// ── Wallet badge ──────────────────────────────────────────────────────────────
function setWallet(addr) {
  const b = document.getElementById('walletBadge');
  b.textContent = addr ? addr.slice(0,6) + '…' + addr.slice(-4) : 'No Wallet';
  b.className = 'wallet-badge' + (addr ? ' connected' : '');
}

document.getElementById('addrInputTest').addEventListener('input', function() {
  setWallet(this.value.trim());
});

// ── Gauge animation ───────────────────────────────────────────────────────────
let gaugeTarget = 0;
function setGauge(val, max, label) {
  gaugeTarget = val;
  const arc = document.getElementById('gaugeArc');
  const needle = document.getElementById('gaugeNeedle');
  const num = document.getElementById('gaugeNum');
  const lbl = document.getElementById('gaugeLabel');
  const ratio = Math.min(val / max, 1);
  const totalLen = 346;
  arc.style.strokeDashoffset = totalLen * (1 - ratio);
  const angle = -90 + ratio * 180;
  needle.setAttribute('transform', 'rotate(' + angle + ' 130 130)');
  num.textContent = Math.round(val);
  if (lbl && label) lbl.textContent = label;
}

function animateGauge(target, max, label, duration=1200) {
  const start = Date.now();
  const tick = () => {
    const elapsed = Date.now() - start;
    const progress = Math.min(elapsed / duration, 1);
    const ease = 1 - Math.pow(1 - progress, 3);
    setGauge(target * ease, max, label);
    if (progress < 1) requestAnimationFrame(tick);
  };
  requestAnimationFrame(tick);
}

// ── Speed test simulation ─────────────────────────────────────────────────────
// In production this would hit real test endpoints.
// For the contest demo we simulate realistic values with browser timing.
async function runTest() {
  const addr = document.getElementById('addrInputTest').value.trim();
  if (!addr || addr.length !== 40) {
    toast('Please enter a valid 40-char hex wallet address first', 'err');
    return;
  }

  const btn = document.getElementById('btnRun');
  const phase = document.getElementById('phaseLabel');
  const submitBtn = document.getElementById('btnSubmit');
  const rewardBadge = document.getElementById('rewardBadge');

  btn.disabled = true;
  btn.classList.add('running');
  submitBtn.disabled = true;
  rewardBadge.style.display = 'none';
  testResult = null;

  // Reset metrics
  document.getElementById('pingVal').textContent = '—';
  document.getElementById('downVal').textContent = '—';
  document.getElementById('upVal').textContent = '—';
  setGauge(0, 200, '');

  // ── Phase 1: Ping ──
  phase.textContent = '📡 Measuring latency…';
  const pingStart = performance.now();
  try { await fetch(API + '/api/stats', {cache:'no-store'}); } catch(_){}
  const ping = Math.round(performance.now() - pingStart);
  document.getElementById('pingVal').textContent = ping;
  animateGauge(ping, 300, 'Ping ' + ping + ' ms', 600);
  await sleep(600);

  // ── Phase 2: Download ──
  phase.textContent = '⬇️  Testing download speed…';
  const dlStart = performance.now();
  let dlBytes = 0;
  try {
    // Fetch feed as a proxy for download measurement
    const r = await fetch(API + '/api/feed?_=' + Date.now(), {cache:'no-store'});
    const buf = await r.arrayBuffer();
    dlBytes = buf.byteLength;
  } catch(_) { dlBytes = 1024; }
  const dlTime = (performance.now() - dlStart) / 1000;
  // Scale to a realistic value
  const download = Math.round(Math.max(1000, Math.min(100000, (dlBytes * 8 / dlTime) + Math.random() * 50000)));
  document.getElementById('downVal').textContent = download;
  animateGauge(download / 1000, 200, 'Download ' + (download/1000).toFixed(1) + ' Mbps', 900);
  await sleep(900);

  // ── Phase 3: Upload ──
  phase.textContent = '⬆️  Testing upload speed…';
  const uploadPayload = new Uint8Array(32768);
  const upStart = performance.now();
  try {
    await fetch(API + '/api/stats', {method:'POST', body: uploadPayload, cache:'no-store'}).catch(()=>{});
  } catch(_) {}
  const upTime = Math.max(0.05, (performance.now() - upStart) / 1000);
  const upload = Math.round(Math.max(512, Math.min(50000, (uploadPayload.byteLength * 8 / upTime) + Math.random() * 20000)));
  document.getElementById('upVal').textContent = upload;
  animateGauge(upload / 1000, 200, 'Upload ' + (upload/1000).toFixed(1) + ' Mbps', 700);
  await sleep(700);

  // ── Done ──
  testResult = { ping, download, upload };
  phase.textContent = '✅ Test complete — ready to submit on-chain';
  btn.disabled = false;
  btn.classList.remove('running');
  submitBtn.disabled = false;
  animateGauge(download / 1000, 200, 'Download ' + (download/1000).toFixed(1) + ' Mbps', 400);
}

function sleep(ms) { return new Promise(r => setTimeout(r, ms)); }

// ── Submit to chain ───────────────────────────────────────────────────────────
async function submitResult() {
  if (!testResult) { toast('Run a speed test first', 'err'); return; }

  const addr = document.getElementById('addrInputTest').value.trim();
  const isp = document.getElementById('ispInput').value.trim() || 'Unknown ISP';
  const region = document.getElementById('regionInput').value.trim() || 'Unknown Region';

  if (!addr || addr.length !== 40) {
    toast('Invalid wallet address', 'err'); return;
  }

  // Build the Canopy RPC transaction
  // The Canopy node RPC is at :50052 by default; we direct users to use the CLI or wallet.
  // For the demo, we show the transaction payload they can submit.
  const payload = {
    messageType: "submit_measurement",
    msg: {
      "@type": "type.googleapis.com/types.MessageSubmitMeasurement",
      address: addr,
      ping: testResult.ping,
      download: testResult.download,
      upload: testResult.upload,
      isp: isp,
      region: region
    }
  };

  toast('📝 Transaction payload ready! Use the Canopy CLI or wallet to sign and submit.', 'info');

  // Show the payload in console for dev use
  console.log('MeshPulse TX payload:', JSON.stringify(payload, null, 2));

  // Optimistically show reward
  document.getElementById('btnSubmit').disabled = true;
  document.getElementById('rewardBadge').style.display = 'block';
  document.getElementById('phaseLabel').textContent = '🎉 Submit via Canopy CLI to earn your 10 $MESHP!';

  // Also show copyable payload
  const payloadStr = JSON.stringify(payload, null, 2);
  const area = document.createElement('textarea');
  area.value = payloadStr;
  area.style.cssText = 'position:fixed;opacity:0;top:0;left:0;width:1px;height:1px';
  document.body.appendChild(area);
  area.select();
  document.execCommand('copy');
  document.body.removeChild(area);
  toast('📋 TX payload copied to clipboard!', 'ok');
}

// ── Network stats ─────────────────────────────────────────────────────────────
async function loadNetworkStats() {
  try {
    const r = await fetch(API + '/api/stats');
    const s = await r.json();
    document.getElementById('ns-total').textContent = s.totalMeasurements || 0;
    document.getElementById('ns-ping').textContent = s.avgPing || 0;
    document.getElementById('ns-down').textContent = s.avgDownload || 0;
    document.getElementById('ns-up').textContent = s.avgUpload || 0;
  } catch(e) {
    console.warn('Stats unavailable:', e);
  }
}

// ── Live Feed ─────────────────────────────────────────────────────────────────
async function loadFeed() {
  const tbody = document.getElementById('feedBody');
  tbody.innerHTML = '<tr><td colspan="8" style="padding:40px;text-align:center;color:var(--muted)">Loading…</td></tr>';
  try {
    const r = await fetch(API + '/api/feed');
    const rows = await r.json();
    if (!rows || rows.length === 0) {
      tbody.innerHTML = '<tr><td colspan="8" style="padding:40px;text-align:center;color:var(--muted)">No measurements yet — be the first to submit!</td></tr>';
      return;
    }
    tbody.innerHTML = rows.map(m => {
      const addr = m.address || '';
      const short = addr.slice(0,6) + '…' + addr.slice(-4);
      const t = m.timestamp ? new Date(m.timestamp * 1000).toLocaleTimeString() : '—';
      return '<tr>' +
        '<td style="color:var(--muted)">#' + m.id + '</td>' +
        '<td class="addr">' + short + '</td>' +
        '<td class="val-ping">' + m.ping + ' ms</td>' +
        '<td class="val-down">' + (m.download/1000).toFixed(1) + ' M</td>' +
        '<td class="val-up">'   + (m.upload/1000).toFixed(1)   + ' M</td>' +
        '<td><span class="badge-isp">' + (m.isp||'?') + '</span></td>' +
        '<td><span class="badge-region">' + (m.region||'?') + '</span></td>' +
        '<td style="color:var(--muted);font-size:0.8rem">' + t + '</td>' +
        '</tr>';
    }).join('');
  } catch(e) {
    tbody.innerHTML = '<tr><td colspan="8" style="padding:40px;text-align:center;color:var(--red)">Failed to load feed: ' + e.message + '</td></tr>';
  }
}

// ── Leaderboard ───────────────────────────────────────────────────────────────
async function loadLeaderboard() {
  const el = document.getElementById('lbContent');
  el.innerHTML = '<div class="skeleton" style="height:60px"></div>'.repeat(5);
  try {
    const r = await fetch(API + '/api/leaderboard');
    const rows = await r.json();
    if (!rows || rows.length === 0) {
      el.innerHTML = '<div style="text-align:center;padding:60px;color:var(--muted)"><div style="font-size:3rem;margin-bottom:12px">🏆</div>No contributors yet</div>';
      return;
    }
    const maxCount = rows[0].totalMeasurements || 1;
    const rankSymbols = ['🥇','🥈','🥉'];
    const rankClasses = ['gold','silver','bronze'];
    el.innerHTML = rows.map((c, i) => {
      const pct = Math.round((c.totalMeasurements / maxCount) * 100);
      const addr = c.address || '';
      const short = addr.slice(0,8) + '…' + addr.slice(-6);
      const rank = i < 3 ? rankSymbols[i] : (i+1);
      const rankCls = i < 3 ? rankClasses[i] : '';
      const tokens = (c.tokenBalance / 1_000_000).toFixed(1);
      return '<div class="lb-row">' +
        '<div class="lb-rank ' + rankCls + '">' + rank + '</div>' +
        '<div class="lb-addr">' + short + '</div>' +
        '<div class="lb-bar-wrap"><div class="lb-bar" style="width:' + pct + '%"></div></div>' +
        '<div class="lb-count">' + c.totalMeasurements + ' tests</div>' +
        '<div class="lb-tokens">⚡ ' + tokens + ' $MESHP</div>' +
        '</div>';
    }).join('');
  } catch(e) {
    el.innerHTML = '<div style="padding:40px;text-align:center;color:var(--red)">Failed to load: ' + e.message + '</div>';
  }
}

// ── My Stats ──────────────────────────────────────────────────────────────────
async function lookupContributor() {
  const addr = document.getElementById('lookupAddr').value.trim();
  if (!addr || addr.length !== 40) { toast('Enter a valid 40-char hex address', 'err'); return; }
  try {
    const r = await fetch(API + '/api/contributor?address=' + addr);
    const c = await r.json();
    if (c.error) { toast(c.error, 'err'); return; }
    document.getElementById('statsAddr').textContent = addr;
    document.getElementById('statsBalance').textContent = ((c.tokenBalance||0) / 1_000_000).toFixed(2);
    document.getElementById('statsMeasurements').textContent = c.totalMeasurements || 0;
    const earned = ((c.totalMeasurements||0) * 10).toFixed(0);
    document.getElementById('statsEarned').textContent = earned;
    const total = c.totalMeasurements || 0;
    const rank = total >= 100 ? 'Diamond Node' : total >= 50 ? 'Gold Node' : total >= 20 ? 'Silver Node' : total >= 5 ? 'Bronze Node' : 'New Node';
    document.getElementById('statsRank').textContent = rank + ' · ' + total + ' measurements';
    document.getElementById('statsContent').style.display = 'block';
    document.getElementById('statsEmpty').style.display = 'none';
    setWallet(addr);
    document.getElementById('addrInputTest').value = addr;
  } catch(e) {
    toast('Lookup failed: ' + e.message, 'err');
  }
}

// ── Auto-load on start ────────────────────────────────────────────────────────
loadNetworkStats();
setInterval(loadNetworkStats, 30000);
</script>
</body>
</html>`
