import request from '@/utils/request';

// 获取角色的菜单权限
export const getRoleMenus = (roleId: number) => {
    return request({
        url: `/system/role/${roleId}/menus`,
        method: 'get',
    });
};

// 设置角色的菜单权限
export const setRoleMenus = (roleId: number, menuIds: number[]) => {
    return request({
        url: `/system/role/${roleId}/menus`,
        method: 'put',
        data: { menuIds },
    });
};

// 获取当前用户的按钮权限
export const getUserPermissions = () => {
    return request({
        url: '/user/permissions',
        method: 'get',
    });
};

// 获取当前用户可访问的菜单
export const getUserMenus = () => {
    return request({
        url: '/user/menus',
        method: 'get',
    });
};
