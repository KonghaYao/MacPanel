import { getPlatformCapabilities } from '@/api/modules/platform';
import { GlobalStore } from '@/store';

let syncPromise: Promise<void> | null = null;
let loaded = false;

const sameFeatures = (left: Record<string, boolean>, right: Record<string, boolean>) => {
    const leftKeys = Object.keys(left);
    const rightKeys = Object.keys(right);
    if (leftKeys.length !== rightKeys.length) {
        return false;
    }
    return leftKeys.every((key) => left[key] === right[key]);
};

export const resetPlatformCapabilitiesSync = () => {
    loaded = false;
    syncPromise = null;
};

export const syncPlatformCapabilities = async (options?: { force?: boolean }) => {
    if (loaded && !options?.force) {
        return;
    }
    if (syncPromise) {
        return syncPromise;
    }

    syncPromise = (async () => {
        const globalStore = GlobalStore();
        try {
            const res = await getPlatformCapabilities();
            const os = res.data?.os || '';
            const features = res.data?.features || {};
            if (globalStore.platformOS !== os || !sameFeatures(globalStore.platformFeatures, features)) {
                globalStore.setPlatformCapabilities({ os, features });
            }
            globalStore.setPlatformCapabilitiesLoaded(true);
            loaded = true;
        } catch {
            if (globalStore.platformOS !== '' || Object.keys(globalStore.platformFeatures).length > 0) {
                globalStore.setPlatformCapabilities({ os: '', features: {} });
            }
            globalStore.setPlatformCapabilitiesLoaded(false);
            loaded = false;
        }
    })().finally(() => {
        syncPromise = null;
    });

    return syncPromise;
};

export const hasPlatformFeature = (feature?: string) => {
    if (!feature) {
        return true;
    }
    const globalStore = GlobalStore();
    return !!globalStore.platformFeatures[feature];
};
