export const getPanelHost = (): string => window.location.host || '';

export const applyPanelTitle = (): void => {
    const host = getPanelHost();
    if (host) {
        document.title = host;
    }
};
