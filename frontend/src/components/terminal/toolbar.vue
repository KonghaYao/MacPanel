<template>
    <div class="terminal-toolbar" :class="{ 'is-inline': inline }">
        <el-tooltip :content="$t('commons.button.search')" placement="top">
            <el-button text class="terminal-toolbar-btn" icon="Search" @click="emit('search')" />
        </el-tooltip>
        <el-tooltip :content="$t('commons.button.copy')" placement="top">
            <el-button text class="terminal-toolbar-btn" icon="CopyDocument" :disabled="!canCopy" @click="emit('copy')" />
        </el-tooltip>
        <el-tooltip :content="$t('terminal.paste')" placement="top">
            <el-button text class="terminal-toolbar-btn" icon="DocumentCopy" @click="emit('paste')" />
        </el-tooltip>
        <el-tooltip :content="$t('terminal.clearScreen')" placement="top">
            <el-button text class="terminal-toolbar-btn" icon="Brush" @click="emit('clear')" />
        </el-tooltip>
        <el-tooltip v-if="showReconnect" :content="$t('commons.button.reconnect')" placement="top">
            <el-button text class="terminal-toolbar-btn" icon="Refresh" @click="emit('reconnect')" />
        </el-tooltip>
        <el-tooltip v-if="showFullscreen" :content="fullscreenLabel" placement="top">
            <el-button text class="terminal-toolbar-btn" icon="FullScreen" @click="emit('fullscreen')" />
        </el-tooltip>
        <div v-if="slots.default" class="terminal-toolbar-extra">
            <slot />
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, useSlots } from 'vue';
import i18n from '@/lang';

const props = defineProps<{
    canCopy?: boolean;
    showReconnect?: boolean;
    showFullscreen?: boolean;
    fullscreenActive?: boolean;
    inline?: boolean;
}>();

const slots = useSlots();

const emit = defineEmits<{
    search: [];
    copy: [];
    paste: [];
    clear: [];
    reconnect: [];
    fullscreen: [];
}>();

const fullscreenLabel = computed(() =>
    i18n.global.t('commons.button.' + (props.fullscreenActive ? 'quitFullscreen' : 'fullscreen')),
);
</script>

<style scoped lang="scss">
.terminal-toolbar {
    display: flex;
    align-items: center;
    gap: 2px;
    min-width: 0;
    padding: 4px 6px;
    background: var(--el-bg-color);
    border-bottom: 1px solid var(--el-border-color-lighter);
}

.terminal-toolbar-btn {
    width: 30px;
    height: 30px;
    margin: 0;
    padding: 0;
    flex: 0 0 auto;
}

.terminal-toolbar.is-inline {
    flex: 0 0 auto;
    height: auto;
    padding: 0;
    background: transparent;
    border: 0;
}

.terminal-toolbar-extra {
    display: flex;
    align-items: center;
    gap: 8px;
    min-width: 0;
    flex: 1;
}

.terminal-toolbar.is-inline .terminal-toolbar-extra {
    flex: 0 0 auto;
}
</style>
