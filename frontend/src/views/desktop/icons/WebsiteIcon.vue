<template>
    <IconShell name="website" :colors="colors">
        <defs>
            <radialGradient id="compass-face" cx="0.5" cy="0.5" r="0.5">
                <stop offset="0%" stop-color="#ffffff" />
                <stop offset="100%" stop-color="#e8eef8" />
            </radialGradient>
        </defs>
        <circle cx="32" cy="32" r="18" fill="url(#compass-face)" stroke="#b0b8c8" stroke-width="0.8" />
        <circle cx="32" cy="32" r="15" fill="none" stroke="#c8c8cc" stroke-width="0.8" />
        <path d="M32 14v36M14 32h36" stroke="#d0d0d4" stroke-width="0.8" />
        <ellipse cx="32" cy="32" rx="15" ry="7.5" fill="none" stroke="#d0d0d4" stroke-width="0.7" />
        <ellipse cx="32" cy="32" rx="7.5" ry="15" fill="none" stroke="#d0d0d4" stroke-width="0.7" />
        <path d="M32 32L32 14L44 32Z" fill="#ff3333" />
        <path d="M32 32L32 50L20 32Z" fill="#0066cc" />
        <path d="M32 32L44 32L32 14Z" fill="#ff6666" opacity="0.6" />
        <path d="M32 32L20 32L32 50Z" fill="#3388dd" opacity="0.6" />
        <circle cx="32" cy="32" r="3" fill="#fff" stroke="#888" stroke-width="0.5" />
        <circle cx="32" cy="32" r="1.2" fill="#666" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.website.colors;
</script>
