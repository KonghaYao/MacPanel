import assert from 'node:assert/strict';
import test from 'node:test';
import { Base64 } from 'js-base64';
import { createBase64StreamDecoder, decodeBase64 } from '../src/utils/base64.ts';

// The agent sends terminal output as independent base64 chunks of a raw byte stream, so a chunk can
// end in the middle of a multi-byte UTF-8 character.
const text = '中文终端';

function chunksOf(value, cuts) {
    const bytes = Base64.toUint8Array(Base64.encode(value));
    const chunks = [];
    let start = 0;
    for (const end of cuts) {
        chunks.push(Base64.fromUint8Array(bytes.slice(start, end)));
        start = end;
    }
    chunks.push(Base64.fromUint8Array(bytes.slice(start)));
    return chunks;
}

test('decodes a character split across two chunks', () => {
    const [first, second] = chunksOf(text, [7]);
    const decode = createBase64StreamDecoder();

    assert.equal(decode(first) + decode(second), text);
});

test('decodes a character split into single bytes', () => {
    const decode = createBase64StreamDecoder();
    const bytes = Base64.toUint8Array(Base64.encode(text));

    let decoded = '';
    for (const byte of bytes) {
        decoded += decode(Base64.fromUint8Array(Uint8Array.of(byte)));
    }

    assert.equal(decoded, text);
});

test('decodes chunks that end on a character boundary', () => {
    const [first, second] = chunksOf(text, [6]);
    const decode = createBase64StreamDecoder();

    assert.equal(decode(first), '中文');
    assert.equal(decode(second), '终端');
});

test('decoding every chunk on its own corrupts a split character', () => {
    const [first, second] = chunksOf(text, [7]);

    // what the terminal did before it decoded the stream as a whole
    assert.equal(decodeBase64(first).includes('\uFFFD'), true);
    assert.notEqual(decodeBase64(first) + decodeBase64(second), text);
});

test('returns an empty string for a payload that is not base64', () => {
    const decode = createBase64StreamDecoder();

    assert.equal(decode('not base64 ***'), '');
});
