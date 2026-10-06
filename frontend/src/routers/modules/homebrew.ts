import { Layout } from '@/routers/constant';

const homebrewRouter = {
    sort: 8,
    path: '/homebrew',
    name: 'Homebrew-Menu',
    component: Layout,
    redirect: '/homebrew/index',
    meta: {
        title: 'menu.homebrew',
        icon: 'p-appstore',
        permission: 'homebrew_view',
        platformFeature: 'homebrew',
        hideInSidebar: true,
    },
    children: [
        {
            path: '/homebrew/index',
            name: 'Homebrew',
            component: () => import('@/views/homebrew/index.vue'),
            meta: {
                title: 'menu.homebrew',
                icon: 'p-appstore',
                permission: 'homebrew_view',
                platformFeature: 'homebrew',
                activeMenu: '/apps',
                parent: 'menu.app',
            },
        },
    ],
};

export default homebrewRouter;
