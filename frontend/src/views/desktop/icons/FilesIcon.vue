<template>
    <IconShell name="files" :colors="colors">
        <defs>
            <linearGradient id="folder-back" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#a8d4ff" />
                <stop offset="100%" stop-color="#2272c9" />
            </linearGradient>
            <linearGradient id="folder-front" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#d8ecff" />
                <stop offset="100%" stop-color="#4a9ff5" />
            </linearGradient>
        </defs>
        <path d="M14 22h14l3 4h19v22a3 3 0 0 1-3 3H14a3 3 0 0 1-3-3V25a3 3 0 0 1 3-3Z" fill="url(#folder-back)" />
        <path d="M14 28h36v20a3 3 0 0 1-3 3H14a3 3 0 0 1-3-3V28Z" fill="url(#folder-front)" />
        <path d="M14 28h36v6H14Z" fill="#fff" opacity="0.25" />
        <path d="M14 22h14l3 4h19" fill="none" stroke="#fff" stroke-width="0.6" opacity="0.35" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.files.colors;
</script>
