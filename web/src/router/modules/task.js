/** When your routing table is too long, you can split it into small modules **/
import Layout from '../../views/layout/Layout'

const taskRouter = {
    path: '/task',
    name: 'task',
    component: Layout,
    meta: {
        title: '任务管理',
        icon: 'el-icon-set-up'
    },
    children: [
        {
            path: "taskList",
            name: 'taskList',
            meta: {
                title: '任务列表'
            },
            component: () => import('../../views/task/TaskList')
        }
    ]
}
export default taskRouter
