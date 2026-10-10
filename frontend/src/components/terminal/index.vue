<template>
    <!-- Right click is inert here: no custom menu and no native browser menu. -->
    <div class="terminal-shell" @contextmenu.prevent @mousedown.right.prevent>
        <div ref="terminalElement" class="terminal-container" :style="shellStyle"></div>
        <div v-if="searchOpen" class="terminal-search">
            <el-input
                ref="searchInputRef"
                v-model="searchText"
                size="small"
                :placeholder="$t('terminal.searchPlaceholder')"
                @keydown.enter.exact.prevent="find(true)"
                @keydown.shift.enter.prevent="find(false)"
                @keydown.esc.prevent="closeSearch"
                @input="find(true)"
            />
            <el-button text size="small" @click="find(false)">{{ $t('terminal.findPrev') }}</el-button>
            <el-button text size="small" @click="find(true)">{{ $t('terminal.findNext') }}</el-button>
            <el-button text size="small" @click="closeSearch">{{ $t('commons.button.close') }}</el-button>
        </div>
        <transition name="ai-mask-fade">
            <div v-if="aiNotice.loading" class="ai-notice-mask"></div>
        </transition>
        <transition name="ai-notice-fade">
            <div
                v-if="aiNotice.visible"
                class="ai-notice"
                :class="[`ai-notice--${aiNotice.level}`, { 'ai-notice--loading': aiNotice.loading }]"
            >
                {{ aiNotice.message }}
            </div>
        </transition>
    </div>
</template>

<script lang="ts" setup>
import { ref, shallowRef, watch, onActivated, onBeforeUnmount, nextTick, computed, onMounted, reactive } from 'vue';
import { Terminal } from '@xterm/xterm';
import type { SearchAddon } from '@xterm/addon-search';
import '@xterm/xterm/css/xterm.css';
import { FitAddon } from '@xterm/addon-fit';
import {
    TERMINAL_DEFAULTS,
    buildTerminalOptions,
    buildTerminalTheme,
    createSearchAddon,
    cursorStyleOf,
    loadClipboard,
    loadUnicode11,
    loadWebLinks,
    loadWebgl,
    searchDecorationOptions,
    type TerminalLinkTarget,
} from '@/components/terminal/options';
import { createBase64StreamDecoder, encodeBase64 } from '@/utils/base64';
import { TerminalStore } from '@/store';
import { MsgError } from '@/utils/message';
import { canReadClipboard, readClipboardText, writeClipboardText } from '@/utils/clipboard-api';
import { checkStreamAuth } from '@/utils/stream-auth';
import { useGlobalStore } from '@/composables/useGlobalStore';
import i18n from '@/lang';
const { currentNode } = useGlobalStore();

// session: agent side session id known (fresh or reattached)
// expired: the agent no longer has the session; a reconnect must open a new one
const emit = defineEmits(['session', 'expired']);

// Close codes of the agent's session protocol (agent/utils/terminal/session.go).
const CLOSE_SESSION_NOT_FOUND = 4404;
const CLOSE_ATTACHED_ELSEWHERE = 4409;
const CLOSE_REVALIDATE = 4410;

const terminalElement = ref<HTMLDivElement | null>(null);
const searchInputRef = ref<{ focus: () => void } | null>(null);
const fitAddon = new FitAddon();
const searchAddon = shallowRef<SearchAddon>();
const linkTarget: TerminalLinkTarget = { getElement: () => terminalElement.value };
const ui = reactive({ canCopy: false });
const searchOpen = ref(false);
const searchText = ref('');
let pending: { kind: 'connect'; endpoint: string; args: string } | { kind: 'error'; message: string } | null = null;
const termReady = ref(false);
const webSocketReady = ref(false);
const term = shallowRef<Terminal>();
const terminalSocket = ref<WebSocket>();
const heartbeatTimer = ref<NodeJS.Timer>();
let initWebSocketToken = 0;
// Output arrives as independent base64 chunks; decode them as one byte stream so that a multi-byte
// UTF-8 character split across two reads is not turned into a replacement character.
let decodeOutput = createBase64StreamDecoder();
const latency = ref(0);
// Reconnect state. Only terminals that received a session hello reconnect;
// the agent keeps a dirty-disconnected session alive for a short grace period.
const sessionId = ref('');
let wsEndpoint = '';
let wsArgs = '';
let closing = false;
let reconnecting = false;
let reconnectNoticeShown = false;
let revalidating = false;
let reconnectStartedAt = 0;
let reconnectDelay = 1000;
let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
// Must match graceTimeout in agent/utils/terminal/session.go: past it the agent has dropped the shell.
const reconnectWindow = 30 * 60 * 1000;
const initCmd = ref('');
const hideInitCmdEcho = ref(false);
const initCmdEchoBuffer = ref('');
const waitForPrompt = ref('');
const waitForPromptBuffer = ref('');
const aiNotice = ref({
    visible: false,
    loading: false,
    level: 'info',
    message: '',
});
let aiNoticeTimer: ReturnType<typeof setTimeout> | null = null;
let resizeFrame: number | undefined;
let lastResizeColumns = 0;
let lastResizeRows = 0;

