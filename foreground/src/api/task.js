import service from '../utils/request'

// 任务列表
export const getTaskList = (data) => {
    return service({
        url: "/sys/getTaskList",
        method: 'post',
        data: data
    })
}

//保存任务
export const saveTask = (data) => {
    return service({
        url: "/sys/saveTask",
        method: 'post',
        data: data
    })
}