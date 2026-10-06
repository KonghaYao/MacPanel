<template>
    <IconShell name="logs" :colors="colors">
        <defs>
            <linearGradient id="notepad-body" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#fffde0" />
                <stop offset="100%" stop-color="#ffeb3b" />
            </linearGradient>
            <linearGradient id="notepad-lines" x1="0" y1="0" x2="1" y2="0">
                <stop offset="0%" stop-color="#e6c800" stop-opacity="0.15" />
                <stop offset="100%" stop-color="#e6c800" stop-opacity="0" />
            </linearGradient>
        </defs>
        <rect x="16" y="12" width="32" height="40" rx="2" fill="url(#notepad-body)" stroke="#c8a800" stroke-width="0.6" />
        <rect x="16" y="12" width="32" height="8" rx="2" fill="#fff" opacity="0.35" />
        <line x1="20" y1="26" x2="44" y2="26" stroke="#c8a800" stroke-width="0.6" opacity="0.4" />
        <line x1="20" y1="32" x2="44" y2="32" stroke="#c8a800" stroke-width="0.6" opacity="0.35" />
        <line x1="20" y1="38" x2="44" y2="38" stroke="#c8a800" stroke-width="0.6" opacity="0.35" />
        <line x1="20" y1="44" x2="40" y2="44" stroke="#c8a800" stroke-width="0.6" opacity="0.35" />
        <rect x="20" y="28" width="18" height="2.2" rx="1" fill="#8a7800" opacity="0.55" />
        <rect x="20" y="34" width="14" height="2.2" rx="1" fill="#8a7800" opacity="0.45" />
        <rect x="20" y="40" width="16" height="2.2" rx="1" fill="#8a7800" opacity="0.45" />
        <rect x="16" y="12" width="32" height="40" rx="2" fill="url(#notepad-lines)" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.logs.colors;
</script>
