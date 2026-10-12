import assert from 'node:assert/strict';
import test from 'node:test';
import { applyPanelTitle, getPanelHost } from '../src/utils/panel-title.ts';

const originalWindow = Object.getOwnPropertyDescriptor(globalThis, 'window');
const originalDocument = Object.getOwnPropertyDescriptor(globalThis, 'document');

function stubHost(host, title = '') {
    const doc = { title };
    Object.defineProperty(globalThis, 'window', { value: { location: { host } }, configurable: true, writable: true });
    Object.defineProperty(globalThis, 'document', { value: doc, configurable: true, writable: true });
    return doc;
}

function restore() {
    if (originalWindow) Object.defineProperty(globalThis, 'window', originalWindow);
    else delete globalThis.window;
    if (originalDocument) Object.defineProperty(globalThis, 'document', originalDocument);
    else delete globalThis.document;
}

test('uses the access address as the tab title', () => {
    const doc = stubHost('192.168.1.20:9999');
    applyPanelTitle();
    assert.equal(doc.title, '192.168.1.20:9999');
    assert.equal(getPanelHost(), '192.168.1.20:9999');
    restore();
});

test('keeps the port so different environments differ', () => {
    const doc = stubHost('localhost:4004');
    applyPanelTitle();
    assert.equal(doc.title, 'localhost:4004');
    restore();
});

test('leaves the title untouched when there is no host', () => {
    const doc = stubHost('', 'loading...');
    applyPanelTitle();
    assert.equal(doc.title, 'loading...');
    restore();
});
