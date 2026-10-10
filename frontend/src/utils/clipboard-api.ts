/**
 * Clipboard read/write helpers used by every copy/paste entry point of the panel.
 *
 * `navigator.clipboard` (the async Clipboard API) only exists in secure contexts, so it is
 * undefined when the panel is opened over `http://<lan-ip>`. The helpers therefore fall back to the
 * clipboard APIs that work on plain HTTP origins: the native copy/paste events (terminal shortcuts
 * rely on those) and a hidden textarea + `document.execCommand('copy')` for buttons.
 */

export function canReadClipboard(): boolean {
    return typeof navigator !== 'undefined' && typeof navigator.clipboard?.readText === 'function';
}

function canWriteClipboard(): boolean {
    return typeof navigator !== 'undefined' && typeof navigator.clipboard?.writeText === 'function';
}

export async function readClipboardText(): Promise<string> {
    if (!canReadClipboard()) {
        throw new Error('clipboard read is unavailable in this browser context');
    }
    return navigator.clipboard.readText();
}

export async function writeClipboardText(text: string): Promise<void> {
    if (!text) {
        return;
    }
    if (canWriteClipboard()) {
        try {
            await navigator.clipboard.writeText(text);
            return;
        } catch {
            // denied permission or an unfocused document, keep going with the legacy path
        }
    }
    if (!copyWithExecCommand(text)) {
        throw new Error('clipboard write is unavailable in this browser context');
    }
}

function copyWithExecCommand(text: string): boolean {
    if (typeof document === 'undefined' || typeof document.execCommand !== 'function') {
        return false;
    }
    // selecting the textarea steals the focus, put it back where it was afterwards
    const focused = document.activeElement as HTMLElement | null;
    const textarea = document.createElement('textarea');
    textarea.value = text;
    textarea.setAttribute('readonly', 'readonly');
    textarea.style.position = 'fixed';
    textarea.style.top = '-1000px';
    textarea.style.left = '-1000px';
    textarea.style.opacity = '0';
    document.body.appendChild(textarea);

    const selection = document.getSelection();
    const ranges: Range[] = [];
    if (selection) {
        for (let index = 0; index < selection.rangeCount; index++) {
            ranges.push(selection.getRangeAt(index).cloneRange());
        }
    }

    let copied = false;
    try {
        textarea.select();
        textarea.setSelectionRange(0, text.length);
        copied = document.execCommand('copy');
    } catch {
        copied = false;
    }
    textarea.remove();

    // restore the previous selection so the text that was copied stays highlighted
    if (selection && ranges.length > 0) {
        selection.removeAllRanges();
        ranges.forEach((range) => selection.addRange(range));
    }
    focused?.focus?.();
    return copied;
}
