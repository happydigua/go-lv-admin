import type { App, Directive } from 'vue';
import { useUserStore } from '@/store/user';

// 按钮级权限指令
// 使用方式: v-permission="'system:user:add'" 或 v-permission="['system:user:add', 'system:user:edit']"
export const permission: Directive = {
    mounted(el, binding) {
        const userStore = useUserStore();
        const permissions = userStore.permissions;
        const value = binding.value;

        // admin 拥有所有权限
        if (permissions.includes('*')) {
            return;
        }

        // 没有传值则不需要权限
        if (!value) {
            return;
        }

        // 单个权限
        if (typeof value === 'string') {
            if (!permissions.includes(value)) {
                el.parentNode?.removeChild(el);
            }
            return;
        }

        // 多个权限（满足其一即可）
        if (Array.isArray(value)) {
            const hasPermission = value.some((p: string) => permissions.includes(p));
            if (!hasPermission) {
                el.parentNode?.removeChild(el);
            }
        }
    }
};

// 注册指令
export function setupPermissionDirective(app: App) {
    app.directive('permission', permission);
}
