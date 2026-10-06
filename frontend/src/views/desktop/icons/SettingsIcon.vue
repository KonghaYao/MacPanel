<template>
    <IconShell name="settings" :colors="colors">
        <circle cx="32" cy="32" r="7.5" fill="#fff" />
        <path
            d="M32 11v7M32 46v7M11 32h7M46 32h7M17.5 17.5l5 5M41.5 41.5l5 5M46.5 17.5l-5 5M22.5 41.5l-5 5"
            stroke="#fff"
            stroke-width="5.2"
            stroke-linecap="round"
        />
        <circle cx="32" cy="32" r="7.5" fill="none" stroke="#fff" stroke-width="1.5" opacity="0.35" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.settings.colors;
</script>
