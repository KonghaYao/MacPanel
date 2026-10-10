import type { ILinkHandler, ITerminalOptions, ITheme, Terminal } from '@xterm/xterm';
import { Unicode11Addon } from '@xterm/addon-unicode11';
import { WebglAddon } from '@xterm/addon-webgl';
import { ClipboardAddon, type IClipboardProvider } from '@xterm/addon-clipboard';
import { SearchAddon, type ISearchOptions } from '@xterm/addon-search';
import { WebLinksAddon } from '@xterm/addon-web-links';
import { canReadClipboard, readClipboardText, writeClipboardText } from '@/utils/clipboard-api';

export const TERMINAL_FONT_FAMILY = "'JetBrains Mono', monospace";

export const TERMINAL_DEFAULTS = {
    fontSize: 14,
    lineHeight: 1.1,
    letterSpacing: 0,
    scrollSensitivity: 1,
    cursorStyle: 'block' as const,
    cursorBlink: 'Enable',
    fontFamily: TERMINAL_FONT_FAMILY,
    backgroundColor: '#111827',
    foregroundColor: '#e5e7eb',
    scrollback: 1000,
};

export const TERMINAL_ANSI_THEME: ITheme = {
    background: '#111827',
    foreground: '#e5e7eb',
    cursor: '#e5e7eb',
    cursorAccent: '#111827',
    selectionBackground: 'rgba(102, 178, 255, 0.30)',
    selectionInactiveBackground: 'rgba(102, 178, 255, 0.20)',
    black: '#111827',
    red: '#f87171',
    green: '#34d399',
    yellow: '#fbbf24',
    blue: '#60a5fa',
    magenta: '#c084fc',
    cyan: '#22d3ee',
    white: '#e5e7eb',
    brightBlack: '#6b7280',
    brightRed: '#fca5a5',
    brightGreen: '#6ee7b7',
    brightYellow: '#fde68a',
    brightBlue: '#93c5fd',
    brightMagenta: '#d8b4fe',
    brightCyan: '#67e8f9',
    brightWhite: '#f9fafb',
};

