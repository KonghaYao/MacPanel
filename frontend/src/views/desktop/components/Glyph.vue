<template>
    <div class="glyph">
        <img :src="SYMBOL_URLS[name]" class="symbol" alt="" />
    </div>
</template>

<script setup lang="ts">
import type { DockGlyph } from '@/views/desktop/apps';
import { SYMBOL_URLS } from '@/views/desktop/icons';

defineProps<{ name: DockGlyph }>();
</script>

<style scoped>
.glyph {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
}

.symbol {
    width: 100%;
    height: 100%;
    object-fit: contain;
    pointer-events: none;
}
</style>
