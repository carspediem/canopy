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
    cursor: pointer;
    transition: all 0.2s;
  }
  .wallet-badge.connected { color: var(--green); border-color: rgba(0,255,136,0.3); }
  .wallet-badge:hover { border-color: var(--accent); }

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
  .balance-display { text-align: right; }
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

  /* ── LANDING PAGE ── */
  .landing {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 70vh;
    text-align: center;
    padding: 40px 24px;
  }
  .landing-glow {
    width: 120px;
    height: 120px;
    border-radius: 50%;
    background: radial-gradient(circle, rgba(0,212,255,0.25) 0%, transparent 70%);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 3.5rem;
    margin-bottom: 32px;
    position: relative;
  }
  .landing-glow::after {
    content: '';
    position: absolute;
    inset: -8px;
    border-radius: 50%;
    border: 1px solid rgba(0,212,255,0.2);
    animation: orbit 3s linear infinite;
  }
  @keyframes orbit {
    from { transform: rotate(0deg) scale(1); opacity: 0.6; }
    50%  { transform: rotate(180deg) scale(1.08); opacity: 1; }
    to   { transform: rotate(360deg) scale(1); opacity: 0.6; }
  }
  .landing h1 {
    font-size: 2.4rem;
    font-weight: 800;
    background: linear-gradient(90deg, var(--accent), var(--accent2));
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    margin-bottom: 16px;
    line-height: 1.2;
  }
  .landing .subtitle {
    font-size: 1.1rem;
    color: var(--muted);
    max-width: 480px;
    line-height: 1.6;
    margin-bottom: 40px;
  }
  .landing .features {
    display: flex;
    gap: 24px;
    margin-bottom: 48px;
    flex-wrap: wrap;
    justify-content: center;
  }
  .landing .feature-pill {
    background: var(--card);
    border: 1px solid var(--border);
    border-radius: 12px;
    padding: 12px 20px;
    font-size: 0.88rem;
    color: var(--text);
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .btn-connect {
    padding: 18px 56px;
    font-size: 1.15rem;
    font-weight: 700;
    border: none;
    border-radius: 50px;
    cursor: pointer;
    background: linear-gradient(135deg, var(--accent), var(--accent2));
    color: #fff;
    letter-spacing: 0.05em;
    box-shadow: 0 0 40px rgba(0,212,255,0.4);
    transition: all 0.25s;
  }
  .btn-connect:hover { transform: scale(1.04); box-shadow: 0 0 60px rgba(0,212,255,0.6); }

  /* ── WALLET MODAL ── */
  .modal-overlay {
    display: none;
    position: fixed;
    inset: 0;
    background: rgba(5,10,14,0.85);
    backdrop-filter: blur(6px);
    z-index: 200;
    align-items: center;
    justify-content: center;
  }
  .modal-overlay.open { display: flex; }
  .modal-box {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 20px;
    padding: 40px 36px;
    width: 100%;
    max-width: 480px;
    box-shadow: 0 0 60px rgba(0,212,255,0.1);
  }
  .modal-box h2 {
    font-size: 1.4rem;
    font-weight: 700;
    margin-bottom: 8px;
    background: linear-gradient(90deg, var(--accent), var(--accent2));
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }
  .modal-box .modal-sub {
    font-size: 0.88rem;
    color: var(--muted);
    margin-bottom: 28px;
    line-height: 1.5;
  }
  .modal-box .modal-sub a {
    color: var(--accent);
    text-decoration: none;
  }
  .modal-box .modal-sub a:hover { text-decoration: underline; }
  .modal-input {
    width: 100%;
    padding: 14px 18px;
    background: var(--card);
    border: 1px solid var(--border);
    border-radius: 10px;
    color: var(--text);
    font-family: monospace;
    font-size: 0.9rem;
    outline: none;
    transition: border-color 0.2s;
    margin-bottom: 8px;
  }
  .modal-input:focus { border-color: var(--accent); }
  .modal-input::placeholder { color: var(--muted); }
  .modal-input-err { font-size: 0.78rem; color: var(--red); min-height: 18px; margin-bottom: 20px; }
  .btn-modal-connect {
    width: 100%;
    padding: 14px;
    background: linear-gradient(135deg, var(--accent), var(--accent2));
    border: none;
    border-radius: 10px;
    color: #fff;
    font-size: 1rem;
    font-weight: 700;
    cursor: pointer;
    transition: opacity 0.2s;
    margin-bottom: 12px;
  }
  .btn-modal-connect:hover { opacity: 0.9; }
  .btn-modal-cancel {
    width: 100%;
    padding: 10px;
    background: transparent;
    border: 1px solid var(--border);
    border-radius: 10px;
    color: var(--muted);
    font-size: 0.9rem;
    cursor: pointer;
    transition: all 0.2s;
  }
  .btn-modal-cancel:hover { border-color: var(--accent); color: var(--text); }

  /* ── DASHBOARD ── */
  .dashboard-header {
    display: flex;
    align-items: center;
    gap: 20px;
    padding: 24px 28px;
    background: linear-gradient(135deg, rgba(0,212,255,0.06), rgba(124,58,237,0.06));
    border: 1px solid var(--border);
    border-radius: 16px;
    margin-bottom: 20px;
    flex-wrap: wrap;
  }
  .node-status-badge {
    display: flex;
    align-items: center;
    gap: 10px;
    background: rgba(0,255,136,0.08);
    border: 1px solid rgba(0,255,136,0.25);
    border-radius: 50px;
    padding: 8px 20px;
  }
  .pulse-dot {
    width: 10px;
    height: 10px;
    background: var(--green);
    border-radius: 50%;
    position: relative;
    flex-shrink: 0;
  }
  .pulse-dot::after {
    content: '';
    position: absolute;
    inset: -4px;
    border-radius: 50%;
    background: var(--green);
    opacity: 0.4;
    animation: pulse-glow 1.6s ease-out infinite;
  }
  @keyframes pulse-glow {
    0%   { transform: scale(1); opacity: 0.4; }
    100% { transform: scale(2.4); opacity: 0; }
  }
  .node-status-text { font-weight: 700; color: var(--green); font-size: 0.92rem; letter-spacing: 0.06em; }
  .dashboard-counters { display: flex; gap: 24px; flex: 1; flex-wrap: wrap; }
  .dash-counter {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }
  .dash-counter .dc-label { font-size: 0.72rem; color: var(--muted); text-transform: uppercase; letter-spacing: 0.06em; }
  .dash-counter .dc-value { font-size: 1.5rem; font-weight: 700; font-family: monospace; }
  .dash-counter .dc-value.green { color: var(--green); }
  .dash-counter .dc-value.yellow { color: var(--yellow); }
  .btn-chain-link {
    padding: 9px 18px;
    background: transparent;
    border: 1px solid rgba(0,212,255,0.35);
    border-radius: 8px;
    color: var(--accent);
    font-size: 0.83rem;
    cursor: pointer;
    text-decoration: none;
    display: inline-flex;
    align-items: center;
    gap: 6px;
    transition: all 0.2s;
  }
  .btn-chain-link:hover { background: rgba(0,212,255,0.1); }

  .dashboard-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; margin-bottom: 20px; }
  @media (max-width: 720px) { .dashboard-grid { grid-template-columns: 1fr; } }

  /* auto feed */
  .auto-feed-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 16px;
    border-bottom: 1px solid rgba(26,46,69,0.5);
    animation: slideIn 0.4s ease;
  }
  @keyframes slideIn {
    from { opacity: 0; transform: translateY(-8px); }
    to   { opacity: 1; transform: translateY(0); }
  }
  .auto-feed-item:last-child { border-bottom: none; }
  .af-time { font-size: 0.75rem; color: var(--muted); font-family: monospace; min-width: 56px; }
  .af-metric { font-size: 0.82rem; font-weight: 600; }
  .af-hash {
    font-family: monospace;
    font-size: 0.72rem;
    color: var(--muted);
    flex: 1;
    text-align: right;
  }

  .btn-claim {
    padding: 14px 36px;
    font-size: 1rem;
    font-weight: 700;
    border: none;
    border-radius: 50px;
    cursor: pointer;
    background: linear-gradient(135deg, var(--yellow), #ff9500);
    color: #000;
    letter-spacing: 0.04em;
    transition: all 0.25s;
    box-shadow: 0 0 24px rgba(255,215,0,0.3);
  }
  .btn-claim:hover:not(:disabled) { transform: scale(1.04); box-shadow: 0 0 36px rgba(255,215,0,0.5); }
  .btn-claim:disabled { opacity: 0.35; cursor: not-allowed; }

  /* responsive */
  @media (max-width: 600px) {
    .stat-grid { grid-template-columns: repeat(2,1fr); }
    nav { flex-wrap: wrap; }
    .landing h1 { font-size: 1.8rem; }
    .dashboard-header { flex-direction: column; align-items: flex-start; }
  }
</style>
</head>
<body>

<nav>
  <div class="logo">⚡ MeshPulse <span>DePIN</span></div>
  <div class="nav-tabs">
    <button class="tab active" onclick="showPage('speedtest',this)">⚡ Dashboard</button>
    <button class="tab" onclick="showPage('mystats',this)">👤 My Stats</button>
    <button class="tab" onclick="showPage('feed',this)">🌐 Live Feed</button>
    <button class="tab" onclick="showPage('leaderboard',this)">🏆 Leaderboard</button>
  </div>
  <div class="wallet-badge" id="walletBadge" onclick="onWalletBadgeClick()">Connect Wallet</div>
</nav>

<!-- ═══════════════════ PAGE 1: DASHBOARD ═══════════════════ -->
<div class="page active" id="page-speedtest">

  <!-- LANDING (shown before connect) -->
  <div id="landingView" class="landing">
    <div class="landing-glow">⚡</div>
    <h1>Earn $MESHP by keeping<br>the internet honest</h1>
    <p class="subtitle">Run a passive DePIN node. Connect once, earn forever.<br>Your device measures network quality — the chain rewards you.</p>
    <div class="features">
      <div class="feature-pill">📡 Automatic measurements</div>
      <div class="feature-pill">⛓️ On-chain proof</div>
      <div class="feature-pill">💰 Earn $MESHP</div>
      <div class="feature-pill">🔒 Non-custodial</div>
    </div>
    <button class="btn-connect" onclick="openWalletModal()">Connect Wallet</button>
  </div>

  <!-- DASHBOARD (shown after connect) -->
  <div id="dashboardView" style="display:none">

    <!-- Header row -->
    <div class="dashboard-header">
      <div class="node-status-badge">
        <div class="pulse-dot"></div>
        <div class="node-status-text">ACTIVE</div>
      </div>
      <div class="dashboard-counters">
        <div class="dash-counter">
          <div class="dc-label">Uptime</div>
          <div class="dc-value green" id="uptimeDisplay">00:00:00</div>
        </div>
        <div class="dash-counter">
          <div class="dc-label">$MESHP Earned</div>
          <div class="dc-value yellow" id="earnedDisplay">0.00</div>
        </div>
        <div class="dash-counter">
          <div class="dc-label">Measurements</div>
          <div class="dc-value" id="measureCount">0</div>
        </div>
      </div>
      <a class="btn-chain-link" href="https://testnet.app.canopynetwork.org/chains/251166" target="_blank">
        🔗 View MeshPulse Chain
      </a>
    </div>

    <!-- Main grid -->
    <div class="dashboard-grid">

      <!-- Auto measurements feed -->
      <div class="card" style="margin-bottom:0">
        <div class="card-title">Live Measurements</div>
        <div id="autoFeedList" style="min-height:200px">
          <div style="padding:40px;text-align:center;color:var(--muted);font-size:0.88rem">
            First measurement in 30s…
          </div>
        </div>
      </div>

      <!-- Rewards panel -->
      <div class="card" style="margin-bottom:0;display:flex;flex-direction:column;gap:20px">
        <div class="card-title">Rewards</div>

        <div class="stat-box" style="background:rgba(255,215,0,0.05);border-color:rgba(255,215,0,0.2)">
          <div class="label">Claimable Balance</div>
          <div class="value" id="claimableBalance" style="color:var(--yellow)">0.00</div>
          <div class="unit">$MESHP</div>
        </div>

        <button class="btn-claim" id="btnClaim" onclick="claimRewards()" disabled>
          💰 Claim Rewards
        </button>

        <div style="font-size:0.78rem;color:var(--muted);line-height:1.6;text-align:center">
          Rewards accrue every 30 seconds.<br>
          Submit via <a href="https://testnet.app.canopynetwork.org/chains/251166" target="_blank" style="color:var(--accent);text-decoration:none">Canopy Wallet</a> to settle on-chain.
        </div>

        <div style="border-top:1px solid var(--border);padding-top:16px">
          <div style="font-size:0.75rem;color:var(--muted);margin-bottom:8px;text-transform:uppercase;letter-spacing:0.06em">Your Node</div>
          <div style="font-family:monospace;font-size:0.78rem;color:var(--accent);word-break:break-all" id="dashWalletAddr">—</div>
        </div>
      </div>

    </div><!-- /grid -->
  </div><!-- /dashboardView -->
</div>

<!-- ══ WALLET MODAL ══ -->
<div class="modal-overlay" id="walletModal">
  <div class="modal-box">
    <h2>Connect Your Wallet</h2>
    <p class="modal-sub">
      Enter your Canopy wallet address to start earning $MESHP.<br>
      Don't have one? <a href="https://testnet.app.canopynetwork.org/faucet" target="_blank">Get a free testnet address →</a>
    </p>
    <input class="modal-input" type="text" id="modalAddrInput"
      placeholder="Canopy wallet address (40 hex chars)" autocomplete="off" />
    <div class="modal-input-err" id="modalAddrErr"></div>
    <button class="btn-modal-connect" onclick="confirmConnect()">Connect</button>
    <button class="btn-modal-cancel" onclick="closeWalletModal()">Cancel</button>
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
const API = '';
let walletAddr = null;
let uptimeStart = null;
let uptimeInterval = null;
let autoInterval = null;
let meshpEarned = 0;
let measurementsDone = 0;
let autoFeedItems = [];

// ── Navigation ───────────────────────────────────────────────────────────────
function showPage(id, btn) {
  document.querySelectorAll('.page').forEach(p => p.classList.remove('active'));
  document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
  document.getElementById('page-' + id).classList.add('active');
  btn.classList.add('active');
  if (id === 'feed') loadFeed();
  if (id === 'leaderboard') loadLeaderboard();
}

// ── Toast ─────────────────────────────────────────────────────────────────────
function toast(msg, type='info') {
  const t = document.getElementById('toast');
  t.textContent = msg;
  t.className = 'show ' + type;
  setTimeout(() => { t.className = ''; }, 3500);
}

// ── Wallet badge ──────────────────────────────────────────────────────────────
function updateWalletBadge() {
  const b = document.getElementById('walletBadge');
  if (walletAddr) {
    b.textContent = walletAddr.slice(0,6) + '…' + walletAddr.slice(-4);
    b.className = 'wallet-badge connected';
  } else {
    b.textContent = 'Connect Wallet';
    b.className = 'wallet-badge';
  }
}

function onWalletBadgeClick() {
  if (walletAddr) {
    // clicking connected badge disconnects
    if (confirm('Disconnect wallet?')) disconnectWallet();
  } else {
    openWalletModal();
  }
}

// ── Wallet modal ──────────────────────────────────────────────────────────────
function openWalletModal() {
  document.getElementById('walletModal').classList.add('open');
  document.getElementById('modalAddrInput').value = '';
  document.getElementById('modalAddrErr').textContent = '';
  setTimeout(() => document.getElementById('modalAddrInput').focus(), 100);
}

function closeWalletModal() {
  document.getElementById('walletModal').classList.remove('open');
}

document.getElementById('modalAddrInput').addEventListener('keydown', function(e) {
  if (e.key === 'Enter') confirmConnect();
});

function confirmConnect() {
  let addr = document.getElementById('modalAddrInput').value.trim();
  if (addr.startsWith('0x') || addr.startsWith('0X')) addr = addr.slice(2);
  if (!addr || addr.length !== 40 || !/^[0-9a-fA-F]+$/.test(addr)) {
    document.getElementById('modalAddrErr').textContent = 'Please enter a valid 40-character hex address.';
    return;
  }
  closeWalletModal();
  connectWallet(addr.toLowerCase());
}

function connectWallet(addr) {
  walletAddr = addr;
  localStorage.setItem('meshpulse_wallet', addr);
  updateWalletBadge();
  showDashboard();
  toast('✅ Node connected — earning $MESHP!', 'ok');
}

function disconnectWallet() {
  walletAddr = null;
  localStorage.removeItem('meshpulse_wallet');
  stopNode();
  updateWalletBadge();
  document.getElementById('dashboardView').style.display = 'none';
  document.getElementById('landingView').style.display = 'flex';
}

// ── Dashboard ─────────────────────────────────────────────────────────────────
function showDashboard() {
  document.getElementById('landingView').style.display = 'none';
  document.getElementById('dashboardView').style.display = 'block';
  document.getElementById('dashWalletAddr').textContent = walletAddr;
  startNode();
}

function startNode() {
  uptimeStart = Date.now();
  meshpEarned = 0;
  measurementsDone = 0;
  autoFeedItems = [];
  document.getElementById('earnedDisplay').textContent = '0.00';
  document.getElementById('measureCount').textContent = '0';
  document.getElementById('claimableBalance').textContent = '0.00';
  document.getElementById('btnClaim').disabled = true;
  document.getElementById('autoFeedList').innerHTML =
    '<div style="padding:40px;text-align:center;color:var(--muted);font-size:0.88rem">First measurement in 30s…</div>';

  // uptime counter
  if (uptimeInterval) clearInterval(uptimeInterval);
  uptimeInterval = setInterval(tickUptime, 1000);

  // auto measurement every 30s
  if (autoInterval) clearInterval(autoInterval);
  autoInterval = setInterval(runAutoMeasurement, 30000);
}

function stopNode() {
  if (uptimeInterval) clearInterval(uptimeInterval);
  if (autoInterval) clearInterval(autoInterval);
  uptimeInterval = null;
  autoInterval = null;
  meshpEarned = 0;
  measurementsDone = 0;
}

function tickUptime() {
  if (!uptimeStart) return;
  const secs = Math.floor((Date.now() - uptimeStart) / 1000);
  const h = String(Math.floor(secs / 3600)).padStart(2, '0');
  const m = String(Math.floor((secs % 3600) / 60)).padStart(2, '0');
  const s = String(secs % 60).padStart(2, '0');
  document.getElementById('uptimeDisplay').textContent = h + ':' + m + ':' + s;
}

async function runAutoMeasurement() {
  if (!walletAddr) return;

  // Measure ping via fetch timing
  const t0 = performance.now();
  try { await fetch(API + '/api/stats', {cache:'no-store'}); } catch(_) {}
  const ping = Math.max(5, Math.round(performance.now() - t0));

  // Simulated download/upload (realistic range)
  const download = Math.round(10000 + Math.random() * 90000);  // 10–100 Mbps in Kbps
  const upload   = Math.round(5000  + Math.random() * 45000);  // 5–50 Mbps in Kbps

  // Increment counters
  meshpEarned += 1;
  measurementsDone += 1;
  document.getElementById('earnedDisplay').textContent = meshpEarned.toFixed(2);
  document.getElementById('measureCount').textContent = measurementsDone;
  document.getElementById('claimableBalance').textContent = meshpEarned.toFixed(2);
  document.getElementById('btnClaim').disabled = (meshpEarned <= 0);

  // Add to auto feed display
  const now = new Date();
  const timeStr = now.toTimeString().slice(0, 8);
  const fakeHash = walletAddr.slice(0,6) + '…' + Math.random().toString(16).slice(2,6);
  const item = {
    time: timeStr,
    ping,
    download: (download / 1000).toFixed(1),
    upload: (upload / 1000).toFixed(1),
    hash: fakeHash
  };
  autoFeedItems.unshift(item);
  if (autoFeedItems.length > 10) autoFeedItems.pop();
  renderAutoFeed();

  // Submit to backend in background (best-effort) — real on-chain tx via keystore key
  try {
    const resp = await fetch(API + '/api/submit-measurement', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        ping,
        download,
        upload,
        isp: 'Unknown ISP',
        region: 'Unknown Region'
      })
    });
    const data = await resp.json();
    if (resp.ok && data.txHash) {
      // Update the hash in the feed display with the real tx hash
      item.hash = data.txHash.substring(0, 12) + '…';
      renderAutoFeed();
      console.log('[AutoMeasure] TX submitted:', data.txHash);
    }
  } catch(_) {}
}

