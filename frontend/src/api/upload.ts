import request from '@/utils/request';

// 上传图片
export const uploadImage = (file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    return request({
        url: '/upload/image',
        method: 'post',
        data: formData,
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    });
};

// 上传文件
export const uploadFile = (file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    return request({
        url: '/upload/file',
        method: 'post',
        data: formData,
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    });
};

// 删除文件
export const deleteFile = (key: string) => {
    return request({
        url: '/upload/file',
        method: 'delete',
        data: { key }
    });
};
