<template>
    <IconShell name="containers" :colors="colors">
        <rect x="13" y="28" width="22" height="16" rx="3" fill="#fff" opacity="0.88" />
        <rect x="22" y="20" width="22" height="16" rx="3" fill="#fff" opacity="0.95" />
        <rect x="31" y="12" width="22" height="16" rx="3" fill="#fff" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.containers.colors;
</script>