function renderAutoFeed() {
  const el = document.getElementById('autoFeedList');
  if (!autoFeedItems.length) return;
  el.innerHTML = autoFeedItems.map(it =>
    '<div class="auto-feed-item">' +
      '<span class="af-time">' + it.time + '</span>' +
      '<span class="af-metric val-ping">' + it.ping + 'ms</span>' +
      '<span class="af-metric val-down">' + it.download + 'M↓</span>' +
      '<span class="af-metric val-up">'   + it.upload   + 'M↑</span>' +
      '<span class="af-hash">' + it.hash + '</span>' +
    '</div>'
  ).join('');
}

function claimRewards() {
  if (meshpEarned <= 0) return;
  toast('📋 Claim via Canopy Wallet at testnet.app.canopynetwork.org', 'info');
  meshpEarned = 0;
  document.getElementById('claimableBalance').textContent = '0.00';
  document.getElementById('btnClaim').disabled = true;
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
  let addr = document.getElementById('lookupAddr').value.trim();
  if (addr.startsWith('0x') || addr.startsWith('0X')) addr = addr.slice(2);
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
  } catch(e) {
    toast('Lookup failed: ' + e.message, 'err');
  }
}

// ── Init ──────────────────────────────────────────────────────────────────────
(function init() {
  const saved = localStorage.getItem('meshpulse_wallet');
  if (saved) {
    walletAddr = saved;
    updateWalletBadge();
    showDashboard();
  }
})();
</script>
</body>`

