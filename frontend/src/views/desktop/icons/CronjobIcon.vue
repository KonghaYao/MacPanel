<template>
    <IconShell name="cronjob" :colors="colors">
        <circle cx="32" cy="32" r="16" fill="#fff" opacity="0.95" />
        <circle cx="32" cy="32" r="14" fill="none" stroke="#cc7700" stroke-width="2" opacity="0.35" />
        <rect x="31" y="19" width="2.2" height="14" rx="1.1" fill="#1c1c1e" />
        <rect x="31" y="30" width="11" height="2.2" rx="1.1" fill="#1c1c1e" />
        <circle cx="32" cy="32" r="2.2" fill="#1c1c1e" />
        <circle cx="32" cy="32" r="16" fill="none" stroke="#fff" stroke-width="1.5" opacity="0.25" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.cronjob.colors;
</script>
