<template>
    <IconShell name="appStore" :colors="colors">
        <path d="M32 14L18 46h8.5l2.8-7h11.4l2.8 7H51L37 14h-5Z" fill="#fff" />
        <rect x="27" y="33" width="10" height="3.5" rx="1.2" fill="#1a9cff" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.appStore.colors;
</script>
