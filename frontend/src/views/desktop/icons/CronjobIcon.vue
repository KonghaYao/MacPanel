<template>
    <IconShell name="cronjob" :colors="colors">
        <circle cx="32" cy="32" r="17" fill="#1c1c1e" />
        <circle cx="32" cy="32" r="15" fill="none" stroke="#48484a" stroke-width="1.5" />
        <rect x="31" y="20" width="2" height="13" rx="1" fill="#fff" />
        <rect x="31" y="31" width="10" height="2" rx="1" fill="#fff" />
        <circle cx="32" cy="32" r="2.5" fill="#fff" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.cronjob.colors;
</script>
