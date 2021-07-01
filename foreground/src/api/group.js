import service from '../utils/request'

// 分组列表
export const getGroupList = (data) => {
    console.log(data)
    return service({
        url: "/group/getGroupList",
        method: 'post',
        data: data
    })
}

//保存分组
export const saveGroup = (data) => {
    return service({
        url: "/group/saveGroup",
        method: 'post',
        data: data
    })
}

//删除分组
export const delGroup = (data) => {
    return service({
        url: "/group/delGroup",
        method: 'post',
        params: data
    })
}