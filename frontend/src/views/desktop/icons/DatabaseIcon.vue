<template>
    <IconShell name="database" :colors="colors">
        <ellipse cx="32" cy="17" rx="15" ry="5.5" fill="#fff" />
        <path d="M17 17v30c0 3 6.7 5.5 15 5.5s15-2.5 15-5.5V17" fill="#fff" opacity="0.9" />
        <ellipse cx="32" cy="31" rx="15" ry="5.5" fill="#248a3d" opacity="0.3" />
        <ellipse cx="32" cy="42" rx="15" ry="5.5" fill="#248a3d" opacity="0.35" />
        <ellipse cx="32" cy="17" rx="15" ry="5.5" fill="#fff" opacity="0.95" />
        <path d="M22 17c3.5 2 8.5 2 12 0M22 31c3.5 2 8.5 2 12 0" stroke="#34c759" stroke-width="1.2" opacity="0.4" fill="none" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.database.colors;
</script>
