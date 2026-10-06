<template>
    <IconShell name="website" :colors="colors">
        <circle cx="32" cy="32" r="17" fill="#f5f5f7" />
        <circle cx="32" cy="32" r="14.5" fill="none" stroke="#aeaeb2" stroke-width="1.2" />
        <path d="M32 15v34M15 32h34" stroke="#c7c7cc" stroke-width="1" />
        <ellipse cx="32" cy="32" rx="14.5" ry="7" fill="none" stroke="#c7c7cc" stroke-width="1" />
        <ellipse cx="32" cy="32" rx="7" ry="14.5" fill="none" stroke="#c7c7cc" stroke-width="1" />
        <path d="M32 15c-4.5 5.5-4.5 18.5 0 24M32 15c4.5 5.5 4.5 18.5 0 24" stroke="#c7c7cc" stroke-width="1" fill="none" />
        <path d="M32 32L32 15L43 32Z" fill="#ff3b30" />
        <path d="M32 32L32 49L21 32Z" fill="#007aff" />
        <circle cx="32" cy="32" r="2.8" fill="#fff" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.website.colors;
</script>
