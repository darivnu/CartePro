import { useLogout } from '../../auth/useAuth'

export function ClientHome() {
  const logout = useLogout() 

  return (
    <div className="p-4">
      Client area
      <button onClick={() => logout.mutate()}>Log out</button>
    </div>
  )
}
