import http from '@/api';

export const fetchDesktopWallpaper = () => {
    return http.get<Blob>('/desktop/wallpaper', {
        responseType: 'blob',
        skipErrorMessage: true,
    });
};
