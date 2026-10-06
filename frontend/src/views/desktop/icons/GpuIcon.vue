<template>
    <IconShell name="gpu" :colors="colors">
        <rect x="17" y="19" width="30" height="26" rx="4.5" fill="#fff" opacity="0.94" />
        <rect x="21" y="23" width="9" height="9" rx="2" fill="#6d28d9" />
        <rect x="34" y="23" width="9" height="9" rx="2" fill="#6d28d9" />
        <rect x="21" y="36" width="22" height="5" rx="1.5" fill="#6d28d9" opacity="0.65" />
        <rect x="23" y="13" width="3.5" height="6" rx="1.2" fill="#64d2ff" />
        <rect x="30" y="13" width="3.5" height="6" rx="1.2" fill="#64d2ff" />
        <rect x="37" y="13" width="3.5" height="6" rx="1.2" fill="#64d2ff" />
        <rect x="23" y="45" width="3.5" height="6" rx="1.2" fill="#64d2ff" />
        <rect x="30" y="45" width="3.5" height="6" rx="1.2" fill="#64d2ff" />
        <rect x="37" y="45" width="3.5" height="6" rx="1.2" fill="#64d2ff" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.gpu.colors;
</script>
