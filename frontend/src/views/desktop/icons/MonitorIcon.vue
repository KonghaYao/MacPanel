<template>
    <IconShell name="monitor" :colors="colors">
        <rect x="16" y="22" width="7" height="22" rx="2" fill="#30d158" />
        <rect x="28.5" y="16" width="7" height="28" rx="2" fill="#64d2ff" />
        <rect x="41" y="26" width="7" height="18" rx="2" fill="#ffd60a" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.monitor.colors;
</script>
