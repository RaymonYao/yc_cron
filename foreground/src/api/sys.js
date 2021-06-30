import service from '../utils/request'

export const getSysInfo = (data) => {
    return service({
        url: "/sys/getSysInfo",
        method: 'get',
        params: data
    })
}
