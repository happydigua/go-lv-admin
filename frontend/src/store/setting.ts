import { defineStore } from 'pinia';
import { ref } from 'vue';
import { getPublicSettings } from '@/api/system/setting';

export const useSettingStore = defineStore('setting', () => {
    const siteName = ref('Go Lv Admin');
    const siteLogo = ref('');
    const siteFooter = ref('© 2024 Go Lv Admin');

    const fetchSettings = async () => {
        try {
            const res: any = await getPublicSettings();
            if (res) {
                siteName.value = res.site_name || 'Go Lv Admin';
                siteLogo.value = res.site_logo || '';
                siteFooter.value = res.site_footer || '© 2024 Go Lv Admin';

                // 更新网页标题
                document.title = siteName.value;
            }
        } catch (error) {
            console.error('Failed to fetch settings:', error);
        }
    };

    return {
        siteName,
        siteLogo,
        siteFooter,
        fetchSettings
    };
});