const readyWatcher = watch(
    () => webSocketReady.value && termReady.value,
    (ready) => {
        if (ready) {
            changeTerminalSize();
            readyWatcher(); // unwatch self
        }
    },
);

const terminalStore = TerminalStore();
const lineHeight = computed(() => terminalStore.lineHeight);
const fontSize = computed(() => terminalStore.fontSize);
const fontFamily = computed(() => terminalStore.fontFamily);
const backgroundColor = computed(() => terminalStore.backgroundColor);
const foregroundColor = computed(() => terminalStore.foregroundColor);
const letterSpacing = computed(() => terminalStore.letterSpacing);
watch(
    [lineHeight, fontSize, letterSpacing, fontFamily],
    ([newLineHeight, newFontSize, newLetterSpacing, newFontFamily]) => {
        if (!term.value) return;
        term.value.options.lineHeight = newLineHeight;
        term.value.options.letterSpacing = newLetterSpacing;
        term.value.options.fontSize = newFontSize;
        term.value.options.fontFamily = newFontFamily;
        changeTerminalSize();
    },
);
watch([backgroundColor, foregroundColor], ([newBackgroundColor, newForegroundColor]) => {
    if (!term.value) return;
    term.value.options.theme = buildTerminalTheme(newBackgroundColor, newForegroundColor);
    applyTerminalBackground(newBackgroundColor);
});
const cursorStyle = computed(() => terminalStore.cursorStyle);
watch(cursorStyle, () => {
    if (!term.value) return;
    term.value.options.cursorStyle = getStyle();
});
const cursorBlink = computed(() => terminalStore.cursorBlink);
watch(cursorBlink, (newCursorBlink) => {
    if (!term.value) return;
    term.value.options.cursorBlink = String(newCursorBlink).toLowerCase() === 'enable';
});
const scrollback = computed(() => terminalStore.scrollback);
watch(scrollback, (newScrollback) => {
    if (!term.value) return;
    term.value.options.scrollback = newScrollback;
});
const scrollSensitivity = computed(() => terminalStore.scrollSensitivity);
watch(scrollSensitivity, (newScrollSensitivity) => {
    if (!term.value) return;
    term.value.options.scrollSensitivity = newScrollSensitivity;
});

interface WsProps {
    endpoint: string;
    args: string;
    error: string;
    initCmd: string;
    waitForPrompt?: string;
    sessionId?: string;
}

interface TerminalBufferLine {
    isWrapped?: boolean;
    translateToString(trimRight?: boolean, startColumn?: number, endColumn?: number): string;
}
const acceptParams = (props: WsProps) => {
    nextTick(() => {
        if (props.error.length !== 0) {
            initError(props.error);
        } else {
            initCmd.value = props.initCmd || '';
            waitForPrompt.value = props.waitForPrompt || '';
            waitForPromptBuffer.value = '';
            sessionId.value = props.sessionId || '';
            init(props.endpoint, props.args);
        }
    });
};

const shellStyle = computed(() => ({
    backgroundColor: terminalStore.backgroundColor || TERMINAL_DEFAULTS.backgroundColor,
    '--terminal-bg': terminalStore.backgroundColor || TERMINAL_DEFAULTS.backgroundColor,
}));

const newTerm = () => {
    term.value = new Terminal(
        buildTerminalOptions(
            {
                lineHeight: terminalStore.lineHeight,
                letterSpacing: terminalStore.letterSpacing,
                fontSize: terminalStore.fontSize,
                fontFamily: terminalStore.fontFamily,
                backgroundColor: terminalStore.backgroundColor,
                foregroundColor: terminalStore.foregroundColor,
                cursorBlink: String(terminalStore.cursorBlink).toLowerCase() === 'enable',
                cursorStyle: cursorStyleOf(terminalStore.cursorStyle),
                scrollback: terminalStore.scrollback,
                scrollSensitivity: terminalStore.scrollSensitivity,
            },
            linkTarget,
        ),
    );
};

