<template>
    <el-button v-if="isIcon" link @click="copyText(content)" icon="DocumentCopy" class="ml-1.5"></el-button>
    <el-button @click="copyText(content)" v-else>{{ label || $t('commons.button.copy') }}</el-button>
</template>

<script lang="ts" setup>
import { copyText } from '@/utils/clipboard';
defineOptions({ name: 'CopyButton' });

defineProps({
    content: String,
    label: String,
    isIcon: {
        type: Boolean,
        default: true,
    },
});
</script>
