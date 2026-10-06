<template>
    <div
        class="mac-desktop"
        :class="{ 'is-dark': isDarkTheme, 'is-interacting': interacting }"
        @pointerdown="closeMenus"
    >
        <div
            class="wallpaper"
            :class="{ 'has-image': !!wallpaperUrl }"
            :style="wallpaperStyle"
            aria-hidden="true"
        />
        <MenuBar
            ref="menuRef"
            :active-title="activeTitle"
            :active-window-id="activeWindowId"
            :windows="menuWindows"
            @about="openAbout"
            @settings="openApp('settings')"
            @back="backToPanel"
            @help="openHelp"
            @close="closeActive"
            @minimize="minimizeActive"
            @zoom="zoomActive"
            @focus="focusWindow"
        />
        <main class="mac-stage">
            <div class="desktop-icons">
                <button
                    v-for="app in desktopIcons"
                    :key="app.id"
                    type="button"
                    class="desktop-icon"
                    :class="{ selected: selectedIcon === app.id }"
                    @click.stop="openApp(app.id)"
                >
                    <span class="tile">
                        <Glyph :name="app.glyph" />
                    </span>
                    <span>{{ t(app.titleKey) }}</span>
                </button>
            </div>
            <Window
                v-for="item in windows"
                v-show="!item.minimized"
                :key="item.id"
                :state="item"
                :title="t(item.titleKey)"
                :active="item.id === activeWindowId"
                :src="item.kind === 'iframe' ? desktopEmbedSrc(item.route) : ''"
                :version="systemVersion"
                @focus="focusWindow(item.id)"
                @close="closeWindow(item.id)"
                @minimize="minimizeWindow(item.id)"
                @zoom="zoomWindow(item.id)"
                @interact="interacting = $event"
                @change="patchWindow(item.id, $event)"
            />
        </main>
        <Dock :apps="dockApps" :running-ids="runningIds" @open="openApp" />
    </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { fetchDesktopWallpaper } from '@/api/modules/desktop';
import { getSettingBaseInfo } from '@/api/modules/setting';
import { useGlobalStore } from '@/composables/useGlobalStore';
import { syncPlatformCapabilities } from '@/utils/platform';
import { desktopEmbedSrc } from '@/utils/desktop-embed';
import {
    DESKTOP_APPS,
    DESKTOP_ICON_IDS,
    findDesktopApp,
    type DesktopAppDef,
    type DesktopWindowState,
} from '@/views/desktop/apps';
import MenuBar from '@/views/desktop/components/MenuBar.vue';
import Dock from '@/views/desktop/components/Dock.vue';
import Window from '@/views/desktop/components/Window.vue';
import Glyph from '@/views/desktop/components/Glyph.vue';

const { t } = useI18n();
const router = useRouter();
const { isDarkTheme, platformFeatures, docsUrl } = useGlobalStore();

const windows = ref<DesktopWindowState[]>([]);
const activeWindowId = ref('');
const interacting = ref(false);
const selectedIcon = ref('');
const systemVersion = ref('');
const wallpaperUrl = ref('');
const menuRef = ref<InstanceType<typeof MenuBar>>();
let zTop = 10;
let windowSeed = 1;

const dockApps = computed(() =>
    DESKTOP_APPS.filter((app) => !app.platformFeature || platformFeatures.value[app.platformFeature]),
);
const desktopIcons = computed(() =>
    DESKTOP_ICON_IDS.map((id) => findDesktopApp(id)).filter((app): app is DesktopAppDef => !!app),
);
const runningIds = computed(() => windows.value.map((item) => item.appId));
const activeWindow = computed(() => windows.value.find((item) => item.id === activeWindowId.value));
const activeTitle = computed(() => (activeWindow.value ? t(activeWindow.value.titleKey) : ''));
const menuWindows = computed(() =>
    windows.value.map((item) => ({
        id: item.id,
        title: t(item.titleKey),
    })),
);
const wallpaperStyle = computed(() =>
    wallpaperUrl.value ? { backgroundImage: `url(${wallpaperUrl.value})` } : undefined,
);

const closeMenus = () => {
    selectedIcon.value = '';
    menuRef.value?.closeMenu();
};

const focusWindow = (id: string) => {
    const target = windows.value.find((item) => item.id === id);
    if (!target) {
        return;
    }
    target.minimized = false;
    zTop += 1;
    target.z = zTop;
    activeWindowId.value = id;
};

const openWindow = (partial: Omit<DesktopWindowState, 'id' | 'z' | 'minimized' | 'maximized'>) => {
    const existing = windows.value.find((item) => item.appId === partial.appId);
    if (existing) {
        focusWindow(existing.id);
        return;
    }
    const offset = (windows.value.length % 6) * 26;
    const stageWidth = Math.max(window.innerWidth - 48, 520);
    const stageHeight = Math.max(window.innerHeight - 150, 360);
    const width = Math.min(partial.width, stageWidth);
    const height = Math.min(partial.height, stageHeight);
    zTop += 1;
    const created: DesktopWindowState = {
        ...partial,
        id: `win-${windowSeed++}`,
        x: Math.max(24, Math.min(partial.x + offset, stageWidth - width)),
        y: Math.max(16, Math.min(partial.y + offset, stageHeight - 80)),
        width,
        height,
        z: zTop,
        minimized: false,
        maximized: false,
    };
    windows.value.push(created);
    activeWindowId.value = created.id;
};

