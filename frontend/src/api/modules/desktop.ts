import http from '@/api';

export const fetchDesktopWallpaper = () => {
    return http.get<Blob>(
        '/desktop/wallpaper',
        { _t: Date.now() },
        {
            responseType: 'blob',
            skipErrorMessage: true,
        },
    );
};
