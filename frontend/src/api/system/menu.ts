import request from '@/utils/request';

// 获取菜单列表
export const getMenuList = () => {
    return request({
        url: '/system/menu/list',
        method: 'get',
    });
};

// 创建菜单
export const createMenu = (data: any) => {
    return request({
        url: '/system/menu',
        method: 'post',
        data,
    });
};

// 更新菜单
export const updateMenu = (id: number, data: any) => {
    return request({
        url: `/system/menu/${id}`,
        method: 'put',
        data,
    });
};

// 删除菜单
export const deleteMenu = (id: number) => {
    return request({
        url: `/system/menu/${id}`,
        method: 'delete',
    });
};
