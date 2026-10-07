import { StreamParser } from '@codemirror/language';

interface DotenvState {
    inValue: boolean;
    stringQuote: string | null;
}

export const dotenv: StreamParser<DotenvState> = {
    name: 'dotenv',

    startState(): DotenvState {
        return {
            inValue: false,
            stringQuote: null,
        };
    },

    token(stream, state) {
        if (state.stringQuote) {
            let escaped = false;
            let ch: string | void;
            while ((ch = stream.next())) {
                if (ch === state.stringQuote && !escaped) {
                    state.stringQuote = null;
                    return 'string';
                }
                escaped = !escaped && ch === '\\';
            }
            return 'string';
        }

        if (stream.sol()) {
            state.inValue = false;
            stream.eatSpace();
            if (stream.eol()) {
                return null;
            }

            const ch = stream.peek();
            if (ch === '#') {
                stream.skipToEnd();
                return 'comment';
            }
        }

        const ch = stream.next();
        if (!ch) {
            return null;
        }

        if (!state.inValue) {
            if (ch === 'e' && stream.match(/^xport(?=\s)/)) {
                stream.eatSpace();
                return 'keyword';
            }
            if (/[\w.-]/.test(ch as string)) {
                stream.eatWhile(/[\w.-]/);
                return 'property';
            }
            if (ch === '=') {
                state.inValue = true;
                return 'operator';
            }
            return null;
        }

        if (ch === '=') {
            return 'operator';
        }

        if (ch === '"' || ch === "'") {
            state.stringQuote = ch;
            return 'string';
        }

        if (ch === '#') {
            stream.skipToEnd();
            return 'comment';
        }

        stream.eatWhile(/[^\s#]/);
        return 'string';
    },

    languageData: {
        commentTokens: { line: '#' },
    },
};
