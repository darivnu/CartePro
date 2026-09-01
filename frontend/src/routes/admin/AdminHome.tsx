import { useLogout } from '../../auth/useAuth'

export function AdminHome() {
  const logout = useLogout()

  return (
    <div className="p-4">
      Admin area
      <button onClick={() => logout.mutate()}>Log out</button>
    </div>
  )
}
