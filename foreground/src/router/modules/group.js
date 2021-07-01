/** When your routing table is too long, you can split it into small modules **/
import Layout from '../../views/layout/Layout'

const groupRouter = {
    path: '/group',
    name: 'group',
    component: Layout,
    redirect: '/group/groupList',
    meta: {
        title: '分组管理',
        icon: 'el-icon-set-up'
    },
    children: [
        {
            path: "groupList",
            name: 'groupList',
            meta: {
                title: '分组列表'
            },
            component: () => import('../../views/group/GroupList')
        }
    ]
}
export default groupRouter
