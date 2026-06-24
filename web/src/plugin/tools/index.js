export default {
  name: 'tools',
  description: 'Tools 插件入口（环境配置、用户数据、粉丝/关注/好友、发送公屏消息）',
  version: '1.0.0',
  routes: [],
  apis: [
    () => import('@/plugin/tools/api/environment'),
    () => import('@/plugin/tools/api/userRelation'),
    () => import('@/plugin/tools/api/fanFollow'),
    () => import('@/plugin/tools/api/sendChat')
  ]
}