import http from '@/api';

export interface PlatformCapabilities {
    os: string;
    arch: string;
    features: Record<string, boolean>;
}

export const getPlatformCapabilities = () => {
    return http.get<PlatformCapabilities>('/platform/capabilities');
};
