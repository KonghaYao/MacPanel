<template>
    <IconShell name="containers" :colors="colors">
        <rect x="14" y="30" width="20" height="15" rx="2.5" fill="#fff" opacity="0.82" />
        <rect x="22" y="22" width="20" height="15" rx="2.5" fill="#fff" opacity="0.92" />
        <rect x="30" y="14" width="20" height="15" rx="2.5" fill="#fff" />
        <path d="M34 19h12M34 23h8" stroke="#2496ed" stroke-width="1.5" stroke-linecap="round" opacity="0.5" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.containers.colors;
</script>
