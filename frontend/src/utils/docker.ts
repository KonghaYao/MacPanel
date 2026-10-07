import i18n from '@/lang';

export interface ParsedImageReference {
    name: string;
    tag: string;
    full: string;
}

function findTagSeparator(ref: string): number {
    const colon = ref.lastIndexOf(':');
    if (colon <= 0) {
        return -1;
    }
    const slash = ref.lastIndexOf('/');
    if (colon <= slash) {
        return -1;
    }
    const firstSlash = ref.indexOf('/');
    if (firstSlash > 0 && colon < firstSlash) {
        const nextColon = ref.indexOf(':', firstSlash);
        return nextColon >= 0 ? nextColon : -1;
    }
    return colon;
}

export function parseImageReference(ref: string): ParsedImageReference {
    const full = (ref || '').trim();
    if (!full || full.includes('<none>')) {
        return { name: full, tag: '', full };
    }
    const tagSep = findTagSeparator(full);
    if (tagSep < 0) {
        return { name: full, tag: '', full };
    }
    return {
        name: full.slice(0, tagSep),
        tag: full.slice(tagSep + 1),
        full,
    };
}

export function getPrimaryImageName(tags: string[]): string {
    const valid = (tags || []).filter((tag) => tag && !tag.includes('<none>'));
    if (valid.length === 0) {
        return '-';
    }
    const names = [...new Set(valid.map((tag) => parseImageReference(tag).name))];
    return names[0];
}

export function getImageTagLabel(ref: string): string {
    const parsed = parseImageReference(ref);
    return parsed.tag || parsed.full;
}

export function getDockerRestartHelper(dockerRuntime?: string, platformOS?: string): string {
    if (platformOS === 'darwin') {
        if (dockerRuntime === 'orbstack') {
            return i18n.global.t('container.restartHelperOrbStack');
        }
        return i18n.global.t('container.restartHelperDockerDesktop');
    }
    return i18n.global.t('container.restartHelper');
}
