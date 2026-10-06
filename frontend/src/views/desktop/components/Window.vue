<template>
    <section
        class="mac-window"
        :class="{ active, maximized: state.maximized }"
        :style="frameStyle"
        @pointerdown="focus"
    >
        <header class="titlebar" @pointerdown="startDrag" @dblclick.stop="emit('zoom')">
            <div class="traffic" @pointerdown.stop @dblclick.stop>
                <button
                    type="button"
                    class="light close"
                    :aria-label="t('desktop.menus.closeWindow')"
                    @click="emit('close')"
                >
                    <svg viewBox="0 0 12 12" aria-hidden="true"><path d="M3.2 3.2 8.8 8.8M8.8 3.2 3.2 8.8" /></svg>
                </button>
                <button
                    type="button"
                    class="light min"
                    :aria-label="t('desktop.menus.minimize')"
                    @click="emit('minimize')"
                >
                    <svg viewBox="0 0 12 12" aria-hidden="true"><path d="M2.6 6h6.8" /></svg>
                </button>
                <button type="button" class="light zoom" :aria-label="t('desktop.menus.zoom')" @click="emit('zoom')">
                    <svg viewBox="0 0 12 12" aria-hidden="true"><path d="M3.2 5.2V3.2H5.2M8.8 6.8v2H6.8" /></svg>
                </button>
            </div>
            <h2 class="title">{{ title }}</h2>
        </header>
        <div class="body">
            <div v-if="state.kind === 'about'" class="about">
                <span class="about-mark" />
                <strong>{{ t('desktop.aboutTitle') }}</strong>
                <p>{{ t('desktop.aboutBody') }}</p>
                <p v-if="version" class="version">{{ version }}</p>
            </div>
            <template v-else>
                <div v-if="loading" class="loading" />
                <iframe :src="src" :title="title" @load="loading = false" />
            </template>
        </div>
        <template v-if="!state.maximized">
            <i class="resize n" @pointerdown="startResize('n', $event)" />
            <i class="resize s" @pointerdown="startResize('s', $event)" />
            <i class="resize e" @pointerdown="startResize('e', $event)" />
            <i class="resize w" @pointerdown="startResize('w', $event)" />
            <i class="resize ne" @pointerdown="startResize('ne', $event)" />
            <i class="resize nw" @pointerdown="startResize('nw', $event)" />
            <i class="resize se" @pointerdown="startResize('se', $event)" />
            <i class="resize sw" @pointerdown="startResize('sw', $event)" />
        </template>
    </section>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import type { DesktopWindowState } from '@/views/desktop/apps';

type ResizeEdge = 'n' | 's' | 'e' | 'w' | 'ne' | 'nw' | 'se' | 'sw';

const props = defineProps<{
    state: DesktopWindowState;
    title: string;
    active: boolean;
    src: string;
    version: string;
}>();

const emit = defineEmits<{
    focus: [];
    close: [];
    minimize: [];
    zoom: [];
    interact: [active: boolean];
    change: [patch: Pick<DesktopWindowState, 'x' | 'y' | 'width' | 'height'>];
}>();

const { t } = useI18n();
const loading = ref(props.state.kind === 'iframe');

const frameStyle = computed(() => {
    if (props.state.maximized) {
        return { zIndex: props.state.z };
    }
    return {
        zIndex: props.state.z,
        left: `${props.state.x}px`,
        top: `${props.state.y}px`,
        width: `${props.state.width}px`,
        height: `${props.state.height}px`,
    };
});

const focus = () => {
    emit('focus');
};

const stageBox = () => {
    const stage = document.querySelector('.mac-stage');
    if (!(stage instanceof HTMLElement)) {
        return { width: window.innerWidth, height: window.innerHeight };
    }
    const rect = stage.getBoundingClientRect();
    return { width: rect.width, height: rect.height };
};

const trackPointer = (event: PointerEvent, move: (ev: PointerEvent) => void) => {
    const handle = event.currentTarget;
    if (handle instanceof HTMLElement) {
        handle.setPointerCapture(event.pointerId);
    }
    emit('interact', true);
    const end = () => {
        window.removeEventListener('pointermove', move);
        window.removeEventListener('pointerup', end);
        emit('interact', false);
    };
    window.addEventListener('pointermove', move);
    window.addEventListener('pointerup', end);
};

const startDrag = (event: PointerEvent) => {
    if (props.state.maximized || event.button !== 0) {
        return;
    }
    const target = event.target;
    if (target instanceof HTMLElement && target.closest('.traffic')) {
        return;
    }
    focus();
    const originX = event.clientX;
    const originY = event.clientY;
    const startX = props.state.x;
    const startY = props.state.y;
    trackPointer(event, (ev) => {
        const box = stageBox();
        const nextX = startX + ev.clientX - originX;
        const nextY = startY + ev.clientY - originY;
        emit('change', {
            x: Math.min(Math.max(nextX, 80 - props.state.width), box.width - 80),
            y: Math.min(Math.max(nextY, 0), Math.max(box.height - 36, 0)),
            width: props.state.width,
            height: props.state.height,
        });
    });
};

