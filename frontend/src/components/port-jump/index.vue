<template>
    <div></div>
</template>
<script lang="ts" setup>
import { getAgentSettingInfo } from '@/api/modules/setting';
import i18n from '@/lang';
import { MsgError, MsgWarning } from '@/utils/message';

interface DialogProps {
    port: any;
    ip: string;
    protocol?: string;
    path?: string;
    query?: string;
    hash?: string;
}

const acceptParams = async (params: DialogProps): Promise<void> => {
    if (Number(params.port) === 0) {
        MsgError(i18n.global.t('commons.msg.errPort'));
        return;
    }
    const res = await getAgentSettingInfo();
    const useWebAddress = !res.data.systemIP;
    const host = useWebAddress ? window.location.hostname : res.data.systemIP;
    if (!host) {
        MsgWarning(i18n.global.t('setting.systemIPWarning'));
        return;
    }
    let protocol = 'http';
    if (params.protocol === 'https') {
        protocol = 'https';
    } else if (params.protocol === 'http') {
        protocol = 'http';
    } else if (useWebAddress) {
        protocol = window.location.protocol === 'https:' ? 'https' : 'http';
    }
    const buildUrl = (targetHost: string) => {
        let url = `${protocol}://${targetHost}:${params.port}`;
        if (params.path) {
            url += params.path.startsWith('/') ? params.path : `/${params.path}`;
        }
        if (params.query) {
            url += params.query.startsWith('?') ? params.query : `?${params.query}`;
        }
        if (params.hash) {
            url += params.hash.startsWith('#') ? params.hash : `#${params.hash}`;
        }
        return url;
    };
    if (host.indexOf(':') === -1) {
        if (params.ip && params.ip === 'ipv6') {
            MsgWarning(i18n.global.t('setting.systemIPWarning1', ['IPv4']));
            return;
        }
        window.open(buildUrl(host), '_blank');
    } else {
        if (params.ip && params.ip === 'ipv4') {
            MsgWarning(i18n.global.t('setting.systemIPWarning1', ['IPv6']));
            return;
        }
        window.open(buildUrl(`[${host}]`), '_blank');
    }
};

defineExpose({ acceptParams });
</script>
