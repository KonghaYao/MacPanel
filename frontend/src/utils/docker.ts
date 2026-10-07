import i18n from '@/lang';

export function getDockerRestartHelper(dockerRuntime?: string, platformOS?: string): string {
    if (platformOS === 'darwin') {
        if (dockerRuntime === 'orbstack') {
            return i18n.global.t('container.restartHelperOrbStack');
        }
        return i18n.global.t('container.restartHelperDockerDesktop');
    }
    return i18n.global.t('container.restartHelper');
}
