<template>
    <IconShell name="terminal" :colors="colors">
        <defs>
            <linearGradient id="term-screen" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#2a2a2c" />
                <stop offset="100%" stop-color="#0a0a0a" />
            </linearGradient>
        </defs>
        <rect x="12" y="16" width="40" height="32" rx="4" fill="url(#term-screen)" stroke="#48484a" stroke-width="0.8" />
        <rect x="12" y="16" width="40" height="6" rx="4" fill="#fff" opacity="0.06" />
        <path
            d="M17 27l8 8-8 8"
            stroke="#39ff14"
            stroke-width="3.8"
            stroke-linecap="round"
            stroke-linejoin="round"
            fill="none"
        />
        <rect x="30" y="38" width="16" height="3.8" rx="1.5" fill="#39ff14" opacity="0.9" />
        <circle cx="18" cy="19.5" r="1.2" fill="#ff5f57" opacity="0.7" />
        <circle cx="22.5" cy="19.5" r="1.2" fill="#febc2e" opacity="0.7" />
        <circle cx="27" cy="19.5" r="1.2" fill="#28c840" opacity="0.7" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.terminal.colors;
</script>
