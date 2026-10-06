<template>
    <IconShell name="appStore" :colors="colors">
        <defs>
            <linearGradient id="app-a" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#ffffff" />
                <stop offset="100%" stop-color="#d0e8ff" />
            </linearGradient>
        </defs>
        <path d="M32 12L16 48h10l3-8h14l3 8h10L32 12Z" fill="url(#app-a)" />
        <rect x="26" y="33.5" width="12" height="4" rx="1.2" fill="#0066cc" />
        <path d="M32 12l-5 12h10L32 12Z" fill="#fff" opacity="0.45" />
        <path d="M32 12L16 48h10l3-8h14l3 8h10L32 12Z" fill="none" stroke="#fff" stroke-width="0.5" opacity="0.25" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.appStore.colors;
</script>
