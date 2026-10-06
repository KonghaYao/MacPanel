<template>
    <IconShell name="gpu" :colors="colors">
        <defs>
            <linearGradient id="chip-body" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#e8d8ff" />
                <stop offset="100%" stop-color="#5b21b6" />
            </linearGradient>
            <linearGradient id="chip-die" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#c4a0ff" />
                <stop offset="100%" stop-color="#7c3aed" />
            </linearGradient>
            <linearGradient id="pin-metal" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#d0d0d4" />
                <stop offset="100%" stop-color="#888" />
            </linearGradient>
        </defs>
        <rect x="16" y="18" width="32" height="28" rx="4" fill="url(#chip-body)" stroke="#5b21b6" stroke-width="0.6" />
        <rect x="20" y="22" width="10" height="10" rx="2" fill="url(#chip-die)" />
        <rect x="34" y="22" width="10" height="10" rx="2" fill="url(#chip-die)" />
        <rect x="20" y="36" width="24" height="6" rx="1.5" fill="url(#chip-die)" opacity="0.7" />
        <rect x="16" y="18" width="32" height="8" rx="4" fill="#fff" opacity="0.15" />
        <rect v-for="x in [20, 27, 34, 41]" :key="'t' + x" :x="x" :y="12" width="3.5" height="6" rx="1" fill="url(#pin-metal)" />
        <rect v-for="x in [20, 27, 34, 41]" :key="'b' + x" :x="x" :y="46" width="3.5" height="6" rx="1" fill="url(#pin-metal)" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.gpu.colors;
</script>
