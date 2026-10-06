<template>
    <IconShell name="toolbox" :colors="colors">
        <path
            d="M40 20l-8 8a5.5 5.5 0 1 1-7.8 7.8l-8 8 3.5 3.5 8-8a5.5 5.5 0 1 1 7.8-7.8l8-8L40 20Z"
            fill="#fff"
            opacity="0.95"
        />
        <circle cx="27" cy="37" r="4.5" fill="none" stroke="#fff" stroke-width="3.5" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.toolbox.colors;
</script>
