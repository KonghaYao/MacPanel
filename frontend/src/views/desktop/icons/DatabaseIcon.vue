<template>
    <IconShell name="database" :colors="colors">
        <defs>
            <linearGradient id="db-body" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#e8ffe8" />
                <stop offset="100%" stop-color="#34c759" />
            </linearGradient>
            <linearGradient id="db-top" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#ffffff" />
                <stop offset="100%" stop-color="#7ee08a" />
            </linearGradient>
        </defs>
        <path d="M17 17v30c0 3 6.7 5.5 15 5.5s15-2.5 15-5.5V17" fill="url(#db-body)" />
        <ellipse cx="32" cy="17" rx="15" ry="5.5" fill="url(#db-top)" stroke="#1e8a38" stroke-width="0.5" />
        <ellipse cx="32" cy="31" rx="15" ry="5.5" fill="none" stroke="#1e8a38" stroke-width="0.8" opacity="0.35" />
        <ellipse cx="32" cy="42" rx="15" ry="5.5" fill="none" stroke="#1e8a38" stroke-width="0.8" opacity="0.35" />
        <ellipse cx="32" cy="17" rx="15" ry="5.5" fill="#fff" opacity="0.3" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.database.colors;
</script>