export const HTTP_URL_REGEX = /\bhttps?:\/\/[^\s<>"'`]+/;

export const TERMINAL_PREVIEW_SAMPLE = [
    '\x1b[1mMacPanel\x1b[0m  \x1b[31mred \x1b[32mgreen \x1b[33myellow \x1b[34mblue \x1b[35mmagenta \x1b[36mcyan\x1b[0m',
    '\x1b[90mgray \x1b[91mred \x1b[92mgreen \x1b[93myellow \x1b[94mblue \x1b[95mmagenta \x1b[96mcyan \x1b[97mwhite\x1b[0m',
    '┌──────┬──────┐',
    '│ box  │ 线框 │',
    '└──────┴──────┘',
    '中文示例 Chinese sample',
].join('\r\n');

export interface TerminalOptionSettings {
    fontSize: number;
    lineHeight: number;
    letterSpacing: number;
    fontFamily: string;
    backgroundColor: string;
    foregroundColor: string;
    cursorBlink: boolean;
    cursorStyle: 'block' | 'underline' | 'bar';
    scrollback: number;
    scrollSensitivity: number;
}

export interface TerminalLinkTarget {
    getElement(): HTMLElement | null | undefined;
}

export function terminalNumber(value: unknown, fallback: number): number {
    if (value === '' || value === null || value === undefined) {
        return fallback;
    }
    const parsed = Number(value);
    return Number.isFinite(parsed) ? parsed : fallback;
}

export function cursorStyleOf(value: string | undefined): 'block' | 'underline' | 'bar' {
    if (value === 'underline' || value === 'bar' || value === 'block') {
        return value;
    }
    return 'block';
}

export function buildTerminalTheme(background?: string, foreground?: string): ITheme {
    const nextBackground = background || TERMINAL_ANSI_THEME.background;
    const nextForeground = foreground || TERMINAL_ANSI_THEME.foreground;
    return {
        ...TERMINAL_ANSI_THEME,
        background: nextBackground,
        foreground: nextForeground,
        cursor: nextForeground,
        cursorAccent: nextBackground,
    };
}

export function isHttpUrl(text: string): boolean {
    try {
        const url = new URL(text);
        return url.protocol === 'http:' || url.protocol === 'https:';
    } catch {
        return false;
    }
}

function openHttpUrl(text: string) {
    if (!isHttpUrl(text)) {
        return;
    }
    window.open(text, '_blank', 'noopener,noreferrer');
}

function tooltipElement(root: HTMLElement): HTMLDivElement {
    const existing = root.querySelector<HTMLDivElement>('.terminal-link-tooltip');
    if (existing) {
        return existing;
    }
    const tip = document.createElement('div');
    tip.className = 'xterm-hover terminal-link-tooltip';
    root.appendChild(tip);
    return tip;
}

export function showLinkTooltip(root: HTMLElement | null | undefined, event: MouseEvent, text: string) {
    if (!root || !isHttpUrl(text)) {
        return;
    }
    const tip = tooltipElement(root);
    const rect = root.getBoundingClientRect();
    tip.textContent = text;
    tip.style.display = 'block';
    tip.style.left = `${event.clientX - rect.left + 8}px`;
    tip.style.top = `${event.clientY - rect.top + 16}px`;
}

export function hideLinkTooltip(root: HTMLElement | null | undefined) {
    const tip = root?.querySelector<HTMLDivElement>('.terminal-link-tooltip');
    if (tip) {
        tip.style.display = 'none';
    }
}

export function createLinkHandler(target?: TerminalLinkTarget): ILinkHandler {
    return {
        allowNonHttpProtocols: false,
        activate(_event, text) {
            openHttpUrl(text);
        },
        hover(event, text) {
            showLinkTooltip(target?.getElement(), event, text);
        },
        leave() {
            hideLinkTooltip(target?.getElement());
        },
    };
}

export function buildTerminalOptions(
    settings: TerminalOptionSettings,
    target?: TerminalLinkTarget,
    extra?: ITerminalOptions,
): ITerminalOptions {
    return {
        fontSize: settings.fontSize,
        lineHeight: settings.lineHeight,
        letterSpacing: settings.letterSpacing,
        fontFamily: settings.fontFamily || TERMINAL_FONT_FAMILY,
        fontWeight: 'normal',
        theme: buildTerminalTheme(settings.backgroundColor, settings.foregroundColor),
        cursorBlink: settings.cursorBlink,
        cursorStyle: settings.cursorStyle,
        scrollback: settings.scrollback,
        scrollSensitivity: settings.scrollSensitivity,
        rightClickSelectsWord: false,
        macOptionClickForcesSelection: true,
        macOptionIsMeta: false,
        windowOptions: {
            getWinSizePixels: true,
            getCellSizePixels: true,
            getWinSizeChars: true,
            pushTitle: true,
            popTitle: true,
        },
        linkHandler: createLinkHandler(target),
        ...extra,
        allowProposedApi: true,
    };
}

export function loadUnicode11(term: Terminal) {
    term.loadAddon(new Unicode11Addon());
    term.unicode.activeVersion = '11';
}

export function loadWebgl(term: Terminal) {
    term.loadAddon(new WebglAddon());
}

/**
 * Clipboard provider for the OSC 52 addon. Programs inside the terminal (tmux, vim, ...) use it to
 * copy to the system clipboard. The bundled provider calls `navigator.clipboard` directly, which
 * does not exist on plain HTTP origins, so route it through the terminal clipboard helpers instead.
 */
class TerminalClipboardProvider implements IClipboardProvider {
    readText(selection: string): string | Promise<string> {
        if (selection !== 'c' || !canReadClipboard()) {
            return '';
        }
        return readClipboardText();
    }

    async writeText(selection: string, text: string): Promise<void> {
        if (selection !== 'c' || !text) {
            return;
        }
        try {
            await writeClipboardText(text);
        } catch {
            // an OSC 52 write has no user gesture behind it, browsers may refuse it outside HTTPS
        }
    }
}

export function loadClipboard(term: Terminal) {
    term.loadAddon(new ClipboardAddon(undefined, new TerminalClipboardProvider()));
}

export function loadWebLinks(term: Terminal, target?: TerminalLinkTarget) {
    term.loadAddon(
        new WebLinksAddon((_event, uri) => openHttpUrl(uri), {
            urlRegex: HTTP_URL_REGEX,
            hover(event, uri) {
                showLinkTooltip(target?.getElement() ?? term.element, event, uri);
            },
            leave() {
                hideLinkTooltip(target?.getElement() ?? term.element);
            },
        }),
    );
}

export function createSearchAddon() {
    return new SearchAddon({ highlightLimit: 1000 });
}

export function searchDecorationOptions(): ISearchOptions {
    return {
        decorations: {
            matchBackground: '#1e3a5f',
            matchBorder: '#60a5fa',
            matchOverviewRuler: '#60a5fa',
            activeMatchBackground: '#2563eb',
            activeMatchBorder: '#93c5fd',
            activeMatchColorOverviewRuler: '#93c5fd',
        },
    };
}

export function patchTerminalStore(
    store: { $patch: (partial: object) => void },
    data: {
        showTerminalButton?: string;
        lineHeight?: string | number;
        letterSpacing?: string | number;
        fontSize?: string | number;
        fontFamily?: string;
        backgroundColor?: string;
        foregroundColor?: string;
        cursorBlink?: string;
        cursorStyle?: string;
        scrollback?: string | number;
        scrollSensitivity?: string | number;
    },
) {
    store.$patch({
        showTerminalButton: data.showTerminalButton !== 'Disable',
        lineHeight: terminalNumber(data.lineHeight, TERMINAL_DEFAULTS.lineHeight),
        letterSpacing: terminalNumber(data.letterSpacing, TERMINAL_DEFAULTS.letterSpacing),
        fontSize: terminalNumber(data.fontSize, TERMINAL_DEFAULTS.fontSize),
        fontFamily: data.fontFamily || TERMINAL_DEFAULTS.fontFamily,
        backgroundColor: data.backgroundColor || TERMINAL_DEFAULTS.backgroundColor,
        foregroundColor: data.foregroundColor || TERMINAL_DEFAULTS.foregroundColor,
        cursorBlink: data.cursorBlink,
        cursorStyle: data.cursorStyle || TERMINAL_DEFAULTS.cursorStyle,
        scrollback: terminalNumber(data.scrollback, TERMINAL_DEFAULTS.scrollback),
        scrollSensitivity: terminalNumber(data.scrollSensitivity, TERMINAL_DEFAULTS.scrollSensitivity),
    });
}
