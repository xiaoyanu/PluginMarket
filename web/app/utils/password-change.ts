interface PasswordChangeCompletion {
  logout: () => void
  redirectToLogin: () => void | Promise<unknown>
}

export const completePasswordChange = async ({ logout, redirectToLogin }: PasswordChangeCompletion) => {
  logout()
  await redirectToLogin()
}
