<template>
    <header class="menubar" @pointerdown.stop>
        <div class="leading">
            <div class="menu-slot">
                <button
                    type="button"
                    class="menu-btn apple"
                    :class="{ open: openMenu === 'apple' }"
                    :aria-label="t('desktop.menus.about')"
                    @mouseenter="hoverMenu('apple')"
                    @click="toggleMenu('apple')"
                >
                    <svg viewBox="0 0 16 16" aria-hidden="true">
                        <path
                            d="M11.2 8.4c0-1.7 1.4-2.5 1.5-2.6-.8-1.2-2.1-1.3-2.5-1.4-1.1-.1-2.1.6-2.6.6s-1.4-.6-2.3-.6c-1.2 0-2.3.7-2.9 1.7-1.2 2.1-.3 5.3.9 7 .6.8 1.3 1.8 2.2 1.7.9 0 1.2-.6 2.3-.6s1.4.6 2.3.6 1.5-.8 2.1-1.7c.7-.9 1-1.8 1-1.9-.1 0-1.9-.7-2-2.8ZM9.9 3.6c.5-.6.8-1.4.7-2.2-.7 0-1.5.5-2 .1-.5.6-.9 1.4-.8 2.2.8.1 1.6-.4 2.1-.1Z"
                        />
                    </svg>
                </button>
                <div v-if="openMenu === 'apple'" class="dropdown">
                    <button type="button" class="item" @click="run('about')">{{ t('desktop.menus.about') }}</button>
                    <button type="button" class="item" @click="run('settings')">
                        {{ t('desktop.menus.settings') }}
                    </button>
                    <i class="rule" />
                    <button type="button" class="item" @click="run('back')">{{ t('desktop.backToPanel') }}</button>
                </div>
            </div>
            <strong class="app-name">{{ activeTitle || t('desktop.appName') }}</strong>
            <div v-for="item in menus" :key="item.id" class="menu-slot">
                <button
                    type="button"
                    class="menu-btn"
                    :class="{ open: openMenu === item.id }"
                    @mouseenter="hoverMenu(item.id)"
                    @click="toggleMenu(item.id)"
                >
                    {{ t(item.label) }}
                </button>
                <div v-if="openMenu === item.id" class="dropdown" :class="{ 'window-menu': item.id === 'window' }">
                    <template v-if="item.id === 'file'">
                        <button type="button" class="item" :disabled="!activeWindowId" @click="run('close')">
                            {{ t('desktop.menus.closeWindow') }}
                        </button>
                    </template>
                    <template v-else-if="item.id === 'window'">
                        <button type="button" class="item" :disabled="!activeWindowId" @click="run('minimize')">
                            {{ t('desktop.menus.minimize') }}
                        </button>
                        <button type="button" class="item" :disabled="!activeWindowId" @click="run('zoom')">
                            {{ t('desktop.menus.zoom') }}
                        </button>
                        <i class="rule" />
                        <p v-if="windows.length === 0" class="empty">{{ t('desktop.noWindows') }}</p>
                        <button
                            v-for="win in windows"
                            :key="win.id"
                            type="button"
                            class="item"
                            :class="{ current: win.id === activeWindowId }"
                            @click="focus(win.id)"
                        >
                            <span class="mark" />
                            {{ win.title }}
                        </button>
                    </template>
                    <template v-else>
                        <button type="button" class="item" @click="run('help')">
                            {{ t('desktop.menus.helpDocs') }}
                        </button>
                    </template>
                </div>
            </div>
        </div>
        <div class="trailing">
            <time :datetime="isoTime">{{ clock }}</time>
        </div>
    </header>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useGlobalStore } from '@/composables/useGlobalStore';

interface MenuWindow {
    id: string;
    title: string;
}

defineProps<{
    activeTitle: string;
    activeWindowId: string;
    windows: MenuWindow[];
}>();

const emit = defineEmits<{
    about: [];
    settings: [];
    back: [];
    help: [];
    close: [];
    minimize: [];
    zoom: [];
    focus: [id: string];
}>();

const { t } = useI18n();
const { language } = useGlobalStore();
const openMenu = ref('');
const now = ref(new Date());
let timer = 0;

const menus = [
    { id: 'file', label: 'desktop.menus.file' },
    { id: 'window', label: 'desktop.menus.window' },
    { id: 'help', label: 'desktop.menus.help' },
];

const locale = computed(() => (language.value === 'zh' ? 'zh-CN' : language.value || 'en'));
const isoTime = computed(() => now.value.toISOString());
const clock = computed(() =>
    new Intl.DateTimeFormat(locale.value, {
        weekday: 'short',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
    }).format(now.value),
);

