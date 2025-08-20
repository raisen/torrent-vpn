package assets

import (
	"fyne.io/fyne/v2"
)

// IconResource is the embedded application icon
var IconResource = &fyne.StaticResource{
	StaticName: "icon.svg",
	StaticContent: []byte(`<?xml version="1.0" encoding="UTF-8"?>
<svg width="256" height="256" viewBox="0 0 256 256" xmlns="http://www.w3.org/2000/svg">
  <defs>
    <radialGradient id="bgGradient" cx="50%" cy="50%" r="50%">
      <stop offset="0%" style="stop-color:#1e293b;stop-opacity:1" />
      <stop offset="100%" style="stop-color:#0f172a;stop-opacity:1" />
    </radialGradient>

    <linearGradient id="vpnGradient" x1="0%" y1="0%" x2="100%" y2="0%">
      <stop offset="0%" style="stop-color:#3b82f6;stop-opacity:1" />
      <stop offset="50%" style="stop-color:#1d4ed8;stop-opacity:1" />
      <stop offset="100%" style="stop-color:#1e40af;stop-opacity:1" />
    </linearGradient>

    <linearGradient id="torrentGradient" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" style="stop-color:#10b981;stop-opacity:1" />
      <stop offset="100%" style="stop-color:#047857;stop-opacity:1" />
    </linearGradient>

    <filter id="glow">
      <feGaussianBlur stdDeviation="3" result="coloredBlur"/>
      <feMerge>
        <feMergeNode in="coloredBlur"/>
        <feMergeNode in="SourceGraphic"/>
      </feMerge>
    </filter>
  </defs>

  <circle cx="128" cy="128" r="120" fill="url(#bgGradient)" stroke="#334155" stroke-width="4"/>

  <ellipse cx="128" cy="128" rx="90" ry="25" fill="url(#vpnGradient)" opacity="0.8"/>
  <ellipse cx="128" cy="128" rx="90" ry="15" fill="#60a5fa" opacity="0.6"/>

  <circle cx="70" cy="128" r="20" fill="#1e40af" stroke="#3b82f6" stroke-width="2"/>
  <rect x="63" y="125" width="14" height="10" rx="2" fill="white"/>
  <path d="M61 125 C61 121 64 118 70 118 C76 118 79 121 79 125"
        stroke="white" stroke-width="2" fill="none"/>

  <circle cx="186" cy="128" r="20" fill="#047857" stroke="#10b981" stroke-width="2"/>

  <g fill="white">
    <rect x="178" y="120" width="6" height="6" rx="1"/>
    <rect x="186" y="120" width="6" height="6" rx="1"/>
    <rect x="178" y="128" width="6" height="6" rx="1"/>
    <rect x="186" y="128" width="6" height="6" rx="1"/>
    <rect x="178" y="136" width="6" height="6" rx="1"/>
    <rect x="186" y="136" width="6" height="6" rx="1"/>
  </g>

  <g fill="#60a5fa" filter="url(#glow)">
    <circle cx="100" cy="128" r="3" opacity="0.8"/>
    <circle cx="115" cy="125" r="2" opacity="0.6"/>
    <circle cx="130" cy="131" r="2.5" opacity="0.7"/>
    <circle cx="145" cy="126" r="2" opacity="0.5"/>
    <circle cx="160" cy="130" r="3" opacity="0.8"/>
  </g>

  <g fill="url(#torrentGradient)" opacity="0.9">
    <path d="M110 190 L110 170 L118 170 L118 190 L126 190 L118 205 L110 205 L102 190 Z"/>
    <path d="M130 195 L130 175 L138 175 L138 195 L146 195 L138 210 L130 210 L122 195 Z"/>
    <path d="M150 190 L150 170 L158 170 L158 190 L166 190 L158 205 L150 205 L142 190 Z"/>
  </g>

  <g>
    <circle cx="40" cy="60" r="6" fill="#10b981" opacity="0.8"/>
    <circle cx="216" cy="60" r="6" fill="#10b981" opacity="0.8"/>
    <circle cx="40" cy="196" r="6" fill="#10b981" opacity="0.8"/>
    <circle cx="216" cy="196" r="6" fill="#10b981" opacity="0.8"/>
  </g>

  <path d="M 80 50 Q 128 30 176 50" stroke="#64748b" stroke-width="2" fill="none" opacity="0.6"/>
  <text x="128" y="45" text-anchor="middle" fill="#cbd5e1" font-family="Arial, sans-serif" font-size="12" font-weight="bold">VPN</text>
</svg>`),
}
