<template>
    <IconShell name="settings" :colors="colors">
        <circle cx="32" cy="32" r="8" fill="#fff" />
        <path
            d="M32 12v6M32 46v6M12 32h6M46 32h6M18.8 18.8l4.2 4.2M41 41l4.2 4.2M45.2 18.8 41 23M23 41l-4.2 4.2"
            stroke="#fff"
            stroke-width="5"
            stroke-linecap="round"
        />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.settings.colors;
</script>