const applyTerminalBackground = (color: string) => {
    if (!terminalElement.value) return;
    terminalElement.value.style.backgroundColor = color || TERMINAL_DEFAULTS.backgroundColor;
    terminalElement.value.style.backgroundImage = '';
    terminalElement.value.style.backgroundSize = '';
    terminalElement.value.style.backgroundPosition = '';
    terminalElement.value.style.backgroundRepeat = '';
    terminalElement.value.style.imageRendering = '';
};

const getStyle = (): 'underline' | 'block' | 'bar' => cursorStyleOf(terminalStore.cursorStyle);

const init = (endpoint: string, args: string) => {
    pending = { kind: 'connect', endpoint, args };
    tryMount();
};

const initError = (errorInfo: string) => {
    pending = { kind: 'error', message: errorInfo };
    tryMount();
};

function tryMount() {
    if (termReady.value || !pending || !terminalElement.value) {
        return;
    }
    if (terminalElement.value.clientWidth <= 0 || terminalElement.value.clientHeight <= 0) {
        return;
    }
    const job = pending;
    if (!openTerminal(job.kind === 'connect')) {
        return;
    }
    pending = null;
    if (job.kind === 'error') {
        term.value?.write(job.message);
        return;
    }
    initWebSocket(job.endpoint, job.args);
}

function onClose(isKeepShow: boolean = false) {
    initWebSocketToken++;
    pending = null;
    closing = true;
    closeSearch();
    stopReconnect();
    window.removeEventListener('resize', changeTerminalSize);
    if (resizeFrame !== undefined) {
        cancelAnimationFrame(resizeFrame);
        resizeFrame = undefined;
    }
    lastResizeColumns = 0;
    lastResizeRows = 0;
    clearAINotice();
    webSocketReady.value = false;
    try {
        // 1000 tells the agent this is deliberate: close the shell now, no grace period
        terminalSocket.value?.close(1000);
    } catch {}
    if (heartbeatTimer.value) {
        clearInterval(Number(heartbeatTimer.value));
        heartbeatTimer.value = undefined;
    }
    terminalSocket.value = undefined;
    if (!isKeepShow) {
        try {
            term.value?.dispose();
        } catch {}
        term.value = undefined;
        searchAddon.value = undefined;
        termReady.value = false;
        ui.canCopy = false;
        if (terminalElement.value) {
            terminalElement.value.innerHTML = '';
        }
    }
}

// terminal 相关代码 start

function openTerminal(online: boolean): boolean {
    if (!terminalElement.value) {
        return false;
    }
    decodeOutput = createBase64StreamDecoder();
    newTerm();
    const current = term.value;
    if (!current) {
        return false;
    }
    lastResizeColumns = 0;
    lastResizeRows = 0;
    current.open(terminalElement.value);
    applyTerminalBackground(terminalStore.backgroundColor);
    loadUnicode11(current);
    loadWebgl(current);
    current.loadAddon(fitAddon);
    fitAddon.fit();
    searchAddon.value = createSearchAddon();
    current.loadAddon(searchAddon.value);
    loadWebLinks(current, linkTarget);
    loadClipboard(current);
    current.attachCustomKeyEventHandler(onKeyDown);
    current.onSelectionChange(() => {
        ui.canCopy = current.hasSelection();
    });
    window.addEventListener('resize', changeTerminalSize);
    if (online) {
        current.onData((data) => onTermData(data));
    }
    termReady.value = true;
    return true;
}

function changeTerminalSize() {
    if (resizeFrame !== undefined) {
        return;
    }
    resizeFrame = requestAnimationFrame(() => {
        resizeFrame = undefined;
        resizeTerminal();
    });
}

function resizeTerminal() {
    if (!terminalElement.value || !term.value) return;
    if (terminalElement.value.clientWidth <= 0 || terminalElement.value.clientHeight <= 0) {
        return;
    }

    fitAddon.fit();
    if (isWsOpen()) {
        const { cols, rows } = term.value;
        if (cols === lastResizeColumns && rows === lastResizeRows) {
            return;
        }
        lastResizeColumns = cols;
        lastResizeRows = rows;
        terminalSocket.value!.send(
            JSON.stringify({
                type: 'resize',
                cols: cols,
                rows: rows,
            }),
        );
    }
}