const startResize = (edge: ResizeEdge, event: PointerEvent) => {
    if (props.state.maximized || event.button !== 0) {
        return;
    }
    event.stopPropagation();
    focus();
    const originX = event.clientX;
    const originY = event.clientY;
    const start = {
        x: props.state.x,
        y: props.state.y,
        width: props.state.width,
        height: props.state.height,
    };
    const minWidth = 480;
    const minHeight = 300;
    trackPointer(event, (ev) => {
        const box = stageBox();
        const dx = ev.clientX - originX;
        const dy = ev.clientY - originY;
        let { x, y, width, height } = start;
        if (edge.includes('e')) {
            width = Math.max(minWidth, start.width + dx);
        }
        if (edge.includes('s')) {
            height = Math.max(minHeight, start.height + dy);
        }
        if (edge.includes('w')) {
            width = Math.max(minWidth, start.width - dx);
            x = start.x + start.width - width;
        }
        if (edge.includes('n')) {
            height = Math.max(minHeight, start.height - dy);
            y = start.y + start.height - height;
        }
        width = Math.min(width, box.width - Math.max(x, 0));
        height = Math.min(height, box.height - Math.max(y, 0));
        if (x < 0) {
            width += x;
            x = 0;
        }
        if (y < 0) {
            height += y;
            y = 0;
        }
        emit('change', {
            x,
            y,
            width: Math.max(minWidth, width),
            height: Math.max(minHeight, height),
        });
    });
};
</script>

<style scoped lang="scss">
.mac-window {
    position: absolute;
    display: flex;
    flex-direction: column;
    min-width: 0;
    min-height: 0;
    border-radius: 12px;
    overflow: hidden;
    background: rgba(246, 246, 246, 0.86);
    backdrop-filter: blur(28px) saturate(1.6);
    box-shadow:
        0 28px 80px rgba(0, 0, 0, 0.32),
        0 0 0 1px rgba(0, 0, 0, 0.18);
    color: #1d1d1f;
}

.mac-window.maximized {
    inset: 8px;
}

.mac-window:not(.active) {
    box-shadow:
        0 16px 40px rgba(0, 0, 0, 0.22),
        0 0 0 1px rgba(0, 0, 0, 0.12);
}

.titlebar {
    position: relative;
    display: flex;
    align-items: center;
    height: 40px;
    flex: none;
    background: rgba(236, 236, 238, 0.78);
    border-bottom: 1px solid rgba(0, 0, 0, 0.08);
    cursor: grab;
    user-select: none;
}

.titlebar:active {
    cursor: grabbing;
}

.traffic {
    display: flex;
    gap: 8px;
    padding-left: 14px;
    z-index: 1;
}

.light {
    width: 12px;
    height: 12px;
    border: 0;
    border-radius: 50%;
    padding: 0;
    display: grid;
    place-items: center;
    background: #c6c6c8;
    box-shadow: inset 0 0 0 0.5px rgba(0, 0, 0, 0.18);
}

.light svg {
    width: 8px;
    height: 8px;
    opacity: 0;
    stroke: rgba(0, 0, 0, 0.55);
    stroke-width: 1.4;
    fill: none;
    stroke-linecap: round;
}

.traffic:hover svg {
    opacity: 1;
}

.active .light.close,
.traffic:hover .light.close {
    background: #ff5f57;
}

.active .light.min,
.traffic:hover .light.min {
    background: #febc2e;
}

.active .light.zoom,
.traffic:hover .light.zoom {
    background: #28c840;
}

.title {
    position: absolute;
    left: 78px;
    right: 78px;
    margin: 0;
    text-align: center;
    font-size: 13px;
    font-weight: 590;
    letter-spacing: -0.01em;
    pointer-events: none;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.body {
    position: relative;
    flex: 1;
    min-height: 0;
    background: #f5f5f7;
}

.body iframe {
    width: 100%;
    height: 100%;
    border: 0;
    background: #f5f5f7;
}

.loading {
    position: absolute;
    inset: 0;
    background:
        linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.55), transparent) 0 0 / 40% 100% no-repeat,
        #f5f5f7;
    animation: mac-window-sheen 1.1s ease-in-out infinite;
    z-index: 1;
    pointer-events: none;
}

.about {
    height: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 8px;
    padding: 24px;
    text-align: center;
}

.about-mark {
    width: 64px;
    height: 64px;
    border-radius: 16px;
    margin-bottom: 6px;
    background:
        radial-gradient(circle at 30% 30%, rgba(255, 255, 255, 0.45), transparent 40%),
        linear-gradient(160deg, #8eb6ff, #f3c6a5);
    box-shadow:
        inset 0 1px 0 rgba(255, 255, 255, 0.45),
        0 8px 18px rgba(0, 0, 0, 0.18);
}

.about strong {
    font-size: 18px;
    font-weight: 650;
}

.about p {
    margin: 0;
    max-width: 320px;
    font-size: 13px;
    line-height: 1.45;
    color: #3a3a3c;
}

.version {
    color: #6e6e73 !important;
}

.resize {
    position: absolute;
    z-index: 3;
}

.resize.n,
.resize.s {
    left: 10px;
    right: 10px;
    height: 8px;
    cursor: ns-resize;
}

.resize.n {
    top: -3px;
}

.resize.s {
    bottom: -3px;
}

.resize.e,
.resize.w {
    top: 10px;
    bottom: 10px;
    width: 8px;
    cursor: ew-resize;
}

.resize.e {
    right: -3px;
}

.resize.w {
    left: -3px;
}

.resize.ne,
.resize.nw,
.resize.se,
.resize.sw {
    width: 14px;
    height: 14px;
}

.resize.ne {
    top: -3px;
    right: -3px;
    cursor: nesw-resize;
}

.resize.nw {
    top: -3px;
    left: -3px;
    cursor: nwse-resize;
}

.resize.se {
    right: -3px;
    bottom: -3px;
    cursor: nwse-resize;
}

.resize.sw {
    left: -3px;
    bottom: -3px;
    cursor: nesw-resize;
}

@keyframes mac-window-sheen {
    from {
        background-position: -40% 0;
    }
    to {
        background-position: 140% 0;
    }
}
</style>
