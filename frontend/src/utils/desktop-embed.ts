const DESKTOP_EMBED_KEY = 'macpanel-desktop-embed';

type EmbedQuery = {
    embed?: unknown;
};

const hasEmbedQuery = (embed: unknown) => {
    if (embed === '1') {
        return true;
    }
    return Array.isArray(embed) && embed.includes('1');
};

export const markDesktopEmbed = (to: { query: EmbedQuery }) => {
    if (window.self === window.top) {
        return;
    }
    if (hasEmbedQuery(to.query.embed)) {
        sessionStorage.setItem(DESKTOP_EMBED_KEY, '1');
    }
};

export const isDesktopEmbed = () => {
    if (typeof window === 'undefined' || window.self === window.top) {
        return false;
    }
    if (sessionStorage.getItem(DESKTOP_EMBED_KEY) === '1') {
        return true;
    }
    return new URLSearchParams(window.location.search).get('embed') === '1';
};

export const desktopEmbedSrc = (route: string) => {
    const url = new URL(route, window.location.origin);
    url.searchParams.set('embed', '1');
    return `${url.pathname}${url.search}${url.hash}`;
};