// terminal 相关代码 end

// websocket 相关代码 start

const initWebSocket = async (endpoint_: string, args: string = '') => {
    const token = ++initWebSocketToken;
    closing = false;
    wsEndpoint = endpoint_;
    wsArgs = args;
    const href = window.location.href;
    const protocol = href.split('//')[0] === 'http:' ? 'ws' : 'wss';
    const host = href.split('//')[1].split('/')[0];
    const endpoint = endpoint_.replace(/^\/+/, '');
    let node = args.indexOf('id=') !== -1 ? 'local' : currentNode.value;
    let conn = `${protocol}://${host}/${endpoint}?cols=${term.value.cols}&rows=${term.value.rows}&${args}&operateNode=${node}`;
    if (args.indexOf('operateNode=') !== -1) {
        conn = `${protocol}://${host}/${endpoint}?cols=${term.value.cols}&rows=${term.value.rows}&${args}`;
    }
    if (sessionId.value) {
        conn += `&session=${encodeURIComponent(sessionId.value)}`;
    }
    if (revalidating) {
        conn += '&terminalRevalidate=1';
    }
    const authError = await checkStreamAuth(conn);
    if (token !== initWebSocketToken || !termReady.value) {
        return;
    }
    if (authError) {
        reconnecting = false;
        revalidating = false;
        sessionId.value = '';
        showWebSocketAuthError(authError);
        emit('expired');
        return;
    }
    if (heartbeatTimer.value) {
        clearInterval(Number(heartbeatTimer.value));
    }
    terminalSocket.value = new WebSocket(conn);
    terminalSocket.value.onopen = runRealTerminal;
    terminalSocket.value.onmessage = onWSReceive;
    terminalSocket.value.onclose = closeRealTerminal;
    terminalSocket.value.onerror = errorRealTerminal;
    heartbeatTimer.value = setInterval(() => {
        if (isWsOpen()) {
            terminalSocket.value!.send(
                JSON.stringify({
                    type: 'heartbeat',
                    timestamp: `${new Date().getTime()}`,
                }),
            );
        }
    }, 1000 * 10);
};

const showWebSocketAuthError = (message: string) => {
    clearAINotice();
    MsgError(message);
    term.value?.write(`\x1b[31m${message}\x1b[m\r\n`);
};

const runRealTerminal = () => {
    webSocketReady.value = true;
    changeTerminalSize();
    term.value?.focus();
    // a reattached shell already ran its init command
    if (initCmd.value !== '' && !sessionId.value) {
        hideInitCmdEcho.value = true;
        initCmdEchoBuffer.value = '';
        sendMsg(initCmd.value);
    }
};

const stripInitCmdEchoLine = (message: string) => {
    if (!hideInitCmdEcho.value) {
        return message;
    }
    initCmdEchoBuffer.value += message;
    const lineBreakIndex = initCmdEchoBuffer.value.search(/\r?\n/);
    if (lineBreakIndex === -1) {
        return '';
    }

    const lineBreakLength = initCmdEchoBuffer.value[lineBreakIndex] === '\r' ? 2 : 1;
    const remaining = initCmdEchoBuffer.value.slice(lineBreakIndex + lineBreakLength);
    hideInitCmdEcho.value = false;
    initCmdEchoBuffer.value = '';
    initCmd.value = '';
    return remaining;
};

const flushPromptBuffer = (message: string) => {
    if (!waitForPrompt.value) {
        return message;
    }
    waitForPromptBuffer.value += message;
    const promptIndex = waitForPromptBuffer.value.indexOf(waitForPrompt.value);
    if (promptIndex === -1) {
        return '';
    }

    const visible = waitForPromptBuffer.value.slice(promptIndex);
    waitForPrompt.value = '';
    waitForPromptBuffer.value = '';
    return visible;
};

