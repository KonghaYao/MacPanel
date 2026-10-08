<template>
    <div class="footer" :class="{ 'footer--mobile': isMobile }">
        <div class="footer-content">
            <div class="footer-copyright">
                <span>
                    Copyright © 2014-{{ year }} {{ $t('commons.fit2cloud') }} · MacPanel{{
                        version ? ` ${version}` : ''
                    }}
                </span>
            </div>
            <FooterNavigation class="footer-navigation-panel" />
            <SystemUpgrade class="footer-upgrade" :show-version="false" />
        </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import SystemUpgrade from '@/components/system-upgrade/index.vue';
import FooterNavigation from '@/components/footer-navigation/index.vue';
import { getSettingBaseInfo } from '@/api/modules/setting';
import { useGlobalStore } from '@/composables/useGlobalStore';

const { isMobile } = useGlobalStore();

const year = new Date().getFullYear();
const version = ref('');

onMounted(async () => {
    const res = await getSettingBaseInfo();
    version.value = res.data.systemVersion;
});
</script>

<style scoped lang="scss">
.footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    min-height: 48px;
    height: auto;
    background: var(--panel-footer-bg);
    border-top: 1px solid var(--panel-footer-border);
    box-sizing: border-box;
    padding: 10px 12px;
    a {
        font-size: 12px;
        color: #858585;
        text-decoration: none;
        letter-spacing: 0.5px;
    }
    span {
        font-size: 12px;
        color: #858585;
        text-decoration: none;
        letter-spacing: 0.5px;
    }
}

.footer--mobile {
    min-height: 76px;
}

.footer-content {
    display: flex;
    width: 100%;
    flex-wrap: wrap;
    align-items: center;
    justify-content: center;
    gap: 8px 12px;
}

.footer-copyright {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
}

.footer-navigation-panel {
    order: -1;
    flex-basis: 100%;
}

.footer-upgrade {
    white-space: nowrap;
}

@media (min-width: 768px) {
    .footer-content {
        flex-wrap: nowrap;
        justify-content: flex-start;
        gap: 0;
    }

    .footer {
        padding: 10px 20px;
    }

    .footer-copyright {
        justify-content: flex-start;
    }

    .footer-navigation-panel {
        order: 0;
        flex-basis: auto;
        margin-left: auto;
    }
}
</style>
