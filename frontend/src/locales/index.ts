import { createI18n } from 'vue-i18n';
import zhCN from './zh-CN';
import enUS from './en-US';

// 获取存储的语言，默认中文
const getDefaultLocale = () => {
    const stored = localStorage.getItem('locale');
    if (stored) return stored;

    // 根据浏览器语言判断
    const browserLang = navigator.language.toLowerCase();
    if (browserLang.startsWith('zh')) return 'zh-CN';
    return 'en-US';
};

const i18n = createI18n({
    legacy: false, // 使用 Composition API 模式
    locale: getDefaultLocale(),
    fallbackLocale: 'zh-CN',
    messages: {
        'zh-CN': zhCN,
        'en-US': enUS
    }
});

export default i18n;

// 切换语言
export const setLocale = (locale: string) => {
    (i18n.global.locale as any).value = locale;
    localStorage.setItem('locale', locale);
    document.querySelector('html')?.setAttribute('lang', locale);
};

// 获取当前语言
export const getLocale = () => {
    return (i18n.global.locale as any).value;
};
