import assert from 'node:assert/strict';
import test from 'node:test';
import { canReadClipboard, readClipboardText, writeClipboardText } from '../src/utils/clipboard-api.ts';

const originalNavigator = Object.getOwnPropertyDescriptor(globalThis, 'navigator');
const originalDocument = Object.getOwnPropertyDescriptor(globalThis, 'document');

function stubNavigator(clipboard) {
    Object.defineProperty(globalThis, 'navigator', { value: { clipboard }, configurable: true, writable: true });
}

function stubDocument(execCommand) {
    const calls = { selected: 0, removed: 0, copied: 0, focused: 0 };
    const textarea = {
        value: '',
        style: {},
        setAttribute() {},
        select() {
            calls.selected += 1;
        },
        setSelectionRange() {},
        remove() {
            calls.removed += 1;
        },
    };
    const doc = {
        activeElement: {
            focus: () => {
                calls.focused += 1;
            },
        },
        createElement: () => textarea,
        body: { appendChild: () => {} },
        getSelection: () => null,
        execCommand: () => {
            calls.copied += 1;
            return execCommand();
        },
    };
    Object.defineProperty(globalThis, 'document', { value: doc, configurable: true, writable: true });
    return { calls, textarea };
}

function stubMissingDocument() {
    Object.defineProperty(globalThis, 'document', { value: undefined, configurable: true, writable: true });
}

test.afterEach(() => {
    const originals = [
        ['navigator', originalNavigator],
        ['document', originalDocument],
    ];
    for (const [name, descriptor] of originals) {
        if (descriptor) {
            Object.defineProperty(globalThis, name, descriptor);
        } else {
            delete globalThis[name];
        }
    }
});

test('reports the clipboard API as unavailable without navigator.clipboard', async () => {
    stubNavigator(undefined);

    assert.equal(canReadClipboard(), false);
    await assert.rejects(readClipboardText(), /clipboard read is unavailable/);
});

test('reads through the Clipboard API when it exists', async () => {
    stubNavigator({ readText: async () => 'from clipboard' });

    assert.equal(canReadClipboard(), true);
    assert.equal(await readClipboardText(), 'from clipboard');
});

test('writes through the Clipboard API when it exists', async () => {
    const written = [];
    stubNavigator({
        writeText: async (text) => {
            written.push(text);
        },
    });
    const { calls } = stubDocument(() => {
        throw new Error('execCommand must not be used when the Clipboard API works');
    });

    await writeClipboardText('中文 text');

    assert.deepEqual(written, ['中文 text']);
    assert.equal(calls.selected, 0);
});

test('falls back to the textarea copy path when the Clipboard API is missing', async () => {
    stubNavigator(undefined);
    const { calls, textarea } = stubDocument(() => true);

    await writeClipboardText('中文 text');

    assert.equal(textarea.value, '中文 text');
    assert.equal(calls.selected, 1);
    assert.equal(calls.removed, 1);
    assert.equal(calls.focused, 1);
});

test('falls back to the textarea copy path when the Clipboard API rejects', async () => {
    stubNavigator({
        writeText: async () => {
            throw new Error('not allowed');
        },
    });
    const { calls, textarea } = stubDocument(() => true);

    await writeClipboardText('line1\nline2');

    assert.equal(textarea.value, 'line1\nline2');
    assert.equal(calls.copied, 1);
});

test('skips empty content', async () => {
    stubNavigator(undefined);
    const { calls } = stubDocument(() => true);

    await writeClipboardText('');

    assert.equal(calls.selected, 0);
});

test('rejects when neither the Clipboard API nor execCommand can copy', async () => {
    stubNavigator(undefined);
    stubDocument(() => false);

    await assert.rejects(writeClipboardText('abc'), /clipboard write is unavailable/);
});

test('rejects when there is no DOM to fall back to', async () => {
    stubNavigator(undefined);
    stubMissingDocument();

    await assert.rejects(writeClipboardText('abc'), /clipboard write is unavailable/);
});
