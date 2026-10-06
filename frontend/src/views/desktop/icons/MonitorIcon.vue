<template>
    <IconShell name="monitor" :colors="colors">
        <defs>
            <linearGradient id="chart-bg" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#ffffff" stop-opacity="0.92" />
                <stop offset="100%" stop-color="#e8f8ea" stop-opacity="0.85" />
            </linearGradient>
            <linearGradient id="bar-green" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#7ee08a" />
                <stop offset="100%" stop-color="#1a9e3a" />
            </linearGradient>
            <linearGradient id="bar-blue" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#7ec8ff" />
                <stop offset="100%" stop-color="#007aff" />
            </linearGradient>
            <linearGradient id="bar-red" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#ff8a80" />
                <stop offset="100%" stop-color="#ff3333" />
            </linearGradient>
        </defs>
        <rect x="14" y="17" width="36" height="30" rx="3.5" fill="url(#chart-bg)" stroke="#1a9e3a" stroke-width="0.6" opacity="0.9" />
        <rect x="17" y="26" width="7" height="18" rx="1.5" fill="url(#bar-green)" />
        <rect x="28.5" y="20" width="7" height="24" rx="1.5" fill="url(#bar-blue)" />
        <rect x="40" y="30" width="7" height="14" rx="1.5" fill="url(#bar-red)" />
        <rect x="14" y="17" width="36" height="8" rx="3.5" fill="#fff" opacity="0.3" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.monitor.colors;
</script>
