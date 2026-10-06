<template>
    <IconShell name="firewall" :colors="colors">
        <path d="M32 12L48 19v13c0 10.5-7 17.5-16 20-9-2.5-16-9.5-16-20V19L32 12Z" fill="#fff" opacity="0.95" />
        <path
            d="M24 30l6 6 12-13"
            stroke="#d70015"
            stroke-width="4"
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