const onWSReceive = (message: MessageEvent) => {
    const wsMsg = JSON.parse(message.data);
    switch (wsMsg.type) {
        case 'cmd': {
            if (wsMsg.data) {
                let receiveMsg = decodeOutput(wsMsg.data);
                if (hideInitCmdEcho.value) {
                    receiveMsg = stripInitCmdEchoLine(receiveMsg);
                }
                if (receiveMsg && waitForPrompt.value) {
                    receiveMsg = flushPromptBuffer(receiveMsg);
                }
                if (!receiveMsg) {
                    break;
                }
                term.value.write(receiveMsg);
            }
            break;
        }
        case 'heartbeat': {
            latency.value = new Date().getTime() - wsMsg.timestamp;
            break;
        }
        case 'session': {
            const wasReconnect = reconnecting;
            const wasRevalidate = revalidating;
            reconnecting = false;
            revalidating = false;
            reconnectDelay = 1000;
            sessionId.value = wsMsg.id || '';
            decodeOutput = createBase64StreamDecoder();
            if (wasReconnect && !wasRevalidate) {
                // replay is a tail of recent output, start from a clean screen
                term.value?.reset();
            }
            emit('session', sessionId.value);
            break;
        }
        case 'ai_notice': {
            const message = wsMsg.message?.trim();
            if (!message) {
                break;
            }
            showAINotice(wsMsg.level || 'info', message);
            break;
        }
    }
};

const errorRealTerminal = (ex: any) => {
    clearAINotice();
    if (reconnecting) return;
    let message = ex.message;
    if (!message) message = 'disconnected';
    term.value.write(`\x1b[31m${message}\x1b[m\r\n`);
};

const closeRealTerminal = (ev: CloseEvent) => {
    clearAINotice();
    webSocketReady.value = false;
    if (heartbeatTimer.value) {
        clearInterval(Number(heartbeatTimer.value));
        heartbeatTimer.value = undefined;
    }
    terminalSocket.value = undefined;
    if (closing || !sessionId.value) {
        term.value?.write('The connection has been disconnected.');
        term.value?.write(ev.reason);
        return;
    }
    switch (ev.code) {
        case 1000: // the shell exited or the agent closed it
        case CLOSE_SESSION_NOT_FOUND:
            sessionId.value = '';
            reconnecting = false;
            writeNotice(
                '31',
                ev.code === 1000 ? 'The connection has been disconnected.' : i18n.global.t('terminal.sessionExpired'),
            );
            emit('expired');
            return;
        case CLOSE_ATTACHED_ELSEWHERE:
            reconnecting = false;
            writeNotice('31', i18n.global.t('terminal.sessionKicked'));
            return;
        case CLOSE_REVALIDATE:
            revalidating = true;
            scheduleReconnect(true);
            return;
        default:
            scheduleReconnect();
    }
};

const writeNotice = (color: string, message: string) => {
    term.value?.write(`\r\n\x1b[${color}m${message}\x1b[m\r\n`);
};

// scheduleReconnect retries with backoff for as long as the agent keeps a detached session.
const scheduleReconnect = (forRevalidation = false) => {
    const now = Date.now();
    if (!reconnecting) {
        reconnecting = true;
        reconnectStartedAt = now;
        reconnectDelay = 1000;
        reconnectNoticeShown = false;
    } else if (now - reconnectStartedAt > reconnectWindow) {
        reconnecting = false;
        sessionId.value = '';
        writeNotice('31', i18n.global.t('terminal.sessionExpired'));
        emit('expired');
        return;
    }
    if (!forRevalidation && !reconnectNoticeShown) {
        writeNotice('33', i18n.global.t('terminal.sessionReconnecting'));
        reconnectNoticeShown = true;
    }
    reconnectTimer = setTimeout(
        () => {
            reconnectTimer = null;
            if (closing || !sessionId.value) return;
            initWebSocket(wsEndpoint, wsArgs);
        },
        forRevalidation ? 0 : reconnectDelay,
    );
    if (!forRevalidation) {
        reconnectDelay = Math.min(reconnectDelay * 2, 8000);
    }
};

const stopReconnect = () => {
    reconnecting = false;
    revalidating = false;
    if (reconnectTimer) {
        clearTimeout(reconnectTimer);
        reconnectTimer = null;
    }
};

const isWsOpen = () => {
    const readyState = terminalSocket.value && terminalSocket.value.readyState;
    return readyState === 1;
};

function isEnterInputData(data: string): boolean {
    return data === '\r' || data === '\n' || data === '\r\n';
}

