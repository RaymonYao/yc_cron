import service from '../utils/request'

// 任务列表
export const getTaskList = (data) => {
    return service({
        url: "/task/getTaskList",
        method: 'post',
        data: data
    })
}

//保存任务
export const saveTask = (data) => {
    return service({
        url: "/task/saveTask",
        method: 'post',
        data: data
    })
}

//删除任务
export const delTask = (data) => {
    return service({
        url: "/task/delTask",
        method: 'post',
        params: data
    })
}