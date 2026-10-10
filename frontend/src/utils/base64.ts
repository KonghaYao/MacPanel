import { Base64 } from 'js-base64';

type StringLikeRecord = object;

export const encodeBase64 = (value: string) => Base64.encode(value);

export const decodeBase64 = (value: string) => Base64.decode(value);

export type Base64StreamDecoder = (value: string) => string;

// Chunked payloads (the terminal output stream) must be decoded as one byte stream: decoding each
// chunk on its own turns a multi-byte UTF-8 character split over two chunks into replacement chars.
export const createBase64StreamDecoder = (): Base64StreamDecoder => {
    const decoder = new TextDecoder();
    return (value: string) => {
        try {
            return decoder.decode(Base64.toUint8Array(value), { stream: true });
        } catch {
            return '';
        }
    };
};

export const encodeBase64Fields = <T extends StringLikeRecord>(target: T, fields: Array<Extract<keyof T, string>>) => {
    fields.forEach((field) => {
        const record = target as Record<string, unknown>;
        const value = record[field];
        if (typeof value === 'string' && value) {
            record[field] = encodeBase64(value);
        }
    });

    return target;
};
