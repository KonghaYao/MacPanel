import i18n from '@/lang';
import { writeClipboardText } from '@/utils/clipboard-api';

import { MsgError, MsgSuccess } from '@/utils/message';

export async function copyText(content: string) {
    try {
        await writeClipboardText(content);
        MsgSuccess(i18n.global.t('commons.msg.copySuccess'));
    } catch (e) {
        MsgError(i18n.global.t('commons.msg.copyFailed'));
    }
}
