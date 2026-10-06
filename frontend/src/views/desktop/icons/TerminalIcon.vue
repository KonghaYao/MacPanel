<template>
    <IconShell name="terminal" :colors="colors">
        <path
            d="M18 24l8 8-8 8"
            stroke="#30d158"
            stroke-width="4.5"
            stroke-linecap="round"
            stroke-linejoin="round"
            fill="none"
        />
        <rect x="30" y="38" width="18" height="4.5" rx="2" fill="#e5e5ea" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.terminal.colors;
</script>
