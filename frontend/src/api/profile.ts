import request from '@/utils/request';

// 获取个人信息
export const getProfile = () => {
    return request({
        url: '/profile',
        method: 'get',
    });
};

// 更新个人资料
export const updateProfile = (data: { nickname: string; email: string; phone: string }) => {
    return request({
        url: '/profile',
        method: 'put',
        data,
    });
};

// 修改密码
export const changePassword = (data: { oldPassword: string; newPassword: string }) => {
    return request({
        url: '/profile/password',
        method: 'put',
        data,
    });
};
