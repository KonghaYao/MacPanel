<template>
    <IconShell name="logs" :colors="colors">
        <path d="M18 14h18l10 10v24a2 2 0 0 1-2 2H18a2 2 0 0 1-2-2V16a2 2 0 0 1 2-2Z" fill="#e5e5ea" />
        <path d="M36 14v10h10" fill="#c7c7cc" />
        <rect x="22" y="32" width="20" height="2.5" rx="1.2" fill="#636366" />
        <rect x="22" y="38" width="16" height="2.5" rx="1.2" fill="#636366" />
        <rect x="22" y="44" width="18" height="2.5" rx="1.2" fill="#636366" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.logs.colors;
</script>
