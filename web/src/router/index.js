import Vue from 'vue'
import VueRouter from 'vue-router'
import {store} from "../store";
/**
 *  导入子模块路由
 *  !!! 路由name 不要用中划线命名, 菜单输出时用了-作为分隔符
 */
import sysRouter from "./modules/sys"
import groupRouter from "./modules/group"
import taskRouter from "./modules/task"
import logRouter from "./modules/log"
import workerRouter from "./modules/worker"

Vue.use(VueRouter)

const routes = [
    {
        path: '/',
        name: 'login',
        component: () => import('../views/sys/Login'), //异步执行，懒加载，建议这样写
        meta: {
            needLogin: false,
            title: '登录',
            menu: false
        }
    },
    sysRouter,
    groupRouter,
    taskRouter,
    logRouter,
    workerRouter,
    {
        path: '*',
        name: '404',
        meta: {
            title: 'page not found',
            menu: false
        },
        component: () => import('../views/error/Index')
    }
]

const router = new VueRouter({
    mode: 'history',
    routes
})

router.beforeEach(async (to, from, next) => {
    if (store.getters['user/expiresAt'] && store.getters['user/expiresAt'] <= (new Date()).getTime()) {
        await store.dispatch('user/Logout', false) //store.dispatch等价于this.$store.dispatch
    }
    const token = store.getters['user/token']
    document.title = process.env.VUE_APP_TITLE + (to.meta.title ? '-' + to.meta.title : '')
    if (to.meta.needLogin !== false) {  //需要登录
        if (token) {
            next()
        } else {
            next({
                name: "login",
                query: {redirect: to.fullPath}
            })
        }
    } else {    //无需登录
        if (token) {
            //replace: true只是一个设置信息，告诉VUE本次操作后，不能通过浏览器后退按钮，返回前一个路由
            next({name: 'sys', replace: true});
        } else {
            next()
        }
    }
});

export default router
