<template>
    <IconShell name="containers" :colors="colors">
        <defs>
            <linearGradient id="box-back" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#e8f4ff" />
                <stop offset="100%" stop-color="#b8d4f0" />
            </linearGradient>
            <linearGradient id="box-mid" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#ffffff" />
                <stop offset="100%" stop-color="#d0e8ff" />
            </linearGradient>
            <linearGradient id="box-front" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#ffffff" />
                <stop offset="100%" stop-color="#e0f0ff" />
            </linearGradient>
        </defs>
        <rect x="12" y="32" width="20" height="16" rx="2.5" fill="url(#box-back)" stroke="#1565c0" stroke-width="0.6" opacity="0.85" />
        <rect x="20" y="24" width="20" height="16" rx="2.5" fill="url(#box-mid)" stroke="#2272c9" stroke-width="0.6" />
        <rect x="28" y="16" width="20" height="16" rx="2.5" fill="url(#box-front)" stroke="#3d9cf0" stroke-width="0.7" />
        <path d="M32 21h12M32 25h8" stroke="#1565c0" stroke-width="1.4" stroke-linecap="round" opacity="0.45" />
        <rect x="28" y="16" width="20" height="5" rx="2.5" fill="#fff" opacity="0.35" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.containers.colors;
</script>
