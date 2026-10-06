<template>
    <component :is="DOCK_ICON_MAP[name]" class="dock-icon" />
</template>

<script setup lang="ts">
import type { DockGlyph } from '@/views/desktop/apps';
import { DOCK_ICON_MAP } from '@/views/desktop/icons';

defineProps<{ name: DockGlyph }>();
</script>

<style scoped>
.dock-icon {
    width: 100%;
    height: 100%;
    display: block;
}
</style>