function getCurrentTerminalLine(): string {
    const xterm = term.value;
    if (!xterm?.buffer?.active) return '';
    const buffer = xterm.buffer.active;
    const cursorRow = buffer.baseY + buffer.cursorY;
    let startRow = cursorRow;
    let endRow = cursorRow;

    for (let row = cursorRow; row > 0; row--) {
        const line = buffer.getLine(row) as TerminalBufferLine | undefined;
        if (!line?.isWrapped) {
            startRow = row;
            break;
        }
        startRow = row - 1;
    }

    for (let row = cursorRow + 1; row < buffer.length; row++) {
        const line = buffer.getLine(row) as TerminalBufferLine | undefined;
        if (!line?.isWrapped) {
            break;
        }
        endRow = row;
    }

    let content = '';
    for (let row = startRow; row <= endRow; row++) {
        const line = buffer.getLine(row) as TerminalBufferLine | undefined;
        if (!line) continue;
        content += line.translateToString(false);
    }
    return content.trimEnd();
}

function sendMsg(data: string, line: string = '') {
    if (isWsOpen()) {
        terminalSocket.value!.send(
            JSON.stringify({
                type: 'cmd',
                data: encodeBase64(data),
                line,
            }),
        );
    }
}

function onTermData(data: string) {
    if (!data) return;
    if (aiNotice.value.loading) return;
    sendMsg(data, isEnterInputData(data) ? getCurrentTerminalLine() : '');
}

function clearAINotice() {
    if (aiNoticeTimer) {
        clearTimeout(aiNoticeTimer);
        aiNoticeTimer = null;
    }
    aiNotice.value = {
        ...aiNotice.value,
        visible: false,
        loading: false,
    };
}

function showAINotice(level: string, message: string) {
    if (aiNoticeTimer) {
        clearTimeout(aiNoticeTimer);
        aiNoticeTimer = null;
    }
    const resolvedLevel = ['success', 'error', 'info'].includes(level) ? level : 'info';
    aiNotice.value = {
        visible: true,
        loading: resolvedLevel === 'info',
        level: resolvedLevel,
        message,
    };
    if (resolvedLevel === 'info') {
        return;
    }
    aiNoticeTimer = setTimeout(() => {
        aiNotice.value = {
            ...aiNotice.value,
            visible: false,
            loading: false,
        };
        aiNoticeTimer = null;
    }, 2600);
}

// websocket 相关代码 end

const resizeObserver = ref<ResizeObserver>();

function onKeyDown(event: KeyboardEvent): boolean {
    if (event.type !== 'keydown' || !term.value) {
        return true;
    }
    const key = event.key.toLowerCase();
    const command = event.metaKey;
    // ⌘C and ⌘V are deliberately left to the browser: xterm picks up the native copy/paste events
    // from its own textarea, and unlike the async Clipboard API those also work on plain HTTP.
    if (command && key === 'f') {
        event.preventDefault();
        openSearch();
        return false;
    }
    const mac = /Mac|iPhone|iPad/.test(navigator.userAgent);
    if (!mac && event.ctrlKey && !event.metaKey && !event.altKey && key === 'c' && term.value.hasSelection()) {
        event.preventDefault();
        copySelection().catch(reportClipboardError);
        return false;
    }
    return true;
}

async function copySelection() {
    const text = term.value?.getSelection() ?? '';
    if (!text) {
        return;
    }
    await writeClipboardText(text);
}

async function pasteClipboard() {
    if (!canReadClipboard()) {
        // the browser handles ⌘V / Ctrl+V itself, keep the terminal focused for it
        term.value?.focus();
        MsgError(i18n.global.t('terminal.pasteUnavailable'));
        return;
    }
    const text = await readClipboardText();
    if (!text || !term.value) {
        return;
    }
    term.value.paste(text);
}

function reportClipboardError(error: unknown) {
    const message = error instanceof Error ? error.message : String(error);
    if (message) {
        MsgError(message);
    }
}

function openSearch() {
    searchOpen.value = true;
    nextTick(() => searchInputRef.value?.focus());
}

function closeSearch() {
    searchOpen.value = false;
    searchText.value = '';
    searchAddon.value?.clearDecorations();
}

function find(forward: boolean) {
    const addon = searchAddon.value;
    if (!addon) {
        return;
    }
    if (!searchText.value) {
        addon.clearDecorations();
        return;
    }
    if (forward) {
        addon.findNext(searchText.value, searchDecorationOptions());
    } else {
        addon.findPrevious(searchText.value, searchDecorationOptions());
    }
}

onMounted(() => {
    resizeObserver.value = new ResizeObserver(() => {
        if (!termReady.value) {
            tryMount();
            return;
        }
        changeTerminalSize();
    });

    if (terminalElement.value) {
        resizeObserver.value.observe(terminalElement.value);
    }
    tryMount();
});

