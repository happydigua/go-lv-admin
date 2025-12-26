import { defineStore } from 'pinia';
import { ref } from 'vue';
import { login } from '@/api/user';
import { getUserMenus, getUserPermissions } from '@/api/permission';

export const useUserStore = defineStore('user', () => {
    const token = ref(localStorage.getItem('token') || '');
    const userInfo = ref<any>(null);
    const menus = ref<any[]>([]);
    const permissions = ref<string[]>([]);

    const handleLogin = async (loginForm: any) => {
        try {
            const res: any = await login(loginForm);
            token.value = res.token;
            userInfo.value = res.user;
            localStorage.setItem('token', res.token);
            // 登录成功后获取菜单和权限
            await fetchMenus();
            await fetchPermissions();
            return true;
        } catch (error) {
            return false;
        }
    };

    const fetchMenus = async () => {
        try {
            const res: any = await getUserMenus();
            menus.value = res || [];
        } catch (error) {
            console.error('Failed to fetch menus:', error);
        }
    };

    const fetchPermissions = async () => {
        try {
            const res: any = await getUserPermissions();
            permissions.value = res || [];
        } catch (error) {
            console.error('Failed to fetch permissions:', error);
        }
    };

    const logout = () => {
        token.value = '';
        userInfo.value = null;
        menus.value = [];
        permissions.value = [];
        localStorage.removeItem('token');
    };

    return {
        token,
        userInfo,
        menus,
        permissions,
        handleLogin,
        fetchMenus,
        fetchPermissions,
        logout
    };
});