const toggleMenu = (id: string) => {
    openMenu.value = openMenu.value === id ? '' : id;
};

const hoverMenu = (id: string) => {
    if (openMenu.value) {
        openMenu.value = id;
    }
};

const closeMenu = () => {
    openMenu.value = '';
};

const run = (action: 'about' | 'settings' | 'back' | 'help' | 'close' | 'minimize' | 'zoom') => {
    closeMenu();
    if (action === 'about') {
        emit('about');
        return;
    }
    if (action === 'settings') {
        emit('settings');
        return;
    }
    if (action === 'back') {
        emit('back');
        return;
    }
    if (action === 'help') {
        emit('help');
        return;
    }
    if (action === 'close') {
        emit('close');
        return;
    }
    if (action === 'minimize') {
        emit('minimize');
        return;
    }
    emit('zoom');
};

const focus = (id: string) => {
    closeMenu();
    emit('focus', id);
};

onMounted(() => {
    timer = window.setInterval(() => {
        now.value = new Date();
    }, 10000);
    window.addEventListener('blur', closeMenu);
});

onBeforeUnmount(() => {
    window.clearInterval(timer);
    window.removeEventListener('blur', closeMenu);
});

defineExpose({ closeMenu });
</script>

<style scoped lang="scss">
.menubar {
    position: relative;
    z-index: 20;
    height: 28px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 10px 0 8px;
    color: #1d1d1f;
    background: rgba(246, 246, 246, 0.5);
    border-bottom: 1px solid rgba(0, 0, 0, 0.06);
    backdrop-filter: blur(22px) saturate(1.8);
    user-select: none;
}

.leading,
.trailing {
    display: flex;
    align-items: center;
    min-width: 0;
    gap: 2px;
}

.app-name {
    margin: 0 8px 0 4px;
    font-size: 13px;
    font-weight: 680;
    letter-spacing: -0.01em;
}

.menu-btn {
    height: 22px;
    padding: 0 8px;
    border: 0;
    border-radius: 4px;
    background: transparent;
    color: inherit;
    font-size: 13px;
    font-weight: 450;
    cursor: pointer;
}

.menu-btn.apple {
    width: 28px;
    padding: 0;
    display: grid;
    place-items: center;
}

.menu-btn.apple svg {
    width: 14px;
    height: 14px;
    fill: currentColor;
}

.menu-btn.open,
.menu-btn:hover {
    background: rgba(0, 0, 0, 0.08);
}

.trailing time {
    font-size: 13px;
    font-variant-numeric: tabular-nums;
    padding-left: 8px;
}

.menu-slot {
    position: relative;
}

.dropdown {
    position: absolute;
    top: 24px;
    left: 0;
    min-width: 220px;
    padding: 4px;
    border-radius: 10px;
    background: rgba(246, 246, 246, 0.9);
    border: 1px solid rgba(0, 0, 0, 0.08);
    box-shadow: 0 16px 40px rgba(0, 0, 0, 0.22);
    backdrop-filter: blur(24px) saturate(1.6);
}

.window-menu {
    max-height: 320px;
    overflow: auto;
}

.item {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 8px;
    border: 0;
    background: transparent;
    color: inherit;
    text-align: left;
    border-radius: 5px;
    padding: 4px 10px;
    font-size: 13px;
    cursor: pointer;
}

.item:hover:not(:disabled) {
    background: #0a84ff;
    color: #fff;
}

.item:disabled {
    opacity: 0.4;
    cursor: default;
}

.mark {
    width: 12px;
    height: 12px;
    flex: none;
    display: grid;
    place-items: center;
}

.item.current .mark::before {
    content: '';
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: currentColor;
}

.rule {
    display: block;
    height: 1px;
    margin: 4px 6px;
    background: rgba(0, 0, 0, 0.1);
}

.empty {
    margin: 4px 10px 6px;
    font-size: 12px;
    color: #6e6e73;
}

:global(.mac-desktop.is-dark) .menubar,
:global(.mac-desktop.is-dark) .dropdown {
    color: #f5f5f7;
    background: rgba(32, 32, 34, 0.62);
    border-color: rgba(255, 255, 255, 0.08);
}

:global(.mac-desktop.is-dark) .menu-btn.open,
:global(.mac-desktop.is-dark) .menu-btn:hover {
    background: rgba(255, 255, 255, 0.12);
}

:global(.mac-desktop.is-dark) .empty {
    color: #aeaeb2;
}
</style>
