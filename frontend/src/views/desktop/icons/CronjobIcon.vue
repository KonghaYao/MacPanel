<template>
    <IconShell name="cronjob" :colors="colors">
        <defs>
            <radialGradient id="clock-face" cx="0.5" cy="0.5" r="0.5">
                <stop offset="0%" stop-color="#ffffff" />
                <stop offset="100%" stop-color="#ffe8c0" />
            </radialGradient>
            <linearGradient id="clock-rim" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#fff8e8" />
                <stop offset="100%" stop-color="#cc7700" />
            </linearGradient>
        </defs>
        <circle cx="32" cy="32" r="17" fill="url(#clock-rim)" />
        <circle cx="32" cy="32" r="15" fill="url(#clock-face)" stroke="#cc7700" stroke-width="0.6" />
        <circle cx="32" cy="32" r="15" fill="#fff" opacity="0.2" />
        <rect x="31" y="18" width="2.4" height="15" rx="1.2" fill="#333" />
        <rect x="31" y="30" width="12" height="2.4" rx="1.2" fill="#333" />
        <circle cx="32" cy="32" r="2.5" fill="#cc7700" />
        <circle cx="32" cy="32" r="1" fill="#fff" opacity="0.6" />
        <circle v-for="i in 12" :key="i" :cx="32 + 12 * Math.sin((i * Math.PI) / 6)" :cy="32 - 12 * Math.cos((i * Math.PI) / 6)" r="0.8" fill="#cc7700" opacity="0.5" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.cronjob.colors;
</script>
