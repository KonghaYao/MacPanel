<template>
    <IconShell name="website" :colors="colors">
        <circle cx="32" cy="32" r="18" fill="#f2f2f7" />
        <circle cx="32" cy="32" r="15" fill="none" stroke="#c7c7cc" stroke-width="1.5" />
        <path d="M32 14v36M14 32h36" stroke="#c7c7cc" stroke-width="1.2" />
        <path d="M32 14c-5 6-5 20 0 26M32 14c5 6 5 20 0 26" stroke="#c7c7cc" stroke-width="1.2" fill="none" />
        <path d="M32 32L32 14L44 32Z" fill="#ff3b30" />
        <path d="M32 32L32 50L20 32Z" fill="#007aff" />
        <circle cx="32" cy="32" r="3" fill="#fff" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.website.colors;
</script>
