<template>
    <IconShell name="logs" :colors="colors">
        <path d="M17 13h19l11 11v25a2.5 2.5 0 0 1-2.5 2.5H17a2.5 2.5 0 0 1-2.5-2.5V15.5A2.5 2.5 0 0 1 17 13Z" fill="#f2f2f7" />
        <path d="M36 13v11h11" fill="#d1d1d6" />
        <rect x="22" y="32" width="21" height="2.6" rx="1.3" fill="#636366" />
        <rect x="22" y="38.5" width="17" height="2.6" rx="1.3" fill="#636366" />
        <rect x="22" y="45" width="19" height="2.6" rx="1.3" fill="#636366" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.logs.colors;
</script>
