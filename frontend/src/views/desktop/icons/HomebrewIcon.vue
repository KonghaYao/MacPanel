<template>
    <IconShell name="homebrew" :colors="colors">
        <defs>
            <linearGradient id="glass-body" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#fff8e8" />
                <stop offset="100%" stop-color="#f0d080" />
            </linearGradient>
            <linearGradient id="beer-liquid" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#f5c840" />
                <stop offset="100%" stop-color="#c87800" />
            </linearGradient>
            <linearGradient id="foam" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#ffffff" />
                <stop offset="100%" stop-color="#f0e8d8" />
            </linearGradient>
        </defs>
        <path d="M20 16h24l-2 34H22L20 16Z" fill="url(#glass-body)" stroke="#a07828" stroke-width="0.6" />
        <ellipse cx="32" cy="16" rx="13" ry="4.5" fill="url(#foam)" stroke="#a07828" stroke-width="0.5" />
        <rect x="19" y="26" width="26" height="5" rx="1.5" fill="#fff" opacity="0.85" />
        <ellipse cx="32" cy="36" rx="9" ry="11" fill="url(#beer-liquid)" opacity="0.75" />
        <path
            d="M25 12c1.5 2.5 2.8 2.5 3.2 4.5M32 10c1.5 2.5 2.8 2.5 3.2 4.5"
            stroke="#fff8e8"
            stroke-width="2.5"
            stroke-linecap="round"
            fill="none"
        />
        <path d="M39 32h8a3 3 0 0 1 0 6h-8" stroke="#a07828" stroke-width="3" stroke-linecap="round" fill="none" />
        <ellipse cx="32" cy="16" rx="13" ry="4.5" fill="#fff" opacity="0.3" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.homebrew.colors;
</script>
