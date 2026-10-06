<template>
    <IconShell name="gpu" :colors="colors">
        <rect x="18" y="20" width="28" height="24" rx="4" fill="#fff" opacity="0.92" />
        <rect x="22" y="24" width="8" height="8" rx="1.5" fill="#5e35b1" />
        <rect x="34" y="24" width="8" height="8" rx="1.5" fill="#5e35b1" />
        <rect x="22" y="36" width="20" height="4" rx="1.5" fill="#5e35b1" opacity="0.7" />
        <rect x="24" y="14" width="3" height="6" rx="1.2" fill="#64d2ff" />
        <rect x="30.5" y="14" width="3" height="6" rx="1.2" fill="#64d2ff" />
        <rect x="37" y="14" width="3" height="6" rx="1.2" fill="#64d2ff" />
        <rect x="24" y="44" width="3" height="6" rx="1.2" fill="#64d2ff" />
        <rect x="30.5" y="44" width="3" height="6" rx="1.2" fill="#64d2ff" />
        <rect x="37" y="44" width="3" height="6" rx="1.2" fill="#64d2ff" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.gpu.colors;
</script>
