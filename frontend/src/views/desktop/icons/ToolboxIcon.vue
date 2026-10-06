<template>
    <IconShell name="toolbox" :colors="colors">
        <defs>
            <linearGradient id="wrench-body" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#ffffff" />
                <stop offset="100%" stop-color="#ffd080" />
            </linearGradient>
            <linearGradient id="hammer-head" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#e8e8e8" />
                <stop offset="100%" stop-color="#888" />
            </linearGradient>
        </defs>
        <rect x="38" y="14" width="10" height="7" rx="1.5" fill="url(#hammer-head)" transform="rotate(45 43 17.5)" />
        <rect x="14" y="38" width="7" height="10" rx="1.5" fill="url(#hammer-head)" transform="rotate(45 17.5 43)" />
        <path
            d="M40 18l-8 8a5.5 5.5 0 1 1-7.8 7.8l-8 8 3.5 3.5 8-8a5.5 5.5 0 1 1 7.8-7.8l8-8-3.5-3.5Z"
            fill="url(#wrench-body)"
            stroke="#cc6600"
            stroke-width="0.5"
        />
        <circle cx="26" cy="38" r="5.5" fill="none" stroke="#fff" stroke-width="3.2" />
        <circle cx="26" cy="38" r="5.5" fill="none" stroke="#cc6600" stroke-width="0.6" opacity="0.4" />
        <path d="M40 18l-4 4" stroke="#fff" stroke-width="1.5" stroke-linecap="round" opacity="0.4" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.toolbox.colors;
</script>
