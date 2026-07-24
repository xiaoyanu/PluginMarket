export default defineNuxtRouteMiddleware(async (to) => {
  if (to.path !== '/manage' && !to.path.startsWith('/manage/')) return

  const userStore = useUserStore()
  const redirectToLogin = () => navigateTo({
    path: '/user/auth',
    query: {
      mode: 'login',
      redirect: to.fullPath,
      notice: 'login-required'
    }
  })

  if (!userStore.token) return redirectToLogin()

  try {
    const res = await useApiFetch<{ code: number; msg: string; data?: any }>('/user/info', {
      suppressErrorMessage: true,
    })
    if (res.data) userStore.setUserInfo(res.data)
  } catch {
    userStore.logout()
    return redirectToLogin()
  }
})
