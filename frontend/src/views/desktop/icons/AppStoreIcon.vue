<template>
    <IconShell name="appStore" :colors="colors">
        <path d="M32 13L17 47h9l3-7.5h14l3 7.5h9L32 13Z" fill="#fff" />
        <rect x="26.5" y="33" width="11" height="3.8" rx="1.2" fill="#0a84ff" />
        <path d="M32 13l-4 10h8L32 13Z" fill="#fff" opacity="0.35" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.appStore.colors;
</script>
