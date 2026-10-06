<template>
    <IconShell name="terminal" :colors="colors">
        <rect x="13" y="17" width="38" height="30" rx="4" fill="#1a1a1a" opacity="0.45" />
        <path
            d="M17 26l7.5 7.5L17 41"
            stroke="#32d74b"
            stroke-width="4.2"
            stroke-linecap="round"
            stroke-linejoin="round"
            fill="none"
        />
        <rect x="29" y="37" width="17" height="4.2" rx="2" fill="#f5f5f7" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.terminal.colors;
</script>
