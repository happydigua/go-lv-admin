import request from '@/utils/request';

// 获取角色列表
export const getRoleList = () => {
    return request({
        url: '/system/role/list',
        method: 'get',
    });
};

// 创建角色
export const createRole = (data: any) => {
    return request({
        url: '/system/role',
        method: 'post',
        data,
    });
};

// 更新角色
export const updateRole = (id: number, data: any) => {
    return request({
        url: `/system/role/${id}`,
        method: 'put',
        data,
    });
};

// 删除角色
export const deleteRole = (id: number) => {
    return request({
        url: `/system/role/${id}`,
        method: 'delete',
    });
};
