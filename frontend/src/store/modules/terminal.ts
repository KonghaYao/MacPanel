import { defineStore } from 'pinia';
import type { StoreDefinition } from 'pinia';
import piniaPersistConfig from '@/config/pinia-persist';
import { TerminalState } from '../interface';
import { TERMINAL_DEFAULTS } from '@/components/terminal/options';

export const TerminalStore = defineStore('TerminalState', {
    state: (): TerminalState => ({
        showTerminalButton: true,
        lineHeight: TERMINAL_DEFAULTS.lineHeight,
        letterSpacing: TERMINAL_DEFAULTS.letterSpacing,
        fontSize: TERMINAL_DEFAULTS.fontSize,
        fontFamily: TERMINAL_DEFAULTS.fontFamily,
        backgroundColor: TERMINAL_DEFAULTS.backgroundColor,
        foregroundColor: TERMINAL_DEFAULTS.foregroundColor,
        cursorBlink: TERMINAL_DEFAULTS.cursorBlink,
        cursorStyle: TERMINAL_DEFAULTS.cursorStyle,
        scrollback: TERMINAL_DEFAULTS.scrollback,
        scrollSensitivity: TERMINAL_DEFAULTS.scrollSensitivity,
    }),
    persist: piniaPersistConfig('TerminalState'),
}) as StoreDefinition<'TerminalState', TerminalState, any, any>;

export default TerminalStore;
