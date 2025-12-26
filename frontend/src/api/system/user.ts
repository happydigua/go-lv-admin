import request from '@/utils/request';

// 获取用户列表
export const getUserList = (params: { page: number; pageSize: number; username?: string; phone?: string; status?: number }) => {
    return request({
        url: '/system/user/list',
        method: 'get',
        params,
    });
};

// 创建用户
export const createUser = (data: any) => {
    return request({
        url: '/system/user',
        method: 'post',
        data,
    });
};

// 更新用户
export const updateUser = (id: number, data: any) => {
    return request({
        url: `/system/user/${id}`,
        method: 'put',
        data,
    });
};

// 删除用户
export const deleteUser = (id: number) => {
    return request({
        url: `/system/user/${id}`,
        method: 'delete',
    });
};

// 重置密码
export const resetPassword = (id: number, password?: string) => {
    return request({
        url: `/system/user/${id}/reset-password`,
        method: 'put',
        data: { password },
    });
};

// 获取角色选项
export const getRoleOptions = () => {
    return request({
        url: '/system/user/role-options',
        method: 'get',
    });
};
