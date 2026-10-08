<template>
    <div>
        <div class="flex flex-wrap items-center">
            <div class="flex flex-wrap items-center">
                <el-link v-if="isEE" underline="never" type="primary" @click="toEdition">
                    {{ $t('license.ee') }}
                </el-link>
                <el-link v-else-if="isMasterPro" underline="never" type="primary" @click="toLxware">
                    {{ $t('license.pro') }}
                </el-link>
                <el-link v-else-if="isOffline" underline="never" type="primary" @click="to1Panel">
                    {{ $t('license.offLine') }}
                </el-link>
                <el-link
                    v-if="showVersion"
                    underline="never"
                    class="version"
                    :class="{ 'version-after-edition': isEE || isMasterPro || isOffline }"
                    type="primary"
                    @click="getVersionLog()"
                >
                    {{ version }}
                </el-link>
                <el-link
                    v-if="isAdmin && !isOffline && !isEE"
                    :class="{ 'ml-2': showVersion || isEE || isMasterPro || isOffline }"
                    underline="never"
                    type="primary"
                    @click="onUpgradeByMise"
                >
                    {{ $t('commons.button.update') }}
                </el-link>
                <el-tag v-if="version === 'Waiting'" round class="ml-2.5">{{ $t('setting.upgrading') }}</el-tag>
            </div>
        </div>

        <Releases ref="releasesRef" />
    </div>
</template>

<script setup lang="ts">
import { getSettingBaseInfo, upgradeByMise } from '@/api/modules/setting';
import Releases from '@/components/system-upgrade/releases/index.vue';
import i18n from '@/lang';
import { MsgSuccess } from '@/utils/message';
import { onMounted, ref } from 'vue';
import { useGlobalStore } from '@/composables/useGlobalStore';
import { ElMessageBox } from 'element-plus';

withDefaults(
    defineProps<{
        showVersion?: boolean;
    }>(),
    {
        showVersion: true,
    },
);

const { isOffline, isMasterPro, isEE, isIntl, isAdmin } = useGlobalStore();
const releasesRef = ref();

const version = ref<string>('');

const search = async () => {
    const res = await getSettingBaseInfo();
    version.value = res.data.systemVersion;
};

const getVersionLog = () => {
    if (isOffline.value) {
        return;
    }
    releasesRef.value.acceptParams();
};

const toLxware = () => {
    if (!isIntl.value) {
        window.open('https://www.lxware.cn/1panel' + '', '_blank', 'noopener,noreferrer');
    } else {
        window.open('https://1panel.pro/pricing' + '', '_blank', 'noopener,noreferrer');
    }
};

const to1Panel = () => {
    let url = isIntl.value ? 'https://1panel.pro' : 'https://1panel.cn';
    window.open(url, '_blank', 'noopener,noreferrer');
};

const toEdition = () => {
    if (!isIntl.value) {
        window.open('https://1panel.cn/versions.html' + '', '_blank', 'noopener,noreferrer');
    } else {
        window.open('https://1panel.pro/pricing' + '', '_blank', 'noopener,noreferrer');
    }
};

const onUpgradeByMise = async () => {
    try {
        await ElMessageBox.confirm(
            'mise use -g "github:KonghaYao/MacPanel[bin=macpanel]@latest"\nmise reshim',
            i18n.global.t('commons.button.update'),
            {
                confirmButtonText: i18n.global.t('commons.button.confirm'),
                cancelButtonText: i18n.global.t('commons.button.cancel'),
                type: 'info',
            },
        );
    } catch {
        return;
    }
    try {
        const res = await upgradeByMise();
        const output = res.data?.output?.trim() ?? '';
        if (output) {
            await ElMessageBox.alert(output, i18n.global.t('commons.button.update'), {
                confirmButtonText: i18n.global.t('commons.button.confirm'),
            });
            return;
        }
        MsgSuccess(i18n.global.t('commons.msg.operationSuccess'));
    } catch {
        return;
    }
};

onMounted(() => {
    search();
});
</script>

<style lang="scss" scoped>
.line-height {
    line-height: 25px;
}
:deep(.el-link__inner) {
    font-weight: 400;
}
.version {
    margin-left: 0;
    font-size: 14px;
    color: var(--panel-color-primary-light-4);
    text-decoration: none;
    letter-spacing: 0.5px;
    cursor: pointer;
    font-family: auto;
}
.version-after-edition {
    margin-left: 8px;
}
</style>
