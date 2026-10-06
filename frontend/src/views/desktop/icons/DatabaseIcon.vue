<template>
    <IconShell name="database" :colors="colors">
        <ellipse cx="32" cy="18" rx="16" ry="6" fill="#fff" opacity="0.95" />
        <path d="M16 18v28c0 3.3 7.2 6 16 6s16-2.7 16-6V18" fill="#fff" opacity="0.88" />
        <ellipse cx="32" cy="32" rx="16" ry="6" fill="#00796b" opacity="0.35" />
        <ellipse cx="32" cy="46" rx="16" ry="6" fill="#00796b" opacity="0.35" />
        <ellipse cx="32" cy="18" rx="16" ry="6" fill="#fff" opacity="0.95" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.database.colors;
</script>
