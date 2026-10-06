<template>
    <IconShell name="monitor" :colors="colors">
        <rect x="15" y="18" width="34" height="28" rx="4" fill="#2c2c2e" opacity="0.55" />
        <rect x="17" y="24" width="6.5" height="20" rx="2" fill="#30d158" />
        <rect x="28.5" y="20" width="6.5" height="24" rx="2" fill="#64d2ff" />
        <rect x="40" y="28" width="6.5" height="16" rx="2" fill="#ffd60a" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.monitor.colors;
</script>