const openApp = (id: string) => {
    const app = findDesktopApp(id);
    if (!app) {
        return;
    }
    selectedIcon.value = id;
    openWindow({
        appId: app.id,
        titleKey: app.titleKey,
        kind: 'iframe',
        route: app.route,
        x: 48,
        y: 24,
        width: 1040,
        height: 680,
    });
};

const openAbout = () => {
    openWindow({
        appId: 'about',
        titleKey: 'desktop.aboutTitle',
        kind: 'about',
        route: '',
        x: 180,
        y: 72,
        width: 460,
        height: 280,
    });
};

const closeWindow = (id: string) => {
    const index = windows.value.findIndex((item) => item.id === id);
    if (index < 0) {
        return;
    }
    windows.value.splice(index, 1);
    if (activeWindowId.value === id) {
        const next = [...windows.value].sort((a, b) => b.z - a.z)[0];
        activeWindowId.value = next?.id || '';
    }
};

const minimizeWindow = (id: string) => {
    const target = windows.value.find((item) => item.id === id);
    if (!target) {
        return;
    }
    target.minimized = true;
    if (activeWindowId.value === id) {
        const next = windows.value.filter((item) => !item.minimized).sort((a, b) => b.z - a.z)[0];
        activeWindowId.value = next?.id || '';
    }
};

const zoomWindow = (id: string) => {
    const target = windows.value.find((item) => item.id === id);
    if (!target) {
        return;
    }
    target.maximized = !target.maximized;
    focusWindow(id);
};

const patchWindow = (id: string, patch: Pick<DesktopWindowState, 'x' | 'y' | 'width' | 'height'>) => {
    const target = windows.value.find((item) => item.id === id);
    if (!target || target.maximized) {
        return;
    }
    target.x = patch.x;
    target.y = patch.y;
    target.width = patch.width;
    target.height = patch.height;
};

const closeActive = () => {
    if (activeWindowId.value) {
        closeWindow(activeWindowId.value);
    }
};

const minimizeActive = () => {
    if (activeWindowId.value) {
        minimizeWindow(activeWindowId.value);
    }
};

const zoomActive = () => {
    if (activeWindowId.value) {
        zoomWindow(activeWindowId.value);
    }
};

const backToPanel = () => {
    router.push({ name: 'home' });
};

const openHelp = () => {
    window.open(docsUrl.value, '_blank', 'noopener');
};

onMounted(async () => {
    await syncPlatformCapabilities();
    try {
        const res = await getSettingBaseInfo();
        systemVersion.value = res.data.systemVersion || '';
    } catch {
        systemVersion.value = '';
    }
    try {
        const res = await fetchDesktopWallpaper();
        if (res.data instanceof Blob && res.data.size > 0) {
            wallpaperUrl.value = URL.createObjectURL(res.data);
        }
    } catch {
        wallpaperUrl.value = '';
    }
});

onBeforeUnmount(() => {
    if (wallpaperUrl.value) {
        URL.revokeObjectURL(wallpaperUrl.value);
    }
});
</script>

<style scoped lang="scss">
.mac-desktop {
    position: relative;
    height: 100vh;
    height: 100dvh;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    color: #1d1d1f;
    font-family:
        -apple-system, BlinkMacSystemFont, 'SF Pro Text', 'SF Pro Display', 'Helvetica Neue', 'PingFang SC',
        'Hiragino Sans GB', sans-serif;
    -webkit-font-smoothing: antialiased;
}

.wallpaper {
    position: absolute;
    inset: 0;
    background:
        radial-gradient(ellipse at 18% 18%, rgba(255, 186, 140, 0.95), transparent 42%),
        radial-gradient(ellipse at 78% 12%, rgba(126, 176, 255, 0.95), transparent 40%),
        radial-gradient(ellipse at 72% 78%, rgba(186, 140, 255, 0.8), transparent 46%),
        radial-gradient(ellipse at 12% 86%, rgba(94, 214, 224, 0.75), transparent 38%),
        linear-gradient(155deg, #6f97f8 0%, #c9b7ff 46%, #f2c7ae 100%);
}

.wallpaper.has-image {
    background-size: cover;
    background-position: center;
    background-repeat: no-repeat;
}

.mac-desktop.is-dark .wallpaper {
    filter: saturate(0.9) brightness(0.72);
}

.mac-stage {
    position: relative;
    z-index: 1;
    flex: 1;
    min-height: 0;
}

.desktop-icons {
    position: absolute;
    top: 18px;
    right: 12px;
    display: flex;
    flex-direction: column;
    gap: 18px;
    z-index: 1;
}

.desktop-icon {
    width: 86px;
    border: 0;
    padding: 4px;
    background: transparent;
    color: #fff;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.55);
    font-size: 12px;
    line-height: 1.25;
    cursor: pointer;
}

.desktop-icon .tile {
    width: 56px;
    height: 56px;
    margin: 0 auto 6px;
    border-radius: 14px;
    overflow: hidden;
    display: block;
    box-shadow:
        0 1px 2px rgba(0, 0, 0, 0.18),
        0 8px 16px rgba(0, 0, 0, 0.18);
}

.desktop-icon.selected .tile,
.desktop-icon:focus-visible .tile {
    outline: 2px solid rgba(255, 255, 255, 0.92);
    outline-offset: 3px;
}

.mac-desktop.is-interacting :deep(iframe) {
    pointer-events: none;
}
</style>
