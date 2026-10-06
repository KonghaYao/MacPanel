<template>
    <IconShell name="firewall" :colors="colors">
        <path d="M32 11L49 18.5v14c0 11-8 18.5-17 21-9-2.5-17-10-17-21v-14L32 11Z" fill="#fff" />
        <path d="M32 11L49 18.5v14c0 11-8 18.5-17 21-9-2.5-17-10-17-21v-14L32 11Z" fill="#fff" opacity="0.2" />
        <path
            d="M23 31l7 7 14-15"
            stroke="#ff453a"
            stroke-width="4.2"
            stroke-linecap="round"
            stroke-linejoin="round"
            fill="none"
        />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.firewall.colors;
</script>
