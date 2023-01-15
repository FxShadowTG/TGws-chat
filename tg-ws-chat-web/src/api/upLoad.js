import request from "@/api";

// 上传图片
export function upImage(data) {
    return request({
        url: '/api/v1/upImage',
        method: 'post',
        params: data
    })
}