defineExpose({
    acceptParams,
    onClose,
    isWsOpen,
    sendMsg,
    getLatency: () => latency.value,
    refit: () => changeTerminalSize(),
    ui,
    openSearch,
    copySelection: () => copySelection().catch(reportClipboardError),
    pasteClipboard: () => pasteClipboard().catch(reportClipboardError),
    clearScreen: () => term.value?.clear(),
});

onBeforeUnmount(() => {
    onClose();
    resizeObserver.value?.disconnect();
});

onActivated(() => {
    nextTick(changeTerminalSize);
});
</script>

<style lang="scss" scoped>
.terminal-container {
    width: 100%;
    height: 100%;
    background-color: var(--terminal-bg, #111827);
}

.terminal-shell {
    position: relative;
    width: 100%;
    height: 100%;
}

.ai-notice-mask {
    position: absolute;
    inset: 0;
    z-index: 10;
    background: rgba(8, 10, 14, 0.12);
    backdrop-filter: blur(1.5px);
    pointer-events: auto;
    cursor: progress;
}

.ai-notice {
    position: absolute;
    left: 50%;
    top: 24px;
    transform: translateX(-50%);
    z-index: 12;
    width: fit-content;
    min-width: 240px;
    max-width: min(72%, 560px);
    padding: 9px 14px;
    border-radius: 999px;
    border: 1px solid rgba(255, 255, 255, 0.16);
    background: rgba(16, 18, 24, 0.78);
    color: #f3f4f6;
    font-size: 12px;
    line-height: 1.4;
    text-align: center;
    box-shadow: 0 10px 28px rgba(0, 0, 0, 0.2);
    backdrop-filter: blur(8px);
    pointer-events: none;
    white-space: pre-wrap;
}

.ai-notice--loading {
    top: 50%;
    width: min(72%, 560px);
    padding: 12px 16px;
    border-radius: 12px;
    font-size: 13px;
    line-height: 1.5;
    transform: translate(-50%, -50%);
    background: rgba(16, 18, 24, 0.92);
    box-shadow: 0 16px 40px rgba(0, 0, 0, 0.3);
    backdrop-filter: blur(10px);
}

.ai-notice--success {
    border-color: rgba(34, 197, 94, 0.45);
    background: rgba(10, 28, 18, 0.78);
}

.ai-notice--error {
    border-color: rgba(248, 113, 113, 0.45);
    background: rgba(40, 16, 16, 0.8);
}

.ai-notice-fade-enter-active,
.ai-notice-fade-leave-active {
    transition:
        opacity 180ms ease,
        transform 180ms ease;
}

.ai-mask-fade-enter-active,
.ai-mask-fade-leave-active {
    transition: opacity 180ms ease;
}

.ai-notice-fade-enter-from,
.ai-notice-fade-leave-to {
    opacity: 0;
    transform: translateX(-50%) translateY(-6px);
}

.ai-notice--loading.ai-notice-fade-enter-from,
.ai-notice--loading.ai-notice-fade-leave-to {
    transform: translate(-50%, calc(-50% + 8px));
}

.ai-mask-fade-enter-from,
.ai-mask-fade-leave-to {
    opacity: 0;
}

.terminal-search {
    position: absolute;
    top: 8px;
    right: 12px;
    z-index: 20;
    display: flex;
    align-items: center;
    gap: 4px;
    max-width: calc(100% - 24px);
    padding: 6px;
    border-radius: 8px;
    background: rgba(17, 24, 39, 0.94);
    border: 1px solid rgba(255, 255, 255, 0.12);
}

:deep(.terminal-link-tooltip) {
    position: absolute;
    z-index: 8;
    max-width: min(480px, 80%);
    padding: 4px 8px;
    border-radius: 4px;
    background: rgba(17, 24, 39, 0.94);
    color: #e5e7eb;
    font-size: 12px;
    line-height: 1.4;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    pointer-events: none;
}

:deep(.xterm) {
    height: 100%;
    padding: 0;
    background-color: var(--terminal-bg, #111827);
}

:deep(.xterm .xterm-viewport) {
    background-color: var(--terminal-bg, #111827) !important;
    scrollbar-width: none;
}

:deep(.xterm .xterm-viewport::-webkit-scrollbar) {
    width: 0;
    height: 0;
    display: none;
}
</style>
